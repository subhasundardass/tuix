package app

import (
	appctx "github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/internal/ui/layout"
	"github.com/subhasundardass/tuix/internal/ui/screen"
	"github.com/subhasundardass/tuix/tuix"
)

// Root is the main application entry point rendered by tuix on every frame.
func Root(props tuix.Props) tuix.Element {

	// ── Context ───────────────────────────────────────────────────────────────
	ctx := &appctx.AppContext{}
	ctx.Set(appctx.AppContextValues{
		AppName:     state.config.AppName,
		DarkMode:    state.darkMode,
		UserName:    state.user.name,
		CurrentPage: state.currentPage,
		NavigateTo: func(p string) {
			state.mu.Lock()
			state.currentPage = p
			state.mu.Unlock()
		},
		ToggleDark: func() {
			state.mu.Lock()
			state.darkMode = !state.darkMode
			state.mu.Unlock()
		},
		PushScreen: PushScreen,
		PopScreen:  PopScreen,
		GetStack:   GetScreenStack,
		GetCurrent: GetCurrentScreen,
	})

	// ── Global keypresses ─────────────────────────────────────────────────────

	// F1 — open About window
	if tuix.CurrentKey.Code == tuix.KeyTab {
		// ⭐ Open window using the factory pattern
		tuix.ShowAt("about", "🏠 About",
			func(focused bool, availHeight int) tuix.Element {
				// ⭐ This function is called when rendering the window
				// It receives focus state and available height
				return screen.AboutPage(ctx, tuix.Props{
					Width:  tuix.Fixed(56),
					Height: tuix.Fixed(30), // ⭐ Use available height
				}, focused) // ⭐ Pass focused state to the page
			},
			30, 3, 60, 10,
		)
	}

	// ESC — close topmost window
	if tuix.HasAnyWindows() && tuix.CurrentKey.Code == tuix.KeyEscape {
		tuix.CloseCurrent()
	}

	// ── Layout ────────────────────────────────────────────────────────────────
	mainLayout := layout.MasterLayout(layout.LayoutProps{
		Ctx:     ctx,
		Title:   "Hello",
		Content: tuix.Box(tuix.Props{}, tuix.NewStyle()),
	})

	// return tuix.Box(
	// 	tuix.Props{Direction: tuix.Column, Width: tuix.Grow(1)},
	// 	tuix.NewStyle(),
	// 	mainLayout,
	// 	tuix.RenderWindows(), // always rendered last so windows appear on top
	// )
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Grow(1),
			Height:    tuix.Grow(1),
		},
		tuix.NewStyle(),
		mainLayout,
		tuix.RenderWindows(),
	)
}
