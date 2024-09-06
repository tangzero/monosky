package monosky

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/bluesky-social/indigo/api/bsky"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var cleanUserNameRegex = regexp.MustCompile("[[:^ascii:]]")

type Post struct {
	Component
	*bsky.FeedDefs_FeedViewPost
	image *Image
}

func NewPost(post *bsky.FeedDefs_FeedViewPost) *Post {
	return &Post{FeedDefs_FeedViewPost: post}
}

func (post *Post) Init() tea.Cmd {
	if post.Post.Embed != nil && post.Post.Embed.EmbedImages_View != nil && len(post.Post.Embed.EmbedImages_View.Images) > 0 {
		post.image = NewImage(post.Post.Embed.EmbedImages_View.Images[0].Fullsize)
		return post.image.Load
	}
	return nil
}

func (post *Post) View() string {
	switch record := post.Post.Record.Val.(type) {
	case *bsky.FeedPost:
		return lipgloss.NewStyle().
			Width(60).
			Border(lipgloss.RoundedBorder(), true, false, false, true).
			Padding(0, 1).
			Render(lipgloss.JoinVertical(lipgloss.Top, post.author(), record.Text, post.embed()))
	}
	return reflect.TypeOf(post.Post.Record.Val).String()
}

func (post *Post) displayName() string {
	if post.Post.Author.DisplayName == nil {
		return ""
	}
	displayName := strings.TrimSpace(cleanUserNameRegex.ReplaceAllLiteralString(*post.Post.Author.DisplayName, ""))
	if displayName == "" {
		return displayName
	}
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Margin(0, 1, 1, 0).
		Render(displayName)
}

func (post *Post) handle() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AEBBC9")).
		Render("@" + post.Post.Author.Handle)
}

func (post *Post) author() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, post.displayName(), post.handle())
}

func (post *Post) embed() string {
	if post.image != nil {
		return post.image.View()
	}
	return ""
}
