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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

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
			if m.screen == screenLoading {
				m.screen = screenDashboard
			}
		}

	case tickMsg:
		return m, tea.Batch(fetchStatusCmd(), tickCmd())
	}

	return m, nil
}
