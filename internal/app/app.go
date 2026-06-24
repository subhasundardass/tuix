package app

import (
	"github.com/subhasundardass/tuix/internal/context"
	appctx "github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/internal/ui"
	"github.com/subhasundardass/tuix/internal/ui/layout"
	"github.com/subhasundardass/tuix/tuix"
)

// Root is the main application component
// It's exported so main can use it
func Root(props tuix.Props) tuix.Element {

	// Local state for UI
	page, setPage := tuix.UseState(state.currentPage)

	ctx := &appctx.AppContext{}

	ctx.Set(appctx.AppContextValues{
		CurrentPage: page,
		AppName:     state.config.AppName,
		DarkMode:    state.darkMode,
		UserName:    state.user.name,
		NavigateTo: func(p string) {
			setPage(p)
			state.mu.Lock()
			state.currentPage = p
			state.mu.Unlock()
		},
		ToggleDark: func() {
			state.mu.Lock()
			state.darkMode = !state.darkMode
			state.mu.Unlock()
		},
		PushScreen:    PushScreen,
		ReplaceScreen: ReplaceScreen,
		PopScreen:     PopScreen,
		GetStack:      GetScreenStack,
		GetCurrent:    GetCurrentScreen,
	})

	// Handle Escape
	if tuix.CurrentKey.Code == tuix.KeyEscape {
		popped := ctx.PopScreen()
		setPage(popped)
	}

	// Sync with stack
	if current := ctx.GetCurrent(); current != page {
		setPage(current)
	}

	var content tuix.Element
	currentID := ctx.GetCurrent()
	screen, ok := ui.GetScreen(currentID)

	if !ok {
		// Fallback for 404
		content = tuix.Text("404 - Page Not Found", tuix.NewStyle().Foreground(tuix.Red))
	} else {
		// Render the screen with context
		content = screen.Render(ctx, tuix.Props{})
	}

	//Get screen title
	title := "Unknown"
	if screen, ok := ui.GetScreen(currentID); ok {
		title = screen.Title
	}

	// switch ctx.GetCurrent() {
	// case "home":
	// 	content = screen.HomePage(ctx, tuix.Props{})
	// case "settings":
	// 	content = screen.SettingsPage(tuix.Props{})
	// case "about":
	// 	content = screen.AboutPage(tuix.Props{})
	// default:
	// 	content = tuix.Text("404 - Page Not Found", tuix.NewStyle())
	// }

	// Provide context
	return context.DefaultContext.Provide(ctx, func() tuix.Element {
		return layout.MasterLayout(layout.LayoutProps{
			Title:   title,
			Content: content,
		})
	})
}
