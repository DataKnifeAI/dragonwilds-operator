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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PhasePending         = "Pending"
	PhaseDeploying       = "Deploying"
	PhaseRunning         = "Running"
	PhaseVersionMismatch = "VersionMismatch"
	PhaseFailed          = "Failed"

	ConditionReady             = "Ready"
	ConditionCredentialsValid  = "CredentialsValid"
	ConditionWorldLoaded       = "WorldLoaded"
	ConditionGatewayConfigured = "GatewayConfigured"

	DefaultImage           = "ghcr.io/runescape/rsdw-dedicated"
	DefaultGamePort        = int32(7777)
	DefaultBeaconOffset    = int32(1111)
	DefaultPlayers         = int32(6)
	DefaultPersistenceSize = "20Gi"
	DefaultGraceSeconds    = int64(60)

	SecretKeyAdminPassword = "admin-password"
	SecretKeyWorldPassword = "world-password"
)

// GatewayConfig configures exposure for Dragonwilds game traffic.
type GatewayConfig struct {
	// Enabled creates Gateway / UDPRoute resources (Envoy Gateway pattern).
	// Default false — Dragonwilds (Unreal) joins fail behind Envoy's UDP proxy
	// because the game sees the Envoy pod IP. Prefer serviceType=LoadBalancer
	// with address set to a kube-vip VIP instead.
	// +kubebuilder:default=false
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Address is the dedicated external IP (Kube-VIP / MetalLB).
	// Required when enabled is true; also used as LoadBalancerIP when
	// enabled is false and serviceType is LoadBalancer.
	// +optional
	// +kubebuilder:validation:Pattern=`^$|^([0-9]{1,3}\.){3}[0-9]{1,3}$`
	Address string `json:"address,omitempty"`

	// ClassName is the GatewayClass name.
	// +kubebuilder:default="envoy"
	// +optional
	ClassName string `json:"className,omitempty"`

	// GatewayName overrides the Gateway resource name.
	// Default: {name}-gateway (strips a trailing "-server" suffix from name).
	// +optional
	GatewayName string `json:"gatewayName,omitempty"`

	// EnvoyProxyName overrides the EnvoyProxy resource name.
	// Default: game-{base}-kubevip.
	// +optional
	EnvoyProxyName string `json:"envoyProxyName,omitempty"`

	// ServiceType overrides the backing Service type when gateway is disabled.
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	// +kubebuilder:default=ClusterIP
	// +optional
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
}

// UpdateConfig controls opt-in update behaviour for the official image.
// The Jagex image runs SteamCMD on each container start; autoUpdateImage maps
// to RSDW_AUTO_STOP_ON_UPDATE so a detected Steam build stop+restarts the pod.
type UpdateConfig struct {
	// AutoUpdateImage when true sets RSDW_AUTO_STOP_ON_UPDATE=true so the
	// container exits when Steam publishes a newer build (orchestrator restarts
	// and SteamCMD updates into the PVC). Default false — opt-in only.
	// +optional
	AutoUpdateImage bool `json:"autoUpdateImage,omitempty"`

	// ValidateOnStart when true sets STEAMAPPVALIDATE=1 on the game container.
	// +optional
	ValidateOnStart bool `json:"validateOnStart,omitempty"`
}

// PersistenceConfig configures the world/save PVC.
type PersistenceConfig struct {
	// Size is the PVC capacity. Default 20Gi (official minimum disk guidance).
	// +kubebuilder:default="20Gi"
	// +optional
	Size string `json:"size,omitempty"`

	// StorageClassName selects the StorageClass. Empty uses the cluster default.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// MountPath is the container path for the PVC. Official image default:
	// /home/steam/rsdw-dedicated
	// +kubebuilder:default="/home/steam/rsdw-dedicated"
	// +optional
	MountPath string `json:"mountPath,omitempty"`
}

// ServerManagerConfig is an optional authenticated admin UI sidecar on the
// game pod. When enabled, TCP port (default 8088) is published on the same
// Service as game UDP (LoadBalancer VIP:port). Basic auth uses username
// "admin" and the admin-password Secret key.
type ServerManagerConfig struct {
	// Enabled starts the Server Manager sidecar and exposes its TCP port.
	// Default false. Enabling Recreate-rolls the game Deployment.
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// Port is the sidecar listen port (and Service TCP port). Default 8088.
	// +kubebuilder:default=8088
	// +kubebuilder:validation:Minimum=1024
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port int32 `json:"port,omitempty"`

	// Image is the container image that provides /server-manager (the operator
	// image). Default: harbor.dataknife.net/library/dragonwilds-operator:latest
	// +optional
	Image string `json:"image,omitempty"`
}

