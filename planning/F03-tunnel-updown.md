# F03 — Tunnel Up / Down

**Priority:** P0 | **Status:** todo

Toggle the Tailscale tunnel with a single key. Accessible from the status dashboard.

## Behavior

Press `u` on the dashboard:

- If tunnel is **up** → confirm prompt → run `tailscale down`
- If tunnel is **down** → run `tailscale up` (no confirm needed — non-destructive)

Confirmation prompt (inline, not a modal):

```
  Bring tunnel down? [y/N]
```

Shown in the bottom bar. Any key other than `y` cancels.

## Command

Shell out to `tailscale` CLI (no local API equivalent for up/down):

```
tailscale up
tailscale down
```

`tailscale up` may open a browser for auth on first run — detect this in stdout and surface a message: `"Auth required — check your browser"`.

## States

| Tunnel state | `u` press result |
|-------------|-----------------|
| Connected | Confirm → `tailscale down` |
| Down (daemon running) | `tailscale up` immediately |
| Daemon not running | Show error: "tailscaled not running" |
| Connecting / transitioning | Disable key, show spinner |

## UI Feedback

- Status indicator on dashboard updates optimistically on keypress
- Spinner shown while command runs
- On error: show CLI stderr in status bar (red), revert optimistic state

## Implementation Notes

- Run as `tea.Cmd` (async) — never block the TUI event loop.
- Capture both stdout and stderr. Surface stderr on non-zero exit.
- After command completes, trigger an immediate status refresh (don't wait for the 5s auto-refresh tick).
- `tailscale up` accepts flags (`--exit-node`, `--accept-routes`, etc.) — for now pass none; add flags as needed when other features require it.
