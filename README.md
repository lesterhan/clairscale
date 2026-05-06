# clairscale

Interactive terminal UI for Tailscale — status at a glance, bring up/down, switch exit nodes, without touching the CLI directly.

## Features

- Live Tailscale status dashboard (peers, IPs, connection state)
- Start/stop the Tailscale daemon
- Bring the tunnel up and down (`tailscale up` / `tailscale down`)
- Browse and connect to exit nodes
- Peer list with ping/latency info

## Installation

**Requirements:** Go 1.21+

```sh
go install ./cmd/clairscale
```

Then add Go's bin directory to your PATH so the binary is found.

**bash** — add to `~/.bashrc` or `~/.bash_profile`:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

**fish** — add to `~/.config/fish/config.fish`:
```fish
fish_add_path (go env GOPATH)/bin
```

After editing, reload your shell (`source ~/.bashrc` / `source ~/.config/fish/config.fish`) or open a new terminal.

## Usage

```
clairscale
```

See [planning/features.md](planning/features.md) for roadmap and status.

## Development

Run directly without installing:

```sh
go run ./cmd/clairscale
```

Or build a local binary:

```sh
go build -o clairscale ./cmd/clairscale
./clairscale
```

See [CLAUDE.md](CLAUDE.md) for project conventions and architecture notes.
