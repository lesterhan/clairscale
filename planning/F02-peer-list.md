# F02 — Peer List

**Priority:** P0 | **Status:** todo

Scrollable list of all Tailscale peers. Accessible from the status dashboard.

## Layout

```
┌─ Peers ───────────────────────────────────────────────┐
│  9 online · 3 offline                     [/] filter  │
│                                                        │
│  ● server-1        100.x.x.1    linux                 │
│  ● laptop-work     100.x.x.2    macOS                 │
│▶ ● phone           100.x.x.3    android               │
│  ○ nas             100.x.x.4    linux   (2h ago)      │
│  ○ old-desktop     100.x.x.5    windows (3d ago)      │
│                                                        │
│  [enter] ping   [y] copy IP   [esc] back              │
└────────────────────────────────────────────────────────┘
```

- ● green = online, ○ gray = offline
- Last-seen shown for offline peers
- Selected row highlighted

## Data Source

`tailscale status --json` → `.Peer` map

Relevant fields per peer:
```
.HostName
.TailscaleIPs[0]
.OS
.Online
.LastSeen          (time.Time)
.ExitNodeOption    (bool — can be used as exit node)
```

## Behavior

- Sorted: online peers first (alpha), then offline (alpha)
- `/` opens inline filter — narrows list as you type, `esc` clears
- `enter` on selected peer triggers ping (F06) — shows latency inline
- `y` copies peer IP to clipboard (F09) — silent confirmation in status bar
- `esc` / `backspace` returns to status dashboard (F01)
- List uses `bubbles/list` component

## Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `↓` / `j` / `k` | Navigate |
| `enter` | Ping selected peer |
| `y` | Copy IP to clipboard |
| `/` | Filter |
| `esc` | Back to dashboard |

## Implementation Notes

- Reuse the same `Status` struct fetched for F01 — no extra API call needed.
- Ping (F06) runs async: fire `tea.Cmd`, update peer row with result when done.
- Clipboard copy: try `wl-copy` (Wayland) then `xclip`/`xsel` (X11). Fail silently with error in status bar.
