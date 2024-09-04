package monosky

import (
	"reflect"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/charmbracelet/lipgloss"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

type Post struct {
	Component
	*bsky.FeedDefs_FeedViewPost
}

func NewPost(post *bsky.FeedDefs_FeedViewPost) *Post {
	return &Post{FeedDefs_FeedViewPost: post}
}

func (post *Post) View() string {
	switch record := post.Post.Record.Val.(type) {
	case *bsky.FeedPost:
		return lipgloss.NewStyle().
			Width(60).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Render(lipgloss.JoinVertical(lipgloss.Top, post.author(), record.Text, post.embed()))
	}
	return reflect.TypeOf(post.Post.Record.Val).String()
}

func (post *Post) displayName() string {
	if post.Post.Author.DisplayName == nil {
		return ""
	}
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Margin(0, 1, 1, 0).
		Render(*post.Post.Author.DisplayName)
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
	if post.Post.Embed == nil || post.Post.Embed.EmbedImages_View == nil || len(post.Post.Embed.EmbedImages_View.Images) == 0 {
		return ""
	}
	image := post.Post.Embed.EmbedImages_View.Images[0]

	flags := aic_package.Flags{
		Width:               58,
		Colored:             true,
		CharBackgroundColor: true,
		Braille:             false,
		Dither:              false,
	}

	asciiArt, err := aic_package.Convert(image.Fullsize, flags)
	if err != nil {
		return lipgloss.NewStyle().
			Bold(true).
			Margin(1, 0).
			Foreground(lipgloss.Color("#FF0000")).
			Render(err.Error())
	}
	return asciiArt

	// return lipgloss.NewStyle().
	// 	Bold(true).
	// 	Margin(1, 0).
	// 	Foreground(lipgloss.Color("#FF00FF")).
	// 	Render()
}
