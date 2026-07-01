package window

import (
	"github.com/subhasundardass/tuix/tuix"
)

// OverlayRenderer handles rendering windows as overlays on top of main content.
// It uses tuix.Overlay() for absolute positioning of windows.
type OverlayRenderer struct {
	screenWidth  int
	screenHeight int
}

// NewOverlayRenderer creates a new overlay renderer with the given screen dimensions.
func NewOverlayRenderer(screenWidth, screenHeight int) *OverlayRenderer {
	return &OverlayRenderer{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// Render returns a full-screen element with windows properly overlaid on main content.
// The render order is: Main Content → Windows (in Z-order)
// Windows are rendered on top of main content using absolute positioning.
func (or *OverlayRenderer) Render(mainContent tuix.Element) tuix.Element {
	mgr := GetManager()

	// If no windows, just return the main content
	if mgr.Count() == 0 {
		return mainContent
	}

	// Start with main content as the base layer
	elements := []tuix.Element{mainContent}

	// Add all visible windows as overlays in Z-order (bottom to top)
	for _, winID := range mgr.GetZOrder() {
		win := mgr.GetWindow(winID)
		if win == nil || !win.visible {
			continue
		}
		elements = append(elements, or.renderWindowAsOverlay(win))
	}

	// Stack everything in a full-screen container
	// Children are rendered in order: earlier elements are behind later ones
	return tuix.Box(
		tuix.Props{
			Width:  tuix.Grow(1),
			Height: tuix.Grow(1),
		},
		tuix.NewStyle(),
		elements...,
	)
}

// renderWindowAsOverlay renders a single window as an absolute-positioned overlay.
// Uses tuix.Overlay() for true absolute positioning at (w.X, w.Y).
func (or *OverlayRenderer) renderWindowAsOverlay(w *Window) tuix.Element {
	// Build the complete window UI
	windowContent := or.buildWindowContent(w)

	// Place the window at its exact screen coordinates using Overlay
	return tuix.Overlay(w.X, w.Y, windowContent)
}

// buildTitleBarWithColor creates title bar with custom color
func (or *OverlayRenderer) buildWindowContent(w *Window) tuix.Element {

	// width := w.Width
	// height := w.Height

	// Simple window with no border - just title bar and content
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Fit(),
			Height:    tuix.Fit(),
		},
		tuix.NewStyle(), // NO Background here - remove it
		// Title bar at the top - this has its own background
		or.buildTitleBar(w),
		// Content with padding - NO background
		tuix.Box(
			tuix.Props{
				Padding: [4]int{1, 1, 1, 1},
			},
			tuix.NewStyle().Background(tuix.Color{Type: tuix.ColorRGB,
				R: 40,
				G: 40,
				B: 40,
			}),
			w.Content,
		),
	)
}

// buildTitleBar creates the window title bar with title text and close indicator.
func (or *OverlayRenderer) buildTitleBar(w *Window) tuix.Element {
	// Set default title if empty
	title := w.Title
	if title == "" {
		title = "Window"
	}

	// Add [MODAL] indicator for modal windows
	if w.Modal {
		title = title + " [MODAL]"
	}

	// Title bar with blue background and white text
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Width:     tuix.Fixed(w.Width),
			Height:    tuix.Fixed(1),
		},
		tuix.NewStyle().
			Background(tuix.Blue).
			Foreground(tuix.White),
		// Title with a space prefix for padding
		tuix.Text(" "+title, tuix.NewStyle().Bold(true)),
		// Flexible spacer to push close button to the right
		tuix.Box(
			tuix.Props{Width: tuix.Grow(1)},
			tuix.NewStyle(),
		),
		// Close button indicator (visual only for now)
		tuix.Text("[Esc]", tuix.NewStyle().Foreground(tuix.BrightRed)),
	)
}

// RenderWindowsOverlay is a convenience function that renders windows as overlays.
// This is the recommended way to integrate window management into your app.
//
// Usage:
//
//	func App(props tuix.Props) tuix.Element {
//	    mainContent := buildYourMainUI()
//	    return window.RenderWindowsOverlay(140, 40, mainContent)
//	}
func RenderWindowsOverlay(screenWidth, screenHeight int, mainContent tuix.Element) tuix.Element {
	renderer := NewOverlayRenderer(screenWidth, screenHeight)
	return renderer.Render(mainContent)
}
