// internal/ui/layout.go
package ui

import (
	"strings"

	"github.com/subhasundardass/tuix/tuix"
)

type LayoutProps struct {
	Title   string
	Content tuix.Element
	Dark    bool
}

func Layout(props LayoutProps) tuix.Element {
	width := 80
	height := 24

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Fixed(width),
			Height:    tuix.Fixed(height),
			Gap:       0,
		},
		tuix.NewStyle().Background(tuix.Black),

		// Header
		renderHeader(props.Title, props.Dark),

		// Content fills remaining height
		tuix.Box(
			tuix.Props{
				Width:  tuix.Grow(1),
				Height: tuix.Grow(1),
			},
			tuix.NewStyle(),
			props.Content,
		),

		// Footer
		tuix.Box(
			tuix.Props{
				Width:  tuix.Grow(1),
				Height: tuix.Fixed(1),
			},
			tuix.NewStyle(),

			tuix.Text("F1 Help  F2 Menu  F10 Exit", tuix.NewStyle()),
		),
	)
}

func renderHeader(title string, dark bool) tuix.Element {
	bgColor := tuix.Blue
	textColor := tuix.White

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Padding:   [4]int{0, 2, 0, 2},
		},
		tuix.NewStyle().Background(bgColor),
		tuix.Text(" "+title+" ",
			tuix.NewStyle().Bold(true).Foreground(textColor)),
		tuix.Text("  [1]Home  [2]Settings  [3]About  [t]Theme  [q]Quit",
			tuix.NewStyle().Foreground(textColor).Italic(true)),
	)
}

func renderFooter(dark bool) tuix.Element {
	textColor := tuix.White

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Padding:   [4]int{1, 2, 1, 2},
		},
		tuix.NewStyle().Background(tuix.BrightBlack),
		tuix.Text(" 👤 Guest | Press Ctrl+C to quit ",
			tuix.NewStyle().Foreground(textColor)),
	)
}

func renderEmptyScreen(width, height int) tuix.Element {
	// Create a box filled with spaces
	spaceRow := strings.Repeat(" ", width)
	children := []tuix.Element{}
	for i := 0; i < height; i++ {
		children = append(children, tuix.Text(spaceRow, tuix.NewStyle()))
	}

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Fixed(width),
			Height:    tuix.Fixed(height),
		},
		tuix.NewStyle(),
		children...,
	)
}
