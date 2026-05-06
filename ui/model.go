package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lesterhan/clairscale/tailscale"
)

type screen int

const (
	screenLoading screen = iota
	screenDaemonDown
	screenDashboard
	screenPeerList
)

type socketCheckMsg struct{ err error }
type statusMsg struct {
	status *tailscale.Status
	err    error
}
type tickMsg time.Time

type Model struct {
	screen    screen
	width     int
	height    int
	status    *tailscale.Status
	lastFetch time.Time
	fetchErr  error
	peerList  peerListModel
}

func Initial() Model {
	return Model{screen: screenLoading}
}

func (m Model) Init() tea.Cmd {
	return checkSocketCmd()
}

func checkSocketCmd() tea.Cmd {
	return func() tea.Msg {
		return socketCheckMsg{err: tailscale.CheckSocket()}
	}
}

func fetchStatusCmd() tea.Cmd {
	return func() tea.Msg {
		s, err := tailscale.FetchStatus()
		return statusMsg{status: s, err: err}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global messages handled before screen delegation.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.peerList.viewHeight = msg.Height
		return m, nil

	case backMsg:
		m.screen = screenDashboard
		return m, nil

	case socketCheckMsg:
		if msg.err != nil {
			m.screen = screenDaemonDown
			return m, nil
		}
		return m, tea.Batch(fetchStatusCmd(), tickCmd())

	case statusMsg:
		m.lastFetch = time.Now()
		if msg.err != nil {
			m.fetchErr = msg.err
		} else {
			m.status = msg.status
			m.fetchErr = nil
			switch m.screen {
			case screenLoading:
				m.screen = screenDashboard
			case screenPeerList:
				m.peerList = m.peerList.refreshPeers(m.status)
			}
		}
		return m, nil

	case tickMsg:
		return m, tea.Batch(fetchStatusCmd(), tickCmd())
	}

	// Screen-specific key handling.
	switch m.screen {
	case screenPeerList:
		var cmd tea.Cmd
		m.peerList, cmd = m.peerList.update(msg)
		return m, cmd

	default:
		key, ok := msg.(tea.KeyMsg)
		if !ok {
			return m, nil
		}
		switch key.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "p":
			if m.screen == screenDashboard && m.status != nil {
				m.peerList = newPeerListModel(m.status, m.height)
				m.screen = screenPeerList
			}
		}
	}

	return m, nil
}