// BackupConfig configures scheduled tar backups of the data PVC.
type BackupConfig struct {
	// Enabled creates a backup PVC + CronJob. Default false.
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// Schedule is a Cron expression (UTC). Default "0 4 * * *" (04:00 daily).
	// +kubebuilder:default="0 4 * * *"
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// Retain is how many compressed archives to keep. Default 7.
	// +kubebuilder:default=7
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=90
	// +optional
	Retain int32 `json:"retain,omitempty"`

	// Size is the backup PVC capacity. Default matches persistence size.
	// +optional
	Size string `json:"size,omitempty"`

	// StorageClassName for the backup PVC. Empty inherits persistence.storageClassName.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// Image is the backup job container image. Default busybox:1.36.
	// +kubebuilder:default="busybox:1.36"
	// +optional
	Image string `json:"image,omitempty"`
}

// DragonwildsServerSpec defines the desired state of a Dragonwilds dedicated server.
// Configuration maps to the official image env vars (RSDW_*) and optionally a
// rendered DedicatedServer.ini ConfigMap for observability / SteamCMD parity.
type DragonwildsServerSpec struct {
	// OwnerId is the RuneScape: Dragonwilds Player ID (EOS Online ID).
	// Required — the server will not start without it.
	// +kubebuilder:validation:MinLength=1
	OwnerId string `json:"ownerId"`

	// ServerName is the display name for the dedicated server (RSDW_SERVER_NAME).
	// +optional
	ServerName string `json:"serverName,omitempty"`

	// WorldName is the visible world name (RSDW_WORLD_NAME / DefaultWorldName).
	// +optional
	WorldName string `json:"worldName,omitempty"`

	// Players is the intended player cap (1–6). Drives default memory requests:
	// max(4Gi, 2Gi + 1Gi×players).
	// +kubebuilder:default=6
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=6
	// +optional
	Players int32 `json:"players,omitempty"`

	// GamePort is the primary UDP game port (RSDW_PORT). Host and container
	// ports must match for the official image.
	// +kubebuilder:default=7777
	// +kubebuilder:validation:Minimum=1024
	// +kubebuilder:validation:Maximum=65535
	// +optional
	GamePort int32 `json:"gamePort,omitempty"`

	// BeaconPort is the UDP beacon port (RSDW_BEACON_PORT). Default is
	// gamePort + 1111 (7777 → 8888).
	// +kubebuilder:validation:Minimum=1024
	// +kubebuilder:validation:Maximum=65535
	// +optional
	BeaconPort int32 `json:"beaconPort,omitempty"`

	// Admins is a comma-separated list of additional owner/admin user IDs
	// (RSDW_ADMINS).
	// +optional
	Admins string `json:"admins,omitempty"`

	// AdditionalArgs are extra CLI args for RSDragonwildsServer.sh
	// (RSDW_ADDITIONAL_ARGS). Do not pass -FORCELOGFLUSH.
	// +optional
	AdditionalArgs string `json:"additionalArgs,omitempty"`

	// Image is the dedicated server container image repository.
	// +kubebuilder:default="ghcr.io/runescape/rsdw-dedicated"
	// +optional
	Image string `json:"image,omitempty"`

	// ImageTag pins an image tag. Empty uses the image as given (often :latest).
	// Prefer digests in production.
	// +optional
	ImageTag string `json:"imageTag,omitempty"`

	// ImagePullPolicy for the game server container.
	// +kubebuilder:default=IfNotPresent
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets for private registries.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// GenerateSecrets when true creates an Opaque Secret with random
	// admin-password and world-password keys if missing. Existing non-empty
	// keys are never overwritten. Passwords are never written into the CR.
	// +optional
	GenerateSecrets bool `json:"generateSecrets,omitempty"`

	// CredentialsSecretRef names an existing Secret with keys admin-password
	// and world-password. Required when generateSecrets is false.
	// +optional
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`

	// AllowEmptyWorldPassword when true sets an empty join password
	// (RSDW_PASSWORD=""). Without this, an unset world password is treated as
	// an error because the official image default is "random" each start.
	// +optional
	AllowEmptyWorldPassword bool `json:"allowEmptyWorldPassword,omitempty"`

	// Update configures SteamCMD / auto-stop-on-update behaviour.
	// +optional
	Update UpdateConfig `json:"update,omitempty"`

	// Gateway configures Envoy Gateway exposure (or NodePort fallback).
	// +optional
	Gateway GatewayConfig `json:"gateway,omitempty"`

	// Persistence configures the PVC mounted at the official data path.
	// +optional
	Persistence PersistenceConfig `json:"persistence,omitempty"`

	// Backup configures an optional CronJob that tars the data PVC into a
	// separate backup PVC (retain N archives).
	// +optional
	Backup BackupConfig `json:"backup,omitempty"`

	// ServerManager is an optional authenticated admin UI sidecar.
	// +optional
	ServerManager ServerManagerConfig `json:"serverManager,omitempty"`

	// Resources overrides auto-selected CPU/memory.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// TerminationGracePeriodSeconds is the SIGTERM grace window (default 60)
	// so the server can flush the world save on stop.
	// +kubebuilder:default=60
	// +kubebuilder:validation:Minimum=30
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`

	// NodeSelector pins the game server pod to specific nodes.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// DragonwildsServerStatus defines the observed state of DragonwildsServer.
type DragonwildsServerStatus struct {
	// Phase is Pending, Deploying, Running, VersionMismatch, or Failed.
	// +optional
	Phase string `json:"phase,omitempty"`

	// Ready is true when the game server Deployment is available.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// ConnectionAddress is the IP clients should use.
	// +optional
	ConnectionAddress string `json:"connectionAddress,omitempty"`

	// ConnectionPort is the UDP game port clients should use.
	// +optional
	ConnectionPort int32 `json:"connectionPort,omitempty"`

	// BeaconPort is the observed UDP beacon port.
	// +optional
	BeaconPort int32 `json:"beaconPort,omitempty"`

	// Message is a human-readable status detail.
	// +optional
	Message string `json:"message,omitempty"`

	// CredentialsSecretName is the Secret that holds join/admin passwords.
	// Never contains plaintext passwords.
	// +optional
	CredentialsSecretName string `json:"credentialsSecretName,omitempty"`

	// CredentialsGenerated is true when generateSecrets manages the Secret.
	// +optional
	CredentialsGenerated bool `json:"credentialsGenerated,omitempty"`

	// ImageVersion is the running container image reference.
	// +optional
	ImageVersion string `json:"imageVersion,omitempty"`

	// RunningBuildId is the game BuildId parsed from server session logs when Ready.
	// +optional
	RunningBuildId string `json:"runningBuildId,omitempty"`

	// WorldName is the observed/recorded world name after first successful start.
	// +optional
	WorldName string `json:"worldName,omitempty"`

	// WorldSaveGuid is recorded after first successful Ready so later recreates
	// can detect accidental new-world creation (palworld DedicatedServerName analogue).
	// +optional
	WorldSaveGuid string `json:"worldSaveGuid,omitempty"`

	// ServerManagerAddress is the IP for the optional admin UI (usually the VIP).
	// +optional
	ServerManagerAddress string `json:"serverManagerAddress,omitempty"`

	// ServerManagerPort is the TCP port for the optional admin UI.
	// +optional
	ServerManagerPort int32 `json:"serverManagerPort,omitempty"`

	// ObservedGeneration is the last reconciled generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represent the latest available observations.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=dws;dragonwilds
// +kubebuilder:printcolumn:name="Ready",type=boolean,JSONPath=`.status.ready`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Address",type=string,JSONPath=`.status.connectionAddress`
// +kubebuilder:printcolumn:name="Port",type=integer,JSONPath=`.status.connectionPort`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// DragonwildsServer is the Schema for the dragonwildsservers API.
type DragonwildsServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DragonwildsServerSpec   `json:"spec,omitempty"`
	Status DragonwildsServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DragonwildsServerList contains a list of DragonwildsServer.
type DragonwildsServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DragonwildsServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DragonwildsServer{}, &DragonwildsServerList{})
}
