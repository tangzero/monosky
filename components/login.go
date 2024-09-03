package components

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Login struct {
	form  *huh.Form
	width int
}

func NewLogin() Login {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	return Login{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("username").
					Title("Username"),
				huh.NewInput().
					Key("password").
					Title("Password").
					EchoMode(huh.EchoModePassword),
			),
		),
		width: width,
	}
}

func (login Login) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg { return tea.ClearScreen() }, login.form.Init())
}

func (login Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == "ctrl+c" {
		return login, tea.Quit
	}

	form, cmd := login.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		login.form = f
	}

	return login, cmd
}

func (login Login) View() string {
	if login.form.State == huh.StateCompleted {
		username := login.form.GetString("username")
		password := login.form.GetString("password")
		return fmt.Sprintf("username: %s\npassword: %s", username, password)
	}

	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	form := lipgloss.NewStyle().Width(80).Render(login.form.View())
	title := lipgloss.NewStyle().Width(40).Background(subtle).Align(lipgloss.Center).Render("Login")
	ui := lipgloss.JoinVertical(lipgloss.Center, title, form)

	return lipgloss.Place(login.width, 20,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		// lipgloss.WithWhitespaceChars("😄"),
		// lipgloss.WithWhitespaceForeground(subtle),
	)
}

var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
