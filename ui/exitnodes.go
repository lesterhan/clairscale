package ui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lesterhan/clairscale/tailscale"
)

type setExitNodeMsg struct{ err error }

type exitNodeRow struct {
	hostName string
	ip       string
	active   bool
	online   bool
}

type exitNodeModel struct {
	rows       []exitNodeRow
	cursor     int
	offset     int
	viewHeight int
	setting    bool
	err        string
}

func newExitNodeModel(status *tailscale.Status, height int) exitNodeModel {
	return exitNodeModel{
		rows:       buildExitNodeRows(status),
		viewHeight: height,
	}
}

func buildExitNodeRows(status *tailscale.Status) []exitNodeRow {
	activeIP := ""
	for _, p := range status.Peer {
		if p.ExitNode && len(p.TailscaleIPs) > 0 {
			activeIP = p.TailscaleIPs[0]
			break
		}
	}

	rows := []exitNodeRow{
		{hostName: "none", ip: "", active: activeIP == "", online: true},
	}

	var candidates []tailscale.PeerStatus
	for _, p := range status.Peer {
		if p.ExitNodeOption {
			candidates = append(candidates, p)
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].HostName < candidates[j].HostName
	})

	for _, p := range candidates {
		ip := ""
		if len(p.TailscaleIPs) > 0 {
			ip = p.TailscaleIPs[0]
		}
		rows = append(rows, exitNodeRow{
			hostName: p.HostName,
			ip:       ip,
			active:   ip != "" && ip == activeIP,
			online:   p.Online,
		})
	}
	return rows
}

func setExitNodeCmd(ip string) tea.Cmd {
	return func() tea.Msg {
		return setExitNodeMsg{err: tailscale.SetExitNode(ip)}
	}
}

func (m exitNodeModel) listHeight() int {
	h := m.viewHeight - 7
	if h < 5 {
		return 5
	}
	return h
}

func (m exitNodeModel) clampOffset() exitNodeModel {
	lh := m.listHeight()
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+lh {
		m.offset = m.cursor - lh + 1
	}
	return m
}

func (m exitNodeModel) update(msg tea.Msg) (exitNodeModel, tea.Cmd) {
	if m.setting {
		return m, nil
	}
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m = m.clampOffset()
		}
	case "down", "j":
		if m.cursor < len(m.rows)-1 {
			m.cursor++
			m = m.clampOffset()
		}
	case "enter":
		row := m.rows[m.cursor]
		m.setting = true
		m.err = ""
		return m, setExitNodeCmd(row.ip)
	case "esc":
		return m, func() tea.Msg { return backMsg{} }
	}
	return m, nil
}

func viewExitNodes(m exitNodeModel) string {
	activeLabel := "none"
	for _, r := range m.rows {
		if r.active && r.ip != "" {
			activeLabel = r.hostName
			break
		}
	}

	var b strings.Builder
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "  Active: %s\n", activeLabel)
	fmt.Fprintln(&b)

	lh := m.listHeight()
	end := m.offset + lh
	if end > len(m.rows) {
		end = len(m.rows)
	}

	for i := m.offset; i < end; i++ {
		r := m.rows[i]

		cursor := "  "
		if i == m.cursor {
			cursor = "▶ "
		}

		dot := "○"
		if r.online {
			dot = "●"
		}

		check := ""
		if r.active {
			check = "  ✓"
		}

		if r.ip == "" {
			fmt.Fprintf(&b, "%s%s %s\n", cursor, dot, r.hostName)
		} else {
			fmt.Fprintf(&b, "%s%s %-28s %-16s%s\n", cursor, dot, r.hostName, r.ip, check)
		}
	}

	fmt.Fprintln(&b)
	if m.setting {
		fmt.Fprintln(&b, "  Connecting...")
	} else if m.err != "" {
		fmt.Fprintf(&b, "  Error: %s\n", m.err)
	}

	fmt.Fprintln(&b, "  [enter] connect   [esc] back")
	return b.String()
}
