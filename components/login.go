package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const Logo = `
_|      _|                                          _|                
_|_|  _|_|    _|_|    _|_|_|      _|_|      _|_|_|  _|  _|    _|    _|
_|  _|  _|  _|    _|  _|    _|  _|    _|  _|_|      _|_|      _|    _|
_|      _|  _|    _|  _|    _|  _|    _|      _|_|  _|  _|    _|    _|
_|      _|    _|_|    _|    _|    _|_|    _|_|_|    _|    _|    _|_|_|
                                                                    _|
                                                                _|_|  

                                    Bluesky for monospaced environemts
`

// Login is a component that handles user login
type Login struct {
	BaseComponent
	LoggedIn      bool
	TokenRequired bool
	form          *huh.Form
	username      huh.Field
	password      huh.Field
	token         huh.Field
}

// NewLogin creates a new Login component
func NewLogin() *Login {
	username := huh.NewInput().
		Key("username").
		Title("Username").
		Validate(required("Username"))

	password := huh.NewInput().
		Key("password").
		Title("Password").
		Validate(required("Password")).
		EchoMode(huh.EchoModePassword)

	token := huh.NewInput().
		Key("token").
		Title("Auth Token").
		Validate(required("Auth Token"))

	return &Login{
		form:     huh.NewForm(huh.NewGroup(username, password)), // initialy only username and password are shown
		username: username,
		password: password,
		token:    token,
	}
}

// OnResize is called when the terminal is resized
func (login *Login) Init() tea.Cmd {
	return login.form.Init()
}

// Update is called when a message is received
func (login *Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := login.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		login.form = f
	}

	return login, cmd
}

// View is called when the component should render
func (login *Login) View() string {
	if login.form.State == huh.StateCompleted {
		username := login.form.GetString("username")
		password := login.form.GetString("password")
		return fmt.Sprintf("username: %s\npassword: %s", username, password)
	}

	logo := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#874BF4")).
		Render(Logo)
	width, _ := lipgloss.Size(logo)

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Width(width).
		Padding(1).
		Render(login.form.View())

	return lipgloss.Place(
		login.Width, login.Height,
		lipgloss.Center, 0.8,
		lipgloss.JoinVertical(lipgloss.Center, logo, form),
	)
}

func required(field string) func(value string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", field)
		}
		return nil
	}
}
