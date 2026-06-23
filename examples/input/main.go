package main

import (
	"fmt"
	"strings"

	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	formData, setFormData := tuix.UseState(make(map[string]string))
	submitted, setSubmitted := tuix.UseState(false)

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	body := tuix.NewStyle().Foreground(tuix.BrightWhite)
	success := tuix.NewStyle().Foreground(tuix.BrightGreen).Bold(true)

	handleSubmit := func(data map[string]string) {
		setFormData(data)
		setSubmitted(true)
	}

	return tuix.Box(

		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("◆ Registration Form", title),

		tuix.Box(
			tuix.Props{Direction: tuix.Column, Width: tuix.Grow(1)},
			tuix.NewStyle().Border(tuix.Border{
				Top: true, Right: true, Bottom: true, Left: true,
				Chars: tuix.BorderSharp,
				Title: "User Information",
			}),
			func() tuix.Element {
				return ExampleForm(ExampleFormProps{
					OnSubmit: handleSubmit,
				})
			}(),
		),

		//---
		// ExampleBasicLegend(),

		tuix.Box(
			tuix.Props{Direction: tuix.Column, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Grow(1)},
			tuix.NewStyle().Border(tuix.Border{
				Top: true, Right: true, Bottom: true, Left: true,
				Chars: tuix.BorderRounded, Color: tuix.BrightCyan,
			}),
			tuix.Text("📋 Form Preview", tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)),

			func() tuix.Element {
				if submitted {
					return tuix.Box(
						tuix.Props{Direction: tuix.Column, Gap: 0},
						tuix.NewStyle(),
						tuix.Text("✅ Submitted successfully!", success),
						tuix.Text(fmt.Sprintf("  Name: %s", formData["name"]), body),
						tuix.Text(fmt.Sprintf("  Email: %s", formData["email"]), body),
						tuix.Text(fmt.Sprintf("  Age: %s", formData["age"]), body),
						tuix.Text(fmt.Sprintf("  Salary: $%s", formData["salary"]), body),
						tuix.Text(fmt.Sprintf("  Birth Date: %s", formData["birthdate"]), body),
						tuix.Text(fmt.Sprintf("  Joining Date: %s", formData["joiningdate"]), body),
						tuix.Text(fmt.Sprintf("  Agreed: %s", formData["agree"]), body),
					)
				}
				return tuix.Text("  Waiting for submission...", dim)
			}(),
		),

		tuix.Text("ctrl-c to quit", dim),
	)
}

type ExampleFormProps struct {
	OnSubmit func(map[string]string)
}

