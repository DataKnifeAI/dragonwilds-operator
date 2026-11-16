# Connect

## Direct connect

1. Start the game → join via **Direct Connect** (or equivalent world browser entry).
2. Enter `status.connectionAddress` and `status.connectionPort` (UDP game port, default 7777).
3. Enter the world password from the credentials Secret (`world-password`), unless you
   set `allowEmptyWorldPassword: true`.

```bash
kubectl get dragonwildsserver ashenfall-home -o wide
kubectl get secret ashenfall-home-secrets -o jsonpath='{.data.world-password}' | base64 -d; echo
kubectl get secret ashenfall-home-secrets -o jsonpath='{.data.admin-password}' | base64 -d; echo
```

## Homelab without Gateway

```yaml
spec:
  gateway:
    enabled: false
    serviceType: NodePort
```

Then use a node IP + the allocated NodePorts (game + beacon). Both UDP ports must be
reachable; missing beacon forwarding often shows the world listed then “connection lost”.

## Crossplay

Console players can join PC-hosted dedicated servers per Jagex docs — confirm current
client build matches the server (version mismatch removes the server from the browser).

## Admin

Use the admin password from the Secret in the in-game Server Management UI
(Pause → Settings).
