package layout

import (
	"fmt"

	"github.com/subhasundardass/tuix/tuix"
)

func Footer(props tuix.Props) tuix.Element {

	footer := tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     tuix.Grow(1),
			Justify:   tuix.JustifySpaceBetween, // ← Pushes left/right apart
		},
		tuix.NewStyle(),
		tuix.Box(
			tuix.Props{},
			tuix.NewStyle(),
			tuix.Text("Ready", tuix.Style{}),
		),
		tuix.Box(
			tuix.Props{},
			tuix.NewStyle(),
			tuix.Text(fmt.Sprintf("v%s", "1.0.1"), tuix.Style{}),
		),
	)

	return footer
}
