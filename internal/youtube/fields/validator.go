package fields

import (
	"github.com/Arvind215271/askito/internal/api"
	"github.com/Arvind215271/askito/internal/youtube"
)

// ValidateFields ensures all requested fields are valid
func ValidateFields(fields []string) *api.AppError {
	err := youtube.Err.Export.InvalidField()
	var foundInvalid bool
	for _, f := range fields {
		if _, ok := ValidFields[f]; !ok {
			err.AddField(f, "Invalid field")
			foundInvalid = true
		}
	}
	if foundInvalid {
		return err
	}
	return nil
}