// ExampleForm - Manual rendering with Box + Label + Input
// NO NAVIGATION HERE - Form handles everything!
func ExampleForm(props ExampleFormProps) tuix.Element {
	// Fields definition ONLY - NO duplicate state!
	fields := []components.Field{
		{ID: "name", Label: "Name", Type: components.FieldTypeText},
		{ID: "email", Label: "Email", Type: components.FieldTypeText},
		{ID: "age", Label: "Age", Type: components.FieldTypeNumber},
		{ID: "salary", Label: "Salary", Type: components.FieldTypeNumber},
		{ID: "birthdate", Label: "Birth Date", Type: components.FieldTypeDate},
		{ID: "joiningdate", Label: "Joining Date", Type: components.FieldTypeDate},
		{ID: "agree", Label: "I agree to terms", Type: components.FieldTypeCheckbox},
		{ID: "submit", Label: "Submit", Type: components.FieldTypeButton},
	}

	return components.Form(components.FormProps{
		Fields: fields,
		Width:  55,
		OnSubmit: func(data map[string]string) {
			if props.OnSubmit != nil {
				props.OnSubmit(data)
			}
		},
		OnValidate: func(data map[string]string) map[string]string {
			errors := make(map[string]string)

			if strings.TrimSpace(data["name"]) == "" {
				errors["name"] = "Name is required"
			}
			if data["email"] == "" || !strings.Contains(data["email"], "@") {
				errors["email"] = "Valid email is required"
			}
			if data["age"] == "" || data["age"] == "0" {
				errors["age"] = "Age is required"
			}
			if data["salary"] == "" || data["salary"] == "0.00" {
				errors["salary"] = "Salary is required"
			}
			if data["birthdate"] == "" {
				errors["birthdate"] = "Birth date is required"
			}
			if data["joiningdate"] == "" {
				errors["joiningdate"] = "Joining date is required"
			}
			if data["agree"] != "true" {
				errors["agree"] = "You must agree to terms"
			}

			return errors
		},
		// Custom render function - Form manages state internally
		RenderFunc: func(focusedIndex int, setFocused func(int), formData map[string]string, setFormData func(map[string]string)) tuix.Element {
			isFocused := func(idx int) bool { return focusedIndex == idx }

			// Get values from formData (Form's internal state)
			name := formData["name"]
			email := formData["email"]
			age := formData["age"]
			salary := formData["salary"]
			birthdate := formData["birthdate"]
			joiningdate := formData["joiningdate"]
			// agree := formData["agree"] == "true"

			return tuix.Box(
				tuix.Props{
					Direction: tuix.Column,
					Gap:       0,
					Padding:   [4]int{1, 2, 2, 2},
				},
				tuix.NewStyle().Background(tuix.Black),

				// 0: Name
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Width: tuix.Fixed(30)},
						tuix.NewStyle(),
						tuix.Text("Name:", tuix.NewStyle().Foreground(tuix.White)), // ← Colon at end
					),
					components.Input(
						isFocused(0),
						name,
						func(v string) {
							formData["name"] = v
							setFormData(formData)
						},
					),
				),

				// 1: Email
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Direction: tuix.Row, Gap: 1, Width: tuix.Fixed(30)},
						tuix.NewStyle(),
						tuix.Text("Email", tuix.NewStyle().Foreground(tuix.White)),
					),
					components.Input(
						isFocused(1),
						email,
						func(v string) {
							formData["email"] = v
							setFormData(formData)
						},
					),
				),

				// 2: Age
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Direction: tuix.Row, Gap: 1, Width: tuix.Fixed(30)},
						tuix.NewStyle(),
						tuix.Text("Age", tuix.NewStyle().Foreground(tuix.White)),
					),
					components.NumberInput(components.NumberInputProps{
						Value:       age,
						Focused:     isFocused(2),
						Decimal:     0,
						Min:         nil,
						Max:         nil,
						Placeholder: "0",
						Width:       10,
						OnChange: func(v string) {
							formData["age"] = v
							setFormData(formData)
						},
					}),
				),

				// 3: Salary
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Direction: tuix.Row, Gap: 1, Width: tuix.Fixed(30)},
						tuix.NewStyle(),
						tuix.Text("Salary", tuix.NewStyle().Foreground(tuix.White)),
					),
					components.NumberInput(components.NumberInputProps{
						Value:       salary,
						Focused:     isFocused(3),
						Decimal:     2,
						Min:         nil,
						Max:         nil,
						Placeholder: "0.00",
						Width:       15,
						OnChange: func(v string) {
							formData["salary"] = v
							setFormData(formData)
						},
					}),
				),

				// 4: Birth Date
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Direction: tuix.Row, Gap: 1, Width: tuix.Fixed(30)},
						tuix.NewStyle(),
						tuix.Text("Birth Date", tuix.NewStyle().Foreground(tuix.White)),
					),
					components.DateInput(components.DateInputProps{
						Value:       birthdate,    // ✅ Passing state value
						Focused:     isFocused(4), // ✅ Proper focus control
						Mask:        "DD-MM-YYYY", // ✅ Valid mask
						Placeholder: "DD-MM-YYYY", // ✅ Helpful placeholder
						OnChange: func(v string) { // ✅ Callback works
							formData["birthdate"] = v
							setFormData(formData)
						},
					}),
				),

				// 5: Joining Date
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Direction: tuix.Row, Gap: 1, Width: tuix.Fixed(30)},
						tuix.NewStyle(),
						tuix.Text("Joining Date", tuix.NewStyle().Foreground(tuix.White)),
					),
					components.DateInput(components.DateInputProps{
						Value:       joiningdate,  // ✅ Passing state value
						Focused:     isFocused(5), // ✅ Proper focus control
						Mask:        "DD-MM-YYYY", // ✅ Valid mask
						Placeholder: "DD-MM-YYYY", // ✅ Helpful placeholder
						OnChange: func(v string) { // ✅ Callback works
							formData["joiningdate"] = v
							setFormData(formData)
						},
					}),
				),

				// 6: Checkbox
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					components.Checkbox(
						isFocused(6),
						func(checked bool) {
							if checked {
								formData["agree"] = "true"
							} else {
								formData["agree"] = "false"
							}
							setFormData(formData)
						},
					),
					tuix.Text("I agree to terms", tuix.NewStyle().Foreground(tuix.White)),
				),

				// 7: Submit Button
				tuix.Box(
					tuix.Props{Direction: tuix.Row},
					tuix.NewStyle(),
					components.Button("Submit", isFocused(7)),
				),

				tuix.Text("", tuix.NewStyle()),
				tuix.Text("↓/↑: Navigate  |  Enter: Next/Submit  |  Space: Toggle",
					tuix.NewStyle().Foreground(tuix.BrightBlack)),
			)
		},
	})
}

func main() {
	app := tuix.NewApp(80, 28)
	app.Run(App, tuix.Props{})
}
