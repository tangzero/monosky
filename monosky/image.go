package monosky

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
)

const (
	ImageWidth                 = 58
	BackgroundTransparentColor = "\x1b[0;39;49m"
	BackgroundRGBColor         = "\x1b[48;2;%d;%d;%dm"
	ForegroundTransparentColor = "\x1b[0m "
	ForegroundRGBColor         = "\x1b[38;2;%d;%d;%dm▄"
	Reset                      = "\x1b[0m"
)

type Image struct {
	image.Image
	URL  string
	Data string
}

func NewImage(url string) *Image {
	return &Image{URL: url}
}

func (img *Image) View() string {
	return lipgloss.NewStyle().
		Margin(1, 0).
		Render(img.Data)
}

func (img *Image) Load() tea.Msg {
	data, err := img.download()
	if err != nil {
		return err
	}
	img.Image, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	img.escape(img.scale())
	return ImageLoadedMsg{}
}

func (image *Image) download() ([]byte, error) {
	resp, err := http.Get(image.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}
	return data, nil
}

func (img *Image) scale() image.Image {
	size := img.Bounds().Size()
	scale := float64(ImageWidth) / float64(size.X)
	w := ImageWidth
	h := int(float64(size.Y) * scale)
	return imaging.Fit(img, w, h, imaging.Lanczos)
}

func (img *Image) escape(in image.Image) {
	var sb strings.Builder
	size := in.Bounds().Size()
	for y := 0; y < size.Y; y += 2 {
		for x := 0; x < size.X; x++ {
			r, g, b, a := in.At(x, y).RGBA()
			if a>>8 < 128 {
				sb.WriteString(BackgroundTransparentColor)
			} else {
				sb.WriteString(fmt.Sprintf(BackgroundRGBColor, r>>8, g>>8, b>>8))
			}
			r, g, b, a = in.At(x, y+1).RGBA()
			if a>>8 < 128 {
				sb.WriteString(ForegroundTransparentColor)
			} else {
				sb.WriteString(fmt.Sprintf(ForegroundRGBColor, r>>8, g>>8, b>>8))
			}
		}
		sb.WriteString(Reset)
		sb.WriteString("\n")
	}
	img.Data = sb.String()
	img.Data = img.Data[:len(img.Data)-1] // remove the last newline
}
