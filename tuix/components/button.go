package components

import "github.com/subhasundardass/tuix/tuix"

type ButtonConfig struct {
	ID       string
	Label    string
	Focused  bool
	Disabled bool
	Style    tuix.Style
	OnPress  func(id string)
}

type ButtonOption func(*ButtonConfig)

func defaultButtonConfig() ButtonConfig {
	return ButtonConfig{
		Style: tuix.NewStyle(),
	}
}

func WithButtonID(id string) ButtonOption {
	return func(c *ButtonConfig) {
		c.ID = id
	}
}

func WithLabel(label string) ButtonOption {
	return func(c *ButtonConfig) {
		c.Label = label
	}
}

func WithDisabled(disabled bool) ButtonOption {
	return func(c *ButtonConfig) {
		c.Disabled = disabled
	}
}

func WithButtonStyle(style tuix.Style) ButtonOption {
	return func(c *ButtonConfig) {
		c.Style = style
	}
}

func WithOnPress(fn func(id string)) ButtonOption {
	return func(c *ButtonConfig) {
		c.OnPress = fn
	}
}

func Button(focused bool, opts ...ButtonOption) tuix.Element {
	cfg := defaultButtonConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	cfg.Focused = focused

	style := cfg.Style

	if cfg.Disabled {
		style = style.Foreground(tuix.Cyan)
	} else if cfg.Focused {
		style = style.
			Foreground(tuix.Black).
			Background(tuix.Cyan).
			Bold(true)
	} else {
		style = style.Foreground(tuix.White)
	}

	if cfg.Focused &&
		tuix.CurrentKey.Code == tuix.KeyEnter &&
		cfg.OnPress != nil {
		cfg.OnPress(cfg.ID)
	}

	return tuix.Text("[ "+cfg.Label+" ]", style)
}
