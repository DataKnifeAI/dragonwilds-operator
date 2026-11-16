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
	"fmt"
	"strings"

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
)

// ValidateSpec returns a user-facing error if the CR cannot be reconciled.
func ValidateSpec(spec dragonwildsv1alpha1.DragonwildsServerSpec) error {
	if strings.TrimSpace(spec.OwnerId) == "" {
		return fmt.Errorf("spec.ownerId is required (RuneScape: Dragonwilds Player ID)")
	}
	players := Players(spec)
	if players < 1 || players > 6 {
		return fmt.Errorf("spec.players must be between 1 and 6, got %d", players)
	}
	if !spec.GenerateSecrets && (spec.CredentialsSecretRef == nil || spec.CredentialsSecretRef.Name == "") {
		return fmt.Errorf("provide spec.credentialsSecretRef or set spec.generateSecrets=true")
	}
	if spec.GenerateSecrets && spec.CredentialsSecretRef != nil && spec.CredentialsSecretRef.Name != "" {
		return fmt.Errorf("specify only one of spec.generateSecrets or spec.credentialsSecretRef")
	}
	if gatewayEnabled(spec) && strings.TrimSpace(spec.Gateway.Address) == "" {
		return fmt.Errorf("spec.gateway.address is required when gateway is enabled")
	}
	return nil
}

func gatewayEnabled(spec dragonwildsv1alpha1.DragonwildsServerSpec) bool {
	if spec.Gateway.Enabled == nil {
		return true
	}
	return *spec.Gateway.Enabled
}

// Players returns the effective player count.
func Players(spec dragonwildsv1alpha1.DragonwildsServerSpec) int32 {
	if spec.Players == 0 {
		return dragonwildsv1alpha1.DefaultPlayers
	}
	return spec.Players
}

// GamePort returns the effective UDP game port.
func GamePort(spec dragonwildsv1alpha1.DragonwildsServerSpec) int32 {
	if spec.GamePort == 0 {
		return dragonwildsv1alpha1.DefaultGamePort
	}
	return spec.GamePort
}

// BeaconPort returns the effective UDP beacon port (default gamePort+1111).
func BeaconPort(spec dragonwildsv1alpha1.DragonwildsServerSpec) int32 {
	if spec.BeaconPort != 0 {
		return spec.BeaconPort
	}
	return GamePort(spec) + dragonwildsv1alpha1.DefaultBeaconOffset
}

// ServerName returns the effective server display name.
func ServerName(spec dragonwildsv1alpha1.DragonwildsServerSpec, fallback string) string {
	if strings.TrimSpace(spec.ServerName) != "" {
		return spec.ServerName
	}
	if fallback != "" {
		return fallback
	}
	return "dragonwilds"
}

// WorldName returns the effective world name.
func WorldName(spec dragonwildsv1alpha1.DragonwildsServerSpec, fallback string) string {
	if strings.TrimSpace(spec.WorldName) != "" {
		return spec.WorldName
	}
	if fallback != "" {
		return fallback
	}
	return "Dragonwilds World"
}

// ImageRef builds the container image reference from image + optional tag.
func ImageRef(spec dragonwildsv1alpha1.DragonwildsServerSpec) string {
	image := spec.Image
	if image == "" {
		image = dragonwildsv1alpha1.DefaultImage
	}
	if spec.ImageTag == "" {
		return image
	}
	// Allow full refs in Image (repo@sha256:… or repo:tag) without double-tagging.
	if strings.Contains(image, "@") || strings.Contains(image, ":") {
		return image
	}
	return image + ":" + spec.ImageTag
}

// MountPath returns the PVC mount path.
func MountPath(spec dragonwildsv1alpha1.DragonwildsServerSpec) string {
	if spec.Persistence.MountPath != "" {
		return spec.Persistence.MountPath
	}
	return "/home/steam/rsdw-dedicated"
}

// PersistenceSize returns the PVC size string.
func PersistenceSize(spec dragonwildsv1alpha1.DragonwildsServerSpec) string {
	if spec.Persistence.Size != "" {
		return spec.Persistence.Size
	}
	return dragonwildsv1alpha1.DefaultPersistenceSize
}

// GraceSeconds returns terminationGracePeriodSeconds.
func GraceSeconds(spec dragonwildsv1alpha1.DragonwildsServerSpec) int64 {
	if spec.TerminationGracePeriodSeconds != nil {
		return *spec.TerminationGracePeriodSeconds
	}
	return dragonwildsv1alpha1.DefaultGraceSeconds
}

// GatewayEnabled reports whether Envoy Gateway resources should be reconciled.
func GatewayEnabled(spec dragonwildsv1alpha1.DragonwildsServerSpec) bool {
	return gatewayEnabled(spec)
}
