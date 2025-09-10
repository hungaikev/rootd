package models

// Field represents a single, universal field in a form.
// It contains all possible attributes for any field type. The 'type' property
// determines which attributes are used by the frontend renderer and backend validator.
type Field struct {
	// --- Core Properties (All Fields) ---
	ID          string `json:"id"`                    // UUID for the field, generated on creation.
	Type        string `json:"type"`                  // The registered type name (e.g., "text", "email", "rating").
	Label       string `json:"label"`                 // The question or title for the field.
	Required    bool   `json:"required,omitempty"`    // Is this field mandatory?
	ReadOnly    bool   `json:"readOnly,omitempty"`    // Is this field editable by the user?
	Placeholder string `json:"placeholder,omitempty"` // Placeholder text for input fields.
	Value       any    `json:"value,omitempty"`       // Default or static value for the field.

	// --- Logical Properties ---
	Conditional *Conditional `json:"conditional,omitempty"` // Rules to show/hide this field.
	Validation  *Validation  `json:"validation,omitempty"`  // Custom validation rules for this instance.

	// --- Options for Choice-Based Fields ---
	Options []Option `json:"options,omitempty"` // Used for select, radio, checkbox, rank.

	// --- Specific Configuration Properties ---
	// For Number / Slider
	Min      *float64 `json:"min,omitempty"`
	Max      *float64 `json:"max,omitempty"`
	Step     *float64 `json:"step,omitempty"`
	MinLabel string   `json:"minLabel,omitempty"` // e.g., "Dissatisfied"
	MaxLabel string   `json:"maxLabel,omitempty"` // e.g., "Very Satisfied"

	// For Textarea
	Rows int `json:"rows,omitempty"`

	// For Calculation
	Formula string `json:"formula,omitempty"` // e.g., "{field_id_1} + {field_id_2}"
	Prefix  string `json:"prefix,omitempty"`  // e.g., "KES", "$"

	// For Lookup
	DataSource *DataSource `json:"dataSource,omitempty"`

	// For Payment
	Provider string `json:"provider,omitempty"` // e.g., "stripe"
}

// Option represents a single choice for fields like dropdown, radio, or checkboxes.
type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Conditional defines a rule for when a field should be displayed.
type Conditional struct {
	FieldID  string `json:"fieldId"`  // The UUID of the field to check.
	Operator string `json:"operator"` // e.g., "==", "!=", "includes", ">=", "<="
	Value    any    `json:"value"`    // The value to compare against.
}

// Validation defines custom rules to apply to a field's input.
type Validation struct {
	MinLength           int    `json:"minLength,omitempty"`
	MaxLength           int    `json:"maxLength,omitempty"`
	Pattern             string `json:"pattern,omitempty"`             // Regex pattern
	PatternErrorMessage string `json:"patternErrorMessage,omitempty"` // Custom error for regex failure
}

// DataSource defines the source for a Lookup field.
type DataSource struct {
	Type       string `json:"type"`       // e.g., "api", "internal_list"
	Endpoint   string `json:"endpoint"`   // URL for the API endpoint
	ValueField string `json:"valueField"` // The field in the response to use as the value
	LabelField string `json:"labelField"` // The field in the response to use as the label
}
