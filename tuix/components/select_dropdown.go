package components

import (
	"strings"

	"github.com/subhasundardass/tuix/tuix"
)

// ─────────────────────────────────────────────────────────────────────────────
// SelectDropdown - A basic dropdown select with overlay
// ─────────────────────────────────────────────────────────────────────────────

type SelectOption struct {
	Label string
	Value string
}

type SelectOptionFunc func(*SelectConfig)

type SelectConfig struct {
	ID         string
	Label      string
	Options    []SelectOption
	Selected   string
	Width      int
	Height     int // Max visible options before scrolling
	Style      tuix.Style
	OnChange   func(id string, value string)
	OnKeyPress func(id string, key tuix.Key) bool
	OnFocus    func(id string)
	OnBlur     func(id string)
}

func SelectWithID(id string) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.ID = id
	}
}

func SelectWithLabel(label string) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.Label = label
	}
}

func SelectWithOptions(options []SelectOption) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.Options = options
	}
}

func SelectWithSelected(value string) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.Selected = value
	}
}

func SelectWithWidth(width int) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.Width = width
	}
}

func SelectWithHeight(height int) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.Height = height
	}
}

func SelectWithStyle(style tuix.Style) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.Style = style
	}
}

func SelectWithOnChange(fn func(id string, value string)) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.OnChange = fn
	}
}

func SelectWithOnKeyPress(fn func(id string, key tuix.Key) bool) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.OnKeyPress = fn
	}
}

func SelectWithOnFocus(fn func(id string)) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.OnFocus = fn
	}
}

func SelectWithOnBlur(fn func(id string)) SelectOptionFunc {
	return func(c *SelectConfig) {
		c.OnBlur = fn
	}
}

// ─── SelectDropdown ─────────────────────────────────────────────────────────

func SelectDropdown(focused bool, opts ...SelectOptionFunc) tuix.Element {
	config := &SelectConfig{
		ID:       "",
		Label:    "",
		Options:  []SelectOption{},
		Selected: "",
		Width:    30,
		Height:   5,
		Style:    tuix.NewStyle(),
		OnChange: nil,
		OnFocus:  nil,
		OnBlur:   nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	// ⭐ State
	isOpen, setIsOpen := tuix.UseState(false)
	highlightedIndex, setHighlightedIndex := tuix.UseState(0)

	// Find selected index
	selectedIndex := 0
	for i, opt := range config.Options {
		if opt.Value == config.Selected {
			selectedIndex = i
			break
		}
	}

	// ⭐ Get selected label
	selectedLabel := ""
	for _, opt := range config.Options {
		if opt.Value == config.Selected {
			selectedLabel = opt.Label
			break
		}
	}
	if selectedLabel == "" && len(config.Options) > 0 {
		selectedLabel = config.Options[0].Label
	}

	if focused && config.OnFocus != nil && config.ID != "" {
		config.OnFocus(config.ID)
	}

	if !focused && config.OnBlur != nil && config.ID != "" {
		config.OnBlur(config.ID)
	}

	// ⭐ Handle keys
	if focused {
		key := tuix.CurrentKey

		if config.OnKeyPress != nil && config.ID != "" {
			if config.OnKeyPress(config.ID, key) {
				goto render
			}
		}

		switch key.Code {
		case tuix.KeyEnter, tuix.KeySpace:
			if !isOpen {
				setIsOpen(true)
				setHighlightedIndex(selectedIndex)
			} else {
				// Select highlighted option
				if len(config.Options) > 0 && highlightedIndex < len(config.Options) {
					selected := config.Options[highlightedIndex].Value
					config.Selected = selected
					if config.OnChange != nil && config.ID != "" {
						config.OnChange(config.ID, selected)
					}
				}
				setIsOpen(false)
			}

		case tuix.KeyEscape:
			if isOpen {
				setIsOpen(false)
			}

		case tuix.KeyDown:
			if isOpen && highlightedIndex < len(config.Options)-1 {
				setHighlightedIndex(highlightedIndex + 1)
			}

		case tuix.KeyUp:
			if isOpen && highlightedIndex > 0 {
				setHighlightedIndex(highlightedIndex - 1)
			}

		case tuix.KeyTab:
			if isOpen {
				setIsOpen(false)
			}
		}
	}

render:
	// ⭐ Build display
	displayText := selectedLabel
	if displayText == "" && len(config.Options) > 0 {
		displayText = config.Options[0].Label
	}
	if displayText == "" {
		displayText = "Select..."
	}

	// ⭐ Styles
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

	// ⭐ Main field
	arrow := "▼"
	if isOpen {
		arrow = "▲"
	}

	// Pad display to width
	paddedDisplay := displayText
	displayLen := len([]rune(paddedDisplay))
	if displayLen < config.Width-4 {
		padding := strings.Repeat(" ", config.Width-4-displayLen)
		paddedDisplay = paddedDisplay + padding
	}

	mainField := tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Width:     tuix.Fixed(config.Width),
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
		tuix.Text(" "+paddedDisplay+" ", textStyle),
		tuix.Text(" "+arrow+" ", bracketStyle),
	)

	// ⭐ Label
	labelElement := tuix.Text("", tuix.NewStyle())
	if config.Label != "" {
		labelElement = tuix.Text(config.Label+":", tuix.NewStyle().Foreground(tuix.White))
	}

	// ⭐ Dropdown options (overlay)
	dropdownOverlay := tuix.Text("", tuix.NewStyle())
	if isOpen && len(config.Options) > 0 {
		optionElements := []tuix.Element{}

		// Determine visible range (simple scroll)
		startIdx := 0
		endIdx := len(config.Options)
		if len(config.Options) > config.Height {
			if highlightedIndex >= config.Height {
				startIdx = highlightedIndex - config.Height + 1
			}
			endIdx = startIdx + config.Height
			if endIdx > len(config.Options) {
				endIdx = len(config.Options)
				startIdx = endIdx - config.Height
			}
		}

		for i := startIdx; i < endIdx; i++ {
			opt := config.Options[i]
			style := tuix.NewStyle()
			if i == highlightedIndex {
				style = style.Background(tuix.Blue).Foreground(tuix.White).Bold(true)
			} else {
				style = style.Foreground(tuix.White)
			}
			prefix := "  "
			if i == highlightedIndex {
				prefix = "▶ "
			}
			optionElements = append(optionElements, tuix.Text(prefix+opt.Label, style))
		}

		dropdownOverlay = tuix.Overlay(10, 5,
			tuix.Box(
				tuix.Props{
					Direction: tuix.Column,
					Width:     tuix.Fixed(config.Width + 2),
				},
				tuix.NewStyle().
					Background(tuix.Black).
					Border(tuix.Border{
						Top:    true,
						Right:  true,
						Bottom: true,
						Left:   true,
						Chars:  tuix.BorderRounded,
						Color:  tuix.Cyan,
					}),
				optionElements...,
			),
		)
	}

	// ⭐ Combine
	elements := []tuix.Element{}

	if config.Label != "" {
		elements = append(elements,
			tuix.Box(
				tuix.Props{Direction: tuix.Row, Gap: 1},
				tuix.NewStyle(),
				labelElement,
				mainField,
			),
		)
	} else {
		elements = append(elements, mainField)
	}

	// Add dropdown overlay (renders on top)
	elements = append(elements, dropdownOverlay)

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
		},
		tuix.NewStyle(),
		elements...,
	)
}
