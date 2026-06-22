// cmd/bootstrap.go
package main

import (
	"fmt"
	"sync"

	"github.com/subhasundardass/tuix/internal/config"
	"github.com/subhasundardass/tuix/internal/events"
)

// AppState holds all global application state
type AppState struct {
	mu          sync.RWMutex
	CurrentPage string
	DarkMode    bool
	User        *User
	Config      *config.Config
	EventBus    *events.EventBus
}

// User represents the current user
type User struct {
	Name  string
	Email string
	Role  string
}

// Global app state instance
var appState *AppState
var once sync.Once

// Init initializes the application with all dependencies
func InitApp() *AppState {
	once.Do(func() {
		// Load configuration
		cfg := config.Load()

		// Create event bus
		bus := events.NewEventBus()

		// Initialize state
		appState = &AppState{
			CurrentPage: "home",
			DarkMode:    false,
			User: &User{
				Name:  "Guest",
				Email: "guest@example.com",
				Role:  "viewer",
			},
			Config:   cfg,
			EventBus: bus,
		}

		// Register default event handlers
		registerDefaultHandlers(bus)

		fmt.Println("✅ Application initialized successfully")
	})

	return appState
}

// Get returns the global app state (panics if not initialized)
func GetApp() *AppState {
	if appState == nil {
		panic("App not initialized - call bootstrap.Init() first")
	}
	return appState
}

// registerDefaultHandlers sets up core event listeners
func registerDefaultHandlers(bus *events.EventBus) {
	// Handle page change events
	bus.Subscribe("page:change", func(data interface{}) {
		if page, ok := data.(string); ok {
			appState.mu.Lock()
			appState.CurrentPage = page
			appState.mu.Unlock()
			fmt.Printf("📍 Page changed to: %s\n", page)
		}
	})

	// Handle theme change events
	bus.Subscribe("theme:toggle", func(data interface{}) {
		appState.mu.Lock()
		appState.DarkMode = !appState.DarkMode
		mode := appState.DarkMode
		appState.mu.Unlock()
		fmt.Printf("🎨 Theme toggled: %v\n", mode)
	})

	// Handle user logout
	bus.Subscribe("user:logout", func(data interface{}) {
		appState.mu.Lock()
		appState.User = &User{
			Name:  "Guest",
			Email: "guest@example.com",
			Role:  "viewer",
		}
		appState.mu.Unlock()
		fmt.Println("👤 User logged out")
	})
}

// Convenience methods for accessing state
func (s *AppState) SetPage(page string) {
	s.mu.Lock()
	s.CurrentPage = page
	s.mu.Unlock()
	s.EventBus.Publish("page:change", page)
}

func (s *AppState) GetPage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CurrentPage
}

func (s *AppState) ToggleTheme() {
	s.EventBus.Publish("theme:toggle", nil)
}

func (s *AppState) IsDarkMode() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.DarkMode
}
