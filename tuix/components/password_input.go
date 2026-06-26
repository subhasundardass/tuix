package components

import (
	"strings"

	"github.com/subhasundardass/tuix/tuix"
)

type PasswordOption func(*PasswordConfig)

type PasswordConfig struct {
	ID          string
	Value       string
	Placeholder string
	Width       int
	Style       tuix.Style
	Prefix      string
	Suffix      string
	MaskChar    string // Character used to mask input (default: "*")
	MinLength   int
	MaxLength   int
	OnChange    func(id string, value string)
	OnKeyPress  func(id string, key tuix.Key) bool
	OnFocus     func(id string)
	OnBlur      func(id string)
	OnSubmit    func(id string, value string)
}

func PasswordWithID(id string) PasswordOption {
	return func(c *PasswordConfig) {
		c.ID = id
	}
}

func PasswordWithValue(value string) PasswordOption {
	return func(c *PasswordConfig) {
		c.Value = value
	}
}

func PasswordWithPlaceholder(text string) PasswordOption {
	return func(c *PasswordConfig) {
		c.Placeholder = text
	}
}

func PasswordWithWidth(width int) PasswordOption {
	return func(c *PasswordConfig) {
		c.Width = width
	}
}

func PasswordWithStyle(style tuix.Style) PasswordOption {
	return func(c *PasswordConfig) {
		c.Style = style
	}
}

func PasswordWithPrefix(prefix string) PasswordOption {
	return func(c *PasswordConfig) {
		c.Prefix = prefix
	}
}

func PasswordWithSuffix(suffix string) PasswordOption {
	return func(c *PasswordConfig) {
		c.Suffix = suffix
	}
}

func PasswordWithMaskChar(char string) PasswordOption {
	return func(c *PasswordConfig) {
		c.MaskChar = char
	}
}

func PasswordWithMinLength(min int) PasswordOption {
	return func(c *PasswordConfig) {
		c.MinLength = min
	}
}

func PasswordWithMaxLength(max int) PasswordOption {
	return func(c *PasswordConfig) {
		c.MaxLength = max
	}
}

func PasswordWithOnChange(fn func(id string, value string)) PasswordOption {
	return func(c *PasswordConfig) {
		c.OnChange = fn
	}
}

func PasswordWithOnKeyPress(fn func(id string, key tuix.Key) bool) PasswordOption {
	return func(c *PasswordConfig) {
		c.OnKeyPress = fn
	}
}

func PasswordWithOnFocus(fn func(id string)) PasswordOption {
	return func(c *PasswordConfig) {
		c.OnFocus = fn
	}
}

func PasswordWithOnBlur(fn func(id string)) PasswordOption {
	return func(c *PasswordConfig) {
		c.OnBlur = fn
	}
}

func PasswordWithOnSubmit(fn func(id string, value string)) PasswordOption {
	return func(c *PasswordConfig) {
		c.OnSubmit = fn
	}
}

// ─── PasswordInput ──────────────────────────────────────────────────────────

func PasswordInput(focused bool, opts ...PasswordOption) tuix.Element {
	config := &PasswordConfig{
		ID:          "",
		Value:       "",
		Placeholder: "",
		Width:       20,
		Style:       tuix.NewStyle(),
		Prefix:      "[",
		Suffix:      "]",
		MaskChar:    "*",
		MinLength:   0,
		MaxLength:   0,
		OnChange:    nil,
		OnKeyPress:  nil,
		OnFocus:     nil,
		OnBlur:      nil,
		OnSubmit:    nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	pos, setPos := tuix.UseState(len(config.Value))

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
			if pos < len(config.Value) {
				setPos(pos + 1)
			}

		case tuix.KeyBackspace:
			if pos > 0 && len(config.Value) > 0 {
				newValue := config.Value[:pos-1] + config.Value[pos:]
				if config.MaxLength == 0 || len(newValue) <= config.MaxLength {
					config.Value = newValue
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, newValue)
					}
					setPos(pos - 1)
				}
			}

		case tuix.KeyHome:
			setPos(0)

		case tuix.KeyEnd:
			setPos(len(config.Value))

		case tuix.KeyDelete:
			if pos < len(config.Value) {
				newValue := config.Value[:pos] + config.Value[pos+1:]
				if config.MaxLength == 0 || len(newValue) <= config.MaxLength {
					config.Value = newValue
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, newValue)
					}
				}
			}

		case tuix.KeyEnter:
			if config.OnSubmit != nil && config.ID != "" {
				config.OnSubmit(config.ID, config.Value)
			}

		default:
			if key.Rune != 0 && key.Rune >= 32 {
				newValue := config.Value[:pos] + string(key.Rune) + config.Value[pos:]
				if config.MaxLength == 0 || len(newValue) <= config.MaxLength {
					config.Value = newValue
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, newValue)
					}
					setPos(pos + 1)
				}
			}
		}
	}

render:
	//Check min length
	isValid := true
	if config.MinLength > 0 && len(config.Value) < config.MinLength {
		isValid = false
	}

	//Build masked display
	display := ""
	if focused {
		// Show actual text when focused (or use mask)
		display = config.Value
	} else {
		// Mask the value when not focused
		if config.Value != "" {
			display = strings.Repeat(config.MaskChar, len(config.Value))
		}
	}

	if display == "" && config.Placeholder != "" {
		display = config.Placeholder
	}

	textStyle := config.Style
	if focused {
		textStyle = textStyle.Foreground(tuix.White)
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

	//Cursor
	cursorDisplay := display
	if focused {
		runes := []rune(cursorDisplay)
		if pos < len(runes) {
			cursorDisplay = string(runes[:pos]) + "█" + string(runes[pos:])
		} else {
			cursorDisplay = string(runes) + "█"
		}
	}

	//Pad to width
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
