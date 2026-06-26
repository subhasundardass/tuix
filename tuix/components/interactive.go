package components

import (
	"github.com/subhasundardass/tuix/tuix"
)

// ansiSequence matches CSI escape sequences (colors, cursor moves, etc).
// var ansiSequence = regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]`)

// lineEndings normalizes clipboard line breaks to '\n'.
// var lineEndings = strings.NewReplacer("\r\n", "\n", "\r", "\n")

// sanitizePaste filters clipboard text for safe display
// func sanitizePaste(s string) string {
// 	s = lineEndings.Replace(s)
// 	s = ansiSequence.ReplaceAllString(s, "")
// 	return strings.Map(func(r rune) rune {
// 		if r == '\n' || r == '\t' {
// 			return r
// 		}
// 		if r < 0x20 || r == 0x7F {
// 			return -1
// 		}
// 		return r
// 	}, s)
// }

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

// ─── Text Input (NO LABEL) ──────────────────────────────────────────────
// func Input(
// 	focused bool,
// 	value string,
// 	onChange func(value string),
// ) tuix.Element {
// 	runes := []rune(value)

// 	pos, setPos := tuix.UseState(len(runes))
// 	clamped := pos
// 	if clamped < 0 {
// 		clamped = 0
// 	}
// 	if clamped > len(runes) {
// 		clamped = len(runes)
// 	}
// 	if clamped != pos {
// 		setPos(clamped)
// 		pos = clamped
// 	}

// 	if focused {
// 		switch tuix.CurrentKey.Code {
// 		case tuix.KeyLeft:
// 			if pos > 0 {
// 				setPos(pos - 1)
// 			}
// 		case tuix.KeyRight:
// 			if pos < len(runes) {
// 				setPos(pos + 1)
// 			}
// 		case tuix.KeyBackspace:
// 			if pos > 0 {
// 				newRunes := append([]rune{}, runes[:pos-1]...)
// 				newRunes = append(newRunes, runes[pos:]...)
// 				onChange(string(newRunes))
// 				setPos(pos - 1)
// 			}
// 		case tuix.KeySpace:
// 			newRunes := append([]rune{}, runes[:pos]...)
// 			newRunes = append(newRunes, ' ')
// 			newRunes = append(newRunes, runes[pos:]...)
// 			onChange(string(newRunes))
// 			setPos(pos + 1)
// 		case tuix.KeyPaste:
// 			insert := []rune(sanitizePaste(tuix.CurrentKey.Paste))
// 			if len(insert) > 0 {
// 				newRunes := append([]rune{}, runes[:pos]...)
// 				newRunes = append(newRunes, insert...)
// 				newRunes = append(newRunes, runes[pos:]...)
// 				onChange(string(newRunes))
// 				setPos(pos + len(insert))
// 			}
// 		default:
// 			if tuix.CurrentKey.Rune != 0 {
// 				newRunes := append([]rune{}, runes[:pos]...)
// 				newRunes = append(newRunes, tuix.CurrentKey.Rune)
// 				newRunes = append(newRunes, runes[pos:]...)
// 				onChange(string(newRunes))
// 				setPos(pos + 1)
// 			}
// 		}
// 	}

// 	var fieldStyle tuix.Style
// 	if focused {
// 		fieldStyle = tuix.NewStyle().Foreground(tuix.BrightWhite).Bold(true)
// 	} else {
// 		fieldStyle = tuix.NewStyle().Foreground(tuix.BrightBlack).Bold(true)
// 	}

// 	var display string
// 	if focused {
// 		if pos < len(runes) {
// 			display = string(runes[:pos]) + "█" + string(runes[pos+1:])
// 		} else {
// 			display = string(runes) + "█"
// 		}
// 	} else {
// 		display = value
// 	}

// 	return tuix.WrappedText(display, fieldStyle)
// }

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

// ─── Date Input (NO LABEL) ─────────────────────────────────────────────────────
// type DateInputProps struct {
// 	Value       string
// 	Focused     bool
// 	Mask        string
// 	Placeholder string
// 	OnChange    func(string)
// 	OnSubmit    func(string)
// }

