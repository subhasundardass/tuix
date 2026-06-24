package screen

import (
	"github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/tuix"
)

func HomePage(ctx *context.AppContext, props tuix.Props) tuix.Element {

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),
		tuix.Text("This is  Home Page", tuix.NewStyle().Bold(true)),
	)
}

func SettingsPage(ctx *context.AppContext, props tuix.Props) tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),
		tuix.Text("This is  Settings Page", tuix.NewStyle().Bold(true)),
	)
}

func AboutPage(ctx *context.AppContext, props tuix.Props) tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),
		tuix.Text("This is About", tuix.NewStyle().Bold(true)),
	)
}
