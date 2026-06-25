package screen

import (
	"github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/tuix"
)

func HomePage(ctx *context.AppContext, props tuix.Props) tuix.Element {

	// fields := []components.Field{
	// 	{ID: "name", Label: "Name", Type: components.FieldTypeText},
	// 	{ID: "email", Label: "Email", Type: components.FieldTypeText},
	// }

	// fm := components.Form(components.FormProps{
	// 	Fields: fields,
	// 	Width:  55,
	// 	OnSubmit: func(data map[string]string) {

	// 	},

	// 	RenderFunc: func(focusedIndex int, setFocused func(int), formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	// 		isFocused := func(idx int) bool { return focusedIndex == idx }

	// 		// Get values from formData (Form's internal state)
	// 		name := formData["name"]
	// 		email := formData["email"]

	// 		return tuix.Box(
	// 			tuix.Props{
	// 				Direction: tuix.Column,
	// 				Gap:       0,
	// 				Padding:   [4]int{1, 2, 2, 2},
	// 			},
	// 			tuix.NewStyle().Background(tuix.Black),

	// 			tuix.Box(
	// 				tuix.Props{Direction: tuix.Row, Gap: 1},
	// 				tuix.NewStyle(),
	// 				tuix.Box(
	// 					tuix.Props{Width: tuix.Fixed(30)},
	// 					tuix.NewStyle(),
	// 					tuix.Text("Name:", tuix.NewStyle().Foreground(tuix.BrightWhite)), // ← Colon at end
	// 				),
	// 				components.Input(
	// 					isFocused(0),
	// 					name,
	// 					func(v string) {
	// 						formData["name"] = v
	// 						setFormData(formData)
	// 					},
	// 				),
	// 			),
	// 			tuix.Box(
	// 				tuix.Props{Direction: tuix.Row, Gap: 1},
	// 				tuix.NewStyle(),
	// 				tuix.Box(
	// 					tuix.Props{Direction: tuix.Row, Gap: 1, Width: tuix.Fixed(30)},
	// 					tuix.NewStyle(),
	// 					tuix.Text("Email", tuix.NewStyle().Foreground(tuix.White)),
	// 				),
	// 				components.Input(
	// 					isFocused(1),
	// 					email,
	// 					func(v string) {
	// 						formData["email"] = v
	// 						setFormData(formData)
	// 					},
	// 				),
	// 			),
	// 		)
	// 	},
	// })

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),
		tuix.Text("This is  Home Page", tuix.NewStyle().Bold(true)),
		// fm,
	)
}

func SettingsPage(ctx *context.AppContext, props tuix.Props) tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),
		tuix.Text("This is  Settings Page", tuix.NewStyle().Bold(true)),
	)
}

func AboutPage(ctx *context.AppContext, props tuix.Props, focused bool) tuix.Element {
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 1, 1, 1},
			Width:     props.Width,  // ← use passed width
			Height:    props.Height, // ← use passed height
		},
		tuix.NewStyle(),
		tuix.Text("This is About", tuix.NewStyle().Bold(true)),
	)
}
