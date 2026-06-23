package app

import (
	"sync"

	"github.com/subhasundardass/tuix/internal/config"
)

type appState struct {
	mu          sync.RWMutex
	CurrentPage string
	DarkMode    bool
	User        *User
	Config      *config.Config
}

type User struct {
	Name  string
	Email string
	Role  string
}

var (
	state *appState
	once  sync.Once
)

// init() runs automatically when the package is imported
func init() {
	once.Do(func() {
		state = &appState{
			CurrentPage: "home",
			DarkMode:    false,
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
		}
	})
}

// getApp is UNEXPORTED - components CANNOT call this
func getAppState() *appState {
	if state == nil {
		panic("App not initialized")
	}
	return state
}

// These methods are exported - used by Context
func (s *appState) GetPage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CurrentPage
}

func (s *appState) IsDarkMode() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.DarkMode
}

func (s *appState) SetPage(page string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CurrentPage = page
}

func (s *appState) ToggleTheme() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DarkMode = !s.DarkMode
}

func (s *appState) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.User = &User{
		Name:  "Guest",
		Email: "guest@example.com",
		Role:  "viewer",
	}
}