// func DateInput(props DateInputProps) tuix.Element {
// 	// Use separate state for editing
// 	rawDigits, setRawDigits := tuix.UseState(props.Value)
// 	cursorPos, setCursorPos := tuix.UseState(0)

// 	mask := props.Mask
// 	if mask == "" {
// 		mask = "YYYY-MM-DD"
// 	}

// 	maxDigits := countDigits(mask)

// 	// ✅ Initialize or sync from parent
// 	if rawDigits != props.Value && props.Value != "" {
// 		setRawDigits(props.Value)
// 		setCursorPos(len(props.Value))
// 	}

// 	// ✅ Clamp cursor
// 	if cursorPos < 0 {
// 		setCursorPos(0)
// 	}
// 	if cursorPos > len(rawDigits) {
// 		setCursorPos(len(rawDigits))
// 	}

// 	// ✅ Handle focused key input
// 	if props.Focused {
// 		key := tuix.CurrentKey

// 		switch key.Code {
// 		case tuix.KeyLeft:
// 			if cursorPos > 0 {
// 				setCursorPos(cursorPos - 1)
// 			}

// 		case tuix.KeyRight:
// 			if cursorPos < len(rawDigits) {
// 				setCursorPos(cursorPos + 1)
// 			}

// 		case tuix.KeyBackspace:
// 			if cursorPos > 0 && len(rawDigits) > 0 {
// 				newRaw := rawDigits[:cursorPos-1] + rawDigits[cursorPos:]
// 				setRawDigits(newRaw)
// 				newPos := cursorPos - 1
// 				setCursorPos(newPos)
// 				if props.OnChange != nil {
// 					props.OnChange(newRaw)
// 				}
// 			}

// 		case tuix.KeyEnter:
// 			if props.OnSubmit != nil {
// 				formatted := formatDateForSubmit(rawDigits, mask)
// 				props.OnSubmit(formatted)
// 			}

// 		default:
// 			// ✅ Only accept digits
// 			if key.Rune >= '0' && key.Rune <= '9' {
// 				// ✅ Check max length
// 				if len(rawDigits) < maxDigits {
// 					newRaw := rawDigits[:cursorPos] + string(key.Rune) + rawDigits[cursorPos:]
// 					setRawDigits(newRaw)
// 					newPos := cursorPos + 1
// 					setCursorPos(newPos)

// 					// ✅ Send raw digits to parent
// 					if props.OnChange != nil {
// 						props.OnChange(newRaw)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// ✅ Build display with formatting
// 	display := buildDateDisplay(rawDigits, mask, props.Placeholder, props.Focused, cursorPos)

// 	// ✅ Styling
// 	var style tuix.Style
// 	if props.Focused {
// 		style = tuix.NewStyle().
// 			Foreground(tuix.White).
// 			Background(tuix.BrightBlack).Bold(true)
// 	} else {
// 		style = tuix.NewStyle().Foreground(tuix.White).Bold(true)
// 	}

// 	return tuix.Text(display, style)
// }

// // ─── Date Display Helper ──────────────────────────────────────────────────
// func buildDateDisplay(rawDigits string, mask string, placeholder string, focused bool, cursorPos int) string {
// 	// If empty, show placeholder or mask
// 	if rawDigits == "" {
// 		display := placeholder
// 		if display == "" {
// 			display = mask
// 		}
// 		if focused {
// 			runes := []rune(display)
// 			if len(runes) > 0 {
// 				return "█" + string(runes)
// 			}
// 			return "█"
// 		}
// 		return display
// 	}

// 	// Format with mask
// 	formatted := applyMask(rawDigits, mask)

// 	// Add cursor if focused
// 	if focused {
// 		// Calculate cursor position in formatted string
// 		cursorDisplayPos := calculateCursorPos(rawDigits, mask, cursorPos)

// 		runes := []rune(formatted)
// 		if cursorDisplayPos < len(runes) {
// 			return string(runes[:cursorDisplayPos]) + "█" + string(runes[cursorDisplayPos:])
// 		}
// 		return string(runes) + "█"
// 	}

// 	return formatted
// }

// // ✅ Apply mask formatting to raw digits
// func applyMask(rawDigits string, mask string) string {
// 	result := ""
// 	digitIdx := 0

