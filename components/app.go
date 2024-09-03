package components

import (
	tea "github.com/charmbracelet/bubbletea"
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
	DefaultClient.LoadAuthInfo()
	loginCmd := app.login.Init()
	return tea.Batch(app.ClearScreen, loginCmd)
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
	if !DefaultClient.LoggedIn() {
		return app.login.View()
	}
	return DefaultClient.Auth.Handle
}

// ClearScreen clears the screen
func (app *App) ClearScreen() tea.Msg {
	return tea.ClearScreen()
}
