# dragonwilds-operator

<p align="center">
  <img src="docs/assets/hero-dragonwilds-operator.jpg" alt="Dragonwilds Operator — Ashenfall world keeper" width="640">
</p>

<p align="center">
  <a href="https://github.com/DataKnifeAI/dragonwilds-operator/actions/workflows/test.yml"><img src="https://github.com/DataKnifeAI/dragonwilds-operator/actions/workflows/test.yml/badge.svg" alt="Test"></a>
  <a href="https://github.com/DataKnifeAI/dragonwilds-operator/actions/workflows/lint.yml"><img src="https://github.com/DataKnifeAI/dragonwilds-operator/actions/workflows/lint.yml/badge.svg" alt="Lint"></a>
  <a href="https://github.com/DataKnifeAI/dragonwilds-operator/actions/workflows/build.yml"><img src="https://github.com/DataKnifeAI/dragonwilds-operator/actions/workflows/build.yml/badge.svg" alt="Build"></a>
</p>

Kubernetes operator for [RuneScape: Dragonwilds](https://www.runescape.com/dragonwilds) dedicated servers.
One custom resource (`DragonwildsServer`) reconciles Deployment, PVC, ConfigMap, Secret, and Service around Jagex’s official image [`ghcr.io/runescape/rsdw-dedicated`](https://github.com/runescape/rsdw-dedicated).

Modeled on [palworld-operator](https://github.com/DataKnifeAI/palworld-operator) (“World Keeper”).

**Site:** [dataknifeai.github.io/dragonwilds-operator](https://dataknifeai.github.io/dragonwilds-operator/)  
**Repos:** [GitHub](https://github.com/DataKnifeAI/dragonwilds-operator) · [GitLab mirror](https://gitlab.com/dk-raas/dkai/game-servers/dragonwilds-operator)  
**Operator image:** `harbor.dataknife.net/library/dragonwilds-operator` (see [docs/GITLAB_MIRROR.md](docs/GITLAB_MIRROR.md))

## Quick start (cluster)

```bash
make manifests generate
make install          # apply CRD
make run              # or deploy the controller image later

kubectl apply -f config/samples/dragonwilds_v1alpha1_dragonwildsserver.yaml
# edit ownerId first!
```

Homelab exposure: `gateway.enabled: false`, `serviceType: LoadBalancer`, set `gateway.address`
to a kube-vip VIP, then Direct Connect to `status.connectionAddress:7777`.

## Quick start (no cluster)

See [docs/LOCAL.md](docs/LOCAL.md) — Docker Compose with the same official image.

## Docs

| Doc | Topic |
|-----|--------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Resource graph & design rules |
| [GAME_SERVER.md](docs/GAME_SERVER.md) | Ports, env, ini, saves, updates |
| [CONNECT.md](docs/CONNECT.md) | Join / VIP / Recent vs Direct |
| [BACKUP.md](docs/BACKUP.md) | Optional CronJob backups |
| [SERVER_MANAGER.md](docs/SERVER_MANAGER.md) | Optional admin UI sidecar |
| [METRICS.md](docs/METRICS.md) | Prometheus ServiceMonitor |
| [LOCAL.md](docs/LOCAL.md) | Compose path |
| [FAQ.md](docs/FAQ.md) | Common footguns |
| [TASKS.md](TASKS.md) | Phased roadmap |
| [CHANGELOG.md](CHANGELOG.md) | Changes |
| [GITLAB_MIRROR.md](docs/GITLAB_MIRROR.md) | GitHub → GitLab → Harbor |

## Develop

```bash
export PATH="$HOME/sdk/go/bin:$HOME/go/bin:$PATH"   # if using a local Go toolchain
make generate manifests
make test
make build
make lint
```

Requires Go 1.25+ (scaffold targets controller-runtime v0.20).

## License

Apache-2.0 — see [LICENSE](LICENSE).
