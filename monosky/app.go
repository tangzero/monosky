package monosky

import (
	"context"

	"github.com/bluesky-social/indigo/api/bsky"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
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
	posts []*Post
}

// NewApp creates a new App component
func NewApp() *App {
	return &App{}
}

// Init is called when the component is initialized
func (app *App) Init() tea.Cmd {
	return tea.Batch(DefaultClient.LoadAuthInfo(), app.FetchTimeline)
}

// Update is called when a message is received
func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	component, cmd := app.Component.Update(msg)
	app.Component = *component.(*Component)

	switch msg := msg.(type) {
	case error:
		panic(msg)
	case AskToLoginMsg:
		login := NewLogin(app)
		return login, login.Init()
	case *bsky.FeedGetTimeline_Output:
		return app, app.HandleTimelineChange(msg)
	}

	return app, cmd
}

// View is called when the component should render
func (app *App) View() string {
	if app.posts == nil {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#874BF4")).
			Render(Logo)
	}
	return lipgloss.JoinVertical(lipgloss.Top, lo.Map(app.posts, func(post *Post, _ int) string { return post.View() })...)
}

func (app *App) FetchTimeline() tea.Msg {
	output, err := bsky.FeedGetTimeline(context.Background(), DefaultClient.xrpc, "", "", 5)
	if err != nil {
		return err
	}
	return output
}

func (app *App) HandleTimelineChange(output *bsky.FeedGetTimeline_Output) tea.Cmd {
	var cmds tea.Cmd
	app.posts = make([]*Post, len(output.Feed))
	for i, post := range output.Feed {
		app.posts[i] = NewPost(post)
		cmds = tea.Batch(cmds, app.posts[i].Init())
	}
	return cmds
}
