package components

import (
	"strings"

	"github.com/subhasundardass/tuix/tuix"
)

// ─── Table ───────────────────────────────────────────────────────────────────

// Table renders a bordered grid with a header row and navigable data rows.
// Up/Down moves row selection when focused.
func Table(headers []string, rows [][]string, focused bool, onChange func(int)) tuix.Element {
	selected, setSelected := tuix.UseState(0)

	if focused {
		if tuix.CurrentKey.Code == tuix.KeyUp && selected > 0 {
			setSelected(selected - 1)
		}
		if tuix.CurrentKey.Code == tuix.KeyDown && selected < len(rows)-1 {
			setSelected(selected + 1)
		}
	}

	if onChange != nil {
		onChange(selected)
	}

	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len([]rune(h))
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len([]rune(cell)) > colWidths[i] {
				colWidths[i] = len([]rune(cell))
			}
		}
	}

	segments := func(fill string, left, mid, right string) string {
		parts := make([]string, len(colWidths))
		for i, w := range colWidths {
			parts[i] = strings.Repeat(fill, w+2)
		}
		return left + strings.Join(parts, mid) + right
	}

	padCell := func(cells []string) string {
		parts := make([]string, len(colWidths))
		for i, w := range colWidths {
			cell := ""
			if i < len(cells) {
				cell = cells[i]
			}
			r := []rune(cell)
			for len(r) < w {
				r = append(r, ' ')
			}
			parts[i] = " " + string(r) + " "
		}
		return "│" + strings.Join(parts, "│") + "│"
	}

	headerStyle := tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)
	borderStyle := tuix.NewStyle().Foreground(tuix.BrightBlack)

	elems := []tuix.Element{
		tuix.Text(segments("─", "┌", "┬", "┐"), borderStyle),
		tuix.Text(padCell(headers), headerStyle),
		tuix.Text(segments("─", "├", "┼", "┤"), borderStyle),
	}

	for i, row := range rows {
		var style tuix.Style
		if i == selected && focused {
			style = tuix.NewStyle().Background(tuix.Blue).Foreground(tuix.White).Bold(true)
		} else if i == selected {
			style = tuix.NewStyle().Foreground(tuix.White).Bold(true)
		} else {
			style = tuix.NewStyle().Foreground(tuix.BrightBlack)
		}
		elems = append(elems, tuix.Text(padCell(row), style))
	}

	elems = append(elems, tuix.Text(segments("─", "└", "┴", "┘"), borderStyle))
	return tuix.Box(tuix.Props{Direction: tuix.Column}, tuix.NewStyle(), elems...)
}

// ─── Tabs ────────────────────────────────────────────────────────────────────

// Tabs renders a horizontal tab bar. Left/Right arrows switch tabs when focused.
func Tabs(tabs []string, focused bool, onChange func(int)) tuix.Element {
	active, setActive := tuix.UseState(0)

	if focused {
		if tuix.CurrentKey.Code == tuix.KeyLeft && active > 0 {
			setActive(active - 1)
		}
		if tuix.CurrentKey.Code == tuix.KeyRight && active < len(tabs)-1 {
			setActive(active + 1)
		}
	}

	if onChange != nil {
		onChange(active)
	}

	divider := tuix.Text("│", tuix.NewStyle().Foreground(tuix.BrightBlack))
	elems := make([]tuix.Element, 0, len(tabs)*2-1)
	for i, tab := range tabs {
		var style tuix.Style
		if i == active {
			style = tuix.NewStyle().Foreground(tuix.Black).Background(tuix.Cyan).Bold(true)
		} else if focused {
			style = tuix.NewStyle().Foreground(tuix.White)
		} else {
			style = tuix.NewStyle().Foreground(tuix.BrightBlack)
		}
		elems = append(elems, tuix.Text(" "+tab+" ", style))
		if i < len(tabs)-1 {
			elems = append(elems, divider)
		}
	}
	return tuix.Box(tuix.Props{Direction: tuix.Row}, tuix.NewStyle(), elems...)
}

// ─── Modal (flow-positioned) ─────────────────────────────────────────────────

