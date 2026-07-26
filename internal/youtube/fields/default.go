package fields

// DefaultFields builds the list of default fields from the registry where Default == true.
func DefaultFields() []string {
	out := make([]string, 0)
	for name, def := range Registry {
		if def.Default {
			out = append(out, name)
		}
	}
	return out
}
