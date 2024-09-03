package components

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/keybase/go-keychain"
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
	form       *huh.Form
	identifier huh.Field
	password   huh.Field
	token      huh.Field
	error      string
}

// NewLogin creates a new Login component
func NewLogin() *Login {
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
	switch msg := msg.(type) {
	case error:
		return login, login.HandleError(msg)
	case *atproto.ServerCreateSession_Output:
		return login, login.HandlerSucessfulLogin(msg)
	}

	form, cmd := login.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		login.form = f
	}

	return login, cmd
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
		login.Width, login.Height,
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

	output, err := atproto.ServerCreateSession(context.Background(), Client, &atproto.ServerCreateSession_Input{
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

// HandlerSucessfulLogin handles successful login storing the session token
// TODO: add error handling
func (login *Login) HandlerSucessfulLogin(session *atproto.ServerCreateSession_Output) tea.Cmd {
	Client.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Handle:     session.Handle,
		Did:        session.Did,
	}

	data, _ := json.Marshal(Client.Auth)

	authItem := keychain.NewItem()
	authItem.SetSecClass(keychain.SecClassGenericPassword)
	authItem.SetService("Monosky")
	authItem.SetAccount(session.Handle)
	authItem.SetAuthenticationType(keychain.AuthenticationTypeKey)
	authItem.SetData(data)
	authItem.SetAccessible(keychain.AccessibleWhenUnlocked)

	keychain.DeleteItem(authItem)
	keychain.AddItem(authItem)

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