// Modal renders a bordered panel as a normal flow element.
// Place it as the last child of its parent so it paints on top.
// For a true floating overlay, use ModalOverlay instead.
// Esc calls onClose.
func Modal(title string, visible bool, width int, onClose func(), children ...tuix.Element) tuix.Element {
	if !visible {
		return tuix.Box(tuix.Props{}, tuix.Style{})
	}

	if tuix.CurrentKey.Code == tuix.KeyEscape && onClose != nil {
		onClose()
	}

	inner := width - 2
	titlePart := "─ " + title + " "
	remaining := max(inner-len([]rune(titlePart)), 0)
	borderStyle := tuix.NewStyle().Foreground(tuix.Cyan)

	top := "┌" + titlePart + strings.Repeat("─", remaining) + "┐"
	bottom := "└" + strings.Repeat("─", inner) + "┘"

	rows := []tuix.Element{tuix.Text(top, borderStyle)}
	for _, child := range children {
		rows = append(rows, tuix.Box(
			tuix.Props{Direction: tuix.Row},
			tuix.NewStyle(),
			tuix.Text("│ ", borderStyle),
			child,
		))
	}
	rows = append(rows, tuix.Text(bottom, borderStyle))
	rows = append(rows, tuix.Text("  Esc to close", tuix.NewStyle().Foreground(tuix.BrightBlack)))

	return tuix.Box(tuix.Props{Direction: tuix.Column}, tuix.NewStyle(), rows...)
}

// ─── ModalOverlay (true floating overlay) ────────────────────────────────────

// ModalOverlay renders a bordered modal dialog as a true floating overlay —
// it paints at absolute screen coordinates (x, y), on top of whatever flow
// content sits beneath it, without affecting the layout of any sibling.
//
// Parameters:
//
//	title    — text shown in the top border, e.g. "Confirm"
//	visible  — when false, returns an empty zero-size element
//	x, y     — absolute screen column and row for the top-left corner
//	width    — total width including the border characters
//	onClose  — called when the user presses Esc; set state to hide the modal
//	children — content rows painted inside the border
//
// Usage:
//
//	showModal, setShowModal := tuix.UseState(false)
//	...
//	components.ModalOverlay("Confirm", showModal, 4, 6, 40, func() {
//	    setShowModal(false)
//	}, tuix.Text("Are you sure?", tuix.NewStyle()))
//
// Place the ModalOverlay anywhere in the element tree — position in the
// tree does not affect where it appears on screen. The overlay is always
// painted last within its subtree, so it visually covers siblings.
func ModalOverlay(
	title string,
	visible bool,
	x, y int,
	width int,
	onClose func(),
	children ...tuix.Element,
) tuix.Element {
	if !visible {
		// Zero-size invisible element — no layout space consumed.
		return tuix.Overlay(0, 0)
	}

	if tuix.CurrentKey.Code == tuix.KeyEscape && onClose != nil {
		onClose()
	}

	inner := width - 2
	titlePart := "─ " + title + " "
	remaining := max(inner-len([]rune(titlePart)), 0)

	borderStyle := tuix.NewStyle().Foreground(tuix.Cyan)
	dimStyle := tuix.NewStyle().Foreground(tuix.BrightBlack)

	top := "┌" + titlePart + strings.Repeat("─", remaining) + "┐"
	bottom := "└" + strings.Repeat("─", inner) + "┘"

	// Build the modal's visual content as a column of rows, each padded
	// with border characters on the sides.
	rows := []tuix.Element{
		tuix.Text(top, borderStyle),
	}
	for _, child := range children {
		rows = append(rows, tuix.Box(
			tuix.Props{Direction: tuix.Row},
			tuix.NewStyle(),
			tuix.Text("│ ", borderStyle),
			child,
			tuix.Text(" │", borderStyle),
		))
	}
	rows = append(rows,
		tuix.Text(bottom, borderStyle),
		tuix.Text("  Esc to close", dimStyle),
	)

	content := tuix.Box(
		tuix.Props{Direction: tuix.Column},
		tuix.NewStyle(),
		rows...,
	)

	// Wrap in an Overlay so paint() places it at (x, y) absolutely,
	// ignoring whatever rect flow layout would normally assign here.
	return tuix.Overlay(x, y, content)
}
