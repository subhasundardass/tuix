package components

import (
	"strings"

	"github.com/subhasundardass/tuix/tuix"
)

type InputOption func(*InputConfig)

type InputConfig struct {
	ID          string
	Value       string
	Placeholder string
	Width       int
	Style       tuix.Style
	Prefix      string
	Suffix      string
	MinLength   int // Minimum length (0 = no limit)
	MaxLength   int // Maximum length (0 = no limit)
	OnChange    func(id string, value string)
	OnKeyPress  func(id string, key tuix.Key) bool
	OnFocus     func(id string)
	OnBlur      func(id string)
	OnSubmit    func(id string, value string)
}

// ─── Option Functions ──────────────────────────────────────────────────────

func WithID(id string) InputOption {
	return func(c *InputConfig) {
		c.ID = id
	}
}

func WithValue(value string) InputOption {
	return func(c *InputConfig) {
		c.Value = value
	}
}

func WithPlaceholder(text string) InputOption {
	return func(c *InputConfig) {
		c.Placeholder = text
	}
}

func WithWidth(width int) InputOption {
	return func(c *InputConfig) {
		c.Width = width
	}
}

func WithStyle(style tuix.Style) InputOption {
	return func(c *InputConfig) {
		c.Style = style
	}
}

func WithPrefix(prefix string) InputOption {
	return func(c *InputConfig) {
		c.Prefix = prefix
	}
}

func WithSuffix(suffix string) InputOption {
	return func(c *InputConfig) {
		c.Suffix = suffix
	}
}

// Length validators
func WithMinLength(min int) InputOption {
	return func(c *InputConfig) {
		c.MinLength = min
	}
}

func WithMaxLength(max int) InputOption {
	return func(c *InputConfig) {
		c.MaxLength = max
	}
}

func WithOnChange(fn func(id string, value string)) InputOption {
	return func(c *InputConfig) {
		c.OnChange = fn
	}
}

func WithOnKeyPress(fn func(id string, key tuix.Key) bool) InputOption {
	return func(c *InputConfig) {
		c.OnKeyPress = fn
	}
}

func WithOnFocus(fn func(id string)) InputOption {
	return func(c *InputConfig) {
		c.OnFocus = fn
	}
}

func WithOnBlur(fn func(id string)) InputOption {
	return func(c *InputConfig) {
		c.OnBlur = fn
	}
}

func WithOnSubmit(fn func(id string, value string)) InputOption {
	return func(c *InputConfig) {
		c.OnSubmit = fn
	}
}

// ─── TextInput Component ────────────────────────────────────────────────────

func TextInput(focused bool, opts ...InputOption) tuix.Element {

	config := &InputConfig{
		ID:          "",
		Value:       "",
		Placeholder: "",
		Width:       30,
		Style:       tuix.NewStyle(),
		Prefix:      "[",
		Suffix:      "]",
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

		case tuix.KeySpace:
			// Check max length before inserting space
			if config.MaxLength == 0 || len(config.Value) < config.MaxLength {
				newValue := config.Value[:pos] + " " + config.Value[pos:]
				if config.OnChange != nil && config.ID != "" {
					config.OnChange(config.ID, newValue)
				}
				setPos(pos + 1)
			}

		case tuix.KeyBackspace:
			if pos > 0 && len(config.Value) > 0 {
				newValue := config.Value[:pos-1] + config.Value[pos:]
				if config.OnChange != nil && config.ID != "" {
					config.OnChange(config.ID, newValue)
				}
				setPos(pos - 1)
			}

		case tuix.KeyHome:
			setPos(0)

		case tuix.KeyEnd:
			setPos(len(config.Value))

		case tuix.KeyDelete:
			if pos < len(config.Value) {
				newValue := config.Value[:pos] + config.Value[pos+1:]
				if config.OnChange != nil && config.ID != "" {
					config.OnChange(config.ID, newValue)
				}
			}

		case tuix.KeyEnter:
			// Check min length on submit
			if config.MinLength > 0 && len(config.Value) < config.MinLength {
				// Optional: trigger validation error
				return renderError(config, "Minimum length is "+string(rune(config.MinLength))+" characters")
			}
			if config.OnSubmit != nil && config.ID != "" {
				config.OnSubmit(config.ID, config.Value)
			}

		default:
			// Insert printable character
			if key.Rune != 0 && key.Rune >= 32 && key.Rune <= 126 {
				// Check max length
				if config.MaxLength == 0 || len(config.Value) < config.MaxLength {
					newValue := config.Value[:pos] + string(key.Rune) + config.Value[pos:]
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, newValue)
					}
					setPos(pos + 1)
				}
			}
		}
	}

render:
	// Check if min length is met (for visual feedback)
	isValid := true
	if config.MinLength > 0 && len(config.Value) < config.MinLength {
		isValid = false
	}

	display := config.Value
	if display == "" && config.Placeholder != "" && !focused {
		display = config.Placeholder
	}

	// Style
	textStyle := config.Style
	if focused {
		textStyle = textStyle.Foreground(tuix.White).Background(tuix.Blue).
			Bold(true)
	} else {
		textStyle = textStyle.Foreground(tuix.BrightBlack)
	}

	// Border color changes based on validation
	borderColor := tuix.BrightBlack
	if focused {
		if isValid {
			borderColor = tuix.Cyan
		} else {
			borderColor = tuix.Red // Show red when invalid
		}
	}

	bracketStyle := tuix.NewStyle()
	if focused {
		bracketStyle = bracketStyle.Foreground(borderColor).Bold(true)
	} else {
		bracketStyle = bracketStyle.Foreground(tuix.BrightBlack)
	}

	cursorDisplay := display
	if focused {
		runes := []rune(display)
		if pos < len(runes) {
			cursorDisplay = string(runes[:pos]) + "█" + string(runes[pos:])
		} else {
			cursorDisplay = string(runes) + "█"
		}
	}

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

	// Add validation indicator
	textStyleForDisplay := textStyle
	if !isValid && !focused {
		textStyleForDisplay = textStyleForDisplay.Foreground(tuix.Red)
	}

	elements = append(elements, tuix.Text(paddedDisplay, textStyleForDisplay))

	if config.Suffix != "" {
		elements = append(elements, tuix.Text(config.Suffix, bracketStyle))
	}

	// Show length indicator
	if focused && config.MaxLength > 0 {
		lengthText := tuix.Text(
			" "+string(rune(len(config.Value)+48))+"/"+string(rune(config.MaxLength+48)),
			tuix.NewStyle().Foreground(tuix.BrightBlack),
		)
		return tuix.Box(
			tuix.Props{Direction: tuix.Row},
			tuix.NewStyle(),
			tuix.Box(
				tuix.Props{Direction: tuix.Row},
				tuix.NewStyle(),
				elements...,
			),
			lengthText,
		)
	}

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
		},
		tuix.NewStyle(),
		elements...,
	)
}

// Helper: Render error state
func renderError(config *InputConfig, msg string) tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Row},
		tuix.NewStyle(),
		tuix.Text(""+msg, tuix.NewStyle().Foreground(tuix.Red)),
	)
}
