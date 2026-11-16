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
	"strings"
	"testing"
)

func TestRenderDedicatedServerINI(t *testing.T) {
	got := RenderDedicatedServerINI(DedicatedServerSettings{
		OwnerId:       "owner-1",
		ServerName:    "Ashenfall",
		DefaultWorld:  "Home",
		AdminPassword: "admin",
		WorldPassword: "join",
	})
	for _, want := range []string{
		"[/Script/Dominion.DedicatedServerSettings]",
		"OwnerId=owner-1",
		"ServerName=Ashenfall",
		"DefaultWorldName=Home",
		"AdminPassword=admin",
		"WorldPassword=join",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}

func TestConfigHashStable(t *testing.T) {
	s := DedicatedServerSettings{OwnerId: "a", ServerName: "b", DefaultWorld: "c"}
	h1 := ConfigHash(s)
	h2 := ConfigHash(s)
	if h1 != h2 {
		t.Fatal("hash not stable")
	}
	s2 := s
	s2.OwnerId = "other"
	if h1 == ConfigHash(s2) {
		t.Fatal("hash should change when OwnerId changes")
	}
}
