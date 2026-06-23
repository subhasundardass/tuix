package app

import (
	"github.com/subhasundardass/tuix/internal/config"
	"github.com/subhasundardass/tuix/tuix"
)

// AppContext is EXPORTED - main can use this type
type AppContext struct {
	Page        string
	Dark        bool
	User        *User
	Config      *config.Config
	NavigateTo  func(string)
	ToggleTheme func()
	Logout      func()
}

// DefaultContext is the ONLY thing exposed to main
// main can use: app.DefaultContext
var DefaultContext = tuix.CreateContext(AppContext{
	Page: "home",
	Dark: false,
	User: &User{
		Name:  "Guest",
		Email: "guest@example.com",
		Role:  "viewer",
	},
	Config: &config.Config{
		AppName:     "My App",
		Version:     "1.0.0",
		Theme:       "dark",
		APIEndpoint: "http://localhost:8080",
		Debug:       false,
	},
	NavigateTo:  func(string) {},
	ToggleTheme: func() {},
	Logout:      func() {},
})

func buildContext() AppContext { // ← lowercase = unexported
	state := getAppState()
	return AppContext{
		Page:        state.CurrentPage,
		Dark:        state.DarkMode,
		User:        state.User,
		Config:      state.Config,
		NavigateTo:  state.SetPage,
		ToggleTheme: state.ToggleTheme,
		Logout:      state.Logout,
	}
}
