package button

import (
	"fmt"
	"strings"
)

// ButtonData represents a button component's data
type ButtonData struct {
	ID         string
	Label      string
	Variant    string // "primary", "secondary", "danger", "success"
	Size       string // "small", "medium", "large"
	Disabled   bool
	Loading    bool
	ClickCount int
}

// DefaultButton creates a button with default values
func DefaultButton() ButtonData {
	return ButtonData{
		ID:         "demo-button",
		Label:      "Click me",
		Variant:    "primary",
		Size:       "medium",
		Disabled:   false,
		Loading:    false,
		ClickCount: 0,
	}
}

// GenerateHTML generates the HTML for a button
func (b ButtonData) GenerateHTML() string {
	var classes []string
	classes = append(classes, "btn")
	classes = append(classes, "btn-"+b.Variant)
	classes = append(classes, "btn-"+b.Size)

	if b.Disabled {
		classes = append(classes, "disabled")
	}
	if b.Loading {
		classes = append(classes, "loading")
	}

	classStr := strings.Join(classes, " ")

	disabledAttr := ""
	if b.Disabled {
		disabledAttr = "disabled"
	}

	loadingText := ""
	if b.Loading {
		loadingText = `<span class="spinner"></span>`
	}

	return fmt.Sprintf(`<button id="%s" class="%s" %s data-click-count="%d">
		%s%s
	</button>`, b.ID, classStr, disabledAttr, b.ClickCount, loadingText, b.Label)
}

// WithLabel creates a copy with new label
func (b ButtonData) WithLabel(label string) ButtonData {
	b.Label = label
	return b
}

// WithVariant creates a copy with new variant
func (b ButtonData) WithVariant(variant string) ButtonData {
	b.Variant = variant
	return b
}

// WithSize creates a copy with new size
func (b ButtonData) WithSize(size string) ButtonData {
	b.Size = size
	return b
}

// WithDisabled creates a copy with disabled state
func (b ButtonData) WithDisabled(disabled bool) ButtonData {
	b.Disabled = disabled
	return b
}

// WithLoading creates a copy with loading state
func (b ButtonData) WithLoading(loading bool) ButtonData {
	b.Loading = loading
	return b
}

// WithClickCount creates a copy with click count
func (b ButtonData) WithClickCount(count int) ButtonData {
	b.ClickCount = count
	return b
}
