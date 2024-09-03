package components

import (
	"encoding/json"

	"github.com/bluesky-social/indigo/xrpc"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/keybase/go-keychain"
)

// App is the main application component
type App struct {
	BaseComponent
	login *Login
}

// NewApp creates a new App component
func NewApp() *App {
	return &App{
		login: NewLogin(),
	}
}

// Init is called when the component is initialized
func (app *App) Init() tea.Cmd {
	loadCmd := app.LoadAuthInfo()
	loginCmd := app.login.Init()
	return tea.Batch(app.ClearScreen, loadCmd, loginCmd)
}

// OnResize is called when the terminal is resized
func (app *App) OnResize(width, height int) tea.Cmd {
	app.BaseComponent.OnResize(width, height)
	loginCmd := app.login.OnResize(width, height)
	return tea.Batch(app.ClearScreen, loginCmd)
}

// Update is called when a message is received
func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return app, tea.Quit
		}
	case tea.WindowSizeMsg:
		return app, app.OnResize(msg.Width, msg.Height)
	}

	loginModel, loginCmd := app.login.Update(msg)
	app.login = loginModel.(*Login)

	return app, tea.Batch(loginCmd)
}

// View is called when the component should render
func (app *App) View() string {
	if !app.LoggedIn() {
		return app.login.View()
	}
	return Client.Auth.Handle
}

// ClearScreen clears the screen
func (app *App) ClearScreen() tea.Msg {
	return tea.ClearScreen()
}

// LoadAuthInfo loads the authentication info from the keychain
func (app *App) LoadAuthInfo() tea.Cmd {
	queryItem := keychain.NewItem()
	queryItem.SetSecClass(keychain.SecClassGenericPassword)
	queryItem.SetService("Monosky")
	queryItem.SetMatchLimit(keychain.MatchLimitOne)
	queryItem.SetReturnAttributes(true)
	queryItem.SetReturnData(true)

	results, err := keychain.QueryItem(queryItem)
	if err != nil || len(results) == 0 {
		return nil
	}

	var auth xrpc.AuthInfo
	if err := json.Unmarshal(results[0].Data, &auth); err != nil {
		return func() tea.Msg { return err }
	}
	Client.Auth = &auth

	return nil
}

// LoggedIn returns true if we're able to load the authentication info
func (app *App) LoggedIn() bool {
	return Client.Auth != nil
}
