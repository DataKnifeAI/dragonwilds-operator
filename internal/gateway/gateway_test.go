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

package gateway

import (
	"strings"
	"testing"

	dragonwildsv1alpha1 "github.com/DataKnifeAI/dragonwilds-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPlanDisabled(t *testing.T) {
	off := false
	st := Plan(&dragonwildsv1alpha1.DragonwildsServer{
		ObjectMeta: metav1.ObjectMeta{Name: "x"},
		Spec:       dragonwildsv1alpha1.DragonwildsServerSpec{Gateway: dragonwildsv1alpha1.GatewayConfig{Enabled: &off}},
	})
	if !st.Configured {
		t.Fatal("disabled gateway should report configured")
	}
}

func TestPlanEnabledPending(t *testing.T) {
	on := true
	st := Plan(&dragonwildsv1alpha1.DragonwildsServer{
		Spec: dragonwildsv1alpha1.DragonwildsServerSpec{
			Gateway: dragonwildsv1alpha1.GatewayConfig{Enabled: &on, Address: "203.0.113.10"},
		},
	})
	if st.Configured {
		t.Fatal("Phase 3 not implemented yet")
	}
	if !strings.Contains(st.Message, "Phase 3") {
		t.Fatalf("message = %q", st.Message)
	}
}
