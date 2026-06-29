package main

import (
	"fmt"

	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

type Theme struct {
	Name   string
	Fg     tuix.Color
	Bg     tuix.Color
	Accent tuix.Color
	Muted  tuix.Color
}

var lightTheme = Theme{
	Name:   "light",
	Fg:     tuix.BrightBlack,
	Bg:     tuix.BrightWhite,
	Accent: tuix.Blue,
	Muted:  tuix.Black,
}

var darkTheme = Theme{
	Name:   "dark",
	Fg:     tuix.BrightWhite,
	Bg:     tuix.Black,
	Accent: tuix.BrightCyan,
	Muted:  tuix.BrightBlack,
}

var ThemeContext = tuix.CreateContext(lightTheme)

// ============================================================================
// TITLE BANNER - First impression matters
// ============================================================================
func TitleBanner() tuix.Element {
	t := tuix.UseContext(ThemeContext)
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Padding:   [4]int{1, 0, 1, 0},
			Gap:       0,
		},
		tuix.NewStyle().Foreground(t.Accent).Bold(true),
		tuix.Text("╔═══════════════════════════════════════════════════════╗", tuix.Style{}),
		tuix.Text("║          TUIX - Terminal UI Framework                 ║", tuix.NewStyle().Bold(true).Foreground(t.Accent)),
		tuix.Text("║     React-style components for the terminal           ║", tuix.NewStyle()),
		tuix.Text("╚═══════════════════════════════════════════════════════╝", tuix.Style{}),
	)
}

// ============================================================================
// FEATURE CARDS - Showcase key capabilities
// ============================================================================
func FeatureCard(title, description string) tuix.Element {
	t := tuix.UseContext(ThemeContext)
	return components.Panel(
		title,
		26,
		tuix.WrappedText(description, tuix.NewStyle().Foreground(t.Fg)),
	)
}

func FeaturesSection() tuix.Element {
	t := tuix.UseContext(ThemeContext)
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Fixed(80),
			Gap:       1,
		},
		tuix.NewStyle(),
		tuix.Text("Core Features", tuix.NewStyle().Foreground(t.Accent).Bold(true).Underline(true)),
		tuix.Box(
			tuix.Props{
				Direction: tuix.Row,
				Width:     tuix.Fixed(60),
				Gap:       1,
			},
			tuix.NewStyle(),
			FeatureCard("Hooks", "UseState, UseEffect, UseContext - familiar React patterns"),
			FeatureCard("Components", "Compose UIs declaratively with functional components"),
			FeatureCard("Fast Render", "Minimal diffing - only changed cells re-render"),
		),
	)
}

// ============================================================================
// INTERACTIVE COUNTER - State management showcase
// ============================================================================
func Counter() tuix.Element {
	t := tuix.UseContext(ThemeContext)
	count, setCount := tuix.UseState(0)

	// Increment on '+', decrement on '-'
	if tuix.CurrentKey.Rune == '+' || tuix.CurrentKey.Rune == '=' {
		if count < 100 {
			setCount(count + 1)
		}
	}
	if tuix.CurrentKey.Rune == '-' || tuix.CurrentKey.Rune == '_' {
		if count > 0 {
			setCount(count - 1)
		}
	}

	// Calculate progress 0.0 to 1.0
	progress := float64(count) / 100.0

	return components.Panel(
		"Interactive Counter (0-100)",
		60,
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			tuix.Text("Count:", tuix.NewStyle().Foreground(t.Muted)),
			tuix.Text(fmt.Sprintf("%d / 100", count), tuix.NewStyle().Foreground(t.Accent).Bold(true)),
		),
		components.ProgressBar(progress, 50, t.Accent),
		tuix.Text("Press + / - to change", tuix.NewStyle().Foreground(t.Muted)),
	)
}

// ============================================================================
// TIMER - Auto-updating timer with Spinner
// ============================================================================
func Timer() tuix.Element {
	t := tuix.UseContext(ThemeContext)
	elapsed, setElapsed := tuix.UseState(0)

	// Auto-increment every render (simulates timer tick)
	// In a real app, you'd use UseEffect with a proper time-based update
	// For demo purposes, we increment on each render
	if elapsed < 100 {
		// Use a simple mechanism: increment based on some condition
		// Note: This will need proper timing logic in production
		setElapsed(elapsed + 1)
	}

	// Calculate progress
	progress := float64(elapsed) / 100.0

	return components.Panel(
		"Auto-Updating Timer",
		60,
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			components.Spinner("Running..."),
			tuix.Text(fmt.Sprintf("%d%%", elapsed), tuix.NewStyle().Foreground(t.Accent).Bold(true)),
		),
		components.ProgressBar(progress, 50, t.Accent),
		tuix.Text("Timer updates automatically", tuix.NewStyle().Foreground(t.Muted)),
	)
}

// ============================================================================
// KEYBINDS FOOTER
// ============================================================================
func Footer() tuix.Element {
	t := tuix.UseContext(ThemeContext)
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Width:     tuix.Grow(1),
			Gap:       2,
		},
		tuix.NewStyle(),
		tuix.Text("T", tuix.NewStyle().Foreground(t.Accent).Bold(true)),
		tuix.Text("toggle theme", tuix.NewStyle().Foreground(t.Fg)),
		tuix.Text("•", tuix.NewStyle().Foreground(t.Muted)),
		tuix.Text("+/-", tuix.NewStyle().Foreground(t.Accent).Bold(true)),
		tuix.Text("counter", tuix.NewStyle().Foreground(t.Fg)),
		tuix.Text("•", tuix.NewStyle().Foreground(t.Muted)),
		tuix.Text("Q", tuix.NewStyle().Foreground(t.Accent).Bold(true)),
		tuix.Text("quit", tuix.NewStyle().Foreground(t.Fg)),
	)
}

// ============================================================================
// MAIN APP - Centered layout
// ============================================================================
func App(props tuix.Props) tuix.Element {
	dark, setDark := tuix.UseState(false)

	if tuix.CurrentKey.Rune == 't' || tuix.CurrentKey.Rune == 'T' {
		setDark(!dark)
	}
	if tuix.CurrentKey.Rune == 'q' || tuix.CurrentKey.Rune == 'Q' {
		tuix.Exit()
	}

	current := lightTheme
	if dark {
		current = darkTheme
	}

	return ThemeContext.Provide(current, func() tuix.Element {
		return tuix.Box(
			tuix.Props{
				Direction: tuix.Column,
				Gap:       1,
				Width:     tuix.Grow(1),
			},
			tuix.NewStyle(),
			// Spacer top
			tuix.Box(tuix.Props{Height: tuix.Fixed(1)}, tuix.NewStyle()),
			// Centered content wrapper
			tuix.Box(
				tuix.Props{
					Direction: tuix.Row,
					Width:     tuix.Grow(1),
				},
				tuix.NewStyle(),
				tuix.Box(
					tuix.Props{
						Direction: tuix.Column,
						Gap:       1,
						Width:     tuix.Fixed(62),
					},
					tuix.NewStyle(),
					TitleBanner(),
					FeaturesSection(),
					Counter(),
					Timer(),
				),
			),
			// Spacer grow
			tuix.Box(tuix.Props{Height: tuix.Grow(1)}, tuix.NewStyle()),
			// Footer
			Footer(),
		)
	})
}

func main() {
	app := tuix.NewApp(140, 40)
	app.Run(App, tuix.Props{})
}
