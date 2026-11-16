# Changelog

## Unreleased

### Added

- Initial Kubebuilder v4 scaffold for `DragonwildsServer` (`dragonwilds.dataknife.ai/v1alpha1`).
- Core reconciler owns Secret, ConfigMap, PVC, Deployment (Recreate), and Service.
- Official image defaults: `ghcr.io/runescape/rsdw-dedicated`, UDP 7777/8888, PVC at `/home/steam/rsdw-dedicated`.
- Docker Compose path under `compose/` (`docs/LOCAL.md`).
- Docs: ARCHITECTURE, GAME_SERVER, CONNECT, LOCAL, FAQ, GITLAB_MIRROR.
- GitHub → GitLab mirror + Harbor publish (`harbor.dataknife.net/library/dragonwilds-operator`).
- Envoy Gateway parity with palworld-operator: EnvoyProxy + Gateway + UDPRoutes (game/beacon), dedicated VIP, GatewayClass detection.
- Status: `worldSaveGuid`, `runningBuildId` from pod logs; Ready fails if save GUID drifts.
- Hero graphic in README.
