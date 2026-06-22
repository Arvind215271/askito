package video

import (
	"encoding/json"
	"os"
	"path/filepath"
	"fmt"
)

func ExportVideoTestResults(
	suite VideoTestSuite,
) error {

	baseDir := filepath.Join(
		"testdata",
		"youtube",
		"transcript_tests",
		suite.VideoID,
	)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	for _, result := range suite.Results {
		path := filepath.Join(
			baseDir,
			result.Name+".txt",
		)

		if err := os.WriteFile(
			path,
			[]byte(result.Output),
			0644,
		); err != nil {
			return err
		}
	}

	metricsPath := filepath.Join(
		baseDir,
		"metrics.json",
	)

	payload, err := json.MarshalIndent(
		suite,
		"",
		"  ",
	)
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		metricsPath,
		payload,
		0644,
	); err != nil {
		return err
	}

	fmt.Printf(
		"[DEBUG] Transcript test results exported: %s\n",
		baseDir,
	)

	return nil
}
