# Feature Tracker

## Legend

| Priority | Meaning |
|----------|---------|
| P0 | MVP — ship nothing without this |
| P1 | Important — ship soon after MVP |
| P2 | Nice-to-have |

| Status | Meaning |
|--------|---------|
| todo | Not started |
| in-progress | Active work |
| done | Shipped |

---

## Features

| ID | Feature | Priority | Status | Notes |
|----|---------|----------|--------|-------|
| F01 | Status dashboard — show connection state, local IP, hostname | P0 | todo | Main screen. Poll `tailscale status --json` or local API |
| F02 | Peer list — show all peers with IP, hostname, online state | P0 | todo | Sortable. Show last-seen for offline peers |
| F03 | Bring tunnel up / down | P0 | todo | `tailscale up` / `tailscale down` with confirmation prompt |
| F04 | Exit node browser — list available exit nodes, connect/disconnect | P0 | todo | Show country/location if available. Highlight active node |
| F05 | Start / stop tailscaled daemon | P1 | todo | systemctl or direct. Detect if daemon not running on launch |
| F06 | Ping peer — show latency to selected peer | P1 | todo | Run `tailscale ping` inline, show result in UI |
| F07 | Auto-refresh — live status updates | P1 | todo | Configurable interval, default 5s |
| F08 | Key bindings help overlay | P1 | todo | `?` to show keybindings |
| F09 | Copy peer IP to clipboard | P2 | todo | Useful for SSH etc |
| F10 | SSH into peer directly from UI | P2 | todo | Spawn terminal with `ssh user@<peer-ip>` |
| F11 | Subnet routes — view advertised/accepted routes | P2 | todo | Read-only display first |
| F12 | MagicDNS toggle | P2 | todo | `tailscale set --accept-dns` |

---

## Decisions

- **Tech stack**: Go + bubbletea. lipgloss for styling, bubbles for components.
- **Tailscale integration**: prefer local API socket over CLI where possible.
