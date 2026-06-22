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
	contentHeight := height - 2 // Header + Footer

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       0,
			Width:     tuix.Fixed(width),
			Height:    tuix.Fixed(height),
		},
		tuix.NewStyle().Background(tuix.Black),

		// Header
		renderHeader(props.Title, props.Dark),

		// ⭐ Content area with clear layer
		tuix.Box(
			tuix.Props{
				Direction: tuix.Column,
				Width:     tuix.Fixed(width),
				Height:    tuix.Fixed(contentHeight),
				Padding:   [4]int{1, 2, 1, 2},
			},
			tuix.NewStyle().
				Background(tuix.Black),

			// ⭐ KEY FIX: First render a clear box (spaces)
			// This overwrites old content with empty space
			clearContent(width-4, contentHeight-2),

			// Then render the new content on top
			props.Content,
		),

		// Footer
		renderFooter(props.Dark),
	)
}

// ⭐ NEW: ClearContent creates a box filled with spaces
func clearContent(width, height int) tuix.Element {
	// Create a blank space string for one row
	spaceRow := tuix.Text(strings.Repeat(" ", width), tuix.NewStyle())

	// Build a column of space rows
	children := []tuix.Element{}
	for i := 0; i < height; i++ {
		children = append(children, spaceRow)
	}

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Fixed(width),
			Height:    tuix.Fixed(height),
		},
		tuix.NewStyle().Background(tuix.Black),
		children...,
	)
}

func renderHeader(title string, dark bool) tuix.Element {
	bgColor := tuix.Blue
	textColor := tuix.White

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Padding:   [4]int{1, 2, 1, 2},
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
