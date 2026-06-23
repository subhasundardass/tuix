package components

import (
	"fmt"

	"github.com/subhasundardass/tuix/tuix"
)

// ⭐ Field types
const (
	FieldTypeText     = "text"
	FieldTypeCheckbox = "checkbox"
	FieldTypeSelect   = "select"
	FieldTypeButton   = "button"
	FieldTypeCustom   = "custom"
)

// Field represents a single form field
type Field struct {
	ID       string
	Label    string
	Type     string // "text", "checkbox", "select", "button", "custom"
	Value    string
	Options  []string
	OnChange func(string)
	OnToggle func(bool)
	OnSubmit func()

	// ⭐ For custom fields - you can pass any component
	CustomRender func(focused bool, data map[string]string, setData func(map[string]string)) tuix.Element
}

// FormProps defines form properties
type FormProps struct {
	Fields     []Field
	OnSubmit   func(map[string]string)
	OnValidate func(map[string]string) map[string]string
	Width      int
	Height     int
}

// Form manages a complete form with auto-focus and navigation
func Form(props FormProps) tuix.Element {
	// ⭐ Focus state
	focusIndex, setFocusIndex := tuix.UseState(0)
	totalFields := len(props.Fields)

	// ⭐ Form data
	formData, setFormData := tuix.UseState(make(map[string]string))
	errors, setErrors := tuix.UseState(make(map[string]string))

	// ⭐ Handle keyboard navigation
	switch tuix.CurrentKey.Code {
	case tuix.KeyDown:
		setFocusIndex((focusIndex + 1) % totalFields)

	case tuix.KeyUp:
		setFocusIndex((focusIndex - 1 + totalFields) % totalFields)

	case tuix.KeyEnter:
		currentField := props.Fields[focusIndex]
		if currentField.Type == FieldTypeButton {
			validateAndSubmit(props, formData, setErrors)
		} else {
			setFocusIndex((focusIndex + 1) % totalFields)
		}
	}

	// ⭐ Build fields with proper focus
	fieldElements := []tuix.Element{}
	for i, field := range props.Fields {
		isFocused := focusIndex == i
		fieldElements = append(fieldElements, renderField(field, isFocused, formData, setFormData))
	}

	// ⭐ Error summary
	errorSummary := renderErrors(errors)

	// ⭐ Navigation help
	help := tuix.Text(
		fmt.Sprintf("  ↓: Next  |  ↑: Previous  |  Enter: Next/Submit  (%d/%d)",
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

		// All fields
		tuix.Box(
			tuix.Props{
				Direction: tuix.Column,
				Gap:       1,
			},
			tuix.NewStyle(),
			fieldElements...,
		),

		// Error summary
		errorSummary,

		// Navigation help
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

	// case FieldTypeSelect:
	// 	return renderSelectField(field, focused, formData, setFormData)

	case FieldTypeButton:
		return renderButtonField(field, focused)

	case FieldTypeCustom:
		// ⭐ Custom field - render whatever you want
		if field.CustomRender != nil {
			return field.CustomRender(focused, formData, setFormData)
		}
		return tuix.Text("Custom field: "+field.ID, tuix.NewStyle().Foreground(tuix.Red))

	default:
		return tuix.Text("Unknown field: "+field.Type, tuix.NewStyle().Foreground(tuix.Red))
	}
}

// ⭐ Built-in field renderers
func renderTextField(field Field, focused bool, value string, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	input := Input(
		field.Label+":",
		focused,
		value,
		func(newValue string) {
			formData[field.ID] = newValue
			setFormData(formData)
			if field.OnChange != nil {
				field.OnChange(newValue)
			}
		},
	)

	if err, ok := formData["error_"+field.ID]; ok && err != "" {
		return tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0},
			tuix.NewStyle(),
			input,
			tuix.Text("  ✗ "+err,
				tuix.NewStyle().Foreground(tuix.Red).Italic(true)),
		)
	}

	return input
}

func renderCheckboxField(field Field, focused bool, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
	return Checkbox(
		field.Label,
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
}

// func renderSelectField(field Field, focused bool, formData map[string]string, setFormData func(map[string]string)) tuix.Element {
// 	selected := formData[field.ID]
// 	if selected == "" && len(field.Options) > 0 {
// 		selected = field.Options[0]
// 		formData[field.ID] = selected
// 		setFormData(formData)
// 	}

// 	return SelectPicker(
// 		field.Options,
// 		focused,
// 		func(value string) {
// 			formData[field.ID] = value
// 			setFormData(formData)
// 			if field.OnChange != nil {
// 				field.OnChange(value)
// 			}
// 		},
// 	)
// }

func renderButtonField(field Field, focused bool) tuix.Element {
	return Button(field.Label, focused)
}

// ⭐ Validation and submission
func validateAndSubmit(props FormProps, formData map[string]string, setErrors func(map[string]string)) {
	// ⭐ Validate if validator provided
	if props.OnValidate != nil {
		errors := props.OnValidate(formData)
		setErrors(errors)

		if len(errors) > 0 {
			return
		}
	}

	// ⭐ Clear errors
	setErrors(make(map[string]string))

	// ⭐ Submit
	if props.OnSubmit != nil {
		props.OnSubmit(formData)
	}
}

// ⭐ Error rendering
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
