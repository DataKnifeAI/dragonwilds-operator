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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
)

// DedicatedServerSettings is the [/Script/Dominion.DedicatedServerSettings] section.
// Passwords are included only when rendering a ConfigMap that is never logged;
// prefer env vars (RSDW_*) for the official image.
type DedicatedServerSettings struct {
	OwnerId       string
	ServerName    string
	DefaultWorld  string
	AdminPassword string
	WorldPassword string
}

// RenderDedicatedServerINI renders DedicatedServer.ini content.
func RenderDedicatedServerINI(s DedicatedServerSettings) string {
	var b strings.Builder
	b.WriteString("[/Script/Dominion.DedicatedServerSettings]\n")
	fmt.Fprintf(&b, "OwnerId=%s\n", s.OwnerId)
	fmt.Fprintf(&b, "ServerName=%s\n", s.ServerName)
	fmt.Fprintf(&b, "DefaultWorldName=%s\n", s.DefaultWorld)
	fmt.Fprintf(&b, "AdminPassword=%s\n", s.AdminPassword)
	fmt.Fprintf(&b, "WorldPassword=%s\n", s.WorldPassword)
	return b.String()
}

// SettingsFromSpec builds settings from the CR (passwords supplied separately).
func SettingsFromSpec(server *dragonwildsv1alpha1.DragonwildsServer, adminPassword, worldPassword string) DedicatedServerSettings {
	return DedicatedServerSettings{
		OwnerId:       strings.TrimSpace(server.Spec.OwnerId),
		ServerName:    ServerName(server.Spec, server.Name),
		DefaultWorld:  WorldName(server.Spec, server.Name),
		AdminPassword: adminPassword,
		WorldPassword: worldPassword,
	}
}

// ConfigHash returns a short hash of non-secret config used to annotate pods so
// ConfigMap/env changes only roll the Deployment when intended. Passwords are
// hashed so the annotation never contains plaintext.
func ConfigHash(s DedicatedServerSettings) string {
	sum := sha256.Sum256([]byte(RenderDedicatedServerINI(s)))
	return hex.EncodeToString(sum[:8])
}
