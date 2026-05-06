package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lesterhan/clairscale/tailscale"
)

type backMsg struct{}

type peerListModel struct {
	allPeers   []tailscale.PeerStatus
	cursor     int
	offset     int
	viewHeight int
	filterText string
	filtering  bool
}

func newPeerListModel(status *tailscale.Status, height int) peerListModel {
	peers := tailnetPeers(status)
	sortPeers(peers)
	return peerListModel{allPeers: peers, viewHeight: height}
}

// tailnetPeers excludes peers whose DNS name doesn't match the local tailnet suffix
// (filters out Mullvad and other external exit node providers).
func tailnetPeers(status *tailscale.Status) []tailscale.PeerStatus {
	suffix := status.CurrentTailnet.MagicDNSSuffix
	var out []tailscale.PeerStatus
	for _, p := range status.Peer {
		if suffix == "" || strings.HasSuffix(strings.TrimSuffix(p.DNSName, "."), suffix) {
			out = append(out, p)
		}
	}
	return out
}

func sortPeers(peers []tailscale.PeerStatus) {
	sort.Slice(peers, func(i, j int) bool {
		if peers[i].Online != peers[j].Online {
			return peers[i].Online
		}
		return peers[i].HostName < peers[j].HostName
	})
}

func (m peerListModel) visiblePeers() []tailscale.PeerStatus {
	if m.filterText == "" {
		return m.allPeers
	}
	f := strings.ToLower(m.filterText)
	var out []tailscale.PeerStatus
	for _, p := range m.allPeers {
		if strings.Contains(strings.ToLower(p.HostName), f) {
			out = append(out, p)
		}
	}
	return out
}

func (m peerListModel) update(msg tea.Msg) (peerListModel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if m.filtering {
		return m.updateFilter(key)
	}
	visible := m.visiblePeers()
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m = m.clampOffset()
		}
	case "down", "j":
		if m.cursor < len(visible)-1 {
			m.cursor++
			m = m.clampOffset()
		}
	case "/":
		m.filtering = true
	case "esc":
		return m, func() tea.Msg { return backMsg{} }
	}
	return m, nil
}

func (m peerListModel) updateFilter(msg tea.KeyMsg) (peerListModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.filtering = false
		m.filterText = ""
		m.cursor = 0
		m.offset = 0
	case "backspace":
		if len(m.filterText) > 0 {
			m.filterText = m.filterText[:len(m.filterText)-1]
			m.cursor = 0
			m.offset = 0
		}
	case "enter":
		m.filtering = false
	default:
		if len(msg.Runes) == 1 {
			m.filterText += string(msg.Runes)
			m.cursor = 0
			m.offset = 0
		}
	}
	return m, nil
}

func (m peerListModel) listHeight() int {
	h := m.viewHeight - 6
	if h < 5 {
		return 5
	}
	return h
}

func (m peerListModel) clampOffset() peerListModel {
	lh := m.listHeight()
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+lh {
		m.offset = m.cursor - lh + 1
	}
	return m
}

func (m peerListModel) refreshPeers(status *tailscale.Status) peerListModel {
	peers := tailnetPeers(status)
	sortPeers(peers)
	m.allPeers = peers
	visible := m.visiblePeers()
	if len(visible) > 0 && m.cursor >= len(visible) {
		m.cursor = len(visible) - 1
	}
	return m.clampOffset()
}

func viewPeerList(m peerListModel) string {
	visible := m.visiblePeers()

	online, offline := 0, 0
	for _, p := range m.allPeers {
		if p.Online {
			online++
		} else {
			offline++
		}
	}

	var b strings.Builder
	fmt.Fprintln(&b)

	if m.filtering {
		fmt.Fprintf(&b, "  %d online · %d offline    filter: %s_\n", online, offline, m.filterText)
	} else {
		fmt.Fprintf(&b, "  %d online · %d offline    [/] filter\n", online, offline)
	}
	fmt.Fprintln(&b)

	end := m.offset + m.listHeight()
	if end > len(visible) {
		end = len(visible)
	}

	for i := m.offset; i < end; i++ {
		p := visible[i]

		cursor := "  "
		if i == m.cursor {
			cursor = "▶ "
		}

		dot := "○"
		if p.Online {
			dot = "●"
		}

		ip := ""
		if len(p.TailscaleIPs) > 0 {
			ip = p.TailscaleIPs[0]
		}

		extra := ""
		if !p.Online && !p.LastSeen.IsZero() {
			extra = "  " + formatLastSeen(p.LastSeen)
		}

		fmt.Fprintf(&b, "%s%s %-20s %-16s %-10s%s\n",
			cursor, dot, p.HostName, ip, p.OS, extra)
	}

	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "  [enter] ping   [y] copy IP   [esc] back")
	return b.String()
}

func formatLastSeen(t time.Time) string {
	ago := time.Since(t)
	switch {
	case ago < time.Minute:
		return fmt.Sprintf("(%ds ago)", int(ago.Seconds()))
	case ago < time.Hour:
		return fmt.Sprintf("(%dm ago)", int(ago.Minutes()))
	case ago < 24*time.Hour:
		return fmt.Sprintf("(%dh ago)", int(ago.Hours()))
	default:
		return fmt.Sprintf("(%dd ago)", int(ago.Hours()/24))
	}
}
