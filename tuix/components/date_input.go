package components

import (
	"github.com/subhasundardass/tuix/tuix"
)

type DateOption func(*DateConfig)

type DateConfig struct {
	ID          string
	Value       string
	Format      string
	Placeholder string
	Width       int
	Style       tuix.Style
	Prefix      string
	Suffix      string
	Min         string
	Max         string
	OnChange    func(id string, value string)
	OnKeyPress  func(id string, key tuix.Key) bool
	OnFocus     func(id string)
	OnBlur      func(id string)
	OnSubmit    func(id string, value string)
}

func DateWithID(id string) DateOption {
	return func(c *DateConfig) {
		c.ID = id
	}
}

func DateWithValue(value string) DateOption {
	return func(c *DateConfig) {
		c.Value = value
	}
}

func DateWithFormat(format string) DateOption {
	return func(c *DateConfig) {
		c.Format = format
	}
}

func DateWithPlaceholder(text string) DateOption {
	return func(c *DateConfig) {
		c.Placeholder = text
	}
}

func DateWithWidth(width int) DateOption {
	return func(c *DateConfig) {
		c.Width = width
	}
}

func DateWithStyle(style tuix.Style) DateOption {
	return func(c *DateConfig) {
		c.Style = style
	}
}

func DateWithPrefix(prefix string) DateOption {
	return func(c *DateConfig) {
		c.Prefix = prefix
	}
}

func DateWithSuffix(suffix string) DateOption {
	return func(c *DateConfig) {
		c.Suffix = suffix
	}
}

func DateWithMin(min string) DateOption {
	return func(c *DateConfig) {
		c.Min = min
	}
}

func DateWithMax(max string) DateOption {
	return func(c *DateConfig) {
		c.Max = max
	}
}

func DateWithOnChange(fn func(id string, value string)) DateOption {
	return func(c *DateConfig) {
		c.OnChange = fn
	}
}

func DateWithOnKeyPress(fn func(id string, key tuix.Key) bool) DateOption {
	return func(c *DateConfig) {
		c.OnKeyPress = fn
	}
}

func DateWithOnFocus(fn func(id string)) DateOption {
	return func(c *DateConfig) {
		c.OnFocus = fn
	}
}

func DateWithOnBlur(fn func(id string)) DateOption {
	return func(c *DateConfig) {
		c.OnBlur = fn
	}
}

func DateWithOnSubmit(fn func(id string, value string)) DateOption {
	return func(c *DateConfig) {
		c.OnSubmit = fn
	}
}

// ─── DateInput Component ───────────────────────────────────────────────────

