// internal/ui/layout.go
package layout

import (
	"github.com/subhasundardass/tuix/internal/context"
	appctx "github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/tuix"
)

type LayoutProps struct {
	Ctx     *appctx.AppContext
	Title   string
	Content tuix.Element
}

func MasterLayout(ctx *context.AppContext, props LayoutProps) tuix.Element {

	// appctx := context.Use() // ← one line, done
	// ctx.NavigateTo("home")
	// ctx.ToggleDark()
	// _ = ctx.Page()

	if ctx == nil {
		return tuix.Text("ERROR: ctx is nil in MasterLayout", tuix.NewStyle().Foreground(tuix.Red))
	}

	//Main content with border and title
	borderChars := tuix.BorderRounded
	if tuix.IsFocused("content") {
		borderChars = tuix.BorderDouble
	}
	mainContent := tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     tuix.Grow(1),
			Gap:       1,
		},
		tuix.NewStyle().Border(tuix.Border{
			Top:    true,
			Right:  true,
			Bottom: true,
			Left:   true,
			Chars:  borderChars,
			Color:  tuix.BrightBlack,
			Title:  props.Title,
		}),
		props.Content, // ← Your page content goes here
	)

	//Body: Sidebar + Main content (both should grow)
	body := tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Gap:       0,
			Width:     tuix.Grow(1),
			Height:    tuix.Grow(1), // ← Fill remaining height
		},
		tuix.NewStyle(),
		SidebarTree(ctx, tuix.Props{}), // Sidebar
		mainContent,                    // Main content
	)

	//Full layout: Header + Body + Footer
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       0,
			Padding:   [4]int{0, 0, 0, 0},
			Width:     tuix.Grow(1),
			Height:    tuix.Grow(1), // ← Fill entire screen
		},
		tuix.NewStyle(),
		Header(tuix.Props{}), // Fixed height
		body,                 // Takes remaining space
		Footer(tuix.Props{}), // Fixed height
	)
}
