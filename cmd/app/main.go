// main.go or cmd/app/main.go
package main

import (
	"github.com/subhasundardass/tuix/internal/ui"
	"github.com/subhasundardass/tuix/internal/ui/pages"
	"github.com/subhasundardass/tuix/tuix"
)

func main() {

	// Initialize
	InitApp()

	// Create and run app
	app := tuix.NewApp(80, 24)
	app.Run(Root, tuix.Props{})

}

func Root(props tuix.Props) tuix.Element {

	// Get app state
	appState := GetApp()

	// Build context value from state
	ctx := BuildContext()

	// Local state for UI (syncs with app state)
	// page, setPage := tuix.UseState(appState.GetPage())
	dark, setDark := tuix.UseState(appState.IsDarkMode())

	// Handle keyboard input
	if tuix.CurrentKey.Rune != 0 {
		switch tuix.CurrentKey.Rune {
		case '1':
			// setPage("home")
			appState.SetPage("home")
			ctx.NavigateTo("home") // Also call via context
		case '2':
			// setPage("settings")
			appState.SetPage("settings")
			ctx.NavigateTo("settings")
		case '3':
			// setPage("about")
			appState.SetPage("about")
			ctx.NavigateTo("about")
		case 't':
			setDark(!dark)
			appState.ToggleTheme()
			ctx.ToggleTheme()
		case 'q':
			return tuix.Text("", tuix.NewStyle())
		}

	}

	// if tuix.CurrentKey.Code == tuix.KeyEscape {
	// 	setPage("home")
	// 	appState.SetPage("home")
	// 	ctx.NavigateTo("home")
	// }

	// if tuix.CurrentKey.Code == tuix.KeyCtrlC {
	// 	return tuix.Text("Goodbye!", tuix.NewStyle())
	// }

	// Choose which page to render
	var content tuix.Element
	switch ctx.Page {
	case "home":
		content = pages.HomePage(tuix.Props{})
	case "settings":
		content = pages.SettingsPage(tuix.Props{})
	case "about":
		content = pages.AboutPage(tuix.Props{})
	default:
		content = tuix.Text("404 - Page Not Found", tuix.NewStyle())
	}

	// ⭐ PROVIDE CONTEXT TO ALL CHILDREN ⭐
	// This makes ctx available via UseContext anywhere in the tree
	return DefaultContext.Provide(ctx, func() tuix.Element {
		return ui.Layout(ui.LayoutProps{
			Title:   appState.Config.AppName,
			Content: content,
			Dark:    dark,
			// User:    app.User,
		})
	})
}
