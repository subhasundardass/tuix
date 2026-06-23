// Example: input
// Demonstrates: the Input component (typing, backspace, paste with Cmd+V),
// plus reading the live value back into a "preview" view.
//
// Run: go run ./examples/input
// See: ../../DOCS.md#input

package main

import (
	"strings"

	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {

	// focusIndex, setFocusIndex := tuix.UseState(0)
	// ⭐ Helper to check if a component is focused
	// isFocused := func(idx int) bool {
	// 	return focusIndex == idx
	// }

	// // ⭐ Handle Tab key to move focus
	// if tuix.CurrentKey.Code == tuix.KeyEnter {
	// 	// Tab moves forward
	// 	if tuix.CurrentKey.Code != tuix.KeyEnter {
	// 		setFocusIndex((focusIndex + 1) % 2) // 2 components
	// 	} else {
	// 		// Shift+Tab moves backward
	// 		setFocusIndex((focusIndex - 1 + 2) % 2)
	// 	}
	// }

	// value, setValue := tuix.UseState("")
	// _, setAgree := tuix.UseState(false)

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	// body := tuix.NewStyle().Foreground(tuix.BrightWhite)

	field := tuix.Box(
		tuix.Props{Direction: tuix.Column, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Grow(1)},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: tuix.BorderRounded, Color: tuix.BrightYellow,
		}),
		// Input(label, focused, value, onChange).
		// focused=true here because there's only one field; in a real
		// form you'd track focus with another UseState.
		// components.Input("name", isFocused(0), value, setValue),
		// components.Checkbox(
		// 	"I Agree",
		// 	isFocused(1),
		// 	func(v bool) {
		// 		setAgree(v)
		// 	},
		// ),

		RegistrationPage(tuix.Props{}),
	)

	preview := tuix.Box(
		tuix.Props{Direction: tuix.Column, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("you typed:", dim),
		// tuix.Text(
		// 	fmt.Sprintf("Focus: %d (Tab to navigate)", focusIndex),
		// 	tuix.NewStyle().Foreground(tuix.Blue),
		// ),
		// tuix.WrappedText(value, body),
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("◆ type something (paste works too!)", title),

		field,
		preview,
		tuix.Text("ctrl-c to quit", dim),
	)
}

func RegistrationPage(props tuix.Props) tuix.Element {
	// ⭐ Define fields
	fields := []components.Field{
		{
			ID:    "name",
			Label: "Name",
			Type:  "text",
		},
		{
			ID:    "email",
			Label: "Email",
			Type:  "text",
		},
		{
			ID:    "agree",
			Label: "I agree to terms",
			Type:  "checkbox",
		},

		{
			ID:    "submit",
			Label: "Submit",
			Type:  "button",
			OnSubmit: func() {
				// Handle button click
			},
		},
	}

	// ⭐ Form component handles everything!
	return components.Form(components.FormProps{
		Fields: fields,
		Width:  50,
		OnSubmit: func(data map[string]string) {
			// Handle form submission
			// ctx.NavigateTo("dashboard")
		},
		OnValidate: func(data map[string]string) map[string]string {
			errors := make(map[string]string)

			if data["name"] == "" {
				errors["name"] = "Name is required"
			}
			if data["email"] == "" || !strings.Contains(data["email"], "@") {
				errors["email"] = "Valid email is required"
			}
			if data["agree"] != "true" {
				errors["agree"] = "You must agree to terms"
			}

			return errors
		},
	})
}

func main() {
	app := tuix.NewApp(80, 12)
	app.Run(App, tuix.Props{})
}