// 	for _, maskChar := range mask {
// 		if maskChar == 'Y' || maskChar == 'M' || maskChar == 'D' {
// 			if digitIdx < len(rawDigits) {
// 				result += string(rawDigits[digitIdx])
// 				digitIdx++
// 			} else {
// 				result += string(maskChar)
// 			}
// 		} else {
// 			result += string(maskChar)
// 		}
// 	}

// 	return result
// }

// // ✅ Calculate where cursor should be in formatted string
// func calculateCursorPos(rawDigits string, mask string, digitCursorPos int) int {
// 	displayPos := 0
// 	digitCount := 0

// 	for _, maskChar := range mask {
// 		if maskChar == 'Y' || maskChar == 'M' || maskChar == 'D' {
// 			if digitCount == digitCursorPos {
// 				return displayPos
// 			}
// 			digitCount++
// 		}
// 		displayPos++
// 	}

// 	return displayPos
// }

// // ✅ Format for submission
// func formatDateForSubmit(rawDigits string, mask string) string {
// 	return applyMask(rawDigits, mask)
// }

// ─── Helper ──────────────────────────────────────────────────────
// func countDigits(mask string) int {
// 	count := 0
// 	for _, ch := range mask {
// 		if ch == 'Y' || ch == 'M' || ch == 'D' {
// 			count++
// 		}
// 	}
// 	return count
// }

// ─── Number Input (NO LABEL) ──────────────────────────────────────────────
// type NumberInputProps struct {
// 	Value       string
// 	Focused     bool
// 	Decimal     int
// 	Min         *float64
// 	Max         *float64
// 	Step        float64
// 	Placeholder string
// 	Width       int
// 	OnChange    func(string)
// 	OnSubmit    func(string)
// }

// func NumberInput(props NumberInputProps) tuix.Element {
// 	rawValue, setRawValue := tuix.UseState("")
// 	pos, setPos := tuix.UseState(0)

// 	decimalPlaces := props.Decimal
// 	if decimalPlaces < 0 {
// 		decimalPlaces = 0
// 	}

// 	if props.Focused {
// 		switch tuix.CurrentKey.Code {
// 		case tuix.KeyLeft:
// 			if pos > 0 {
// 				setPos(pos - 1)
// 			}
// 		case tuix.KeyRight:
// 			if pos < len(rawValue) {
// 				setPos(pos + 1)
// 			}
// 		case tuix.KeyBackspace:
// 			if pos > 0 && len(rawValue) > 0 {
// 				newRaw := rawValue[:pos-1] + rawValue[pos:]
// 				setRawValue(newRaw)
// 				setPos(pos - 1)
// 				if props.OnChange != nil {
// 					props.OnChange(newRaw)
// 				}
// 			}
// 		case tuix.KeyEnter:
// 			if props.OnSubmit != nil {
// 				props.OnSubmit(rawValue)
// 			}
// 		default:
// 			if tuix.CurrentKey.Rune != 0 {
// 				ch := string(tuix.CurrentKey.Rune)
// 				if isNumberChar(ch, rawValue, decimalPlaces) {
// 					newRaw := rawValue[:pos] + ch + rawValue[pos:]
// 					if validateNumberRange(newRaw, props.Min, props.Max) {
// 						setRawValue(newRaw)
// 						setPos(pos + 1)
// 						if props.OnChange != nil {
// 							props.OnChange(newRaw)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	display := buildNumberDisplay(rawValue, decimalPlaces, props.Placeholder)

// 	var displayWithCursor string
// 	if props.Focused {
// 		displayWithCursor = insertNumberCursor(display, rawValue, pos, decimalPlaces)
// 	} else {
// 		displayWithCursor = display
// 	}

// 	width := props.Width
// 	if width == 0 {
// 		width = 20
// 	}

// 	displayLen := len([]rune(displayWithCursor))
// 	if displayLen < width {
// 		padding := strings.Repeat(" ", width-displayLen)
// 		displayWithCursor = padding + displayWithCursor
// 	}

