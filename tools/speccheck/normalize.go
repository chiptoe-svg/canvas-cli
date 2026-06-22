package main

import (
	"regexp"
	"strings"
)

// paramRe matches documented Canvas path parameters like :course_id, :id, :user_id, etc.
var paramRe = regexp.MustCompile(`:[a-z_]+`)

// fmtVerbRe matches Go fmt format verbs (%d, %s, %v, etc.) used in Sprintf path strings.
var fmtVerbRe = regexp.MustCompile(`%[a-z]`)

// selfSegRe matches the literal "self" as a path segment, which Canvas uses as an alias
// for the authenticated user. Documented paths use :user_id in the same position.
var selfSegRe = regexp.MustCompile(`/self(/|$)`)

// normalizePath reduces a path to a comparable form:
//
//   - Documented params (:course_id, :id, :user_id ...) → :x
//   - Go format verbs (%d, %s, %v)                     → :x
//   - Canvas "self" user alias (/self/ or /self$)       → /:x
//   - Trailing slash stripped
//   - Literal segments (front_page, latest, activate)   → kept as-is
//   - Query string stripped (everything from ? onward)
func normalizePath(path string) string {
	// Strip query string.
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	// Strip trailing slash.
	path = strings.TrimRight(path, "/")

	// Replace Go format verbs first (CLI paths), then documented params (spec paths).
	path = fmtVerbRe.ReplaceAllString(path, ":x")
	path = paramRe.ReplaceAllString(path, ":x")

	// Replace /self/ and /self (at path end) with /:x/ and /:x — Canvas "self" means
	// "current user" and is documented as /:user_id in the API spec.
	path = selfSegRe.ReplaceAllStringFunc(path, func(m string) string {
		if strings.HasSuffix(m, "/") {
			return "/:x/"
		}
		return "/:x"
	})

	return path
}
