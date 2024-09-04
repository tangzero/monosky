package components

import (
	"context"
	"fmt"
	"strings"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
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
	Component
	parent     tea.Model
	form       *huh.Form
	identifier huh.Field
	password   huh.Field
	token      huh.Field
	error      string
}

// NewLogin creates a new Login component
func NewLogin(parent tea.Model) *Login {
	username := huh.NewInput().
		Key("identifier").
		Title("Identifier").
		Validate(required("Identifier"))

	password := huh.NewInput().
		Key("password").
		Title("Password").
		Validate(required("Password")).
		EchoMode(huh.EchoModePassword)

	token := huh.NewInput().
		Key("token").
		Title("Sign in code").
		Validate(required("Sign in code"))

	return &Login{
		parent:     parent,
		form:       huh.NewForm(huh.NewGroup(username, password)),
		identifier: username,
		password:   password,
		token:      token,
	}
}

// OnResize is called when the terminal is resized
func (login *Login) Init() tea.Cmd {
	login.form.SubmitCmd = login.DoLogin
	return login.form.Init()
}

// Update is called when a message is received
func (login *Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	component, cmd := login.Component.Update(msg)
	login.Component = *component.(*Component)

	switch msg := msg.(type) {
	case error:
		return login, login.HandleError(msg)
	case *atproto.ServerCreateSession_Output:
		return login.parent, login.HandleSucessfulLogin(msg)
	}

	form, formCmd := login.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		login.form = f
	}

	return login, tea.Batch(cmd, formCmd)
}

// View is called when the component should render
func (login *Login) View() string {
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

	error := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Render(login.error)

	return lipgloss.Place(
		WindowWidth, WindowHeight,
		lipgloss.Center, 0.8,
		lipgloss.JoinVertical(lipgloss.Center, logo, form, error),
	)
}

// DoLogin performs the login operation
func (login *Login) DoLogin() tea.Msg {
	login.error = ""

	identifier := login.form.GetString("identifier")
	password := login.form.GetString("password")
	token := login.form.GetString("token")

	output, err := atproto.ServerCreateSession(context.Background(), DefaultClient.xrpc, &atproto.ServerCreateSession_Input{
		Identifier:      identifier,
		Password:        password,
		AuthFactorToken: &token,
	})
	if err != nil {
		return err
	}
	return output
}

// HandleError handles errors that occur during login
// please note that this is a very basic error handling
func (login *Login) HandleError(err error) tea.Cmd {
	login.error = err.Error()
	login.form = huh.NewForm(huh.NewGroup(login.identifier, login.password)) // redisplay the login form
	login.form.SubmitCmd = login.DoLogin
	cmd := login.form.Init()

	if err, ok := err.(*xrpc.Error); ok {
		if err, ok := err.Wrapped.(*xrpc.XRPCError); ok {
			login.error = err.Message
			switch login.error {
			case "A sign in code has been sent to your email address", "Token is invalid":
				login.form = huh.NewForm(huh.NewGroup(login.identifier, login.password, login.token)) // request the auth token
				login.form.SubmitCmd = login.DoLogin
				return tea.Batch(login.form.Init(), login.form.NextField(), login.form.NextField()) // ugly way to focus the token field
			}
			return cmd
		}
		login.error = err.Wrapped.Error()
		return cmd
	}

	return cmd
}

// HandleSucessfulLogin handles successful login storing the session token
func (login *Login) HandleSucessfulLogin(output *atproto.ServerCreateSession_Output) tea.Cmd {
	DefaultClient.SaveAuthInfo(&xrpc.AuthInfo{
		AccessJwt:  output.AccessJwt,
		RefreshJwt: output.RefreshJwt,
		Handle:     output.Handle,
		Did:        output.Did,
	})
	return nil
}

func required(field string) func(value string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", field)
		}
		return nil
	}
}
