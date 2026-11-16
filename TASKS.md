# Tasks

## Phase 0 — Discovery

- [ ] Manual SteamCMD / official image smoke on a VM
- [ ] Confirm save placement / world migration notes in docs/DISCOVERY.md
- [x] Map official env vars + ports (7777/8888) into CRD defaults

## Phase 1 — Scaffold & CRD

- [x] Kubebuilder project + `DragonwildsServer` API
- [x] Sample CR, docs stubs, Compose parity
- [x] `make generate manifests build test`

## Phase 2 — Core reconciler

- [x] Secret / ConfigMap / PVC / Deployment / Service builders
- [x] Validation for ownerId + credentials + gateway address
- [ ] World-name / save presence Ready condition (log or PVC inspect)
- [ ] Status imageVersion from server log header
- [ ] KinD end-to-end: apply sample → connect → delete pod → world survives

## Phase 3 — Exposure

- [ ] Envoy Gateway + UDPRoute wiring (`internal/gateway`)
- [ ] Detect missing GatewayClass → clear condition
- [ ] NodePort / annotation escape hatch docs polish

## Phase 4 — Ops

- [ ] Optional backup CronJob
- [ ] Prometheus metrics + ServiceMonitor
- [ ] Image digests / SBOM-ready publish (GHCR later)

## Phase 5 — Launch

- [ ] README polish + CI badges
- [ ] v0.1.0 tag
