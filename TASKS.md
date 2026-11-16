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

## Phase 3 — Exposure (palworld parity)

- [x] Envoy Gateway + UDPRoute + EnvoyProxy wiring
- [x] Detect missing GatewayClass → clear condition
- [x] NodePort escape hatch when `gateway.enabled: false`

## Phase 4 — Ops

- [ ] Optional backup CronJob
- [x] Controller metrics endpoint (Kubebuilder default)
- [ ] ServiceMonitor polish / scrape docs
- [ ] Image digests / SBOM-ready publish (GHCR later)

## Phase 5 — Launch

- [x] README + hero graphic
- [ ] CI badges
- [ ] v0.1.0 tag
