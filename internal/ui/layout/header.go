package layout

import (
	"github.com/subhasundardass/tuix/internal/context"
	"github.com/subhasundardass/tuix/tuix"
)

func Header(props tuix.Props) tuix.Element {

	appctx := context.Use() // ← one line, done

	header := tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     tuix.Grow(1),
			Justify:   tuix.JustifySpaceBetween,
			Align:     tuix.AlignCenter,
		},
		tuix.NewStyle().Foreground(tuix.BrightCyan).
			Border(tuix.Border{Bottom: true, Left: true, Right: true, Top: true}),
		tuix.Text(appctx.AppName(), tuix.NewStyle()),
		tuix.Text("layout demo", tuix.NewStyle()),
		tuix.Text("v0.0.15", tuix.NewStyle()),
	)

	return header
}
