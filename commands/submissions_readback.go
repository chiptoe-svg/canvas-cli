package commands

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// submissionReadBack is the evidence that a grade or comment write took:
// the submission as it was before the write, as it reads back afterwards,
// what was asked for, and whether the read-back matches the request.
type submissionReadBack struct {
	Before    *api.Submission        `json:"before"` // nil when the pre-read failed
	After     *api.Submission        `json:"after"`
	Requested gradeRequest           `json:"requested"`
	Comment   *api.SubmissionComment `json:"comment,omitempty"` // the comment found on read-back, when one was requested
	// CommentExisted is true when the requested comment was already on the
	// submission before the write, so it was not posted again.
	CommentExisted bool     `json:"comment_existed,omitempty"`
	Verified       bool     `json:"verified"`
	Mismatches     []string `json:"mismatches,omitempty"`
}

// gradeRequest is the part of a grade write that can be checked on read-back.
type gradeRequest struct {
	PostedGrade string                            `json:"posted_grade,omitempty"`
	Comment     string                            `json:"comment,omitempty"`
	Excuse      bool                              `json:"excuse,omitempty"`
	Rubric      map[string]rubricCriterionRequest `json:"rubric,omitempty"` // by criterion id
}

// rubricCriterionRequest is what was asked for on one rubric criterion.
type rubricCriterionRequest struct {
	Points   *float64 `json:"points,omitempty"`
	Rating   string   `json:"rating_id,omitempty"`
	Comments string   `json:"comments,omitempty"`
}

func gradeRequestFromParams(params *api.GradeSubmissionParams) gradeRequest {
	r := gradeRequest{PostedGrade: params.PostedGrade, Excuse: params.Excuse}
	if params.Comment != nil {
		r.Comment = params.Comment.TextComment
	}
	if len(params.RubricAssessment) > 0 {
		r.Rubric = make(map[string]rubricCriterionRequest, len(params.RubricAssessment))
		for id, c := range params.RubricAssessment {
			r.Rubric[id] = rubricCriterionRequest{Points: c.Points, Rating: c.Rating, Comments: c.Comments}
		}
	}
	return r
}

// findExistingComment returns a comment on s with exactly the given text.
func findExistingComment(s *api.Submission, text string) *api.SubmissionComment {
	if s == nil || text == "" {
		return nil
	}
	for i := range s.SubmissionComments {
		if s.SubmissionComments[i].Comment == text {
			return &s.SubmissionComments[i]
		}
	}
	return nil
}

// writesNothing reports whether params, after any comment was dropped, has
// nothing left to send.
func writesNothing(p *api.GradeSubmissionParams) bool {
	return p.PostedGrade == "" && !p.Excuse && p.LatePolicyStatus == "" && p.SecondsLateOverride == nil &&
		p.Comment == nil && len(p.RubricAssessment) == 0
}

// readBackScoreTolerance absorbs Canvas rounding grades to two decimals.
const readBackScoreTolerance = 0.005

// gradeAndReadBack performs the write and evidences it: the submission is
// read before (best effort) and after (required, live) the PUT, and the
// read-back is compared with what was requested. The returned error is
// the write's or the read-back's failure; a mismatch is reported in the
// result, not as an error, so callers can print it first.
//
// With skipExistingComment, a comment whose exact text is already on the
// submission is reported as present and not posted again; the grade still
// is. Bulk grading uses this so a CSV re-run after a partial failure does
// not repeat feedback. A single `grade` or `add-comment` posts what it was
// told to.
func gradeAndReadBack(ctx context.Context, svc *api.SubmissionsService, courseID, assignmentID, userID int64, params *api.GradeSubmissionParams, skipExistingComment bool) (*submissionReadBack, error) {
	include := []string{"submission_comments", "rubric_assessment"}
	rb := &submissionReadBack{Requested: gradeRequestFromParams(params)}
	if before, err := svc.Get(ctx, courseID, assignmentID, userID, include); err == nil {
		rb.Before = before
	}

	if skipExistingComment && params.Comment != nil {
		if existing := findExistingComment(rb.Before, params.Comment.TextComment); existing != nil {
			rb.Comment, rb.CommentExisted = existing, true
			rb.Requested.Comment = ""
			p := *params
			p.Comment = nil
			params = &p
			if writesNothing(params) {
				rb.After = rb.Before
				rb.Verified = true
				return rb, nil
			}
		}
	}

	if _, err := svc.Grade(ctx, courseID, assignmentID, userID, params); err != nil {
		return rb, err
	}

	after, err := svc.Get(ctx, courseID, assignmentID, userID, include)
	if err != nil {
		return rb, fmt.Errorf("written, but the read-back failed: %w", err)
	}
	rb.After = after
	comment, mismatches := verifyGradeReadBack(rb.Before, after, rb.Requested)
	if !rb.CommentExisted {
		rb.Comment = comment
	}
	rb.Mismatches = mismatches
	rb.Verified = len(rb.Mismatches) == 0
	return rb, nil
}

