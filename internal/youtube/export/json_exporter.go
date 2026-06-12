package export

import (
	"encoding/json"
	"github.com/Arvind215271/askito/internal/youtube"

)	

type JSONExporter struct{}

func NewJSONExporter() JSONExporter {
    return JSONExporter{}
}

func (e JSONExporter) Export(data ExportData, ) ([]byte, error) {

	resp := ExportResponse{
		SchemaVersion: SchemaVersion,
		Format:        FormatJSON,
		Data:          data,
	}

	b, err := json.MarshalIndent(resp, "", "  ")

	if err != nil {
		return nil, youtube.Err.Export.MarshalFailed().Wrap(err)
	}

	return b, nil
}

