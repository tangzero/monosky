package components

import tea "github.com/charmbracelet/bubbletea"

// Component is an interface that all components must implement
type Component interface {
	tea.Model
	OnResize(width, height int) tea.Cmd
}

// Ensure BaseComponent implements Component
var _ Component = new(BaseComponent)

// BaseComponent is a base struct for all components
type BaseComponent struct {
	Width  int
	Height int
}

// OnResize is called when the terminal is resized
func (c *BaseComponent) OnResize(width, height int) tea.Cmd {
	c.Width = width
	c.Height = height
	return nil
}

// Init is called when the component is initialized
func (c *BaseComponent) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received
func (c *BaseComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

// View is called when the component should render
func (c *BaseComponent) View() string {
	return ""
}
