package main

import (
	"fmt"

	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/window"
)

// ============================================================================
// EXAMPLE: Complete Window Manager Demo
// ============================================================================
// This example demonstrates:
// 1. Creating windows from button clicks
// 2. Modal vs non-modal windows
// 3. Window positioning and sizing
// 4. Closing and managing windows
// 5. Rendering windows in the correct layer

func main() {
	app := tuix.NewApp(180, 40)
	app.Run(App, tuix.Props{})
}

func App(props tuix.Props) tuix.Element {

	// Handle keyboard shortcuts
	handleKeyboard()

	const screenW, screenH = 140, 40

	// Build main content normally
	mainContent := tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 2, 1, 2},
			Width:     tuix.Grow(1),
			Height:    tuix.Grow(1),
		},
		tuix.NewStyle(), // Ensure main content has background
		Header(),
		Menu(),
		StatusLine(),
	)

	// Use the overlay renderer to render windows on top
	overlayRenderer := window.NewOverlayRenderer(screenW, screenH)
	return overlayRenderer.Render(mainContent)
}

// ============================================================================
// COMPONENTS
// ============================================================================

func Header() tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 0},
		tuix.NewStyle(),
		tuix.Text(
			"╔═══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗",
			tuix.NewStyle().Foreground(tuix.BrightCyan),
		),
		tuix.Text(
			"║                          TUIX WINDOW MANAGER -  Demo                                                          ║",
			tuix.NewStyle().Foreground(tuix.BrightCyan).Bold(true),
		),
		tuix.Text(
			"╚═══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝",
			tuix.NewStyle().Foreground(tuix.BrightCyan),
		),
	)
}

func Menu() tuix.Element {
	openWins := window.Count()

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),

		tuix.Text("📋 MENU", tuix.NewStyle().Bold(true).Foreground(tuix.BrightYellow)),
		tuix.Text("", tuix.NewStyle()),

		tuix.Text("1 - Open simple window (centered)", tuix.NewStyle()),
		tuix.Text("2 - Open modal dialog", tuix.NewStyle()),
		tuix.Text("3 - Open notification window", tuix.NewStyle()),
		tuix.Text("4 - Open form window", tuix.NewStyle()),
		tuix.Text("c - Close all windows", tuix.NewStyle()),
		tuix.Text("ESC - Close focused window", tuix.NewStyle()),
		tuix.Text("q - Quit", tuix.NewStyle()),

		tuix.Text("", tuix.NewStyle()),
		tuix.Text(
			fmt.Sprintf("📊 Windows open: %d", openWins),
			tuix.NewStyle().Foreground(tuix.BrightGreen),
		),
	)
}

func StatusLine() tuix.Element {
	focused := window.GetFocused()
	isModal := window.IsAnyModalOpen()

	statusText := "No windows open"
	if focused != "" {
		statusText = fmt.Sprintf("Focused: %s", focused)
	}
	if isModal {
		statusText += " (Modal active)"
	}

	return tuix.Box(
		tuix.Props{Direction: tuix.Row, Width: tuix.Grow(1)},
		tuix.NewStyle().Foreground(tuix.BrightBlack),
		tuix.Text("─────────────────────────────────────────────────────────────────────", tuix.NewStyle()),
		tuix.Text("Status: "+statusText, tuix.NewStyle().Foreground(tuix.BrightMagenta)),
	)
}

// ============================================================================
// KEYBOARD HANDLING
// ============================================================================
// closeFocusedWindow closes the currently focused window
func closeFocusedWindow() {
	focusedID := window.GetFocused()
	if focusedID != "" {
		win := window.GetManager().GetWindow(focusedID)
		if win != nil {
			win.Close()
		}
	}
}

func handleKeyboard() {
	// Check for ESC key - close focused window
	if tuix.CurrentKey.Code == tuix.KeyEscape {
		closeFocusedWindow()
		return
	}

	switch tuix.CurrentKey.Rune {
	case '1':
		openSimpleWindow()
	case '2':
		openModalDialog()
	case '3':
		openNotificationWindow()
	case '4':
		openFormWindow()
	case 'c':
		window.CloseAll()
	case 'q':
		tuix.Exit()
	}

}

// ============================================================================
// WINDOW CREATORS
// ============================================================================

func openSimpleWindow() {
	content := tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.Text("This is a simple window!", tuix.NewStyle().Foreground(tuix.BrightGreen).Bold(true)),
		tuix.Text("It renders on top of the main content.", tuix.NewStyle()),
		tuix.Text("Windows can contain any Tuix elements.", tuix.NewStyle()),
		tuix.Text("Press '1' again to open another window", tuix.NewStyle().Foreground(tuix.BrightBlack)),
	)

	window.Create(content).
		SetTitle("Simple Window").
		SetSize(50, 1).
		SetModal(false).
		CenterOnScreen(180, 40).
		Show()
}

func openModalDialog() {
	content := tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.MultilineText("⚠ Confirmation Required", tuix.NewStyle().Foreground(tuix.BrightYellow).Bold(true)),
		tuix.MultilineText("Modal windows block input to the background.", tuix.NewStyle()),
		tuix.MultilineText("This window must be closed before interacting", tuix.NewStyle()),
		tuix.MultilineText("with the main application.", tuix.NewStyle()),
		tuix.MultilineText("Press '2' again to open another modal", tuix.NewStyle().Foreground(tuix.BrightBlack)),
	)

	window.Create(content).
		SetTitle("Modal Dialog").
		SetSize(55, 12). // Increased size
		SetModal(false).
		CenterOnScreen(140, 40).
		Show()
}

func openNotificationWindow() {
	content := tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.Text("✓ Success!", tuix.NewStyle().Foreground(tuix.BrightGreen).Bold(true)),
		tuix.Text("This is a notification window.", tuix.NewStyle()),
		tuix.Text("It doesn't block background interaction.", tuix.NewStyle()),
		tuix.Text("Stack multiple notifications!", tuix.NewStyle()),
	)

	window.Create(content).
		SetTitle("Notification").
		SetSize(45, 9).
		SetModal(false).
		SetPosition(85, 5).
		Show()
}

func openFormWindow() {
	content := tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.Text("User Input Form", tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)),
		tuix.Text("Name:  [________________]", tuix.NewStyle()),
		tuix.Text("Email: [________________]", tuix.NewStyle()),
		tuix.Text("", tuix.NewStyle()),
		tuix.Text("[ Submit ]  [ Cancel ]", tuix.NewStyle().Foreground(tuix.BrightYellow)),
	)

	window.Create(content).
		SetTitle("Registration Form").
		SetSize(48, 12).
		SetModal(true).
		CenterOnScreen(140, 40).
		Show()
}

// ============================================================================
// NOTES
// ============================================================================
//
// 1. window.Create() returns a *Window that you can configure
// 2. Windows render automatically via window.RenderWindows()
// 3. Windows layer in Z-order (last opened is on top)
// 4. Modal windows visually indicate blocked background interaction
// 5. Call window.CloseAll() to close all windows at once
//
// Phase 1 Features:
// ✓ Window creation and lifecycle (create, show, hide, close)
// ✓ Sizing and positioning (manual + CenterOnScreen)
// ✓ Z-order stacking and focus
// ✓ Modal window support
// ✓ Thread-safe global registry
//
// Coming in Phase 2:
// - FocusManager integration
// - Tab/Escape keyboard routing
// - Close button (✕)
// - Window resize
// - Draggable title bar
