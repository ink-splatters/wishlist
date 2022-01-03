package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/accesscontrol"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/charmbracelet/wishlist"
	"github.com/gliderlabs/ssh"
)

func main() {
	if err := wishlist.Serve(&wishlist.Config{
		Listen: "127.0.0.1",
		Port:   2222,
		Factory: func(e wishlist.Endpoint) (*ssh.Server, error) {
			return wish.NewServer(
				wish.WithAddress(e.Address),
				wish.WithMiddleware(
					bm.Middleware(e.Handler),
					lm.Middleware(),
					accesscontrol.Middleware(),
					activeterm.Middleware(),
				),
			)
		},
		Endpoints: []*wishlist.Endpoint{
			{
				Name:    "foo bar",
				Address: "some.other.server:2222",
			},
			{
				Name: "entries without handlers and without addresses are ignored",
			},
			{
				Address: "entries without names are ignored",
			},
			{
				Name: "example app",
				Handler: func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
					return initialModel(), []tea.ProgramOption{}
				},
			},
		},
	}); err != nil {
		panic(err)
	}
}

type model struct {
	spinner  spinner.Model
	quitting bool
}

func initialModel() model {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}