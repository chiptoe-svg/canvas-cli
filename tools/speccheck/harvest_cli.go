package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// harvestCLIPaths uses go/ast to extract /api/v1/... paths from non-test Go
// source files under dir. It captures:
//   - Plain string literals containing /api/v1/
//   - The first argument of fmt.Sprintf calls when it is a string literal
//     containing /api/v1/
//
// The returned slice is sorted and deduplicated.
func harvestCLIPaths(dir string) ([]string, error) {
	fset := token.NewFileSet()
	seen := map[string]bool{}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// Only .go source files; skip test files.
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		f, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}

		// NOTE: paths built by string concatenation where the /api/v1/ prefix is
		// in a variable (not a literal) are not captured. Enhance this function if
		// such patterns are added to the service layer.
		ast.Inspect(f, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.BasicLit:
				// Plain string literal.
				if v.Kind == token.STRING {
					s := unquote(v.Value)
					if isAPIPath(s) {
						seen[s] = true
					}
				}
			case *ast.CallExpr:
				// fmt.Sprintf("...", ...) — capture the format string.
				if isFmtSprintf(v) && len(v.Args) > 0 {
					if lit, ok := v.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						s := unquote(lit.Value)
						if isAPIPath(s) {
							seen[s] = true
						}
					}
				}
			case *ast.BinaryExpr:
				// String concatenation: "/api/v1/..." + someVar
				if v.Op == token.ADD {
					if lhs, ok := v.X.(*ast.BasicLit); ok && lhs.Kind == token.STRING {
						s := unquote(lhs.Value)
						if isAPIPath(s) {
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

// isAPIPath reports whether s is a Canvas API path.
func isAPIPath(s string) bool {
	return strings.HasPrefix(s, "/api/v1/") || strings.HasPrefix(s, "/api/lti/")
}

// isFmtSprintf reports whether call is a fmt.Sprintf invocation.
func isFmtSprintf(call *ast.CallExpr) bool {
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

// unquote strips surrounding double or backtick quotes from a Go string literal.
func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		return s[1 : len(s)-1]
	}
	return s
}
