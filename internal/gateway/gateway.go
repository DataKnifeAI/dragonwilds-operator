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

// Package gateway will own Envoy Gateway + UDPRoute wiring (Phase 3).
package gateway

import (
	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
	"github.com/DataKnifeAI/dragonwilds-operator/internal/config"
)

// Status describes gateway reconcile outcome for status conditions.
type Status struct {
	Configured bool
	Address    string
	Message    string
}

// Plan returns the intended gateway posture without creating API objects yet.
// Phase 3 will create Gateway / UDPRoute / EnvoyProxy resources here.
func Plan(server *dragonwildsv1alpha1.DragonwildsServer) Status {
	if !config.GatewayEnabled(server.Spec) {
		return Status{
			Configured: true,
			Address:    "",
			Message:    "gateway disabled; using Service only (set gateway.serviceType=NodePort for homelab)",
		}
	}
	if server.Spec.Gateway.Address == "" {
		return Status{
			Configured: false,
			Message:    "spec.gateway.address is required when gateway.enabled=true",
		}
	}
	className := server.Spec.Gateway.ClassName
	if className == "" {
		className = "envoy"
	}
	return Status{
		Configured: false,
		Address:    server.Spec.Gateway.Address,
		Message:    "Envoy Gateway UDPRoute wiring lands in Phase 3 (class=" + className + ")",
	}
}
