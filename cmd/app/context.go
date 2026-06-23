// bootstrap/context.go
package main

import (
	"github.com/subhasundardass/tuix/internal/config"
	"github.com/subhasundardass/tuix/tuix"
)

// AppContextData holds all data available via context
type AppContext struct {
	// Data
	Page   string
	Dark   bool
	User   *User
	Config *config.Config

	// Behavior (functions that can be called from anywhere)
	NavigateTo  func(string)
	ToggleTheme func()
	Logout      func()
}

// Create the context with a default value
var DefaultContext = tuix.CreateContext(AppContext{
	Page:        "home",
	Dark:        false,
	User:        &User{Name: "Guest", Email: "guest@example.com", Role: "viewer"},
	Config:      &config.Config{AppName: "My App", Version: "1.0.0"},
	NavigateTo:  func(string) {},
	ToggleTheme: func() {},
	Logout:      func() {},
})

// BuildContext creates a context value from the current app state
func BuildContext() AppContext {
	state := GetApp()
	return AppContext{
		Page:        state.CurrentPage,
		Dark:        state.DarkMode,
		User:        state.User,
		Config:      state.Config,
		NavigateTo:  state.SetPage,
		ToggleTheme: state.ToggleTheme,
		// Logout:      state.Logout,
	}
}
