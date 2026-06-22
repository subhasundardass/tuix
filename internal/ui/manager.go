// // internal/manager/page_manager.go
package ui

// import (
// 	"github.com/subhasundardass/tuix/tuix"
// )

// type PageType string

// const (
// 	PageHome     PageType = "home"
// 	PageSettings PageType = "settings"
// 	PageAbout    PageType = "about"
// )

// // ─── Page State ──────────────────────────────────────────────────────────

// type HomePageState struct {
// 	selectedItem int
// 	items        []string
// }

// type SettingsPageState struct {
// 	selectedOption int
// 	theme          string
// }

// type AboutPageState struct {
// 	scrollPos int
// }

// // ─── Page Manager ───────────────────────────────────────────────────────

// type PageManager struct {
// 	currentPage  PageType
// 	previousPage PageType

// 	// Page-specific states
// 	homeState     HomePageState
// 	settingsState SettingsPageState
// 	aboutState    AboutPageState

// 	// Render functions for each page
// 	renderFuncs map[PageType]func() tuix.Element
// }

// func NewPageManager() *PageManager {
// 	pm := &PageManager{
// 		currentPage: PageHome,
// 		renderFuncs: make(map[PageType]func() tuix.Element),

// 		homeState: HomePageState{
// 			items: []string{
// 				"📄 Open Document",
// 				"📝 Create New",
// 				"⚙️  Settings",
// 				"ℹ️  About",
// 				"🚪 Quit",
// 			},
// 		},

// 		settingsState: SettingsPageState{
// 			theme: "dark",
// 		},
// 	}

// 	// Register page renderers
// 	pm.renderFuncs[PageHome] = pm.renderHome
// 	pm.renderFuncs[PageSettings] = pm.renderSettings
// 	pm.renderFuncs[PageAbout] = pm.renderAbout

// 	return pm
// }

// // ─── Navigation ─────────────────────────────────────────────────────────

// func (pm *PageManager) Navigate(page PageType) {
// 	pm.previousPage = pm.currentPage
// 	pm.currentPage = page
// }

// func (pm *PageManager) Back() {
// 	if pm.previousPage != "" {
// 		temp := pm.currentPage
// 		pm.currentPage = pm.previousPage
// 		pm.previousPage = temp
// 	}
// }

// // ─── Render Current Page ────────────────────────────────────────────────

// func (pm *PageManager) RenderCurrentPage() tuix.Element {
// 	renderer, exists := pm.renderFuncs[pm.currentPage]
// 	if !exists {
// 		return renderErrorPage("Page not found")
// 	}
// 	return renderer()
// }

// func (pm *PageManager) GetPageTitle() string {
// 	switch pm.currentPage {
// 	case PageHome:
// 		return "Home"
// 	case PageSettings:
// 		return "Settings"
// 	case PageAbout:
// 		return "About"
// 	default:
// 		return "Unknown"
// 	}
// }

// // ─── Page Renderers ─────────────────────────────────────────────────────

// func (pm *PageManager) renderHome() tuix.Element {
// 	state := &pm.homeState

// 	items := []tuix.Element{}

// 	for i, item := range state.items {
// 		style := tuix.NewStyle()

// 		// Highlight selected item
// 		if i == state.selectedItem {
// 			style = style.
// 				Background(tuix.Blue).
// 				Foreground(tuix.White).
// 				Bold(true)
// 		} else {
// 			style = style.Foreground(tuix.White)
// 		}

// 		items = append(items, tuix.Text("  "+item, style))
// 	}

// 	content := tuix.Box(
// 		tuix.Props{
// 			Direction: tuix.Column,
// 			Gap:       1,
// 		},
// 		tuix.NewStyle(),
// 		append([]tuix.Element{
// 			tuix.Text("Welcome to Home Page", tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)),
// 			tuix.Text("", tuix.NewStyle()),
// 		}, items...)...,
// 	)

// 	return content
// }

// func (pm *PageManager) renderSettings() tuix.Element {
// 	state := &pm.settingsState

// 	options := []tuix.Element{
// 		tuix.Text("  ☀️  Light Theme", tuix.NewStyle().Foreground(tuix.White)),
// 		tuix.Text("  🌙 Dark Theme", tuix.NewStyle().Foreground(tuix.White)),
// 		tuix.Text("  🔔 Notifications", tuix.NewStyle().Foreground(tuix.White)),
// 	}

// 	if state.selectedOption < len(options) {
// 		oldStyle := options[state.selectedOption]
// 		// Highlight selected
// 		options[state.selectedOption] = tuix.Text(
// 			"  > "+state.getOptionLabel(),
// 			tuix.NewStyle().Background(tuix.Blue).Foreground(tuix.White).Bold(true),
// 		)
// 	}

