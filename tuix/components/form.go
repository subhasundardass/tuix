package components

import (
	"fmt"

	"github.com/subhasundardass/tuix/tuix"
)

// Field types
const (
	FieldTypeText     = "text"
	FieldTypeCheckbox = "checkbox"
	FieldTypeSelect   = "select"
	FieldTypeButton   = "button"
	FieldTypeCustom   = "custom"
	FieldTypeNumber   = "number"
	FieldTypeDate     = "date"
)

// Field represents a single form field
type Field struct {
	ID       string
	Label    string
	Type     string
	Value    string
	Options  []string
	OnChange func(string)
	OnToggle func(bool)
	OnSubmit func()

	// Number Input specific
	Decimal     int
	Min         *float64
	Max         *float64
	Step        float64
	Placeholder string
	Width       int

	// Date Input specific
	DateMask        string
	DatePlaceholder string

	// For custom fields
	CustomRender func(focused bool, data map[string]string, setData func(map[string]string)) tuix.Element
}

// FormProps defines form properties
type FormProps struct {
	Fields     []Field
	OnSubmit   func(map[string]string)
	OnValidate func(map[string]string) map[string]string
	Width      int
	Height     int
	RenderFunc func(focusedIndex int, setFocused func(int), formData map[string]string, setFormData func(map[string]string)) tuix.Element
}

// Form manages a complete form with auto-focus and navigation
func Form(props FormProps) tuix.Element {
	// Focus state - managed HERE
	focusIndex, setFocusIndex := tuix.UseState(0)
	totalFields := len(props.Fields)

	// Form data
	formData, setFormData := tuix.UseState(make(map[string]string))
	errors, setErrors := tuix.UseState(make(map[string]string))

	// Get current field
	currentField := props.Fields[focusIndex]

	// Handle keyboard navigation - managed HERE
	switch tuix.CurrentKey.Code {
	case tuix.KeyDown:
		setFocusIndex((focusIndex + 1) % totalFields)

	case tuix.KeyUp:
		setFocusIndex((focusIndex - 1 + totalFields) % totalFields)

	case tuix.KeyEnter:
		if currentField.Type == FieldTypeButton {
			validateAndSubmit(props, formData, setErrors)
		} else {
			setFocusIndex((focusIndex + 1) % totalFields)
		}

	case tuix.KeySpace:
		if currentField.Type != FieldTypeCheckbox {
			setFocusIndex((focusIndex + 1) % totalFields)
		}
	}

	// If custom render function provided, use it
	if props.RenderFunc != nil {
		return props.RenderFunc(focusIndex, setFocusIndex, formData, setFormData)
	}

	// Default rendering (automatic)
	fieldElements := []tuix.Element{}
	for i, field := range props.Fields {
		isFocused := focusIndex == i
		fieldElements = append(fieldElements, renderField(field, isFocused, formData, setFormData))
	}

	// Error summary
	errorSummary := renderErrors(errors)

	// Navigation help
	help := tuix.Text(
		fmt.Sprintf("  ↓: Next  |  ↑: Previous  |  Enter: Next/Submit  |  Space: Toggle  (%d/%d)",
			focusIndex+1, totalFields),
		tuix.NewStyle().Foreground(tuix.BrightBlack),
	)

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{2, 4, 2, 4},
			Width:     tuix.Fixed(props.Width),
		},
		tuix.NewStyle().Background(tuix.Black),

		tuix.Box(
			tuix.Props{
				Direction: tuix.Column,
				Gap:       1,
			},
			tuix.NewStyle(),
			fieldElements...,
		),

		errorSummary,
		help,
	)
}

// renderField renders a single form field
func renderField(field Field, focused bool, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	value := formData[field.ID]

	switch field.Type {
	case FieldTypeText:
		return renderTextField(field, focused, value, formData, setFormData)

	case FieldTypeCheckbox:
		return renderCheckboxField(field, focused, formData, setFormData)

	case FieldTypeButton:
		return renderButtonField(field, focused)

	case FieldTypeNumber:
		return renderNumberField(field, focused, value, formData, setFormData)

	case FieldTypeDate:
		return renderDateField(field, focused, value, formData, setFormData)

	case FieldTypeCustom:
		if field.CustomRender != nil {
			return field.CustomRender(focused, formData, setFormData)
		}
		return tuix.Text("Custom field: "+field.ID, tuix.NewStyle().Foreground(tuix.Red))

	default:
		return tuix.Text("Unknown field: "+field.Type, tuix.NewStyle().Foreground(tuix.Red))
	}
}

