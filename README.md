# clairscale

Interactive terminal UI for Tailscale — status at a glance, bring up/down, switch exit nodes, without touching the CLI directly.

## Why

`tailscale` CLI works but is tedious for day-to-day use. trayscale requires a desktop tray. clairscale lives in the terminal, fast and keyboard-driven.

## Features

- Live Tailscale status dashboard (peers, IPs, connection state)
- Start/stop the Tailscale daemon
- Bring the tunnel up and down (`tailscale up` / `tailscale down`)
- Browse and connect to exit nodes
- Peer list with ping/latency info

## Usage

```
clairscale
```

See [planning/features.md](planning/features.md) for roadmap and status.

## Development

See [CLAUDE.md](CLAUDE.md) for project conventions and architecture notes.
