package api

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// specEndpoint mirrors testdata/spec/canvas_endpoints.json.
type specEndpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type specManifest struct {
	Endpoints []specEndpoint `json:"endpoints"`
}

// knownUndocumented lists CLI paths that have no matching documented endpoint
// and have been triaged individually. Every entry MUST have a real justification.
var knownUndocumented = map[string]string{
	// Canvas's course-level bulk-grade endpoint is POST /api/v1/courses/:id/submissions/update_grades.
	// The CLI additionally calls POST /api/v1/courses/:id/assignments/:id/submissions/update_grades
	// (assignment-scoped bulk grade). Canvas does support this path in practice but it is not
	// listed in the Canvas REST API documentation captured in .ai/canvas-lms-docs.
	"/api/v1/courses/:x/assignments/:x/submissions/update_grades": "assignment-scoped bulk-grade endpoint: undocumented Canvas API path supported in practice (docs only list course-level /courses/:id/submissions/update_grades)",
}

// TestSpecContract_CLIPathsAreDocumented verifies that every path the CLI
// calls in internal/api/ matches at least one documented Canvas API endpoint
// in testdata/spec/canvas_endpoints.json.
//
// Refresh testdata/spec/canvas_endpoints.json with: make spec-sync
//
// Triage process for failures:
//
//	(a) Real CLI bug — path is wrong; fix the implementation, do not allowlist
//	(b) Normalization artifact — fix normalizePath in tools/speccheck/normalize.go
//	(c) Legitimately undocumented — add to knownUndocumented with a real comment
func TestSpecContract_CLIPathsAreDocumented(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "spec", "canvas_endpoints.json"))
	if err != nil {
		t.Fatalf("cannot read manifest (run `make spec-sync` first): %v", err)
	}

	var man specManifest
	if err := json.Unmarshal(data, &man); err != nil {
		t.Fatalf("parsing manifest: %v", err)
	}
	if len(man.Endpoints) == 0 {
		t.Fatal("manifest has no endpoints — run `make spec-sync`")
	}

	// Build a set of normalized documented paths.
	docPaths := map[string]bool{}
	for _, ep := range man.Endpoints {
		docPaths[specNormalizePath(ep.Path)] = true
	}

	// Harvest CLI paths from internal/api non-test source files.
	cliPaths, err := specHarvestCLIPaths(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("harvesting CLI paths: %v", err)
	}

	var failures []string
	for _, rawPath := range cliPaths {
		norm := specNormalizePath(rawPath)
		if docPaths[norm] {
			continue
		}
		if reason, ok := knownUndocumented[norm]; ok {
			t.Logf("SKIP undocumented (allowed): %s — %s", norm, reason)
			continue
		}
		failures = append(failures, rawPath+" -> "+norm)
	}
	sort.Strings(failures)
	for _, f := range failures {
		t.Errorf("CLI path has no documented Canvas endpoint: %s\n"+
			"  Triage: (a) fix the CLI path, (b) fix the normalizer, or (c) add to knownUndocumented with justification", f)
	}
}

// specNormalizePath mirrors tools/speccheck/normalize.go. Kept here inline to
// avoid a dependency on the tool package from the test package.
//
//   - :param_name → :x
//   - %d/%s/%v   → :x
//   - /self/      → /:x/
//   - query string stripped
func specNormalizePath(path string) string {
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	path = strings.TrimRight(path, "/")
	path = specFmtVerbRe.ReplaceAllString(path, ":x")
	path = specParamRe.ReplaceAllString(path, ":x")
	path = specSelfSegRe.ReplaceAllStringFunc(path, func(m string) string {
		if strings.HasSuffix(m, "/") {
			return "/:x/"
		}
		return "/:x"
	})
	return path
}

var (
	specParamRe   = regexp.MustCompile(`:[a-z_]+`)
	specFmtVerbRe = regexp.MustCompile(`%[a-z]`)
	specSelfSegRe = regexp.MustCompile(`/self(/|$)`)
)

// specHarvestCLIPaths walks internal/api/ (non-test .go files) and extracts
// /api/v1/... and /api/lti/... string literals including fmt.Sprintf format strings.
func specHarvestCLIPaths(repoRoot string) ([]string, error) {
	dir := filepath.Join(repoRoot, "internal", "api")
	fset := token.NewFileSet()
	seen := map[string]bool{}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.BasicLit:
				if v.Kind == token.STRING {
					s := specUnquote(v.Value)
					if specIsAPIPath(s) {
						seen[s] = true
					}
				}
			case *ast.CallExpr:
				if specIsFmtSprintf(v) && len(v.Args) > 0 {
					if lit, ok := v.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						s := specUnquote(lit.Value)
						if specIsAPIPath(s) {
							seen[s] = true
						}
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(seen))
	for p := range seen {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths, nil
}

func specIsAPIPath(s string) bool {
	return strings.HasPrefix(s, "/api/v1/") || strings.HasPrefix(s, "/api/lti/")
}

func specIsFmtSprintf(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "fmt" && sel.Sel.Name == "Sprintf"
}

func specUnquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		return s[1 : len(s)-1]
	}
	return s
}
