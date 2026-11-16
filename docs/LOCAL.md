# Local (Docker Compose)

Feature parity without Kubernetes — same official image the operator uses.

```bash
cd compose
cp .env.example .env
# set RSDW_OWNER_ID and RSDW_ADMIN_PASSWORD
docker compose up -d
docker compose logs -f
```

Or from the repo root:

```bash
make compose-up
make compose-logs
make compose-down
```

## Requirements

- Docker with Compose v2
- UDP 7777 + 8888 free on the host (or change `GAME_PORT` / `BEACON_PORT` together)
- ~8 GiB RAM recommended for a full party; official minimum is 4 GiB
- Disk under `compose/data/` for the world (gitignored)

## Persistence

`./data` mounts to `/home/steam/rsdw-dedicated`. Keep this directory to preserve the world
across `compose down` / `up`.
