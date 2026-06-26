package screen

import (
	"github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

func SettingsPage(ctx *context.AppContext, props tuix.Props) tuix.Element {
	setting, setSetting := tuix.UseState("")
	focusIndex, setFocusIndex := tuix.UseState(0)
	totalFields := 1

	//Check if content has focus
	contentFocused := tuix.IsFocused("content")

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

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),

		tuix.Text("⚙️ Settings Page", tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)),
		tuix.Text("Press ESC to go back", tuix.NewStyle().Foreground(tuix.BrightBlack)),

		tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0, Padding: [4]int{1, 2, 2, 2}},
			tuix.NewStyle().Background(tuix.Black),

			components.TextInput(
				contentFocused,
				components.WithID("setting"),
				components.WithValue(setting),
				components.WithWidth(30),
				components.WithPrefix("[ "),
				components.WithSuffix(" ]"),
				components.WithOnChange(func(id, value string) {
					setSetting(value)
				}),
			),
		),

		//Debug info
		tuix.Text(
			"🔵 Focus: ",
			tuix.NewStyle().Foreground(tuix.BrightBlack),
		),
	)
}
