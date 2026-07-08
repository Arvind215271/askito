package fields

// Planner represents the validated export request.
type Planner struct {
	fields   []string
	fieldSet map[string]bool
}
