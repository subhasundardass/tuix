package main

import (
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
	// state
	name, setName := tuix.UseState("Subha")
	email, setEmail := tuix.UseState("subha@email.com")
	age, setAge := tuix.UseState("25")

	tuix.UseEffect(func() func() {
		// set initial focus
		tuix.SetFocus(nameID)
		return nil
	}, []any{})

	title := tuix.NewStyle().
		Bold(true).
		Foreground(tuix.BrightBlack)

	dim := tuix.NewStyle().
		Foreground(tuix.BrightBlack)

	label := tuix.NewStyle().
		Foreground(tuix.BrightWhite)

	// ⭐ Set focus order
	tuix.SetFocusOrder([]string{nameID, ageID, submitID})

	// ⭐ Handle Enter key to move focus
	switch {
	case tuix.IsFocused(nameID) && tuix.CurrentKey.Code == tuix.KeyEnter:
		tuix.Debug("Enter Pressed in Name → focusing Age")
		tuix.SetFocus(ageID)

	case tuix.IsFocused(ageID) && tuix.CurrentKey.Code == tuix.KeyEnter:
		tuix.Debug("Enter Pressed in Age → focusing Submit")
		tuix.SetFocus(submitID)

	case tuix.IsFocused(submitID) && tuix.CurrentKey.Code == tuix.KeyEnter:
		tuix.Debug("Submit pressed!")
		// Handle submit action
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

		// ⭐ Other Input (disabled)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			tuix.Text("Label:", label),
			components.TextInput(
				false, // ⭐ disabled
				components.WithID("other"),
				components.WithValue("disabled"),
				components.WithWidth(20),
				components.WithStyle(tuix.NewStyle().Foreground(tuix.BrightBlack)),
			),
		),

		// ⭐ NAME (focusable)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			tuix.Text("Name:", label),
			components.TextInput(
				tuix.IsFocused(nameID),
				components.WithID(nameID),
				components.WithValue(name),
				components.WithWidth(20),
				components.WithPrefix("["),
				components.WithSuffix("]"),
				components.WithOnChange(func(id, value string) {
					tuix.Debug("Name changed to:", value)
					setName(value)
				}),
			),
			func() tuix.Element {
				if tuix.IsFocused(nameID) {
					return tuix.Text(" █", tuix.NewStyle().Foreground(tuix.Cyan))
				}
				return tuix.Text("", tuix.NewStyle())
			}(),
		),

		// ⭐ EMAIL (NOT focusable)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			tuix.Text("Email:", label),
			components.TextInput(
				false, // ⭐ NOT focusable
				components.WithID(emailID),
				components.WithValue(email),
				components.WithWidth(20),
				components.WithPrefix("["),
				components.WithSuffix("]"),
				components.WithStyle(tuix.NewStyle().Foreground(tuix.BrightBlack)),
				components.WithOnChange(func(id, value string) {
					setEmail(value)
				}),
			),
			func() tuix.Element {
				return tuix.Text(" 🔒", tuix.NewStyle().Foreground(tuix.BrightBlack))
			}(),
		),

		// ⭐ AGE (focusable)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			tuix.Text("Age:", label),
			components.TextInput(
				tuix.IsFocused(ageID),
				components.WithID(ageID),
				components.WithValue(age),
				components.WithWidth(10),
				components.WithPrefix("["),
				components.WithSuffix("]"),
				components.WithOnChange(func(id, value string) {
					setAge(value)
				}),
			),
			func() tuix.Element {
				if tuix.IsFocused(ageID) {
					return tuix.Text(" █", tuix.NewStyle().Foreground(tuix.Cyan))
				}
				return tuix.Text("", tuix.NewStyle())
			}(),
		),

		// ⭐ SUBMIT BUTTON (focusable)
		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			func() tuix.Element {
				if tuix.IsFocused(submitID) {
					return tuix.Text("[ Submit ]", tuix.NewStyle().
						Background(tuix.Blue).
						Foreground(tuix.White).
						Bold(true))
				}
				return tuix.Text("[ Submit ]", tuix.NewStyle().
					Foreground(tuix.White))
			}(),
		),

		tuix.Text("", dim),
		tuix.Text("Tab: Next | Shift+Tab: Prev | Enter: Move Down", dim),
	)
}

func main() {
	app := tuix.NewApp(60, 13)
	app.Run(App, tuix.Props{})
}
