package components

import (
	"strings"

	"github.com/subhasundardass/tuix/tuix"
)

type TextAreaOption func(*TextAreaConfig)

type TextAreaConfig struct {
	ID          string
	Value       string
	Placeholder string
	Width       int
	Height      int
	Style       tuix.Style
	Prefix      string
	Suffix      string
	OnChange    func(id string, value string)
	OnKeyPress  func(id string, key tuix.Key) bool
	OnFocus     func(id string)
	OnBlur      func(id string)
	OnSubmit    func(id string, value string)
}

func TextAreaWithID(id string) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.ID = id
	}
}

func TextAreaWithValue(value string) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.Value = value
	}
}

func TextAreaWithPlaceholder(text string) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.Placeholder = text
	}
}

func TextAreaWithWidth(width int) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.Width = width
	}
}

func TextAreaWithHeight(height int) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.Height = height
	}
}

func TextAreaWithStyle(style tuix.Style) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.Style = style
	}
}

func TextAreaWithPrefix(prefix string) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.Prefix = prefix
	}
}

func TextAreaWithSuffix(suffix string) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.Suffix = suffix
	}
}

func TextAreaWithOnChange(fn func(id string, value string)) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.OnChange = fn
	}
}

func TextAreaWithOnKeyPress(fn func(id string, key tuix.Key) bool) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.OnKeyPress = fn
	}
}

func TextAreaWithOnFocus(fn func(id string)) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.OnFocus = fn
	}
}

func TextAreaWithOnBlur(fn func(id string)) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.OnBlur = fn
	}
}

func TextAreaWithOnSubmit(fn func(id string, value string)) TextAreaOption {
	return func(c *TextAreaConfig) {
		c.OnSubmit = fn
	}
}

// ─── TextArea ──────────────────────────────────────────────────────────────

func TextArea(focused bool, opts ...TextAreaOption) tuix.Element {
	config := &TextAreaConfig{
		ID:          "",
		Value:       "",
		Placeholder: "",
		Width:       40,
		Height:      5,
		Style:       tuix.NewStyle(),
		Prefix:      "",
		Suffix:      "",
		OnChange:    nil,
		OnKeyPress:  nil,
		OnFocus:     nil,
		OnBlur:      nil,
		OnSubmit:    nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	//State for cursor position
	pos, setPos := tuix.UseState(len(config.Value))
	currentLine, setCurrentLine := tuix.UseState(0) //FIXED: renamed to currentLine

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
		case tuix.KeyUp:
			if currentLine > 0 {
				setCurrentLine(currentLine - 1)
			}
		case tuix.KeyDown:
			if currentLine < config.Height-1 {
				setCurrentLine(currentLine + 1)
			}
		case tuix.KeyBackspace:
			if pos > 0 && len(config.Value) > 0 {
				newValue := config.Value[:pos-1] + config.Value[pos:]
				config.Value = newValue
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
				config.Value = newValue
				if config.OnChange != nil && config.ID != "" {
					config.OnChange(config.ID, newValue)
				}
			}
		case tuix.KeyEnter:
			//Insert newline on Enter
			newValue := config.Value[:pos] + "\n" + config.Value[pos:]
			config.Value = newValue
			if config.OnChange != nil && config.ID != "" {
				config.OnChange(config.ID, newValue)
			}
			setPos(pos + 1)
			if config.OnSubmit != nil && config.ID != "" {
				config.OnSubmit(config.ID, newValue)
			}
		default:
			if key.Rune != 0 && key.Rune >= 32 {
				newValue := config.Value[:pos] + string(key.Rune) + config.Value[pos:]
				config.Value = newValue
				if config.OnChange != nil && config.ID != "" {
					config.OnChange(config.ID, newValue)
				}
				setPos(pos + 1)
			}
		}
	}

render:
	display := config.Value
	if display == "" && config.Placeholder != "" {
		display = config.Placeholder
	}

	//Wrap text to width
	lines := wrapTextArea(display, config.Width)

	//Pad to height
	for len(lines) < config.Height {
		lines = append(lines, "")
	}

	textStyle := config.Style
	if focused {
		textStyle = textStyle.Foreground(tuix.White)
	} else {
		textStyle = textStyle.Foreground(tuix.BrightBlack)
	}

	borderColor := tuix.BrightBlack
	if focused {
		borderColor = tuix.Cyan
	}

	bracketStyle := tuix.NewStyle()
	if focused {
		bracketStyle = bracketStyle.Foreground(borderColor).Bold(true)
	} else {
		bracketStyle = bracketStyle.Foreground(tuix.BrightBlack)
	}

	//Build content lines with cursor
	contentLines := []tuix.Element{}
	for lineIdx, lineContent := range lines {
		displayLine := lineContent
		//FIXED: Use currentLine instead of line
		if focused && lineIdx == currentLine {
			runes := []rune(displayLine)
			if pos < len(runes) {
				displayLine = string(runes[:pos]) + "█" + string(runes[pos:])
			} else {
				displayLine = string(runes) + "█"
			}
		}
		contentLines = append(contentLines, tuix.Text(displayLine, textStyle))
	}

	elements := []tuix.Element{}

	if config.Prefix != "" {
		elements = append(elements, tuix.Text(config.Prefix, bracketStyle))
	}

	//Text area content
	elements = append(elements, tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Width:     tuix.Fixed(config.Width + 2),
			Height:    tuix.Fixed(config.Height),
		},
		tuix.NewStyle().
			Border(tuix.Border{
				Top:    true,
				Right:  true,
				Bottom: true,
				Left:   true,
				Chars:  tuix.BorderRounded,
				Color:  borderColor,
			}),
		contentLines...,
	))

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

// wrapTextArea wraps text to max width
func wrapTextArea(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	// Split by newlines first
	paragraphs := strings.Split(text, "\n")
	var allLines []string

	for _, para := range paragraphs {
		if para == "" {
			allLines = append(allLines, "")
			continue
		}

		words := strings.Fields(para)
		if len(words) == 0 {
			allLines = append(allLines, "")
			continue
		}

		var line strings.Builder
		lineWidth := 0

		for _, word := range words {
			wordWidth := len([]rune(word))
			if lineWidth == 0 {
				line.WriteString(word)
				lineWidth = wordWidth
			} else if lineWidth+1+wordWidth <= maxWidth {
				line.WriteByte(' ')
				line.WriteString(word)
				lineWidth += 1 + wordWidth
			} else {
				allLines = append(allLines, line.String())
				line.Reset()
				line.WriteString(word)
				lineWidth = wordWidth
			}
		}

		if line.Len() > 0 {
			allLines = append(allLines, line.String())
		}
	}

	return allLines
}
