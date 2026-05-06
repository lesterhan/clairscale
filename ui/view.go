package ui

import (
	"fmt"
	"strings"
	"time"
)

func (m Model) View() string {
	switch m.screen {
	case screenLoading:
		return "\n  " + styleDim.Render("Loading...") + "\n"
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
	fmt.Fprintln(&b, "  "+styleError.Render("✗ tailscaled is not running"))
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "  "+styleDim.Render("[s] Start daemon   [q] Quit"))
	return b.String()
}

func viewDashboard(m Model) string {
	if m.status == nil {
		return "\n  " + styleDim.Render("Loading...") + "\n"
	}

	s := m.status
	self := s.Self

	ip := "—"
	if len(self.TailscaleIPs) > 0 {
		ip = self.TailscaleIPs[0]
	}

	exitNode := styleDim.Render("none")
	for _, p := range s.Peer {
		if p.ExitNode {
			exitNode = styleActive.Render(p.HostName)
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
	fmt.Fprintf(&b, "  %s  %s\n",
		styleTitle.Render(self.HostName),
		styleDim.Render(ip),
	)
	fmt.Fprintf(&b, "  %s %s  %s\n",
		styleLabel.Render("tailnet:"),
		styleDim.Render(s.CurrentTailnet.Name),
		connStateStr(s.BackendState),
	)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "  %s %s\n", styleLabel.Render("Exit node:"), exitNode)
	fmt.Fprintf(&b, "  %s %s · %s\n",
		styleLabel.Render("Peers:"),
		styleOnline.Render(fmt.Sprintf("%d online", online)),
		styleOffline.Render(fmt.Sprintf("%d offline", offline)),
	)
	fmt.Fprintln(&b)

	if m.fetchErr != nil {
		ago := time.Since(m.lastFetch).Round(time.Second)
		fmt.Fprintf(&b, "  %s\n",
			styleError.Render(fmt.Sprintf("Last updated %s ago — connection error", ago)),
		)
		fmt.Fprintln(&b)
	}

	if !m.canManage {
		fmt.Fprintf(&b, "  %s\n",
			styleWarning.Render("⚠ Not root or in the tailscale group — exit node changes will fail"),
		)
		fmt.Fprintf(&b, "  %s\n",
			styleDim.Render("  Fix: sudo tailscale set --operator=$USER  (one-time, no group needed)"),
		)
		fmt.Fprintf(&b, "  %s\n",
			styleDim.Render("   Or: sudo usermod -aG tailscale $USER  (then log out and back in)"),
		)
		fmt.Fprintln(&b)
	}

	fmt.Fprintln(&b, "  "+styleDim.Render("[p] peers  [e] exit nodes  [u] up/down  [q] quit"))
	return b.String()
}

func connStateStr(state string) string {
	switch state {
	case "Running":
		return styleConnected.Render("● Connected")
	case "Stopped":
		return styleWarning.Render("○ Tunnel down")
	case "NeedsLogin":
		return styleError.Render("○ Needs login")
	case "Starting":
		return styleWarning.Render("◌ Starting")
	default:
		if state == "" {
			return styleOffline.Render("○ Unknown")
		}
		return styleOffline.Render("○ " + state)
	}
}
