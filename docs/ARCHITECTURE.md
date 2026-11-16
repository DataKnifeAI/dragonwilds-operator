# Architecture

## Resource graph

```
DragonwildsServer (CR)
├── Secret          ({name}-secrets)          admin-password, world-password
├── ConfigMap       ({name}-config)           DedicatedServer.ini (hash-annotated)
├── PVC             ({name}-data)             /home/steam/rsdw-dedicated
├── Deployment      ({name})                  Recreate, official image
├── Service         ({name})                  UDP game + beacon (LoadBalancer VIP recommended)
├── PVC             ({name}-backups)          optional tar archives (backup.enabled)
├── CronJob         ({name}-backup)           optional scheduled backup
├── Service         ({name}-envoy)            ClusterIP backend for UDPRoutes (Envoy only)
└── Gateway         EnvoyProxy + UDPRoutes    optional — not recommended for Dragonwilds
```

Default exposure: `gateway.enabled: false` + `serviceType: LoadBalancer` + `address`
(kube-vip VIP). Envoy UDP proxy rewrites client source IPs and breaks Unreal joins;
keep Gateway wiring for experiments / parity with other game operators only.

## Design rules

| Rule | Why |
|------|-----|
| `ownerId` required | Official server refuses to start without it |
| Passwords only in Secrets | Never in CR, status, or logs |
| Deployment strategy `Recreate` | Single persistent world — no rolling pairs |
| `terminationGracePeriodSeconds ≥ 30` (default 60) | Clean save-on-stop |
| Config hash annotation | Ini/env changes only roll the pod intentionally |
| PVC retain (document Delete reclaim) | Deleting the CR must not silently destroy worlds |

## Image & updates

Default game image: [`ghcr.io/runescape/rsdw-dedicated`](https://github.com/runescape/rsdw-dedicated).

The official image runs SteamCMD on each container start. `spec.update.autoUpdateImage`
maps to `RSDW_AUTO_STOP_ON_UPDATE=true` so a newer Steam build stops the container and
Kubernetes restarts it (SteamCMD then updates into the PVC).
