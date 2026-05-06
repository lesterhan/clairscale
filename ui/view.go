package ui

import (
	"fmt"
	"strings"
	"time"
)

func (m Model) View() string {
	switch m.screen {
	case screenLoading:
		return "\n  Loading...\n"
	case screenDaemonDown:
		return viewDaemonDown()
	case screenDashboard:
		return viewDashboard(m)
	case screenPeerList:
		return viewPeerList(m.peerList)
	case screenExitNodes:
		return viewExitNodes(m.exitNodes)
	}
	return ""
}

func viewDaemonDown() string {
	var b strings.Builder
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "  ✗ tailscaled is not running")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "  [s] Start daemon   [q] Quit")
	return b.String()
}

func viewDashboard(m Model) string {
	if m.status == nil {
		return "\n  Loading...\n"
	}

	s := m.status
	self := s.Self

	ip := "—"
	if len(self.TailscaleIPs) > 0 {
		ip = self.TailscaleIPs[0]
	}

	exitNode := "none"
	for _, p := range s.Peer {
		if p.ExitNode {
			exitNode = p.HostName
			break
		}
	}

	online, offline := 0, 0
	for _, p := range tailnetPeers(m.status) {
		if p.Online {
			online++
		} else {
			offline++
		}
	}

	var b strings.Builder
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "  %-24s %s\n", self.HostName, ip)
	fmt.Fprintf(&b, "  tailnet: %-20s %s\n", s.CurrentTailnet.Name, connStateStr(s.BackendState))
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "  Exit node: %s\n", exitNode)
	fmt.Fprintf(&b, "  Peers: %d online · %d offline\n", online, offline)
	fmt.Fprintln(&b)

	if m.fetchErr != nil {
		ago := time.Since(m.lastFetch).Round(time.Second)
		fmt.Fprintf(&b, "  Last updated %s ago — connection error\n", ago)
		fmt.Fprintln(&b)
	}

	fmt.Fprintln(&b, "  [p] peers  [e] exit nodes  [u] up/down  [q] quit")
	return b.String()
}

func connStateStr(backendState string) string {
	switch backendState {
	case "Running":
		return "● Connected"
	case "Stopped":
		return "○ Tunnel down"
	case "NeedsLogin":
		return "○ Needs login"
	case "Starting":
		return "◌ Starting"
	default:
		if backendState == "" {
			return "○ Unknown"
		}
		return "○ " + backendState
	}
}
