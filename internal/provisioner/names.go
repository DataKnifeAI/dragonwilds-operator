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

import dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"

const (
	LabelAppName      = "app.kubernetes.io/name"
	LabelAppInstance  = "app.kubernetes.io/instance"
	LabelAppManagedBy = "app.kubernetes.io/managed-by"
	LabelAppComponent = "app.kubernetes.io/component"

	ManagedBy            = "dragonwilds-operator"
	ComponentGame        = "game-server"
	AnnotationConfigHash = "dragonwilds.dataknife.ai/config-hash"
)

func labelsFor(server *dragonwildsv1alpha1.DragonwildsServer) map[string]string {
	return map[string]string{
		LabelAppName:      "dragonwilds-server",
		LabelAppInstance:  server.Name,
		LabelAppManagedBy: ManagedBy,
		LabelAppComponent: ComponentGame,
	}
}

func secretName(server *dragonwildsv1alpha1.DragonwildsServer) string {
	if server.Spec.CredentialsSecretRef != nil && server.Spec.CredentialsSecretRef.Name != "" {
		return server.Spec.CredentialsSecretRef.Name
	}
	return server.Name + "-secrets"
}

func configMapName(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return server.Name + "-config"
}

func pvcName(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return server.Name + "-data"
}

func deploymentName(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return server.Name
}

func serviceName(server *dragonwildsv1alpha1.DragonwildsServer) string {
	return server.Name
}
