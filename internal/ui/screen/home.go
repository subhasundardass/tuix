package screen

import (
	"github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

func HomePage(ctx *context.AppContext, props tuix.Props) tuix.Element {

	name, setName := tuix.UseState("")
	email, setEmail := tuix.UseState("")

	// Focus management

	focusIndex, setFocusIndex := tuix.UseState(0)
	totalFields := 9

	// Navigation
	if tuix.IsFocused("content") {
		switch tuix.CurrentKey.Code {
		case tuix.KeyEscape:
			tuix.SetFocus("sidebar")
		case tuix.KeyDown:
			setFocusIndex((focusIndex + 1) % totalFields)
		case tuix.KeyUp:
			setFocusIndex((focusIndex - 1 + totalFields) % totalFields)
		}
	}

	isFocused := func(idx int) bool { return focusIndex == idx }

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),
		tuix.Text("This is  Home Page", tuix.NewStyle().Bold(true)),

		tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0, Padding: [4]int{1, 2, 2, 2}},
			tuix.NewStyle().Background(tuix.Black),
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

		tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0, Padding: [4]int{1, 2, 2, 2}},
			tuix.NewStyle().Background(tuix.Black),
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
	)
}