// verifyGradeReadBack compares the read-back with the request: a numeric
// posted grade against the entered score (the score before any late
// policy), a letter/pass grade against the entered grade, an excuse
// against excused, and a comment against the submission's comments —
// preferring one that was not there before. Every difference is returned.
func verifyGradeReadBack(before, after *api.Submission, req gradeRequest) (*api.SubmissionComment, []string) {
	var mismatches []string

	if req.PostedGrade != "" {
		if want, err := strconv.ParseFloat(req.PostedGrade, 64); err == nil {
			got := after.Score
			if after.EnteredScore != 0 {
				got = after.EnteredScore
			}
			if math.Abs(got-want) > readBackScoreTolerance {
				mismatches = append(mismatches, fmt.Sprintf("score read back %s, requested %s", formatScore(got), formatScore(want)))
			}
		} else {
			got := after.EnteredGrade
			if got == "" {
				got = after.Grade
			}
			if !strings.EqualFold(got, req.PostedGrade) {
				mismatches = append(mismatches, fmt.Sprintf("grade read back %q, requested %q", got, req.PostedGrade))
			}
		}
	}

	if req.Excuse && !after.ExcusedTLN {
		mismatches = append(mismatches, "submission is not excused on read-back")
	}

	// A rubric-scored grade may carry no posted grade at all, so the rubric
	// rows are the only thing there is to verify. Every requested criterion
	// must read back with the points (and rating / comment, when given) asked.
	for id, want := range req.Rubric {
		got, ok := after.Rubric[id]
		if !ok {
			mismatches = append(mismatches, fmt.Sprintf("rubric criterion %s not on read-back", id))
			continue
		}
		if want.Points != nil && math.Abs(got.Points-*want.Points) > readBackScoreTolerance {
			mismatches = append(mismatches, fmt.Sprintf("rubric criterion %s read back %s points, requested %s", id, formatScore(got.Points), formatScore(*want.Points)))
		}
		if want.Rating != "" && got.RatingID != want.Rating {
			mismatches = append(mismatches, fmt.Sprintf("rubric criterion %s read back rating %q, requested %q", id, got.RatingID, want.Rating))
		}
		if want.Comments != "" && got.Comments != want.Comments {
			mismatches = append(mismatches, fmt.Sprintf("rubric criterion %s comment differs on read-back", id))
		}
	}

	var comment *api.SubmissionComment
	if req.Comment != "" {
		known := map[int64]bool{}
		if before != nil {
			for _, c := range before.SubmissionComments {
				known[c.ID] = true
			}
		}
		for i := range after.SubmissionComments {
			c := &after.SubmissionComments[i]
			if c.Comment != req.Comment {
				continue
			}
			if comment == nil || (known[comment.ID] && !known[c.ID]) {
				comment = c
			}
		}
		if comment == nil {
			mismatches = append(mismatches, "comment not found on read-back")
		}
	}

	return comment, mismatches
}

func formatScore(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// gradeState renders a submission's grade for the before → after line.
func gradeState(s *api.Submission) string {
	if s == nil {
		return "?"
	}
	if s.ExcusedTLN {
		return "excused"
	}
	if s.Grade == "" {
		return "ungraded"
	}
	if score := formatScore(s.Score); score != s.Grade && strings.TrimSuffix(strings.TrimSuffix(s.Grade, "0"), ".") != score {
		return fmt.Sprintf("%s (%s)", s.Grade, score)
	}
	return s.Grade
}

// commentSummary is "#id by author — first 80 chars".
func commentSummary(c *api.SubmissionComment) string {
	text := c.Comment
	if len(text) > 80 {
		text = text[:80] + "…"
	}
	author := c.AuthorName
	if author == "" && c.Author != nil {
		author = c.Author.Name
	}
	if author == "" {
		author = "author " + strconv.FormatInt(c.AuthorID, 10)
	}
	return fmt.Sprintf("#%d by %s — %q", c.ID, author, text)
}

// printReadBackLines prints the compact before/after evidence.
func printReadBackLines(rb *submissionReadBack) {
	if rb.Requested.PostedGrade != "" || rb.Requested.Excuse {
		printInfo("   grade: %s → %s\n", gradeState(rb.Before), gradeState(rb.After))
	}
	switch {
	case rb.CommentExisted:
		printInfo("   comment: already present, not posted again: %s\n", commentSummary(rb.Comment))
	case rb.Requested.Comment != "" && rb.Comment != nil:
		printInfo("   comment: %s\n", commentSummary(rb.Comment))
	case rb.Requested.Comment != "":
		printInfo("   comment: not found on read-back\n")
	}
	if n := len(rb.Requested.Rubric); n > 0 {
		ids := make([]string, 0, n)
		for id := range rb.Requested.Rubric {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		parts := make([]string, 0, n)
		for _, id := range ids {
			if got, ok := rb.After.Rubric[id]; ok {
				parts = append(parts, fmt.Sprintf("%s=%s", id, formatScore(got.Points)))
			} else {
				parts = append(parts, id+"=missing")
			}
		}
		printInfo("   rubric: %s\n", strings.Join(parts, " "))
	}
	if rb.Verified {
		printInfo("   verified: yes\n")
	} else {
		printInfo("   verified: no — %s\n", strings.Join(rb.Mismatches, "; "))
	}
}

// readBackError is the non-zero exit for a mismatch.
func readBackError(userID int64, rb *submissionReadBack) error {
	return fmt.Errorf("submission for user %d did not read back as requested: %s", userID, strings.Join(rb.Mismatches, "; "))
}
