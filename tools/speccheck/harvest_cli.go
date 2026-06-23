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
	var files []*ast.File

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
		files = append(files, f)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Pass 1: collect path-helper functions — those that RETURN an /api path
	// template (e.g. discussionContextPath -> "/api/v1/%s/%d/discussion_topics").
	// Service methods then build sub-resource paths as fmt.Sprintf("%s/.../x", helper(...))
	// or `base := helper(...); base + "/x"`, which a literal-only scan misses.
	helpers := map[string][]string{}
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || !returnsString(fn) {
				continue
			}
			var rets []string
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if ret, ok := n.(*ast.ReturnStmt); ok {
					for _, r := range ret.Results {
						rets = append(rets, apiStringsInExpr(r, nil)...)
					}
				}
				return true
			})
			if len(rets) > 0 {
				helpers[fn.Name.Name] = dedupe(rets)
			}
		}
	}

	// Pass 2: harvest (verb, path) per function, resolving helper substitution and
	// local variables assigned from helpers/path expressions.
	seen := map[string]Endpoint{}
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			locals := map[string][]string{} // var name -> possible /api path templates
			var paths []string
			verbSet := map[string]bool{}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch v := n.(type) {
				case *ast.AssignStmt:
					// Track `name := <path-expr>` so later "%s" substitutions resolve.
					if len(v.Lhs) >= 1 && len(v.Rhs) >= 1 {
						if id, ok := v.Lhs[0].(*ast.Ident); ok {
							if got := apiStringsInExpr(v.Rhs[0], helpers); len(got) > 0 {
								locals[id.Name] = dedupe(append(locals[id.Name], got...))
							}
						}
					}
				case *ast.CallExpr:
					if verb := callClientVerb(v); verb != "" {
						verbSet[verb] = true
					}
				}
				return true
			})
			// Collect candidate paths from every expression, now that locals are known.
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if e, ok := n.(ast.Expr); ok {
					paths = append(paths, apiStringsInExpr(e, helpers, locals)...)
				}
				return true
			})
			paths = dedupe(paths)
			if len(paths) == 0 || len(verbSet) == 0 {
				continue
			}
			for _, p := range paths {
				if !isAPIPath(p) {
					continue
				}
				for verb := range verbSet {
					seen[verb+" "+p] = Endpoint{Method: verb, Path: p}
				}
			}
		}
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

// returnsString reports whether fn declares at least one string result.
func returnsString(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil {
		return false
	}
	for _, fld := range fn.Type.Results.List {
		if id, ok := fld.Type.(*ast.Ident); ok && id.Name == "string" {
			return true
		}
	}
	return false
}

// dedupe returns ss with duplicates removed (insertion order preserved).
func dedupe(ss []string) []string {
	seen := map[string]bool{}
	out := ss[:0:0]
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// apiStringsInExpr resolves an expression to the set of /api path templates it
// can represent, substituting known path helpers and (optionally) local
// variables assigned from path expressions. It handles:
//
//	"/api/v1/..."                                  string literal
//	fmt.Sprintf("/api/v1/...%d", ...)              literal format string
//	fmt.Sprintf("%s/.../x", helper(...) | var)     %s-prefixed, base resolved
//	helper(...)                                    path-helper call
//	<api-expr> + "literal"                         concatenation
func apiStringsInExpr(e ast.Expr, helpers map[string][]string, localsOpt ...map[string][]string) []string {
	var locals map[string][]string
	if len(localsOpt) > 0 {
		locals = localsOpt[0]
	}
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			if s := unquote(v.Value); isAPIPath(s) {
				return []string{s}
			}
		}
	case *ast.Ident:
		if locals != nil {
			return locals[v.Name]
		}
	case *ast.CallExpr:
		if isFmtSprintf(v) && len(v.Args) > 0 {
			lit, ok := v.Args[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return nil
			}
			format := unquote(lit.Value)
			if isAPIPath(format) {
				return []string{format}
			}
			if strings.HasPrefix(format, "%s") && len(v.Args) >= 2 {
				var out []string
				for _, base := range apiStringsInExpr(v.Args[1], helpers, locals) {
					out = append(out, base+format[len("%s"):])
				}
				return out
			}
			return nil
		}
		if id, ok := v.Fun.(*ast.Ident); ok {
			return helpers[id.Name]
		}
	case *ast.BinaryExpr:
		if v.Op == token.ADD {
			bases := apiStringsInExpr(v.X, helpers, locals)
			if len(bases) == 0 {
				return nil
			}
			suffix := ""
			if lit, ok := v.Y.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				suffix = unquote(lit.Value)
			}
			out := make([]string, 0, len(bases))
			for _, b := range bases {
				out = append(out, b+suffix)
			}
			return out
		}
	}
	return nil
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