func DateInput(focused bool, opts ...DateOption) tuix.Element {
	config := &DateConfig{
		ID:          "",
		Value:       "",
		Format:      "YYYY-MM-DD",
		Placeholder: "",
		Width:       20,
		Style:       tuix.NewStyle(),
		Prefix:      "",
		Suffix:      "",
		Min:         "",
		Max:         "",
		OnChange:    nil,
		OnKeyPress:  nil,
		OnFocus:     nil,
		OnBlur:      nil,
		OnSubmit:    nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	mask := config.Format
	if mask == "" {
		mask = "YYYY-MM-DD"
	}

	// State for raw digits
	rawDigits, setRawDigits := tuix.UseState(extractDigitsOnly(config.Value))
	cursorPos, setCursorPos := tuix.UseState(len(rawDigits))

	maxDigits := countDigitsInMask(mask)

	// Initialize or sync from parent
	if rawDigits != config.Value && config.Value != "" {
		setRawDigits(extractDigitsOnly(config.Value))
		setCursorPos(len(rawDigits))
	}

	// Clamp cursor
	if cursorPos < 0 {
		setCursorPos(0)
	}
	if cursorPos > len(rawDigits) {
		setCursorPos(len(rawDigits))
	}

	// Handle focused key input
	if focused {
		key := tuix.CurrentKey

		switch key.Code {
		case tuix.KeyLeft:
			if cursorPos > 0 {
				setCursorPos(cursorPos - 1)
			}

		case tuix.KeyRight:
			if cursorPos < len(rawDigits) {
				setCursorPos(cursorPos + 1)
			}

		case tuix.KeyBackspace:
			if cursorPos > 0 && len(rawDigits) > 0 {
				newRaw := rawDigits[:cursorPos-1] + rawDigits[cursorPos:]
				setRawDigits(newRaw)
				newPos := cursorPos - 1
				setCursorPos(newPos)
				formatted := applyMaskFormat(newRaw, mask)
				if config.OnChange != nil && config.ID != "" {
					config.OnChange(config.ID, formatted)
				}
			}

		case tuix.KeyEnter:
			if config.OnSubmit != nil && config.ID != "" {
				formatted := applyMaskFormat(rawDigits, mask)
				config.OnSubmit(config.ID, formatted)
			}

		default:
			// Only accept digits
			if key.Rune >= '0' && key.Rune <= '9' {
				if len(rawDigits) < maxDigits {
					newRaw := rawDigits[:cursorPos] + string(key.Rune) + rawDigits[cursorPos:]
					setRawDigits(newRaw)
					newPos := cursorPos + 1
					setCursorPos(newPos)

					formatted := applyMaskFormat(newRaw, mask)
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, formatted)
					}
				}
			}
		}
	}

	// Build display with formatting
	display := buildDateDisplayString(rawDigits, mask, config.Placeholder, focused, cursorPos)

	// Style
	var textStyle tuix.Style
	if focused {
		textStyle = tuix.NewStyle().
			Foreground(tuix.White).Background(tuix.Blue).
			Bold(true)
	} else {
		textStyle = tuix.NewStyle().Foreground(tuix.White)
	}

	// Apply prefix/suffix
	elements := []tuix.Element{}

	if config.Prefix != "" {
		prefixStyle := tuix.NewStyle()
		if focused {
			prefixStyle = prefixStyle.Foreground(tuix.Cyan).Bold(true)
		} else {
			prefixStyle = prefixStyle.Foreground(tuix.BrightBlack)
		}
		elements = append(elements, tuix.Text(config.Prefix, prefixStyle))
	}

	elements = append(elements, tuix.Text(display, textStyle))

	if config.Suffix != "" {
		suffixStyle := tuix.NewStyle()
		if focused {
			suffixStyle = suffixStyle.Foreground(tuix.Cyan).Bold(true)
		} else {
			suffixStyle = suffixStyle.Foreground(tuix.BrightBlack)
		}
		elements = append(elements, tuix.Text(config.Suffix, suffixStyle))
	}

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
		},
		tuix.NewStyle(),
		elements...,
	)
}

// ─── Helper Functions ─────────────────────────────────────────────────────

func extractDigitsOnly(value string) string {
	result := ""
	for _, ch := range value {
		if ch >= '0' && ch <= '9' {
			result += string(ch)
		}
	}
	return result
}

func countDigitsInMask(mask string) int {
	count := 0
	for _, ch := range mask {
		if ch == 'Y' || ch == 'M' || ch == 'D' {
			count++
		}
	}
	return count
}

func applyMaskFormat(rawDigits string, mask string) string {
	if rawDigits == "" {
		return ""
	}

	result := ""
	digitIdx := 0

	for _, maskChar := range mask {
		if maskChar == 'Y' || maskChar == 'M' || maskChar == 'D' {
			if digitIdx < len(rawDigits) {
				result += string(rawDigits[digitIdx])
				digitIdx++
			} else {
				result += string(maskChar)
			}
		} else {
			result += string(maskChar)
		}
	}

	return result
}

func buildDateDisplayString(rawDigits string, mask string, placeholder string, focused bool, cursorPos int) string {
	if rawDigits == "" {
		display := placeholder
		if display == "" {
			display = mask
		}
		if focused {
			return "█" + display
		}
		return display
	}

	formatted := applyMaskFormat(rawDigits, mask)

	if focused {
		cursorDisplayPos := calculateCursorDisplayPos(rawDigits, mask, cursorPos)

		runes := []rune(formatted)
		if cursorDisplayPos < len(runes) {
			return string(runes[:cursorDisplayPos]) + "█" + string(runes[cursorDisplayPos:])
		}
		return string(runes) + "█"
	}

	return formatted
}

func calculateCursorDisplayPos(rawDigits string, mask string, digitCursorPos int) int {
	displayPos := 0
	digitCount := 0

	for _, maskChar := range mask {
		if maskChar == 'Y' || maskChar == 'M' || maskChar == 'D' {
			if digitCount == digitCursorPos {
				return displayPos
			}
			digitCount++
		}
		displayPos++
	}

	return displayPos
}
