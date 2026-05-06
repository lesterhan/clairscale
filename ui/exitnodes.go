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
	country  string
	city     string
	active   bool
	online   bool
}

type exitNodeModel struct {
	allRows    []exitNodeRow
	cursor     int
	offset     int
	viewHeight int
	filterText string
	setting    bool
	err        string
}

func newExitNodeModel(status *tailscale.Status, height int) exitNodeModel {
	return exitNodeModel{
		allRows:    buildExitNodeRows(status),
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
		ci, cj := locationFields(candidates[i]), locationFields(candidates[j])
		if ci[0] != cj[0] {
			return ci[0] < cj[0]
		}
		if ci[1] != cj[1] {
			return ci[1] < cj[1]
		}
		return candidates[i].HostName < candidates[j].HostName
	})

	for _, p := range candidates {
		ip := ""
		if len(p.TailscaleIPs) > 0 {
			ip = p.TailscaleIPs[0]
		}
		country, city := locationFields(p)[0], locationFields(p)[1]
		rows = append(rows, exitNodeRow{
			hostName: p.HostName,
			ip:       ip,
			country:  country,
			city:     city,
			active:   ip != "" && ip == activeIP,
			online:   p.Online,
		})
	}
	return rows
}

func locationFields(p tailscale.PeerStatus) [2]string {
	if p.Location == nil {
		return [2]string{"", ""}
	}
	return [2]string{p.Location.Country, p.Location.City}
}

// filteredRows always keeps the "none" row at top; filter applies to exit node peers only.
func (m exitNodeModel) filteredRows() []exitNodeRow {
	if m.filterText == "" {
		return m.allRows
	}
	f := strings.ToLower(m.filterText)
	out := []exitNodeRow{m.allRows[0]}
	for _, r := range m.allRows[1:] {
		if strings.Contains(strings.ToLower(r.hostName), f) ||
			strings.Contains(strings.ToLower(r.country), f) ||
			strings.Contains(strings.ToLower(r.city), f) {
			out = append(out, r)
		}
	}
	return out
}

func setExitNodeCmd(ip string) tea.Cmd {
	return func() tea.Msg {
		return setExitNodeMsg{err: tailscale.SetExitNode(ip)}
	}
}

func (m exitNodeModel) listHeight() int {
	h := m.viewHeight - 8
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

	visible := m.filteredRows()

	switch key.String() {
	case "up":
		if m.cursor > 0 {
			m.cursor--
			m = m.clampOffset()
		}
	case "down":
		if m.cursor < len(visible)-1 {
			m.cursor++
			m = m.clampOffset()
		}
	case "enter":
		if m.cursor < len(visible) {
			row := visible[m.cursor]
			m.setting = true
			m.err = ""
			return m, setExitNodeCmd(row.ip)
		}
	case "esc":
		if m.filterText != "" {
			m.filterText = ""
			m.cursor = 0
			m.offset = 0
		} else {
			return m, func() tea.Msg { return backMsg{} }
		}
	case "backspace":
		if len(m.filterText) > 0 {
			m.filterText = m.filterText[:len(m.filterText)-1]
			m.cursor = 0
			m.offset = 0
		}
	default:
		if len(key.Runes) == 1 {
			m.filterText += string(key.Runes)
			m.cursor = 0
			m.offset = 0
			m.err = ""
		}
	}
	return m, nil
}

func viewExitNodes(m exitNodeModel) string {
	visible := m.filteredRows()

	activeLabel := styleDim.Render("none")
	for _, r := range m.allRows {
		if r.active && r.ip != "" {
			activeLabel = styleActive.Render(r.hostName)
			break
		}
	}

	var b strings.Builder
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "  %s %s\n", styleLabel.Render("Active:"), activeLabel)
	fmt.Fprintln(&b)

	if m.filterText != "" {
		peers := len(visible) - 1
		suffix := "s"
		if peers == 1 {
			suffix = ""
		}
		fmt.Fprintf(&b, "  %s    %s\n",
			styleFilter.Render("filter: "+m.filterText+"_"),
			styleDim.Render(fmt.Sprintf("%d result%s", peers, suffix)),
		)
	} else {
		fmt.Fprintln(&b, "  "+styleDim.Render("type to filter by name, country, or city"))
	}
	fmt.Fprintln(&b)

	lh := m.listHeight()
	end := m.offset + lh
	if end > len(visible) {
		end = len(visible)
	}

	for i := m.offset; i < end; i++ {
		r := visible[i]
		isSelected := i == m.cursor
		fmt.Fprintln(&b, exitNodeRowStr(r, isSelected))
	}

	fmt.Fprintln(&b)
	if m.setting {
		fmt.Fprintln(&b, "  "+styleWarning.Render("Connecting..."))
	} else if m.err != "" {
		fmt.Fprintf(&b, "  %s\n", styleError.Render("Error: "+m.err))
	}

	fmt.Fprintln(&b, "  "+styleDim.Render("[↑↓] navigate   [enter] connect   [esc] clear/back"))
	return b.String()
}

func exitNodeRowStr(r exitNodeRow, isSelected bool) string {
	dot := styleOffline.Render("○")
	if r.ip == "" {
		// "none" row always gets a plain dot
		name := r.hostName
		if isSelected {
			name = styleBold.Render(name)
		}
		return cursor(isSelected) + dot + " " + name
	}
	if r.online {
		dot = styleOnline.Render("●")
	}

	hn := fmt.Sprintf("%-22s", r.hostName)
	ip := styleDim.Render(fmt.Sprintf("%-16s", r.ip))
	country := fmt.Sprintf("%-18s", r.country)
	city := fmt.Sprintf("%-16s", r.city)
	check := "  "

	if r.active {
		hn = styleActive.Render(hn)
		country = styleActive.Render(country)
		city = styleActive.Render(city)
		check = " " + styleActive.Render("✓")
	} else if isSelected {
		hn = styleBold.Render(hn)
		country = styleBold.Render(country)
		city = styleBold.Render(city)
	} else {
		country = styleDim.Render(country)
		city = styleDim.Render(city)
	}

	return cursor(isSelected) + dot + " " + hn + " " + ip + " " + country + " " + city + check
}
