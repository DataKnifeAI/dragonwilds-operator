# Changelog

## Unreleased

### Added

- Optional `spec.serverManager` admin UI sidecar (`/server-manager` in the operator
  image): Overview, Controls, Settings, Access, Updates, Backups on
  `http://VIP:8088` with basic auth (`admin` + `admin-password`).

## 0.1.0

### Added

- Kubebuilder v4 operator for `DragonwildsServer` (`dragonwilds.dataknife.ai/v1alpha1`).
- Reconciler owns Secret, ConfigMap, PVC, Deployment (Recreate), Service, optional backup CronJob/PVC, optional Envoy Gateway stack.
- Official image defaults: `ghcr.io/runescape/rsdw-dedicated`, UDP 7777/8888, PVC at `/home/steam/rsdw-dedicated`.
- Recommended exposure: LoadBalancer VIP via `gateway.address` (`gateway.enabled: false`).
- Optional Envoy Gateway + UDPRoute wiring (experimental — Unreal joins fail behind Envoy UDP proxy).
- Optional `spec.backup` CronJob (tar.gz retain-N onto `{name}-backups` PVC).
- Status: `worldSaveGuid`, `runningBuildId`; Ready fails if save GUID drifts.
- Docs: ARCHITECTURE, GAME_SERVER, CONNECT, BACKUP, METRICS, LOCAL, FAQ, GITLAB_MIRROR.
- Docker Compose path under `compose/`.
- GitHub → GitLab mirror + Harbor publish (`harbor.dataknife.net/library/dragonwilds-operator`).
- Hero graphic + CI badges in README.

### Notes

- Prefer pinning game and operator images by digest in production.
- GHCR publish for the operator image is optional later; Harbor is the primary registry today.
