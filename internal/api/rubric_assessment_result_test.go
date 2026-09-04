package api

import (
	"encoding/json"
	"testing"
)

// Canvas returns a submission's rubric_assessment as a hash keyed by
// criterion id; an array of entries must decode to the same map.
func TestRubricAssessmentResult_Unmarshal(t *testing.T) {
	var s Submission
	if err := json.Unmarshal([]byte(`{"id": 1, "rubric_assessment": {"_1": {"points": 8, "rating_id": "r1", "comments": "ok"}, "_2": {"points": 0}}}`), &s); err != nil {
		t.Fatal(err)
	}
	if len(s.Rubric) != 2 || s.Rubric["_1"].Points != 8 || s.Rubric["_1"].RatingID != "r1" || s.Rubric["_1"].CriterionID != "_1" {
		t.Errorf("hash form decoded wrong: %+v", s.Rubric)
	}
	var a Submission
	if err := json.Unmarshal([]byte(`{"id": 1, "rubric_assessment": [{"criterion_id": "_1", "points": 8}]}`), &a); err != nil {
		t.Fatal(err)
	}
	if len(a.Rubric) != 1 || a.Rubric["_1"].Points != 8 {
		t.Errorf("array form decoded wrong: %+v", a.Rubric)
	}
	var n Submission
	if err := json.Unmarshal([]byte(`{"id": 1, "rubric_assessment": null}`), &n); err != nil || n.Rubric != nil {
		t.Errorf("null: err=%v rubric=%v", err, n.Rubric)
	}
}
