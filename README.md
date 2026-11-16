# dragonwilds-operator

<p align="center">
  <img src="docs/assets/hero-dragonwilds-operator.jpg" alt="Dragonwilds Operator — Ashenfall world keeper" width="640">
</p>

Kubernetes operator for [RuneScape: Dragonwilds](https://www.runescape.com/dragonwilds) dedicated servers.
One custom resource (`DragonwildsServer`) reconciles Deployment, PVC, ConfigMap, Secret, and Service around Jagex’s official image [`ghcr.io/runescape/rsdw-dedicated`](https://github.com/runescape/rsdw-dedicated).

Modeled on [palworld-operator](https://github.com/DataKnifeAI/palworld-operator) (“World Keeper”).

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

## Quick start (no cluster)

See [docs/LOCAL.md](docs/LOCAL.md) — Docker Compose with the same official image.

## Docs

| Doc | Topic |
|-----|--------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Resource graph & design rules |
| [GAME_SERVER.md](docs/GAME_SERVER.md) | Ports, env, ini, saves, updates |
| [CONNECT.md](docs/CONNECT.md) | Join / passwords / NodePort |
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