// 	content := tuix.Box(
// 		tuix.Props{
// 			Direction: tuix.Column,
// 			Gap:       1,
// 		},
// 		tuix.NewStyle(),
// 		append([]tuix.Element{
// 			tuix.Text("Settings", tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)),
// 			tuix.Text("", tuix.NewStyle()),
// 		}, options...)...,
// 	)

// 	return content
// }

// func (pm *PageManager) renderAbout() tuix.Element {
// 	content := tuix.Box(
// 		tuix.Props{
// 			Direction: tuix.Column,
// 			Gap:       1,
// 		},
// 		tuix.NewStyle(),
// 		tuix.Text("About This Application", tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)),
// 		tuix.Text("", tuix.NewStyle()),
// 		tuix.Text("Version: 1.0.0", tuix.NewStyle().Foreground(tuix.White)),
// 		tuix.Text("Built with: Tuix Framework", tuix.NewStyle().Foreground(tuix.White)),
// 		tuix.Text("Author: Your Name", tuix.NewStyle().Foreground(tuix.White)),
// 		tuix.Text("", tuix.NewStyle()),
// 		tuix.Text("Press [ESC] to go back", tuix.NewStyle().Dim(true)),
// 	)

// 	return content
// }

// // ─── Input Handling ─────────────────────────────────────────────────────

// func (pm *PageManager) HandleInput(key tuix.KeyEvent) {
// 	// Global shortcuts
// 	switch key.Code {
// 	case tuix.KeyChar:
// 		switch key.Rune {
// 		case '1':
// 			pm.Navigate(PageHome)
// 			return
// 		case '2':
// 			pm.Navigate(PageSettings)
// 			return
// 		case '3':
// 			pm.Navigate(PageAbout)
// 			return
// 		}
// 	case tuix.KeyEsc:
// 		pm.Back()
// 		return
// 	}

// 	// Page-specific input
// 	switch pm.currentPage {
// 	case PageHome:
// 		pm.handleHomeInput(key)
// 	case PageSettings:
// 		pm.handleSettingsInput(key)
// 	case PageAbout:
// 		pm.handleAboutInput(key)
// 	}
// }

// func (pm *PageManager) handleHomeInput(key tuix.KeyEvent) {
// 	state := &pm.homeState

// 	switch key.Code {
// 	case tuix.KeyDown:
// 		if state.selectedItem < len(state.items)-1 {
// 			state.selectedItem++
// 		}
// 	case tuix.KeyUp:
// 		if state.selectedItem > 0 {
// 			state.selectedItem--
// 		}
// 	case tuix.KeyEnter:
// 		switch state.selectedItem {
// 		case 0: // Open Document
// 			// Handle open
// 		case 1: // Create New
// 			// Handle create
// 		case 2: // Settings
// 			pm.Navigate(PageSettings)
// 		case 3: // About
// 			pm.Navigate(PageAbout)
// 		case 4: // Quit
// 			// Return special value or signal
// 		}
// 	}
// }

// func (pm *PageManager) handleSettingsInput(key tuix.KeyEvent) {
// 	state := &pm.settingsState

// 	switch key.Code {
// 	case tuix.KeyDown:
// 		if state.selectedOption < 2 {
// 			state.selectedOption++
// 		}
// 	case tuix.KeyUp:
// 		if state.selectedOption > 0 {
// 			state.selectedOption--
// 		}
// 	case tuix.KeyEnter:
// 		// Apply setting
// 		state.applyOption()
// 	}
// }

// func (pm *PageManager) handleAboutInput(key tuix.KeyEvent) {
// 	// About page only responds to navigation keys
// 	switch key.Code {
// 	case tuix.KeyUp, tuix.KeyDown:
// 		// Scroll if needed
// 	}
// }

// // ─── Helper Methods ─────────────────────────────────────────────────────

// func (s *SettingsPageState) getOptionLabel() string {
// 	switch s.selectedOption {
// 	case 0:
// 		return "☀️  Light Theme"
// 	case 1:
// 		return "🌙 Dark Theme"
// 	case 2:
// 		return "🔔 Notifications"
// 	default:
// 		return "Unknown"
// 	}
// }

// func (s *SettingsPageState) applyOption() {
// 	// Handle option application
// }

// func renderErrorPage(message string) tuix.Element {
// 	return tuix.Text("❌ "+message, tuix.NewStyle().Foreground(tuix.Red))
// }
