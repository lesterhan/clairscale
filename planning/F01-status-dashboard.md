# F01 — Status Dashboard

**Priority:** P0 | **Status:** todo

The main screen. First thing shown on launch. All other views are reachable from here.

## Layout

```
┌─ clairscale ──────────────────────────────────────────┐
│                                                        │
│  Status   ● Connected                                  │
│  Node     myhost (100.x.x.x)                          │
│  Account  user@example.com / my-tailnet               │
│  Exit     none                                         │
│                                                        │
│  Peers  12   Online  9   Offline  3                   │
│                                                        │
│  [p] Peers   [e] Exit Nodes   [?] Help   [q] Quit     │
└────────────────────────────────────────────────────────┘
```

Status indicator colors (lipgloss):
- ● green  = Connected
- ● yellow = Connecting / authenticating
- ● red    = Offline / not running

## Data Source

`tailscale status --json` or local API `GET /localapi/v0/status`

Relevant fields from response:
```
.Self.HostName
.Self.TailscaleIPs[0]
.Self.Online
.CurrentTailnet.Name
.User[...].LoginName
.Peer (map)  — count online/offline
.ExitNodeStatus.Active + .ExitNodeStatus.TailscaleIPs
```

## Behavior

- Shown on launch after daemon check passes (see F05)
- Auto-refreshes every 5s (see F07)
- If daemon not reachable on launch → show error state with instructions
- `tailscale up` / `tailscale down` toggle accessible via `u` key from here (see F03)
- Status line updates immediately on any user action (optimistic UI, then real refresh)

## Key Bindings

| Key | Action |
|-----|--------|
| `p` | Open peer list (F02) |
| `e` | Open exit node browser (F04) |
| `u` | Toggle tunnel up/down (F03) |
| `?` | Help overlay (F08) |
| `q` / `ctrl+c` | Quit |

## Implementation Notes

- This is the root bubbletea model. Child screens (peers, exit nodes) are composed as sub-models.
- Daemon check on init: try socket, if fail show "tailscaled not running" + instructions.
- Keep status fetch in `tailscale/` package — UI just reads a `Status` struct.
