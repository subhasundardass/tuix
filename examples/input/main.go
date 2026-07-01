package main

import (
	"fmt"

	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	// State
	name, setName := tuix.UseState("")
	email, setEmail := tuix.UseState("")
	age, setAge := tuix.UseState("")
	salary, setSalary := tuix.UseState("")
	birthdate, setBirthdate := tuix.UseState("")
	joiningdate, setJoiningDate := tuix.UseState("")
	agree, setAgree := tuix.UseState(false)
	// gender, setGender := tuix.UseState("Male")
	submitted, setSubmitted := tuix.UseState(false)

	// Focus management
	focusIndex, setFocusIndex := tuix.UseState(0)
	totalFields := 9

	// Styles
	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	body := tuix.NewStyle().Foreground(tuix.BrightWhite)
	success := tuix.NewStyle().Foreground(tuix.BrightGreen).Bold(true)

	// Navigation
	switch tuix.CurrentKey.Code {
	case tuix.KeyDown:
		setFocusIndex((focusIndex + 1) % totalFields)
	case tuix.KeyUp:
		setFocusIndex((focusIndex - 1 + totalFields) % totalFields)
	case tuix.KeyEnter:
		if focusIndex == totalFields-1 {
			// Submit
			if name != "" && email != "" && age != "" && salary != "" && birthdate != "" && joiningdate != "" && agree {
				setSubmitted(true)
			}
		} else {
			setFocusIndex((focusIndex + 1) % totalFields)
		}
	}

	isFocused := func(idx int) bool { return focusIndex == idx }

	// Helper: render field row
	renderField := func(label string, fieldIndex int, input tuix.Element) tuix.Element {
		return tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.NewStyle(),
			tuix.Box(
				tuix.Props{Width: tuix.Fixed(30)},
				tuix.NewStyle(),
				tuix.Text(label+":", tuix.NewStyle().Foreground(tuix.BrightWhite)),
			),
			input,
		)
	}

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 0, Padding: [4]int{1, 1, 1, 1}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("◆ Registration Form", title),

		// Form Fields
		tuix.Box(
			tuix.Props{Direction: tuix.Column, Width: tuix.Grow(1)},
			tuix.NewStyle().Border(tuix.Border{
				Top: true, Right: true, Bottom: true, Left: true,
				Chars: tuix.BorderSharp,
				Title: "User Information",
			}),
			tuix.Box(
				tuix.Props{Direction: tuix.Column, Gap: 0, Padding: [4]int{1, 2, 2, 2}},
				tuix.NewStyle().Background(tuix.Black),

				// 1. Name (TextInput)
				renderField("Name", 0,
					components.TextInput(
						isFocused(0),
						components.WithID("name"),
						components.WithValue(name),
						components.WithWidth(30),
						components.WithPrefix("[ "),
						components.WithSuffix(" ]"),
						components.WithOnChange(func(id, value string) {
							setName(value)
						}),
					),
				),

				// 2. Email (TextInput)
				renderField("Email", 1,
					components.TextInput(
						isFocused(1),
						components.WithID("email"),
						components.WithValue(email),
						components.WithWidth(30),
						components.WithPrefix("[ "),
						components.WithSuffix(" ]"),
						components.WithOnChange(func(id, value string) {
							setEmail(value)
						}),
					),
				),

				components.Spinner("Spinner...."),

				// 3. Age (NumberInput)
				renderField("Age", 2,
					components.NumberInput(
						isFocused(2),
						components.NumberWithID("age"),
						components.NumberWithValue(age),
						components.NumberWithWidth(10),
						components.NumberWithPlaceholder("0"),
						components.NumberWithPrefix("[ "),
						components.NumberWithSuffix(" ]"),
						components.NumberWithDecimal(0),
						components.NumberWithMin(0),
						components.NumberWithMax(150),
						components.NumberWithOnChange(func(id, value string) {
							setAge(value)
						}),
					),
				),

				// 4. Salary (NumberInput with decimals)
				renderField("Salary", 3,
					components.NumberInput(
						isFocused(3),
						components.NumberWithID("salary"),
						components.NumberWithValue(salary),
						components.NumberWithWidth(15),
						components.NumberWithPlaceholder("0.00"),
						components.NumberWithPrefix("[ $ "),
						components.NumberWithSuffix(" ]"),
						components.NumberWithDecimal(2),
						components.NumberWithMin(0),
						components.NumberWithMax(999999.99),
						components.NumberWithOnChange(func(id, value string) {
							setSalary(value)
						}),
					),
				),

				// 5. Birth Date (DateInput)
				renderField("Birth Date", 4,
					components.DateInput(
						isFocused(4),
						components.DateWithID("birthday"),
						components.DateWithValue(birthdate),
						components.DateWithFormat("DD/MM/YYYY"),
						components.DateWithPrefix("[ "),
						components.DateWithSuffix(" ]"),
						components.DateWithWidth(15),
						components.DateWithPlaceholder("DD/MM/YYYY"),
						components.DateWithOnChange(func(id, value string) {
							setBirthdate(value)
						}),
					),
				),

				// 6. Joining Date (DateInput)
				renderField("Joining Date", 5,
					components.DateInput(
						isFocused(5),
						components.DateWithID("joiningtdate"),
						components.DateWithValue(joiningdate),
						components.DateWithFormat("DD/MM/YYYY"),
						components.DateWithPrefix("[ "),
						components.DateWithSuffix(" ]"),
						components.DateWithWidth(15),
						components.DateWithPlaceholder("DD/MM/YYYY"),
						components.DateWithOnChange(func(id, value string) {
							setJoiningDate(value)
						}),
					),
				),

				// 7. Checkbox
				tuix.Box(
					tuix.Props{Direction: tuix.Row, Gap: 1},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Width: tuix.Fixed(30)},
						tuix.NewStyle(),
						tuix.Text("", tuix.NewStyle()),
					),
					components.Checkbox(
						isFocused(6),
						func(checked bool) {
							setAgree(checked)
						},
					),
					tuix.Text("I agree to terms", tuix.NewStyle().Foreground(tuix.White)),
				),

				components.Spinner("Spinner...."),
				components.Badge("Badge", tuix.Green, tuix.Blue),
				components.ProgressBar(100.00, 100.00, tuix.Green),

				// 8. Submit Button
				tuix.Box(
					tuix.Props{Direction: tuix.Row},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Width: tuix.Fixed(31)},
						tuix.NewStyle(),
						tuix.Text("", tuix.NewStyle()),
					),
					// components.Button("Submit", isFocused(7)),
					components.Button(
						isFocused(7),
						components.WithButtonID("save"),
						components.WithLabel("Save"),
						components.WithOnPress(func(id string) {
							// controller.Save()
						}),
					),
				),

				// 9. Select Option
				tuix.Box(
					tuix.Props{Direction: tuix.Row},
					tuix.NewStyle(),
					tuix.Box(
						tuix.Props{Width: tuix.Fixed(31)},
						tuix.NewStyle(),
						tuix.Text("", tuix.NewStyle()),
					),
					components.SelectPicker([]string{"Male", "Female"}, isFocused(8)),
				),

				tuix.Text("", tuix.NewStyle()),
				tuix.Text("↓/↑: Navigate  |  Enter: Next/Submit ",
					tuix.NewStyle().Foreground(tuix.BrightBlack)),
			),
		),

		// Preview
		tuix.Box(
			tuix.Props{Direction: tuix.Column, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Grow(1)},
			tuix.NewStyle().Border(tuix.Border{
				Top: true, Right: true, Bottom: true, Left: true,
				Chars: tuix.BorderRounded, Title: "Output",
			}),
			tuix.Text("📋 Form Preview", tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)),
			func() tuix.Element {
				if !submitted {
					return tuix.Text("  Waiting for submission...", dim)
				}
				return tuix.Box(
					tuix.Props{Direction: tuix.Column, Gap: 0},
					tuix.NewStyle(),
					tuix.Text("✅ Submitted successfully!", success),
					tuix.Text(fmt.Sprintf("  Name: %s", name), body),
					tuix.Text(fmt.Sprintf("  Email: %s", email), body),
					tuix.Text(fmt.Sprintf("  Age: %s", age), body),
					tuix.Text(fmt.Sprintf("  Salary: $%s", salary), body),
					tuix.Text(fmt.Sprintf("  Birth Date: %s", birthdate), body),
					tuix.Text(fmt.Sprintf("  Joining Date: %s", joiningdate), body),
					tuix.Text(fmt.Sprintf("  Agreed: %v", agree), body),
				)
			}(),
		),

		tuix.Text("ctrl-c to quit", dim),
	)
}

func main() {
	app := tuix.NewApp(80, 28)
	app.Run(App, tuix.Props{})
}
