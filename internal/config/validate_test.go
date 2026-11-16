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

package config

import (
	"testing"

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestValidateSpec(t *testing.T) {
	const testOwnerID = "abc123"
	enabled := true
	disabled := false

	tests := []struct {
		name    string
		spec    dragonwildsv1alpha1.DragonwildsServerSpec
		wantErr bool
	}{
		{
			name: "valid generateSecrets",
			spec: dragonwildsv1alpha1.DragonwildsServerSpec{
				OwnerId:         testOwnerID,
				GenerateSecrets: true,
				Gateway:         dragonwildsv1alpha1.GatewayConfig{Enabled: &disabled},
			},
		},
		{
			name: "missing ownerId",
			spec: dragonwildsv1alpha1.DragonwildsServerSpec{
				GenerateSecrets: true,
				Gateway:         dragonwildsv1alpha1.GatewayConfig{Enabled: &disabled},
			},
			wantErr: true,
		},
		{
			name: "missing credentials",
			spec: dragonwildsv1alpha1.DragonwildsServerSpec{
				OwnerId: testOwnerID,
				Gateway: dragonwildsv1alpha1.GatewayConfig{Enabled: &disabled},
			},
			wantErr: true,
		},
		{
			name: "both credential sources",
			spec: dragonwildsv1alpha1.DragonwildsServerSpec{
				OwnerId:              testOwnerID,
				GenerateSecrets:      true,
				CredentialsSecretRef: &corev1.LocalObjectReference{Name: "creds"},
				Gateway:              dragonwildsv1alpha1.GatewayConfig{Enabled: &disabled},
			},
			wantErr: true,
		},
		{
			name: "gateway enabled without address",
			spec: dragonwildsv1alpha1.DragonwildsServerSpec{
				OwnerId:         testOwnerID,
				GenerateSecrets: true,
				Gateway:         dragonwildsv1alpha1.GatewayConfig{Enabled: &enabled},
			},
			wantErr: true,
		},
		{
			name: "players too high",
			spec: dragonwildsv1alpha1.DragonwildsServerSpec{
				OwnerId:         testOwnerID,
				Players:         7,
				GenerateSecrets: true,
				Gateway:         dragonwildsv1alpha1.GatewayConfig{Enabled: &disabled},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSpec(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBeaconPortDefault(t *testing.T) {
	if got := BeaconPort(dragonwildsv1alpha1.DragonwildsServerSpec{}); got != 8888 {
		t.Fatalf("default beacon = %d, want 8888", got)
	}
	if got := BeaconPort(dragonwildsv1alpha1.DragonwildsServerSpec{GamePort: 7778}); got != 8889 {
		t.Fatalf("beacon for 7778 = %d, want 8889", got)
	}
}

func TestImageRef(t *testing.T) {
	if got := ImageRef(dragonwildsv1alpha1.DragonwildsServerSpec{}); got != dragonwildsv1alpha1.DefaultImage {
		t.Fatalf("default image = %q", got)
	}
	if got := ImageRef(dragonwildsv1alpha1.DragonwildsServerSpec{ImageTag: "v1"}); got != dragonwildsv1alpha1.DefaultImage+":v1" {
		t.Fatalf("tagged image = %q", got)
	}
}

func TestDefaultResources(t *testing.T) {
	res := DefaultResources(dragonwildsv1alpha1.DragonwildsServerSpec{Players: 6})
	mem := res.Requests[corev1.ResourceMemory]
	want := resource.MustParse("8Gi")
	if mem.Cmp(want) != 0 {
		t.Fatalf("6-player memory request = %s, want 8Gi", mem.String())
	}
}
