# Server Manager

Optional authenticated admin UI (Palworld-style, Dragonwilds-scoped).

```yaml
spec:
  serverManager:
    enabled: true
    port: 8088   # TCP on the same VIP as game UDP
    # image: harbor.dataknife.net/library/dragonwilds-operator:v0.1.0
```

## Access

- URL: `http://<status.serverManagerAddress>:<status.serverManagerPort>/`
- Auth: basic — user `admin`, password = Secret `admin-password`
- Enabling Recreate-rolls the game Deployment (sidecar + Service TCP port)

## Tabs

| Tab | Actions |
|-----|---------|
| Overview | Ready/phase, connect string, world GUID / build |
| Controls | Restart (Recreate), download world data `.tar.gz` |
| Settings | serverName, worldName, players, ownerId, admins → CR |
| Access | Rotate world/admin passwords (Secret); kick/ban stay in-game |
| Updates | `autoUpdateImage`, `validateOnStart` |
| Backups | List/download archives; run backup Job (needs `spec.backup`) |

Settings and password changes require a **restart** to apply into the game process.

## Notes

- No Palworld REST kick/ban — use in-game Server Management as owner.
- Prefer Direct connect for play; the UI is ops, not EOS Recent.
- Sidecar binary is `/server-manager` in the operator image.
