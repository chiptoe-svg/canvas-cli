package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeCriteriaFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "criteria.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestLoadRubricCriteria_BareArray(t *testing.T) {
	path := writeCriteriaFile(t, `[
	  {"description": "Thesis", "long_description": "Clear claim", "points": 10,
	   "ratings": [{"description": "Excellent", "points": 10}, {"description": "Weak", "points": 4}]},
	  {"description": "Evidence", "points": 5}
	]`)

	criteria, err := loadRubricCriteria(path)
	require.NoError(t, err)
	require.Len(t, criteria, 2)
	assert.Equal(t, "Thesis", criteria[0].Description)
	assert.Equal(t, "Clear claim", criteria[0].LongDescription)
	assert.Equal(t, 10.0, criteria[0].Points)
	require.Len(t, criteria[0].Ratings, 2)
	assert.Equal(t, "Weak", criteria[0].Ratings[1].Description)
	assert.Equal(t, 4.0, criteria[0].Ratings[1].Points)
	assert.Equal(t, "Evidence", criteria[1].Description)
}

func TestLoadRubricCriteria_WrappedObject(t *testing.T) {
	// The shape `canvas rubrics get -o json` prints, so a rubric can be
	// exported, edited, and sent back.
	for _, key := range []string{"criteria", "data"} {
		path := writeCriteriaFile(t, `{"id": 7, "title": "Essay", "`+key+`": [{"description": "Thesis", "points": 10}]}`)
		criteria, err := loadRubricCriteria(path)
		require.NoError(t, err, key)
		require.Len(t, criteria, 1, key)
		assert.Equal(t, "Thesis", criteria[0].Description, key)
	}
}

func TestLoadRubricCriteria_Errors(t *testing.T) {
	cases := map[string]string{
		"invalid json":               `{not json`,
		"empty array":                `[]`,
		"object without criteria":    `{"title": "Essay"}`,
		"criterion missing desc":     `[{"points": 5}]`,
		"negative criterion points":  `[{"description": "X", "points": -1}]`,
		"rating missing description": `[{"description": "X", "points": 5, "ratings": [{"points": 5}]}]`,
		"negative rating points":     `[{"description": "X", "points": 5, "ratings": [{"description": "R", "points": -2}]}]`,
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := loadRubricCriteria(writeCriteriaFile(t, content))
			assert.Error(t, err)
		})
	}

	t.Run("missing file", func(t *testing.T) {
		_, err := loadRubricCriteria(filepath.Join(t.TempDir(), "nope.json"))
		assert.Error(t, err)
	})
}
