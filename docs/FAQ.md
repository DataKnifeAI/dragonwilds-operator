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
On the same LAN, Direct Connect to the LAN IP on the game port.

## Client can’t see the server after a game patch

Server build can lag the client. Prefer pinning an image digest; opt in to
`spec.update.autoUpdateImage` only when you accept SteamCMD stop+restart behaviour.
See [GAME_SERVER.md](GAME_SERVER.md).

## I deleted the CR and the world is gone

If the StorageClass reclaim policy is Delete, removing the PVC destroys data.
Back up `*.sav` before deleting CRs in non-lab clusters. Prefer Retain for production.
