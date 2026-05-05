# F04 — Exit Node Browser

**Priority:** P0 | **Status:** todo

Browse available exit nodes and connect/disconnect. Accessible from the status dashboard.

## Layout

```
┌─ Exit Nodes ──────────────────────────────────────────┐
│  Active: none                                          │
│                                                        │
│  ○ none  (disable exit node)                          │
│  ● server-home      100.x.x.10   (this tailnet)      │
│▶ ● vpn-us           100.x.x.20   (this tailnet)  ✓  │
│                                                        │
│  [enter] connect   [esc] back                         │
└────────────────────────────────────────────────────────┘
```

- ✓ marks currently active exit node
- First row is always "none" — selecting it disconnects from any active exit node
- Only peers with `ExitNodeOption: true` shown

## Data Source

`tailscale status --json`

Relevant fields:
```
.Peer[*].ExitNodeOption    (bool — available as exit node)
.Peer[*].ExitNode          (bool — currently active)
.Peer[*].HostName
.Peer[*].TailscaleIPs[0]
.ExitNodeStatus.Active
```

## Connect / Disconnect

Shell out to CLI:

```sh
# connect
tailscale set --exit-node=<ip>

# disconnect
tailscale set --exit-node=
```

## Behavior

- `enter` on a peer → set as exit node → immediate status refresh → return to dashboard
- `enter` on "none" row → clear exit node
- Show spinner while command runs
- On error: surface stderr in status bar, do not change active state
- If tunnel is down, show warning: "Bring tunnel up first"

## Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `↓` / `j` / `k` | Navigate |
| `enter` | Connect to selected exit node |
| `esc` | Back to dashboard |

## Implementation Notes

- Exit node list is a subset of the peer list — filter from same `Status` struct, no extra fetch.
- After `tailscale set`, trigger immediate refresh (same pattern as F03).
- `tailscale set --exit-node=` (empty string) clears the exit node — test this on real tailscale CLI to confirm syntax before implementing.
