package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

const criteriaFileUsage = `JSON file of criteria: either a bare array of
{"description", "long_description", "points", "ratings": [{"description", "long_description", "points"}]}
or the object 'canvas rubrics get -o json' prints (its "data" array is used)`

// loadRubricCriteria reads a rubric's criteria from a JSON file. It accepts a
// bare array of criteria, or an object carrying them under "criteria" or
// "data" — the latter being what `canvas rubrics get -o json` prints, so a
// rubric can be exported, edited, and sent back.
func loadRubricCriteria(path string) ([]api.RubricCriterion, error) {
	raw, err := os.ReadFile(path) // #nosec G304 -- path comes from --criteria-file, user-controlled by design
	if err != nil {
		return nil, fmt.Errorf("failed to read criteria file: %w", err)
	}

	var criteria []api.RubricCriterion
	if err := json.Unmarshal(raw, &criteria); err != nil {
		var wrapped struct {
			Criteria []api.RubricCriterion `json:"criteria"`
			Data     []api.RubricCriterion `json:"data"`
		}
		if err2 := json.Unmarshal(raw, &wrapped); err2 != nil {
			return nil, fmt.Errorf("invalid criteria file %s: expected a JSON array of criteria or an object with a \"criteria\" or \"data\" array: %w", path, err)
		}
		criteria = wrapped.Criteria
		if len(criteria) == 0 {
			criteria = wrapped.Data
		}
	}

	if len(criteria) == 0 {
		return nil, fmt.Errorf("criteria file %s contains no criteria", path)
	}

	for i, c := range criteria {
		if c.Description == "" {
			return nil, fmt.Errorf("criteria file %s: criterion %d has no description", path, i+1)
		}
		if c.Points < 0 {
			return nil, fmt.Errorf("criteria file %s: criterion %q has negative points", path, c.Description)
		}
		for j, r := range c.Ratings {
			if r.Description == "" {
				return nil, fmt.Errorf("criteria file %s: criterion %q rating %d has no description", path, c.Description, j+1)
			}
			if r.Points < 0 {
				return nil, fmt.Errorf("criteria file %s: criterion %q rating %q has negative points", path, c.Description, r.Description)
			}
		}
	}

	return criteria, nil
}
