# CLAUDE.md — clairscale

## Project

Terminal UI (TUI) app for managing Tailscale. Wraps `tailscale` CLI and/or the Tailscale local API. Goal: keyboard-driven, fast, no desktop dependency.

## Stack

**Go + bubbletea**

- bubbletea for TUI (Elm architecture — Model/Update/View)
- lipgloss for styling
- bubbles for common components (list, spinner, viewport)
- Tailscale local API via Unix socket (`/var/run/tailscale/tailscaled.sock`)

## Architecture

```
clairscale/
  cmd/clairscale/main.go   # entry point — wires bubbletea program
  ui/                      # TUI models, views, components (bubbletea Elm arch)
  tailscale/               # Tailscale API/CLI wrapper — no UI imports allowed here
  planning/                # feature docs (not shipped)
```

bubbletea pattern: each screen is a `Model` with `Init() / Update() / View()`. Compose screens by embedding or delegating to child models.

## Tailscale Integration

Prefer the Tailscale local API (`/var/run/tailscale/tailscaled.sock`) over shelling out to the CLI where possible. Fall back to CLI for operations not exposed in the API.

Relevant CLI commands:
- `tailscale status --json`
- `tailscale up / down`
- `tailscale set --exit-node=<ip>`
- `tailscale ping <peer>`

## Planning

Feature docs live in `planning/`. See `planning/features.md` for master list.
- Status column: `todo` | `in-progress` | `done`
- Priority: `P0` (must-have MVP) | `P1` (important) | `P2` (nice-to-have)

## Conventions

- No comments explaining what code does — only why (non-obvious constraints, workarounds).
- No unused error handling or defensive fallbacks for impossible cases.
- Keep UI and data layers separate — nothing in `ui/` talks to Tailscale directly.
- Update `planning/features.md` when feature status changes.
