package components

import (
	"strconv"
	"strings"

	"github.com/subhasundardass/tuix/tuix"
)

type NumberOption func(*NumberConfig)

type NumberConfig struct {
	ID          string
	Value       string
	Placeholder string
	Width       int
	Style       tuix.Style
	Prefix      string
	Suffix      string
	Decimal     int
	Min         float64
	Max         float64
	OnChange    func(id string, value string)
	OnKeyPress  func(id string, key tuix.Key) bool
	OnFocus     func(id string)
	OnBlur      func(id string)
	OnSubmit    func(id string, value string)
}

func NumberWithID(id string) NumberOption {
	return func(c *NumberConfig) {
		c.ID = id
	}
}

func NumberWithValue(value string) NumberOption {
	return func(c *NumberConfig) {
		c.Value = value
	}
}

func NumberWithPlaceholder(text string) NumberOption {
	return func(c *NumberConfig) {
		c.Placeholder = text
	}
}

func NumberWithWidth(width int) NumberOption {
	return func(c *NumberConfig) {
		c.Width = width
	}
}

func NumberWithStyle(style tuix.Style) NumberOption {
	return func(c *NumberConfig) {
		c.Style = style
	}
}

func NumberWithPrefix(prefix string) NumberOption {
	return func(c *NumberConfig) {
		c.Prefix = prefix
	}
}

func NumberWithSuffix(suffix string) NumberOption {
	return func(c *NumberConfig) {
		c.Suffix = suffix
	}
}

func NumberWithDecimal(places int) NumberOption {
	return func(c *NumberConfig) {
		c.Decimal = places
	}
}

func NumberWithMin(min float64) NumberOption {
	return func(c *NumberConfig) {
		c.Min = min
	}
}

func NumberWithMax(max float64) NumberOption {
	return func(c *NumberConfig) {
		c.Max = max
	}
}

func NumberWithOnChange(fn func(id string, value string)) NumberOption {
	return func(c *NumberConfig) {
		c.OnChange = fn
	}
}

func NumberWithOnKeyPress(fn func(id string, key tuix.Key) bool) NumberOption {
	return func(c *NumberConfig) {
		c.OnKeyPress = fn
	}
}

func NumberWithOnFocus(fn func(id string)) NumberOption {
	return func(c *NumberConfig) {
		c.OnFocus = fn
	}
}

func NumberWithOnBlur(fn func(id string)) NumberOption {
	return func(c *NumberConfig) {
		c.OnBlur = fn
	}
}

func NumberWithOnSubmit(fn func(id string, value string)) NumberOption {
	return func(c *NumberConfig) {
		c.OnSubmit = fn
	}
}

// ─── NumberInput ───────────────────────────────────────────────────────────

func NumberInput(focused bool, opts ...NumberOption) tuix.Element {
	config := &NumberConfig{
		ID:          "",
		Value:       "",
		Placeholder: "",
		Width:       20,
		Style:       tuix.NewStyle(),
		Prefix:      "[",
		Suffix:      "]",
		Decimal:     0,
		Min:         0,
		Max:         0,
		OnChange:    nil,
		OnKeyPress:  nil,
		OnFocus:     nil,
		OnBlur:      nil,
		OnSubmit:    nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	// Keep raw value, sanitize only for display
	rawValue := config.Value
	pos, setPos := tuix.UseState(len(rawValue))

	if focused && config.OnFocus != nil && config.ID != "" {
		config.OnFocus(config.ID)
	}

	if !focused && config.OnBlur != nil && config.ID != "" {
		config.OnBlur(config.ID)
	}

	if focused {
		key := tuix.CurrentKey

		if config.OnKeyPress != nil && config.ID != "" {
			if config.OnKeyPress(config.ID, key) {
				goto render
			}
		}

		switch key.Code {
		case tuix.KeyLeft:
			if pos > 0 {
				setPos(pos - 1)
			}

		case tuix.KeyRight:
			if pos < len(rawValue) {
				setPos(pos + 1)
			}

		case tuix.KeyBackspace:
			// SIMPLEST BACKSPACE - just delete the character
			if pos > 0 && len(rawValue) > 0 {
				newValue := rawValue[:pos-1] + rawValue[pos:]
				// Only check if it's a valid number format
				if isValidNumberFormat(newValue, config.Decimal) {
					rawValue = newValue
					config.Value = rawValue
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, rawValue)
					}
					setPos(pos - 1)
				}
			}

		case tuix.KeyHome:
			setPos(0)

		case tuix.KeyEnd:
			setPos(len(rawValue))

		case tuix.KeyDelete:
			if pos < len(rawValue) {
				newValue := rawValue[:pos] + rawValue[pos+1:]
				if isValidNumberFormat(newValue, config.Decimal) {
					rawValue = newValue
					config.Value = rawValue
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, rawValue)
					}
				}
			}

		case tuix.KeyEnter:
			if config.OnSubmit != nil && config.ID != "" {
				config.OnSubmit(config.ID, rawValue)
			}

		default:
			if key.Rune != 0 {
				char := string(key.Rune)
				if isValidNumberChar(char, rawValue, config.Decimal) {
					newValue := rawValue[:pos] + char + rawValue[pos:]
					if isValidNumberFormat(newValue, config.Decimal) {
						rawValue = newValue
						config.Value = rawValue
						if config.OnChange != nil && config.ID != "" {
							config.OnChange(config.ID, rawValue)
						}
						setPos(pos + 1)
					}
				}
			}
		}
	}

