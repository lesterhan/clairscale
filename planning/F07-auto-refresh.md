# F07 — Auto Refresh

**Priority:** P1 | **Status:** todo

Keep the status dashboard and peer list current without user action.

## Behavior

- Poll `tailscale status --json` every 5 seconds
- Update displayed data in place — no flicker, no full redraw
- User-triggered actions (up/down, exit node change) fire an immediate refresh on top of the tick

## Implementation

Use bubbletea's tick pattern:

```go
type tickMsg time.Time

func tickCmd() tea.Cmd {
    return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
```

On `tickMsg` in `Update`: fetch status, return next `tickCmd()`.

Fetch runs as a `tea.Cmd` (async) — result comes back as a `statusMsg`. Never block the event loop.

## Refresh Triggers

| Trigger | Delay |
|---------|-------|
| App start | Immediate |
| Tick | Every 5s |
| After `tailscale up/down` | Immediate (skip next tick, reset timer) |
| After exit node change | Immediate (skip next tick, reset timer) |

## Error Handling

If a refresh fails (socket error, command error):
- Keep last-known data displayed
- Show subtle error indicator (e.g., status bar: `"Last updated 12s ago — connection error"`)
- Keep retrying on next tick — do not crash or freeze

## Configuration

No user-facing config for now. Hard-code 5s interval. Can expose later via flag or config file.

## Implementation Notes

- Single fetch path used by both tick and manual triggers — same `tea.Cmd` function.
- Reset the tick timer after a manual refresh to avoid double-fetches within seconds of each other.
- Only refresh the active view — no need to fetch peer data when only the dashboard is visible, and vice versa. (Optimize later if needed; start simple with one fetch per tick.)