func renderTextField(field Field, focused bool, value string, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	input := TextInput(
		focused,
		WithID(field.ID),
		WithValue(value),
		WithWidth(field.Width),
		WithPlaceholder(field.Placeholder),
		WithPrefix("["),
		WithSuffix("]"),
		WithOnChange(func(id, v string) {
			formData[field.ID] = v
			setFormData(formData)
			if field.OnChange != nil {
				field.OnChange(v)
			}
		}),
	)

	if err, ok := formData["error_"+field.ID]; ok && err != "" {
		return tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0},
			tuix.NewStyle(),
			tuix.Text(field.Label+": ", tuix.NewStyle().Foreground(tuix.White)),
			input,
			tuix.Text("  ✗ "+err,
				tuix.NewStyle().Foreground(tuix.Red).Italic(true)),
		)
	}

	return tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1},
		tuix.NewStyle(),
		tuix.Text(field.Label+": ", tuix.NewStyle().Foreground(tuix.White)),
		input,
	)
}

func renderCheckboxField(field Field, focused bool, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	checkbox := Checkbox(
		focused,
		func(value bool) {
			if value {
				formData[field.ID] = "true"
			} else {
				formData[field.ID] = "false"
			}
			setFormData(formData)
			if field.OnToggle != nil {
				field.OnToggle(value)
			}
		},
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1},
		tuix.NewStyle(),
		checkbox,
		tuix.Text(" "+field.Label, tuix.NewStyle().Foreground(tuix.White)),
	)
}

func renderButtonField(field Field, focused bool) tuix.Element {
	return Button(
		focused,
		WithButtonID(field.ID),
		WithLabel(field.Label),
		WithOnPress(func(id string) {
			if field.OnSubmit != nil {
				// field.OnPress(id)
			}
		}),
	)
}

func renderNumberField(field Field, focused bool, value string, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	input := NumberInput(
		focused,
		NumberWithID(field.ID),
		NumberWithValue(value),
		NumberWithWidth(field.Width),
		NumberWithPlaceholder(field.Placeholder),
		NumberWithPrefix("["),
		NumberWithSuffix("]"),
		NumberWithDecimal(field.Decimal),
		NumberWithMin(*field.Min),
		NumberWithMax(*field.Max),
		NumberWithOnChange(func(id, v string) {
			formData[field.ID] = v
			setFormData(formData)
			if field.OnChange != nil {
				field.OnChange(v)
			}
		}),
	)

	if err, ok := formData["error_"+field.ID]; ok && err != "" {
		return tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0},
			tuix.NewStyle(),
			tuix.Text(field.Label+": ", tuix.NewStyle().Foreground(tuix.White)),
			input,
			tuix.Text("  ✗ "+err,
				tuix.NewStyle().Foreground(tuix.Red).Italic(true)),
		)
	}

	return tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1},
		tuix.NewStyle(),
		tuix.Text(field.Label+": ", tuix.NewStyle().Foreground(tuix.White)),
		input,
	)
}

func renderDateField(field Field, focused bool, value string, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	mask := field.DateMask
	if mask == "" {
		mask = "YYYY-MM-DD"
	}

	placeholder := field.DatePlaceholder
	if placeholder == "" {
		placeholder = mask
	}

	// Updated DateInput with options pattern
	input := DateInput(
		focused,
		DateWithID(field.ID),
		DateWithValue(value),
		DateWithFormat(mask),
		DateWithPlaceholder(placeholder),
		DateWithWidth(field.Width),
		DateWithPrefix("["),
		DateWithSuffix("]"),
		DateWithOnChange(func(id, v string) {
			formData[field.ID] = v
			setFormData(formData)
			if field.OnChange != nil {
				field.OnChange(v)
			}
		}),
		DateWithOnSubmit(func(id, v string) {
			// Handle submit if needed
		}),
	)

	if err, ok := formData["error_"+field.ID]; ok && err != "" {
		return tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0},
			tuix.NewStyle(),
			tuix.Text(field.Label+": ", tuix.NewStyle().Foreground(tuix.White)),
			input,
			tuix.Text("  ✗ "+err,
				tuix.NewStyle().Foreground(tuix.Red).Italic(true)),
		)
	}

	return tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1},
		tuix.NewStyle(),
		tuix.Text(field.Label+": ", tuix.NewStyle().Foreground(tuix.White)),
		input,
	)
}

func validateAndSubmit(props FormProps, formData map[string]string, setErrors func(map[string]string)) {
	if props.OnValidate != nil {
		errors := props.OnValidate(formData)
		setErrors(errors)
		if len(errors) > 0 {
			return
		}
	}

	setErrors(make(map[string]string))

	if props.OnSubmit != nil {
		props.OnSubmit(formData)
	}
}

func renderErrors(errors map[string]string) tuix.Element {
	if len(errors) == 0 {
		return tuix.Text("", tuix.NewStyle())
	}

	errElements := []tuix.Element{
		tuix.Text("⚠️ Please fix the following errors:",
			tuix.NewStyle().Foreground(tuix.Red).Bold(true)),
	}

	for field, msg := range errors {
		errElements = append(errElements,
			tuix.Text("  • "+field+": "+msg,
				tuix.NewStyle().Foreground(tuix.Yellow)),
		)
	}

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 0},
		tuix.NewStyle().
			Background(tuix.Black),
		errElements...,
	)
}
