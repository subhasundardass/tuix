package app

import (
	appctx "github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/internal/ui"
	"github.com/subhasundardass/tuix/internal/ui/layout"
	"github.com/subhasundardass/tuix/tuix"
)

// Root is the main application entry point rendered by tuix on every frame.
func Root(props tuix.Props) tuix.Element {

	page, setPage := tuix.UseState(state.currentPage)
	// ── Context ───────────────────────────────────────────────────────────────
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

	// ── Global keypresses ─────────────────────────────────────────────────────

	// // F1 — open About window
	// if tuix.CurrentKey.Code == tuix.KeyTab {
	// 	// Open window using the factory pattern
	// 	tuix.ShowAt("about", "🏠 About",
	// 		func(focused bool, availHeight int) tuix.Element {
	// 			// This function is called when rendering the window
	// 			// It receives focus state and available height
	// 			return screen.AboutPage(ctx, tuix.Props{
	// 				Width:  tuix.Fixed(56),
	// 				Height: tuix.Fixed(30), // Use available height
	// 			}, focused) // Pass focused state to the page
	// 		},
	// 		30, 3, 60, 10,
	// 	)
	// }

	// // ESC — close topmost window
	// if tuix.HasAnyWindows() && tuix.CurrentKey.Code == tuix.KeyEscape {
	// 	tuix.CloseCurrent()
	// }

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

	// Get screen title
	title := "Unknown"
	if screen, ok := ui.GetScreen(page); ok {
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
	// ── Focus ────────────────────────────────────────────────────────────────
	tuix.SetFocusOrder([]string{"sidebar", "content"})
	if ctx.GetCurrent() == "" {
		tuix.Focus("sidebar")
	}
	// Handle Tab key
	if tuix.CurrentKey.Code == tuix.KeyTab {
		tuix.FocusNext()
	}

	// ── Layout ────────────────────────────────────────────────────────────────
	mainLayout := layout.MasterLayout(ctx, layout.LayoutProps{
		Ctx:   ctx,
		Title: title,
		// Content: tuix.Box(tuix.Props{}, tuix.NewStyle()),
		Content: content,
	})

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Grow(1),
			Height:    tuix.Grow(1),
		},
		tuix.NewStyle(),
		mainLayout,
		// tuix.RenderWindows(),
	)
}
