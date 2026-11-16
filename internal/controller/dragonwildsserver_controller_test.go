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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
)

var _ = Describe("DragonwildsServer Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "ashenfall-home"

		ctx := context.Background()

		const resourceNamespace = "default"

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}

		BeforeEach(func() {
			By("creating the custom resource for the Kind DragonwildsServer")
			existing := &dragonwildsv1alpha1.DragonwildsServer{}
			err := k8sClient.Get(ctx, typeNamespacedName, existing)
			if err != nil && errors.IsNotFound(err) {
				gatewayOff := false
				resource := &dragonwildsv1alpha1.DragonwildsServer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: resourceNamespace,
					},
					Spec: dragonwildsv1alpha1.DragonwildsServerSpec{
						OwnerId:                 "test-owner-id",
						ServerName:              "Ashenfall Home",
						WorldName:               "Homelab World",
						Players:                 4,
						GenerateSecrets:         true,
						AllowEmptyWorldPassword: true,
						Gateway: dragonwildsv1alpha1.GatewayConfig{
							Enabled:     &gatewayOff,
							ServiceType: corev1.ServiceTypeClusterIP,
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &dragonwildsv1alpha1.DragonwildsServer{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())
			By("Cleanup the specific resource instance DragonwildsServer")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should create owned child resources", func() {
			controllerReconciler := &DragonwildsServerReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			// Finalizer add requeues once.
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())
			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())

			secret := &corev1.Secret{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: resourceName + "-secrets", Namespace: resourceNamespace}, secret)).To(Succeed())
			Expect(secret.Data).To(HaveKey(dragonwildsv1alpha1.SecretKeyAdminPassword))

			cm := &corev1.ConfigMap{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: resourceName + "-config", Namespace: resourceNamespace}, cm)).To(Succeed())
			Expect(cm.Data).To(HaveKey("DedicatedServer.ini"))

			pvc := &corev1.PersistentVolumeClaim{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: resourceName + "-data", Namespace: resourceNamespace}, pvc)).To(Succeed())

			svc := &corev1.Service{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: resourceNamespace}, svc)).To(Succeed())
			Expect(svc.Spec.Ports).To(HaveLen(2))

			deploy := &appsv1.Deployment{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: resourceNamespace}, deploy)).To(Succeed())
			Expect(deploy.Spec.Strategy.Type).To(Equal(appsv1.RecreateDeploymentStrategyType))

			updated := &dragonwildsv1alpha1.DragonwildsServer{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, updated)).To(Succeed())
			Expect(updated.Status.Phase).To(Equal(dragonwildsv1alpha1.PhaseDeploying))
			Expect(updated.Status.CredentialsSecretName).To(Equal(resourceName + "-secrets"))
		})
	})
})
