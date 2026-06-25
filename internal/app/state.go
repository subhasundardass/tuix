// internal/app/state.go
package app // ← same package name

import (
	"sync"

	"github.com/subhasundardass/tuix/internal/config"
)

type appState struct {
	mu          sync.RWMutex
	currentPage string
	darkMode    bool
	user        *user
	config      *config.Config
	screenStack []string
}

type user struct {
	name  string
	email string
	role  string
}

var (
	state *appState
	once  sync.Once
)

func init() {
	once.Do(func() {
		cfg := config.Load()
		state = &appState{
			currentPage: cfg.DefaultPage,
			darkMode:    false,
			user: &user{
				name:  "Guest",
				email: "guest@example.com",
				role:  "viewer",
			},
			config:      cfg,
			screenStack: []string{cfg.DefaultPage},
		}
	})
}

func SetPage(page string) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.currentPage = page
}

func ToggleTheme() {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.darkMode = !state.darkMode
}

func IsDarkMode() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.darkMode
}

func GetUserName() string {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.user.name
}

func GetConfig() *config.Config {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.config
}