render:
	// Validate min/max
	isValid := true
	if rawValue != "" {
		if !validateNumberRange(rawValue, config) {
			isValid = false
		}
	}

	// Format display with padding
	display := formatNumberDisplay(rawValue, config.Decimal, config.Placeholder, focused)

	// Clamp pos for display
	currentPos := clampPos(pos, display)

	// Styles
	textStyle := config.Style
	if focused {
		textStyle = textStyle.Foreground(tuix.White).Background(tuix.Blue)
	} else {
		textStyle = textStyle.Foreground(tuix.BrightBlack)
	}

	borderColor := tuix.BrightBlack
	if focused {
		if isValid {
			borderColor = tuix.Cyan
		} else {
			borderColor = tuix.Red
		}
	}

	bracketStyle := tuix.NewStyle()
	if focused {
		bracketStyle = bracketStyle.Foreground(borderColor).Bold(true)
	} else {
		bracketStyle = bracketStyle.Foreground(tuix.BrightBlack)
	}

	// Cursor with safe bounds
	cursorDisplay := display
	if focused {
		runes := []rune(cursorDisplay)
		if currentPos < len(runes) {
			cursorDisplay = string(runes[:currentPos]) + "█" + string(runes[currentPos:])
		} else {
			cursorDisplay = string(runes) + "█"
		}
	}

	// Pad to fixed width
	paddedDisplay := cursorDisplay
	displayLen := len([]rune(paddedDisplay))
	if displayLen < config.Width {
		padding := strings.Repeat(" ", config.Width-displayLen)
		paddedDisplay = paddedDisplay + padding
	}

	elements := []tuix.Element{}

	if config.Prefix != "" {
		elements = append(elements, tuix.Text(config.Prefix, bracketStyle))
	}

	textStyleForDisplay := textStyle
	if !isValid && !focused {
		textStyleForDisplay = textStyleForDisplay.Foreground(tuix.Red)
	}

	elements = append(elements, tuix.Text(paddedDisplay, textStyleForDisplay))

	if config.Suffix != "" {
		elements = append(elements, tuix.Text(config.Suffix, bracketStyle))
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

func clampPos(pos int, value string) int {
	maxPos := len(value)
	if pos < 0 {
		return 0
	}
	if pos > maxPos {
		return maxPos
	}
	return pos
}

// Check if the number format is valid (only digits, optional decimal, optional minus)
func isValidNumberFormat(value string, decimal int) bool {
	if value == "" || value == "-" || value == "." || value == "-." {
		return true
	}

	// Check for valid characters
	decimalFound := false
	hasMinus := false

	for i, ch := range value {
		if ch == '-' {
			if i != 0 || hasMinus {
				return false
			}
			hasMinus = true
			continue
		}
		if ch >= '0' && ch <= '9' {
			continue
		}
		if ch == '.' {
			if decimal == 0 || decimalFound {
				return false
			}
			decimalFound = true
			continue
		}
		return false
	}

	// Check decimal places
	if decimalFound {
		parts := strings.Split(value, ".")
		if len(parts) == 2 && len(parts[1]) > decimal {
			return false
		}
	}

	return true
}

// Validate min/max range
func validateNumberRange(value string, config *NumberConfig) bool {
	if value == "" || value == "-" || value == "." || value == "-." {
		return true
	}

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return true
	}

	if config.Min != 0 && num < config.Min {
		return false
	}

	if config.Max != 0 && num > config.Max {
		return false
	}

	return true
}

// Check if character is valid for number input
func isValidNumberChar(char string, current string, decimal int) bool {
	if char >= "0" && char <= "9" {
		return true
	}

	if char == "." && decimal > 0 && !strings.Contains(current, ".") {
		return true
	}

	if char == "-" && len(current) == 0 {
		return true
	}

	return false
}

// Format display with padding and decimal places
func formatNumberDisplay(value string, decimal int, placeholder string, focused bool) string {
	if value == "" {
		if placeholder != "" && !focused {
			return placeholder
		}
		if decimal > 0 && !focused {
			return "0." + strings.Repeat("0", decimal)
		}
		return ""
	}

	// If value has decimal and we need to pad
	if strings.Contains(value, ".") {
		parts := strings.Split(value, ".")
		if len(parts) == 2 {
			if len(parts[1]) < decimal {
				return parts[0] + "." + parts[1] + strings.Repeat("0", decimal-len(parts[1]))
			}
		}
	}

	return value
}
