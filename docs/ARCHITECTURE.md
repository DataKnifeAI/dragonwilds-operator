# Architecture

## Resource graph

```
DragonwildsServer (CR)
├── Secret          ({name}-secrets)          admin-password, world-password
├── ConfigMap       ({name}-config)           DedicatedServer.ini (hash-annotated)
├── PVC             ({name}-data)             /home/steam/rsdw-dedicated
├── Deployment      ({name})                  Recreate, official image
├── Service         ({name})                  UDP game + beacon
├── Service         ({name}-envoy)            ClusterIP backend for UDPRoutes
└── Gateway         EnvoyProxy + UDPRoutes    VIP → game 7777 / beacon 8888
```

When `gateway.enabled: false`, skip Gateway/EnvoyProxy/UDPRoutes and expose via
`gateway.serviceType` (NodePort / LoadBalancer / ClusterIP).

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
