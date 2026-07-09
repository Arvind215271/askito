package export

import (
	"net/http"
	"github.com/Arvind215271/askito/internal/api"
)

var (
	ErrInvalidFormat = api.NewError("invalid_format", "Unsupported export format", http.StatusBadRequest)
    ErrInvalidInputType = api.NewError("invalid_input_type", "Invalid input type", http.StatusBadRequest)
)
