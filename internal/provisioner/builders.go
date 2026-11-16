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

package provisioner

import (
	"strconv"

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
	"github.com/DataKnifeAI/dragonwilds-operator/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// BuildConfigMap renders DedicatedServer.ini (passwords included; never log the object).
func BuildConfigMap(server *dragonwildsv1alpha1.DragonwildsServer, adminPassword, worldPassword string) *corev1.ConfigMap {
	settings := config.SettingsFromSpec(server, adminPassword, worldPassword)
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(server),
			Namespace: server.Namespace,
			Labels:    labelsFor(server),
			Annotations: map[string]string{
				AnnotationConfigHash: config.ConfigHash(settings),
			},
		},
		Data: map[string]string{
			"DedicatedServer.ini": config.RenderDedicatedServerINI(settings),
		},
	}
}

// BuildPVC creates the data volume for /home/steam/rsdw-dedicated.
func BuildPVC(server *dragonwildsv1alpha1.DragonwildsServer) *corev1.PersistentVolumeClaim {
	size := config.PersistenceSize(server.Spec)
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName(server),
			Namespace: server.Namespace,
			Labels:    labelsFor(server),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(size),
				},
			},
		},
	}
	if server.Spec.Persistence.StorageClassName != "" {
		sc := server.Spec.Persistence.StorageClassName
		pvc.Spec.StorageClassName = &sc
	}
	return pvc
}

// BuildService exposes game + beacon UDP ports.
//
// When Envoy Gateway is disabled, ServiceType/Address pin exposure. Dragonwilds
// (Unreal) breaks behind Envoy's UDP proxy because joins appear to come from the
// Envoy pod IP — prefer LoadBalancer+VIP (or NodePort) straight to the game pod.
func BuildService(server *dragonwildsv1alpha1.DragonwildsServer) *corev1.Service {
	gamePort := config.GamePort(server.Spec)
	beaconPort := config.BeaconPort(server.Spec)
	svcType := corev1.ServiceTypeClusterIP
	if !config.GatewayEnabled(server.Spec) && server.Spec.Gateway.ServiceType != "" {
		svcType = server.Spec.Gateway.ServiceType
	}

	spec := corev1.ServiceSpec{
		Type:     svcType,
		Selector: labelsFor(server),
		Ports: []corev1.ServicePort{
			{
				Name:       "game",
				Protocol:   corev1.ProtocolUDP,
				Port:       gamePort,
				TargetPort: intstr.FromInt32(gamePort),
			},
			{
				Name:       "beacon",
				Protocol:   corev1.ProtocolUDP,
				Port:       beaconPort,
				TargetPort: intstr.FromInt32(beaconPort),
			},
		},
	}
	if !config.GatewayEnabled(server.Spec) && svcType == corev1.ServiceTypeLoadBalancer {
		if addr := server.Spec.Gateway.Address; addr != "" {
			spec.LoadBalancerIP = addr
		}
		// Cluster: kube-vip may park the VIP on a node without the game pod.
		// Local would silently drop UDP there (VIP ping still works).
		spec.ExternalTrafficPolicy = corev1.ServiceExternalTrafficPolicyCluster
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(server),
			Namespace: server.Namespace,
			Labels:    labelsFor(server),
		},
		Spec: spec,
	}
}

// BuildDeployment creates the Recreate Deployment for the official image.
func BuildDeployment(server *dragonwildsv1alpha1.DragonwildsServer, configHash, secretName string) *appsv1.Deployment {
	gamePort := config.GamePort(server.Spec)
	beaconPort := config.BeaconPort(server.Spec)
	replicas := int32(1)
	grace := config.GraceSeconds(server.Spec)
	labels := labelsFor(server)
	image := config.ImageRef(server.Spec)
	pullPolicy := server.Spec.ImagePullPolicy
	if pullPolicy == "" {
		pullPolicy = corev1.PullIfNotPresent
	}

	autoStop := "false"
	if server.Spec.Update.AutoUpdateImage {
		autoStop = "true"
	}
	validate := "0"
	if server.Spec.Update.ValidateOnStart {
		validate = "1"
	}

	env := []corev1.EnvVar{
		{Name: "RSDW_OWNER_ID", Value: server.Spec.OwnerId},
		{Name: "RSDW_SERVER_NAME", Value: config.ServerName(server.Spec, server.Name)},
		{Name: "RSDW_WORLD_NAME", Value: config.WorldName(server.Spec, server.Name)},
		{Name: "RSDW_PORT", Value: strconv.FormatInt(int64(gamePort), 10)},
		{Name: "RSDW_BEACON_PORT", Value: strconv.FormatInt(int64(beaconPort), 10)},
		{Name: "RSDW_AUTO_STOP_ON_UPDATE", Value: autoStop},
		{Name: "STEAMAPPVALIDATE", Value: validate},
		{Name: "RSDW_LINE_BUFFERED_LOGS", Value: "true"},
		{
			Name: "RSDW_ADMIN_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
					Key:                  dragonwildsv1alpha1.SecretKeyAdminPassword,
				},
			},
		},
		{
			Name: "RSDW_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
					Key:                  dragonwildsv1alpha1.SecretKeyWorldPassword,
					Optional:             boolPtr(true),
				},
			},
		},
	}
	if server.Spec.Admins != "" {
		env = append(env, corev1.EnvVar{Name: "RSDW_ADMINS", Value: server.Spec.Admins})
	}
	if server.Spec.AdditionalArgs != "" {
		env = append(env, corev1.EnvVar{Name: "RSDW_ADDITIONAL_ARGS", Value: server.Spec.AdditionalArgs})
	}

	mountPath := config.MountPath(server.Spec)
	podSpec := corev1.PodSpec{
		TerminationGracePeriodSeconds: &grace,
		ImagePullSecrets:              server.Spec.ImagePullSecrets,
		NodeSelector:                  server.Spec.NodeSelector,
		Containers: []corev1.Container{
			{
				Name:            "dragonwilds",
				Image:           image,
				ImagePullPolicy: pullPolicy,
				Env:             env,
				Ports: []corev1.ContainerPort{
					{Name: "game", ContainerPort: gamePort, Protocol: corev1.ProtocolUDP},
					{Name: "beacon", ContainerPort: beaconPort, Protocol: corev1.ProtocolUDP},
				},
				Resources: config.DefaultResources(server.Spec),
				VolumeMounts: []corev1.VolumeMount{
					{Name: "data", MountPath: mountPath},
				},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: "data",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName(server),
					},
				},
			},
		},
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(server),
			Namespace: server.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						AnnotationConfigHash: configHash,
					},
				},
				Spec: podSpec,
			},
		},
	}
}

func boolPtr(v bool) *bool { return &v }

// SecretNameFor returns the credentials Secret name for a server.
func SecretNameFor(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return secretName(server)
}

// ConfigMapNameFor returns the ConfigMap name.
func ConfigMapNameFor(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return configMapName(server)
}

// PVCNameFor returns the PVC name.
func PVCNameFor(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return pvcName(server)
}

// DeploymentNameFor returns the Deployment name.
func DeploymentNameFor(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return deploymentName(server)
}

// ServiceNameFor returns the Service name.
func ServiceNameFor(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return serviceName(server)
}
