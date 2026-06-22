package main

import (
	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	showModal, setShowModal := tuix.UseState(false)
	confirmed, setConfirmed := tuix.UseState(false)

	// Open modal on Enter, but only when modal is not already open
	if !showModal && tuix.CurrentKey.Code == tuix.KeyEnter {
		setShowModal(true)
	}

	// Confirm action on 'y' when modal is open
	if showModal && tuix.CurrentKey.Rune == 'y' {
		setConfirmed(true)
		setShowModal(false)
	}

	titleStyle := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	normalStyle := tuix.NewStyle().Foreground(tuix.White)
	dimStyle := tuix.NewStyle().Foreground(tuix.BrightBlack)
	greenStyle := tuix.NewStyle().Foreground(tuix.BrightGreen).Bold(true)
	yellowStyle := tuix.NewStyle().Foreground(tuix.BrightYellow)

	statusText := tuix.If(confirmed,
		tuix.Text("✓ Action confirmed!", greenStyle),
		tuix.Text("No action taken yet.", dimStyle),
	)

	// Main page content — rendered normally in flow layout.
	// The overlay floats above this regardless of where it sits in the tree.
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}},
		tuix.NewStyle(),

		tuix.Text("◆ Overlay Modal Demo", titleStyle),
		tuix.Text("Press Enter to open the confirmation modal.", normalStyle),
		statusText,
		tuix.Text("ctrl-c to quit", dimStyle),

		// ModalOverlay floats at absolute position (6, 4) on screen.
		// It takes zero space in the column layout above — siblings are
		// unaffected and the modal paints on top of everything.
		components.ModalOverlay(
			"Confirm Action",
			showModal,
			6, 4, // x=6 col, y=4 row — centered-ish for a 60x10 terminal
			46, // width
			func() { setShowModal(false) },
			tuix.Text("Do you want to proceed?", normalStyle),
			tuix.Text("Press y to confirm · Esc to cancel", yellowStyle),
		),
	)
}

func main() {
	app := tuix.NewApp(60, 10)
	app.Run(App, tuix.Props{})
}
