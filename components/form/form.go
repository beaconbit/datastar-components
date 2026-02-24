package form

import (
	"fmt"
	"strings"
)

// FieldType represents the type of form field
type FieldType string

const (
	FieldText     FieldType = "text"
	FieldEmail    FieldType = "email"
	FieldPassword FieldType = "password"
	FieldNumber   FieldType = "number"
	FieldSelect   FieldType = "select"
	FieldTextarea FieldType = "textarea"
	FieldCheckbox FieldType = "checkbox"
)

// FormField represents a single form field
type FormField struct {
	ID          string
	Name        string
	Label       string
	Type        FieldType
	Value       string
	Placeholder string
	Required    bool
	Disabled    bool
	Options     []SelectOption // For select fields
	Error       string
}

// SelectOption represents an option for select fields
type SelectOption struct {
	Value string
	Label string
}

// FormData represents a form component's data
type FormData struct {
	ID        string
	Title     string
	Action    string
	Method    string // "GET", "POST"
	Fields    []FormField
	Submitted bool
	Errors    map[string]string
}

// DefaultForm creates a form with default values
func DefaultForm() FormData {
	return FormData{
		ID:     "demo-form",
		Title:  "Contact Form",
		Action: "/api/form/submit",
		Method: "POST",
		Fields: []FormField{
			{
				ID:          "name",
				Name:        "name",
				Label:       "Name",
				Type:        FieldText,
				Value:       "",
				Placeholder: "Enter your name",
				Required:    true,
			},
			{
				ID:          "email",
				Name:        "email",
				Label:       "Email",
				Type:        FieldEmail,
				Value:       "",
				Placeholder: "Enter your email",
				Required:    true,
			},
			{
				ID:          "message",
				Name:        "message",
				Label:       "Message",
				Type:        FieldTextarea,
				Value:       "",
				Placeholder: "Enter your message",
				Required:    true,
			},
			{
				ID:       "newsletter",
				Name:     "newsletter",
				Label:    "Subscribe to newsletter",
				Type:     FieldCheckbox,
				Value:    "true",
				Required: false,
			},
		},
		Submitted: false,
		Errors:    make(map[string]string),
	}
}

// GenerateHTML generates the HTML for a form
func (f FormData) GenerateHTML() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<form id="%s" class="form" action="%s" method="%s">`, f.ID, f.Action, f.Method))
	sb.WriteString(fmt.Sprintf(`<h3>%s</h3>`, f.Title))

	for _, field := range f.Fields {
		sb.WriteString(`<div class="form-field">`)
		sb.WriteString(fmt.Sprintf(`<label for="%s">%s</label>`, field.ID, field.Label))

		switch field.Type {
		case FieldText, FieldEmail, FieldPassword, FieldNumber:
			sb.WriteString(fmt.Sprintf(`<input type="%s" id="%s" name="%s" value="%s" placeholder="%s"`,
				string(field.Type), field.ID, field.Name, field.Value, field.Placeholder))
			if field.Required {
				sb.WriteString(` required`)
			}
			if field.Disabled {
				sb.WriteString(` disabled`)
			}
			sb.WriteString(`>`)

		case FieldTextarea:
			sb.WriteString(fmt.Sprintf(`<textarea id="%s" name="%s" placeholder="%s"`,
				field.ID, field.Name, field.Placeholder))
			if field.Required {
				sb.WriteString(` required`)
			}
			if field.Disabled {
				sb.WriteString(` disabled`)
			}
			sb.WriteString(`>`)
			sb.WriteString(field.Value)
			sb.WriteString(`</textarea>`)

		case FieldSelect:
			sb.WriteString(fmt.Sprintf(`<select id="%s" name="%s"`, field.ID, field.Name))
			if field.Required {
				sb.WriteString(` required`)
			}
			if field.Disabled {
				sb.WriteString(` disabled`)
			}
			sb.WriteString(`>`)
			for _, option := range field.Options {
				sb.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, option.Value, option.Label))
			}
			sb.WriteString(`</select>`)

		case FieldCheckbox:
			sb.WriteString(fmt.Sprintf(`<input type="checkbox" id="%s" name="%s" value="%s"`,
				field.ID, field.Name, field.Value))
			if field.Required {
				sb.WriteString(` required`)
			}
			if field.Disabled {
				sb.WriteString(` disabled`)
			}
			if field.Value == "true" {
				sb.WriteString(` checked`)
			}
			sb.WriteString(`>`)
		}

		if field.Error != "" {
			sb.WriteString(fmt.Sprintf(`<div class="error">%s</div>`, field.Error))
		}

		sb.WriteString(`</div>`)
	}

	sb.WriteString(`<div class="form-actions">`)
	sb.WriteString(`<button type="submit" class="btn btn-primary">Submit</button>`)
	sb.WriteString(`<button type="reset" class="btn btn-secondary">Reset</button>`)
	sb.WriteString(`</div>`)

	if f.Submitted {
		sb.WriteString(`<div class="success">Form submitted successfully!</div>`)
	}

	sb.WriteString(`</form>`)
	return sb.String()
}

// WithFieldValue updates a field's value
func (f FormData) WithFieldValue(fieldID, value string) FormData {
	for i, field := range f.Fields {
		if field.ID == fieldID {
			f.Fields[i].Value = value
			break
		}
	}
	return f
}

// WithSubmitted sets the submitted state
func (f FormData) WithSubmitted(submitted bool) FormData {
	f.Submitted = submitted
	return f
}

// WithError adds an error to a field
func (f FormData) WithError(fieldID, error string) FormData {
	if f.Errors == nil {
		f.Errors = make(map[string]string)
	}
	f.Errors[fieldID] = error

	// Also update the field's error
	for i, field := range f.Fields {
		if field.ID == fieldID {
			f.Fields[i].Error = error
			break
		}
	}

	return f
}

// ClearErrors clears all errors
func (f FormData) ClearErrors() FormData {
	f.Errors = make(map[string]string)
	for i := range f.Fields {
		f.Fields[i].Error = ""
	}
	return f
}
