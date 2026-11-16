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

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// DefaultResources returns CPU/memory for the player count.
// Memory request: max(4Gi, 2Gi + 1Gi×players) per official minimum + plan formula.
func DefaultResources(spec dragonwildsv1alpha1.DragonwildsServerSpec) corev1.ResourceRequirements {
	if spec.Resources != nil {
		return *spec.Resources
	}
	players := Players(spec)
	memGi := 2 + int(players)
	if memGi < 4 {
		memGi = 4
	}
	limitGi := memGi + 2
	if limitGi < 6 {
		limitGi = 6
	}
	cpuReq := "1"
	cpuLim := "2"
	if players >= 4 {
		cpuReq = "2"
		cpuLim = "4"
	}
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuReq),
			corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dGi", memGi)),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLim),
			corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dGi", limitGi)),
		},
	}
}