// 	fieldStyle := tuix.NewStyle()
// 	if props.Focused {
// 		fieldStyle = fieldStyle.Foreground(tuix.White).Bold(true)
// 	} else {
// 		fieldStyle = fieldStyle.Foreground(tuix.BrightBlack).Bold(true)
// 	}

// 	return tuix.WrappedText(displayWithCursor, fieldStyle)
// }

// // ─── Number Helpers ──────────────────────────────────────────────────────
// func buildNumberDisplay(raw string, decimalPlaces int, placeholder string) string {
// 	if raw == "" {
// 		if placeholder != "" {
// 			return placeholder
// 		}
// 		if decimalPlaces > 0 {
// 			return "0." + strings.Repeat("0", decimalPlaces)
// 		}
// 		return "0"
// 	}

// 	display := raw

// 	if !strings.Contains(display, ".") && decimalPlaces > 0 {
// 		display = display + "." + strings.Repeat("0", decimalPlaces)
// 	}

// 	if strings.Contains(display, ".") {
// 		parts := strings.Split(display, ".")
// 		intPart := parts[0]
// 		decPart := parts[1]

// 		if len(decPart) > decimalPlaces {
// 			decPart = decPart[:decimalPlaces]
// 		} else if len(decPart) < decimalPlaces {
// 			decPart = decPart + strings.Repeat("0", decimalPlaces-len(decPart))
// 		}

// 		if decimalPlaces > 0 {
// 			display = intPart + "." + decPart
// 		} else {
// 			display = intPart
// 		}
// 	}

// 	return display
// }

// func insertNumberCursor(display string, raw string, pos int, decimalPlaces int) string {
// 	if raw == "" {
// 		runes := []rune(display)
// 		if len(runes) > 0 {
// 			return "█" + string(runes)
// 		}
// 		return "█"
// 	}

// 	cursorDisplayPos := 0

// 	if !strings.Contains(raw, ".") && decimalPlaces > 0 {
// 		decIndex := strings.Index(display, ".")
// 		if decIndex == -1 {
// 			cursorDisplayPos = pos
// 		} else {
// 			if pos <= len(raw) {
// 				cursorDisplayPos = pos
// 			} else {
// 				cursorDisplayPos = decIndex + 1 + (pos - len(raw))
// 			}
// 		}
// 	} else if strings.Contains(raw, ".") {
// 		rawParts := strings.Split(raw, ".")
// 		rawIntLen := len(rawParts[0])

// 		decIndex := strings.Index(display, ".")
// 		if decIndex == -1 {
// 			cursorDisplayPos = pos
// 		} else {
// 			if pos <= rawIntLen {
// 				cursorDisplayPos = pos
// 			} else {
// 				cursorDisplayPos = decIndex + 1 + (pos - rawIntLen - 1)
// 			}
// 		}
// 	} else {
// 		cursorDisplayPos = pos
// 	}

// 	if cursorDisplayPos < 0 {
// 		cursorDisplayPos = 0
// 	}
// 	if cursorDisplayPos > len(display) {
// 		cursorDisplayPos = len(display)
// 	}

// 	runes := []rune(display)
// 	if cursorDisplayPos < len(runes) {
// 		return string(runes[:cursorDisplayPos]) + "█" + string(runes[cursorDisplayPos:])
// 	}
// 	return string(runes) + "█"
// }

// func isNumberChar(ch string, current string, decimalPlaces int) bool {
// 	if ch >= "0" && ch <= "9" {
// 		return true
// 	}
// 	if ch == "." && decimalPlaces > 0 {
// 		return !strings.Contains(current, ".")
// 	}
// 	if ch == "-" && len(current) == 0 {
// 		return true
// 	}
// 	return false
// }

// func validateNumberRange(raw string, min, max *float64) bool {
// 	if raw == "" || raw == "-" || raw == "." {
// 		return true
// 	}

// 	val, err := strconv.ParseFloat(raw, 64)
// 	if err != nil {
// 		return true
// 	}

// 	if min != nil && val < *min {
// 		return false
// 	}
// 	if max != nil && val > *max {
// 		return false
// 	}

// 	return true
// }

// func GetNumberValue(raw string) (float64, error) {
// 	if raw == "" {
// 		return 0, nil
// 	}
// 	return strconv.ParseFloat(raw, 64)
// }
