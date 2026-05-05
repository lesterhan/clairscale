# F05 — Daemon Start / Stop

**Priority:** P1 | **Status:** todo

Detect if `tailscaled` is running. Start/stop it from within the UI.

## Behavior on Launch

Before showing the dashboard, clairscale checks if tailscaled is reachable:

- Try connecting to `/var/run/tailscale/tailscaled.sock`
- If unreachable → show daemon-down screen instead of dashboard

```
┌─ clairscale ──────────────────────────────────────────┐
│                                                        │
│  ✗ tailscaled is not running                          │
│                                                        │
│  [s] Start daemon   [q] Quit                          │
│                                                        │
└────────────────────────────────────────────────────────┘
```

## Start / Stop

Use systemctl (most Linux systems). Try in order, use first that works:

```sh
# start
systemctl start tailscaled

# stop
systemctl stop tailscaled
```

Fallback if systemctl unavailable: surface error telling user to start tailscaled manually.

Check if user has permission (may need sudo). If `systemctl` fails with permission error, show: `"Permission denied — try: sudo systemctl start tailscaled"`.

## States

| State | UI |
|-------|----|
| Daemon running | Normal dashboard (F01) |
| Daemon not running | Daemon-down screen with `[s]` to start |
| Starting (in progress) | Spinner + "Starting tailscaled…" |
| Start failed | Error message + stderr |
| Daemon stopped by user | Daemon-down screen |

## Key Bindings (daemon-down screen)

| Key | Action |
|-----|--------|
| `s` | Start daemon |
| `q` | Quit |

## Key Bindings (dashboard, daemon running)

Daemon stop not exposed as a single key from dashboard — too destructive. Access via `?` help if needed, or add later.

## Implementation Notes

- Socket path: `/var/run/tailscale/tailscaled.sock`. Make configurable via env var `TAILSCALE_SOCKET` for non-standard installs.
- After successful start, wait for socket to be connectable (poll with timeout ~5s) before transitioning to dashboard.
- Detection runs once on startup, not on every refresh tick.
