package ui

import (
	"github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/tuix"
)

// ─── Screen Definition ────────────────────────────────────────────────────

type Screen struct {
	ID     string
	Title  string
	Render func(ctx *context.AppContext, props tuix.Props) tuix.Element
}

// ─── Helper Functions ─────────────────────────────────────────────────────

// GetScreen returns a screen by ID
func GetScreen(id string) (Screen, bool) {
	screen, ok := Registry[id]
	return screen, ok
}

// Registry holds all available screens
var Registry = map[string]Screen{
	// "home": {
	// 	ID:     "home",
	// 	Title:  "Home",
	// 	Render: screen.HomePage,
	// },
	// "settings": {
	// 	ID:     "settings",
	// 	Title:  "Settings",
	// 	Render: screen.SettingsPage,
	// },
	// "about": {
	// 	ID:     "about",
	// 	Title:  "About",
	// 	Render: screen.AboutPage,
	// },

	// Add new screens here...
}
