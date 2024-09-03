package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// App is the main application component
type App struct {
	BaseComponent
	// login Login
}

// NewApp creates a new App component
func NewApp() *App {
	// width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	// var app App
	// app.OnResize(width, height)
	return &App{}
}

// Init is called when the component is initialized
func (app *App) Init() tea.Cmd {
	// TODO: initialize the child components
	return tea.Batch(app.ClearScreen)
}

// OnResize is called when the terminal is resized
func (app *App) OnResize(width, height int) tea.Cmd {
	baseCmd := app.BaseComponent.OnResize(width, height)
	// TODO: resize the child components
	return tea.Batch(baseCmd, app.ClearScreen)
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

	// TODO: update the child components
	return app, nil
}

// View is called when the component should render
func (app *App) View() string {
	// TODO: render the child components
	return fmt.Sprintf("width: %d, height: %d", app.Width, app.Height)
}

// ClearScreen clears the screen
func (app *App) ClearScreen() tea.Msg {
	return tea.ClearScreen()
}
