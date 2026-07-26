package fields

// FieldDefinition describes a valid field in the registry.
type FieldDefinition struct {
	Name    string
	Group   string
	Default bool
}

// Planner represents the validated export request.
type Planner struct {
	fields   []string
	fieldSet map[string]bool
}
