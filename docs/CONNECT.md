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

## LoadBalancer VIP (recommended for Dragonwilds)

Dragonwilds (Unreal) expects a direct UDP path to the game pod. Envoy’s UDP proxy
rewrites the client source to the Envoy pod IP and joins fail with “network connection
failed”. Prefer a kube-vip LoadBalancer on the game Service:

```yaml
spec:
  gateway:
    enabled: false
    serviceType: LoadBalancer
    address: "192.168.14.188"   # dedicated VIP (Kube-VIP / MetalLB)
```

Clients use `status.connectionAddress:7777` (beacon `8888`). Host and container ports match
the official image requirement.

## Envoy Gateway

```yaml
spec:
  gateway:
    enabled: true
    className: envoy
    address: "192.168.14.188"
```

Works for some UDP games (e.g. Palworld) but is **not** reliable for Dragonwilds joins.

## Homelab NodePort

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
