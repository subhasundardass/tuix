// Example: focus-demo
// Demonstrates:
// - Focusable inputs
// - Non-focusable (disabled) input
// - Focus movement testing
// Run: go run ./examples/focus-demo

package main

import (
	"fmt"

	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

// local state keys
const (
	nameID   = "name"
	emailID  = "email"
	ageID    = "age"
	submitID = "submit"
)

func App(props tuix.Props) tuix.Element {

	//state
	name, setName := tuix.UseState("Subha")
	email, setEmail := tuix.UseState("subha@email.com")
	age, setAge := tuix.UseState("25")

	tuix.UseEffect(func() func() {
		// set initial focus
		tuix.Focus(nameID)
		return nil
	}, []any{})

	title := tuix.NewStyle().
		Bold(true).
		Foreground(tuix.BrightBlack)

	dim := tuix.NewStyle().
		Foreground(tuix.BrightBlack)

	label := tuix.NewStyle().
		Foreground(tuix.BrightWhite)

	switch {
	case tuix.IsFocused(nameID) && tuix.CurrentKey.Code == tuix.KeyEnter:
		tuix.Debug("Enter Pressed in Name")
		tuix.Debug("Focusing: " + ageID) // ← check what ageID actually is
		tuix.Focus(ageID)

	case tuix.IsFocused(ageID) && tuix.CurrentKey.Code == tuix.KeyEnter:
		tuix.Debug("Enter Pressed in Age")
		tuix.Focus(submitID)

	}

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 2, 1, 2},
		},
		tuix.NewStyle(),

		// Header
		tuix.Text("◆ Focus Demo (Tab / Shift+Tab)", title),
		tuix.Text("(email is NOT focusable)", dim),

		// NAME (focusable)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			label,
			tuix.Text("Name:", tuix.NewStyle()),
			components.Input(
				tuix.IsFocused(nameID),
				name,
				func(v string) {
					tuix.Debug("Value Changed")
					setName(v)
				},
			),
			func() tuix.Element {
				// if tuix.IsFocused("name") {
				// 	return tuix.Text("Focussed █", tuix.NewStyle())
				// }

				return tuix.Text("", tuix.NewStyle())
			}(),
		),

		// EMAIL (NOT focusable)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			label,
			tuix.Text("Email:", tuix.NewStyle()),
			func() tuix.Element {
				tuix.Debug("email focused: " + fmt.Sprintf("%v", tuix.IsFocused(emailID)))
				return components.Input(
					false, //tuix.IsFocused(emailID),
					email,
					func(v string) {
						setEmail(v)
					},
				)
			}(),
		),

		// AGE (focusable)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			label,
			tuix.Text("Age:", label),

			func() tuix.Element {
				return components.Input(
					tuix.IsFocused(ageID),
					age,
					func(value string) { setAge(value) },
				)
			}(),
		),

		// BUTTON (focusable)
		func() tuix.Element {
			if tuix.IsFocused(submitID) {
				return tuix.Text("[ Submited ] ", tuix.NewStyle().Foreground(tuix.BrightWhite).Bold(true))
			}
			return tuix.Text("[ Submit ]", tuix.NewStyle())
		}(),

		tuix.Text("", dim),
		tuix.Text("Tab: Next | Shift+Tab: Prev", dim),
	)
}

func main() {
	app := tuix.NewApp(60, 10)
	app.Run(App, tuix.Props{})
}
