package components

import (
	"github.com/subhasundardass/tuix/tuix"
)

// ─── Button ───────────────────────────────────────────────────────────────
func Button(label string, focused bool) tuix.Element {
	var style tuix.Style
	if focused {
		style = tuix.NewStyle().
			Foreground(tuix.Black).
			Background(tuix.Cyan).
			Bold(true)
	} else {
		style = tuix.NewStyle().Foreground(tuix.White)
	}
	return tuix.Text("[ "+label+" ]", style)
}

// ─── Checkbox (NO LABEL) ──────────────────────────────────────────────────
func Checkbox(focused bool, onChange func(bool)) tuix.Element {
	checked, setChecked := tuix.UseState(false)

	if focused && tuix.CurrentKey.Code == tuix.KeySpace {
		newValue := !checked
		setChecked(newValue)
		if onChange != nil {
			onChange(newValue)
		}
	}

	box := "[ ]"
	if checked {
		box = "[x]"
	}
	var style tuix.Style
	if focused {
		style = tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	} else {
		style = tuix.NewStyle().Foreground(tuix.White).Bold(true)
	}
	return tuix.Text(box, style)
}

// ─── List (NO LABEL) ──────────────────────────────────────────────────────
func List(items []string, focused bool) tuix.Element {
	selected, setSelected := tuix.UseState(0)

	if focused {
		if tuix.CurrentKey.Code == tuix.KeyDown && selected < len(items)-1 {
			setSelected(selected + 1)
		}
		if tuix.CurrentKey.Code == tuix.KeyUp && selected > 0 {
			setSelected(selected - 1)
		}
	}

	children := make([]tuix.Element, len(items))
	for i, item := range items {
		prefix := "  "
		var style tuix.Style
		if i == selected {
			prefix = "> "
			if focused {
				style = tuix.NewStyle().
					Background(tuix.Blue).
					Foreground(tuix.Cyan).
					Bold(true)
			} else {
				style = tuix.NewStyle().Foreground(tuix.White).Bold(true)
			}
		} else {
			style = tuix.NewStyle().Foreground(tuix.BrightBlack)
		}
		children[i] = tuix.Text(prefix+item, style)
	}
	return tuix.Box(
		tuix.Props{Direction: tuix.Column},
		tuix.NewStyle(),
		children...)
}

// ─── SelectPicker (NO LABEL) ──────────────────────────────────────────────
func SelectPicker(options []string, focused bool) tuix.Element {
	selected, setSelected := tuix.UseState(0)

	if focused {
		if tuix.CurrentKey.Code == tuix.KeyLeft && selected > 0 {
			setSelected(selected - 1)
		} else if tuix.CurrentKey.Code == tuix.KeyRight && selected < len(options)-1 {
			setSelected(selected + 1)
		}
	}

	label := options[selected]
	const optWidth = 12
	for len([]rune(label)) < optWidth {
		label += " "
	}
	var style tuix.Style
	if focused {
		style = tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	} else {
		style = tuix.NewStyle().Foreground(tuix.White)
	}
	return tuix.Text("< "+label+" >", style)
}
