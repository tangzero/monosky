package components

import tea "github.com/charmbracelet/bubbletea"

var (
	WindowWidth  int
	WindowHeight int
)

// Ensure Component implements tea.Model
var _ tea.Model = new(Component)

// Component is a base struct for all components
type Component struct{}

// Init is called when the component is initialized
func (c *Component) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received
func (c *Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		}
	case tea.WindowSizeMsg:
		WindowWidth = msg.Width
		WindowHeight = msg.Height
		return c, tea.ClearScreen
	}
	return c, nil
}

// View is called when the component should render
func (c *Component) View() string {
	return ""
}
