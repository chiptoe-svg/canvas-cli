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

// clientVerb maps a *Client HTTP-helper method name to its HTTP verb. These are
// the wrappers the service layer calls (s.client.<Method>).
var clientVerb = map[string]string{
	"Get": "GET", "GetJSON": "GET", "GetAllPages": "GET", "GetAllPagesGeneric": "GET",
	"Post": "POST", "PostJSON": "POST",
	"Put": "PUT", "PutJSON": "PUT",
	"Patch": "PATCH", "PatchJSON": "PATCH",
	"Delete": "DELETE", "DeleteJSON": "DELETE",
}

// harvestCLIEndpoints extracts (HTTP method, path) pairs from the service layer
// by pairing, within each function, the /api/... path string(s) with the verb of
// the s.client.<Method> call(s) made in that function. Most service methods are
// exactly one path + one client call, so the pairing is exact; the rare
// multi-path or multi-verb function emits the cross product (an accepted
// approximation for coverage accounting). Skips paths in functions with no
// recognizable client verb (cannot attribute a method).
func harvestCLIEndpoints(dir string) ([]Endpoint, error) {
	fset := token.NewFileSet()
	seen := map[string]Endpoint{}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			var paths []string
			verbSet := map[string]bool{}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch v := n.(type) {
				case *ast.BasicLit:
					if v.Kind == token.STRING {
						if s := unquote(v.Value); isAPIPath(s) {
							paths = append(paths, s)
						}
					}
				case *ast.BinaryExpr:
					if v.Op == token.ADD {
						if lhs, ok := v.X.(*ast.BasicLit); ok && lhs.Kind == token.STRING {
							if s := unquote(lhs.Value); isAPIPath(s) {
								paths = append(paths, s)
							}
						}
					}
				case *ast.CallExpr:
					if isFmtSprintf(v) && len(v.Args) > 0 {
						if lit, ok := v.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
							if s := unquote(lit.Value); isAPIPath(s) {
								paths = append(paths, s)
							}
						}
					}
					if verb := callClientVerb(v); verb != "" {
						verbSet[verb] = true
					}
				}
				return true
			})
			if len(paths) == 0 || len(verbSet) == 0 {
				continue
			}
			for _, p := range paths {
				for verb := range verbSet {
					seen[verb+" "+p] = Endpoint{Method: verb, Path: p}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	eps := make([]Endpoint, 0, len(seen))
	for _, e := range seen {
		eps = append(eps, e)
	}
	sort.Slice(eps, func(i, j int) bool {
		if eps[i].Path != eps[j].Path {
			return eps[i].Path < eps[j].Path
		}
		return eps[i].Method < eps[j].Method
	})
	return eps, nil
}

// callClientVerb returns the HTTP verb if call is an s.client.<Method> (or
// client.<Method> / c.<Method>) HTTP-helper invocation, else "".
func callClientVerb(call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	verb, isVerb := clientVerb[sel.Sel.Name]
	if !isVerb {
		return ""
	}
	switch x := sel.X.(type) {
	case *ast.SelectorExpr:
		if x.Sel.Name == "client" {
			return verb
		}
	case *ast.Ident:
		if x.Name == "client" || x.Name == "c" {
			return verb
		}
	}
	return ""
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
