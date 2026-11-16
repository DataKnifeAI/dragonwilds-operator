/*
Copyright 2026 DataKnifeAI.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
	"github.com/DataKnifeAI/dragonwilds-operator/internal/config"
	"github.com/DataKnifeAI/dragonwilds-operator/internal/gateway"
	"github.com/DataKnifeAI/dragonwilds-operator/internal/provisioner"
)

const finalizerName = "dragonwilds.dataknife.ai/finalizer"

// DragonwildsServerReconciler reconciles a DragonwildsServer object.
type DragonwildsServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dragonwilds.dataknife.ai,resources=dragonwildsservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dragonwilds.dataknife.ai,resources=dragonwildsservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dragonwilds.dataknife.ai,resources=dragonwildsservers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *DragonwildsServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	server := &dragonwildsv1alpha1.DragonwildsServer{}
	if err := r.Get(ctx, req.NamespacedName, server); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(server, finalizerName) {
		controllerutil.AddFinalizer(server, finalizerName)
		if err := r.Update(ctx, server); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	if !server.DeletionTimestamp.IsZero() {
		controllerutil.RemoveFinalizer(server, finalizerName)
		if err := r.Update(ctx, server); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := config.ValidateSpec(server.Spec); err != nil {
		_ = r.patchStatus(ctx, server, func(s *dragonwildsv1alpha1.DragonwildsServer) {
			s.Status.Phase = dragonwildsv1alpha1.PhaseFailed
			s.Status.Ready = false
			s.Status.Message = err.Error()
			meta.SetStatusCondition(&s.Status.Conditions, metav1.Condition{
				Type:               dragonwildsv1alpha1.ConditionReady,
				Status:             metav1.ConditionFalse,
				Reason:             "ValidationFailed",
				Message:            err.Error(),
				ObservedGeneration: s.Generation,
			})
		})
		log.Info("validation failed", "error", err)
		return ctrl.Result{}, nil
	}

	secret, generated, err := r.ensureSecret(ctx, server)
	if err != nil {
		return ctrl.Result{}, err
	}
	adminPassword, worldPassword, err := provisioner.ReadPasswords(secret)
	if err != nil {
		_ = r.patchStatus(ctx, server, func(s *dragonwildsv1alpha1.DragonwildsServer) {
			s.Status.Phase = dragonwildsv1alpha1.PhaseFailed
			s.Status.Ready = false
			s.Status.Message = err.Error()
			meta.SetStatusCondition(&s.Status.Conditions, metav1.Condition{
				Type:               dragonwildsv1alpha1.ConditionCredentialsValid,
				Status:             metav1.ConditionFalse,
				Reason:             "CredentialsInvalid",
				Message:            err.Error(),
				ObservedGeneration: s.Generation,
			})
		})
		return ctrl.Result{}, nil
	}

	cm := provisioner.BuildConfigMap(server, adminPassword, worldPassword)
	if err := controllerutil.SetControllerReference(server, cm, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.applyObject(ctx, cm); err != nil {
		return ctrl.Result{}, err
	}

	pvc := provisioner.BuildPVC(server)
	if err := controllerutil.SetControllerReference(server, pvc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.applyObject(ctx, pvc); err != nil {
		return ctrl.Result{}, err
	}

	svc := provisioner.BuildService(server)
	if err := controllerutil.SetControllerReference(server, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.applyObject(ctx, svc); err != nil {
		return ctrl.Result{}, err
	}

	configHash := cm.Annotations[provisioner.AnnotationConfigHash]
	deploy := provisioner.BuildDeployment(server, configHash, secret.Name)
	if err := controllerutil.SetControllerReference(server, deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.applyObject(ctx, deploy); err != nil {
		return ctrl.Result{}, err
	}

	gw := gateway.Plan(server)

	ready := false
	phase := dragonwildsv1alpha1.PhaseDeploying
	message := "reconciling child resources"
	current := &appsv1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: server.Namespace, Name: deploy.Name}, current); err == nil {
		if current.Status.ReadyReplicas >= 1 {
			ready = true
			phase = dragonwildsv1alpha1.PhaseRunning
			message = "game server Deployment is ready"
		}
	}

	connAddr := server.Spec.Gateway.Address
	if !config.GatewayEnabled(server.Spec) {
		connAddr = ""
		message = message + "; gateway disabled — use Service/NodePort (see docs/CONNECT.md)"
	}

	_ = r.patchStatus(ctx, server, func(s *dragonwildsv1alpha1.DragonwildsServer) {
		s.Status.Phase = phase
		s.Status.Ready = ready
		s.Status.Message = message
		s.Status.ConnectionAddress = connAddr
		s.Status.ConnectionPort = config.GamePort(server.Spec)
		s.Status.BeaconPort = config.BeaconPort(server.Spec)
		s.Status.CredentialsSecretName = secret.Name
		s.Status.CredentialsGenerated = generated
		s.Status.ImageVersion = config.ImageRef(server.Spec)
		s.Status.WorldName = config.WorldName(server.Spec, server.Name)
		s.Status.ObservedGeneration = server.Generation

		meta.SetStatusCondition(&s.Status.Conditions, metav1.Condition{
			Type:               dragonwildsv1alpha1.ConditionCredentialsValid,
			Status:             metav1.ConditionTrue,
			Reason:             "CredentialsReady",
			Message:            fmt.Sprintf("using secret %s", secret.Name),
			ObservedGeneration: server.Generation,
		})
		meta.SetStatusCondition(&s.Status.Conditions, metav1.Condition{
			Type:               dragonwildsv1alpha1.ConditionReady,
			Status:             boolCond(ready),
			Reason:             readyReason(ready),
			Message:            message,
			ObservedGeneration: server.Generation,
		})
		gwStatus := metav1.ConditionFalse
		gwReason := "GatewayPending"
		if gw.Configured {
			gwStatus = metav1.ConditionTrue
			gwReason = "GatewayConfigured"
		}
		meta.SetStatusCondition(&s.Status.Conditions, metav1.Condition{
			Type:               dragonwildsv1alpha1.ConditionGatewayConfigured,
			Status:             gwStatus,
			Reason:             gwReason,
			Message:            gw.Message,
			ObservedGeneration: server.Generation,
		})
	})

	return ctrl.Result{}, nil
}

func (r *DragonwildsServerReconciler) ensureSecret(ctx context.Context, server *dragonwildsv1alpha1.DragonwildsServer) (*corev1.Secret, bool, error) {
	name := provisioner.SecretNameFor(server)
	existing := &corev1.Secret{}
	err := r.Get(ctx, client.ObjectKey{Namespace: server.Namespace, Name: name}, existing)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, false, err
	}
	found := err == nil

	if !server.Spec.GenerateSecrets {
		if !found {
			return nil, false, fmt.Errorf("credentials secret %q not found", name)
		}
		return existing, false, nil
	}

	var base *corev1.Secret
	if found {
		base = existing
	}
	desired, err := provisioner.DesiredSecret(server, base)
	if err != nil {
		return nil, false, err
	}
	if err := controllerutil.SetControllerReference(server, desired, r.Scheme); err != nil {
		return nil, false, err
	}
	if !found {
		if err := r.Create(ctx, desired); err != nil {
			return nil, false, err
		}
		return desired, true, nil
	}
	// Preserve resourceVersion for update when keys were filled in.
	desired.ResourceVersion = existing.ResourceVersion
	if err := r.Update(ctx, desired); err != nil {
		return nil, false, err
	}
	return desired, true, nil
}

func (r *DragonwildsServerReconciler) applyObject(ctx context.Context, obj client.Object) error {
	key := client.ObjectKeyFromObject(obj)
	existing := obj.DeepCopyObject().(client.Object)
	err := r.Get(ctx, key, existing)
	if apierrors.IsNotFound(err) {
		return r.Create(ctx, obj)
	}
	if err != nil {
		return err
	}

	switch desired := obj.(type) {
	case *corev1.ConfigMap:
		cur := existing.(*corev1.ConfigMap)
		cur.Data = desired.Data
		cur.Labels = desired.Labels
		cur.Annotations = desired.Annotations
		return r.Update(ctx, cur)
	case *corev1.Service:
		cur := existing.(*corev1.Service)
		cur.Spec.Type = desired.Spec.Type
		cur.Spec.Ports = desired.Spec.Ports
		cur.Spec.Selector = desired.Spec.Selector
		cur.Labels = desired.Labels
		return r.Update(ctx, cur)
	case *appsv1.Deployment:
		cur := existing.(*appsv1.Deployment)
		cur.Spec = desired.Spec
		cur.Labels = desired.Labels
		return r.Update(ctx, cur)
	case *corev1.PersistentVolumeClaim:
		// PVC spec is mostly immutable; only ensure labels.
		cur := existing.(*corev1.PersistentVolumeClaim)
		cur.Labels = desired.Labels
		return r.Update(ctx, cur)
	default:
		return fmt.Errorf("unsupported apply type %T", obj)
	}
}

func (r *DragonwildsServerReconciler) patchStatus(ctx context.Context, server *dragonwildsv1alpha1.DragonwildsServer, mutate func(*dragonwildsv1alpha1.DragonwildsServer)) error {
	latest := &dragonwildsv1alpha1.DragonwildsServer{}
	if err := r.Get(ctx, client.ObjectKeyFromObject(server), latest); err != nil {
		return err
	}
	mutate(latest)
	return r.Status().Update(ctx, latest)
}

func boolCond(ready bool) metav1.ConditionStatus {
	if ready {
		return metav1.ConditionTrue
	}
	return metav1.ConditionFalse
}

func readyReason(ready bool) string {
	if ready {
		return "DeploymentReady"
	}
	return "DeploymentNotReady"
}

// SetupWithManager sets up the controller with the Manager.
func (r *DragonwildsServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dragonwildsv1alpha1.DragonwildsServer{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Named("dragonwildsserver").
		Complete(r)
}
