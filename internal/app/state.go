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
