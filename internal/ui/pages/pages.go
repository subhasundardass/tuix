package pages

import "github.com/subhasundardass/tuix/tuix"

func HomePage(props tuix.Props) tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.Text("🏠 Home Page", tuix.NewStyle().Bold(true)),
	)
}

func SettingsPage(props tuix.Props) tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.Text("⚙️ Settings Page", tuix.NewStyle().Bold(true)),
	)
}

func AboutPage(props tuix.Props) tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.Text("About", tuix.NewStyle().Bold(true)),
	)
}
