package app

import (
	"github.com/subhasundardass/tuix/internal/ui"
	"github.com/subhasundardass/tuix/internal/ui/pages"
	"github.com/subhasundardass/tuix/tuix"
)

// Root is the main application component
// It's exported so main can use it
func Root(props tuix.Props) tuix.Element {
	// Get current context
	ctx := buildContext()

	// Local state for UI
	page, setPage := tuix.UseState(ctx.Page)
	dark, setDark := tuix.UseState(ctx.Dark)

	// Handle keyboard input
	if tuix.CurrentKey.Rune != 0 {
		switch tuix.CurrentKey.Rune {
		case '1':
			setPage("home")
			ctx.NavigateTo("home")
		case '2':
			setPage("settings")
			ctx.NavigateTo("settings")
		case '3':
			setPage("about")
			ctx.NavigateTo("about")
		case 't':
			setDark(!dark)
			ctx.ToggleTheme()
		case 'q':
			return tuix.Text("", tuix.NewStyle())
		}
	}

	if tuix.CurrentKey.Code == tuix.KeyEscape {
		setPage("home")
		ctx.NavigateTo("home")
	}

	if tuix.CurrentKey.Code == tuix.KeyCtrlC {
		return tuix.Text("Goodbye!", tuix.NewStyle())
	}

	// Choose page
	var content tuix.Element
	switch page {
	case "home":
		content = pages.HomePage(tuix.Props{})
	case "settings":
		content = pages.SettingsPage(tuix.Props{})
	case "about":
		content = pages.AboutPage(tuix.Props{})
	default:
		content = tuix.Text("404 - Page Not Found", tuix.NewStyle())
	}

	// Provide context
	return DefaultContext.Provide(ctx, func() tuix.Element {
		return ui.Layout(ui.LayoutProps{
			Title:   ctx.Config.AppName,
			Content: content,
			Dark:    dark,
			// User:    ctx.User,
		})
	})
}
