# FAQ

## Why won’t the server start?

`spec.ownerId` is required. Find your Player ID in-game (Settings → bottom → copy).
Without it the dedicated server refuses normal operation.

## Why did my password change after every restart?

The official image defaults `RSDW_PASSWORD` / `RSDW_ADMIN_PASSWORD` / `RSDW_WORLD_NAME`
to `random` when unset. Always set them via Secret / Compose `.env`.

## Runtime edits to DedicatedServer.ini disappeared

Expected. Stop the server, edit, then start. The operator only applies config at pod start.

## Server shows in the list but “connection lost”

Usually UDP port forwarding: both **game** and **beacon** ports must reach the host.
On the same LAN, Direct Connect to the LAN VIP/IP on the game port.

Recent / Public browser uses the **WAN IP EOS advertises**, not your kube-vip.
Without router port-forward (+ hairpin), those joins fail even when Direct works.
See [CONNECT.md](CONNECT.md).

## Why not Envoy Gateway like Palworld?

Envoy’s UDP proxy makes the game see the Envoy pod IP. Dragonwilds (Unreal) then
fails the join handshake. Use `gateway.enabled: false` + `serviceType: LoadBalancer`
+ `address` instead.

## Client can’t see the server after a game patch

Server build can lag the client. Prefer pinning an image digest; opt in to
`spec.update.autoUpdateImage` only when you accept SteamCMD stop+restart behaviour.
See [GAME_SERVER.md](GAME_SERVER.md).

## I deleted the CR and the world is gone

If the StorageClass reclaim policy is Delete, removing the PVC destroys data.
Enable `spec.backup` (see [BACKUP.md](BACKUP.md)) and/or copy `*.sav` before deleting
CRs in non-lab clusters. Prefer Retain for production.
