# Game server reference

Official image: [`ghcr.io/runescape/rsdw-dedicated`](https://github.com/runescape/rsdw-dedicated)  
SteamCMD app: `4019830`  
Docs: [Jagex dedicated servers](https://dragonwilds.runescape.com/news/how-to-dedicated-servers), [wiki (Linux)](https://dragonwilds.runescape.wiki/w/Dedicated_Servers/Linux)

## Ports (UDP)

| Port | Env | Role |
|------|-----|------|
| 7777 | `RSDW_PORT` | Game |
| 8888 | `RSDW_BEACON_PORT` | Beacon (default = game + 1111) |

Host and container ports **must match** for the official image.

## Env vars (operator → container)

| Spec field | Env |
|------------|-----|
| `ownerId` | `RSDW_OWNER_ID` |
| `serverName` | `RSDW_SERVER_NAME` |
| `worldName` | `RSDW_WORLD_NAME` |
| Secret `admin-password` | `RSDW_ADMIN_PASSWORD` |
| Secret `world-password` | `RSDW_PASSWORD` |
| `gamePort` / `beaconPort` | `RSDW_PORT` / `RSDW_BEACON_PORT` |
| `update.autoUpdateImage` | `RSDW_AUTO_STOP_ON_UPDATE` |
| `update.validateOnStart` | `STEAMAPPVALIDATE` (`1`/`0`) |

Leaving `RSDW_PASSWORD` unset makes the image generate a **new random password every start**.
The operator always sets it from the Secret (use `allowEmptyWorldPassword: true` for an open world).

## DedicatedServer.ini

Path (on PVC / bare install):

`RSDragonwilds/Saved/Config/LinuxServer/DedicatedServer.ini`

Section: `[/Script/Dominion.DedicatedServerSettings]`

Mandatory fields: `OwnerId`, `ServerName`, `DefaultWorldName`, `AdminPassword`.

**Edits while the server is running are lost.** Stop → edit → start. The operator
renders a ConfigMap for observability/hashing; the official image is driven primarily
via env vars at pod start.

## Saves

World saves live under the data volume (`/home/steam/rsdw-dedicated` in the official
image), typically `RSDragonwilds/Saved/Savegames/*.sav`. Pod recreate must keep the PVC.

## Resources

Official minimum: 1 CPU, 4 GiB RAM, 20 GB disk.  
Operator default memory request: `max(4Gi, 2Gi + 1Gi×players)` (cap 6 players → 8 GiB).
