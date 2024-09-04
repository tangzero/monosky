package monosky

import (
	tea "github.com/charmbracelet/bubbletea"
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

// App is the main application component
type App struct {
	Component
}

// NewApp creates a new App component
func NewApp() *App {
	return &App{}
}

// Init is called when the component is initialized
func (app *App) Init() tea.Cmd {
	return DefaultClient.LoadAuthInfo()
}

// Update is called when a message is received
func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	component, cmd := app.Component.Update(msg)
	app.Component = *component.(*Component)

	switch msg.(type) {
	case AskToLoginMsg:
		login := NewLogin(app)
		return login, login.Init()
	}

	return app, cmd
}

// View is called when the component should render
func (app *App) View() string {
	return Logo
}
