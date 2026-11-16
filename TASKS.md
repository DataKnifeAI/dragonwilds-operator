# Tasks

## Phase 0 — Discovery

- [x] Map official env vars + ports (7777/8888) into CRD defaults
- [x] Live cluster smoke: apply CR → join succeeded (Surreal Wolf)

## Phase 1 — Scaffold & CRD

- [x] Kubebuilder project + `DragonwildsServer` API
- [x] Sample CR, docs stubs, Compose parity
- [x] `make generate manifests build test`

## Phase 2 — Core reconciler

- [x] Secret / ConfigMap / PVC / Deployment / Service builders
- [x] Validation for ownerId + credentials + gateway address
- [x] Status worldSaveGuid / runningBuildId from pod logs
- [x] KinD/cluster: apply sample → connect → world on PVC

## Phase 3 — Exposure

- [x] Envoy Gateway + UDPRoute + EnvoyProxy wiring (optional / experimental)
- [x] Detect missing GatewayClass → clear condition
- [x] Recommended path: LoadBalancer VIP (`gateway.enabled: false`) — Envoy breaks Unreal joins
- [x] NodePort escape hatch

## Phase 4 — Ops

- [x] Optional backup CronJob + backup PVC
- [x] Controller metrics endpoint (Kubebuilder default)
- [x] ServiceMonitor polish / scrape docs (`docs/METRICS.md`)
- [x] Optional Server Manager admin UI sidecar (`docs/SERVER_MANAGER.md`)
- [ ] Image digests / SBOM-ready publish (GHCR — optional later)

## Phase 5 — Launch

- [x] README + hero graphic
- [x] CI badges
- [x] v0.1.0 tag
