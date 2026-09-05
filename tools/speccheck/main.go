// Command speccheck provides two modes for Canvas API spec compliance.
//
//	-sync     fetches the official Canvas Swagger 1.2 spec from a live Canvas
//	          host and writes testdata/spec/canvas_endpoints.json and
//	          testdata/spec/canvas_models.json.
//
//	-coverage prints a coverage gap report against the committed manifest and
//	          writes /tmp/canvas_coverage_gap.md
//
// The source of truth is the official Canvas Swagger 1.2 API, not the gitignored
// .ai/canvas-lms-docs markdown mirror. Configure the host via -host flag or the
// CANVAS_SPEC_HOST environment variable. Any Canvas host that serves /doc/api works;
// canvas.instructure.com does NOT (returns 503).
//
// Examples:
//
//	go run ./tools/speccheck -sync
//	go run ./tools/speccheck -sync -host https://myschool.instructure.com
//	CANVAS_SPEC_HOST=https://learn.canvas.net go run ./tools/speccheck -sync
//	go run ./tools/speccheck -coverage
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// defaultHost is the public Canvas Network host that returns 200 on /doc/api.
// canvas.instructure.com is NOT valid — it returns 503 from datacenter IPs.
const defaultHost = "https://learn.canvas.net"

// Paths relative to the repo root (tool is invoked as go run ./tools/speccheck from repo root).
const (
	manifestPath = "testdata/spec/canvas_endpoints.json"
	modelsPath   = "testdata/spec/canvas_models.json"
	cacheDir     = ".canvas-spec"
	apiSrcDir    = "internal/api"
)

// coverageOut is the path the coverage gap report is written to. It uses the OS
// temp dir so the tool and its tests work on Windows (no hardcoded /tmp).
var coverageOut = filepath.Join(os.TempDir(), "canvas_coverage_gap.md")

// userAgent mimics a browser so Canvas hosts don't block datacenter requests.
const userAgent = "Mozilla/5.0 (compatible; canvas-cli/speccheck; +https://github.com/chiptoe-svg/canvas-cli)"

// minResourceFraction is the minimum fraction of listed resources that must be
// successfully fetched before the manifest is written. Guards against silent
// partial fetches producing a misleadingly small manifest.
const minResourceFraction = 0.80

// Endpoint is a single documented Canvas API endpoint.
type Endpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// Manifest is the committed inventory of documented Canvas API endpoints.
type Manifest struct {
	GeneratedNote string     `json:"generated_note"`
	Source        string     `json:"source"`
	Endpoints     []Endpoint `json:"endpoints"`
}

// ModelField describes a single field in a Canvas API model.
type ModelField struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// ModelDef is a single Canvas API model definition.
type ModelDef struct {
	Resource string                `json:"resource"`
	Fields   map[string]ModelField `json:"fields"`
}

// ModelsManifest is the committed map of Canvas API models.
type ModelsManifest struct {
	GeneratedNote string               `json:"generated_note"`
	Source        string               `json:"source"`
	Models        map[string]*ModelDef `json:"models"`
}

// --- Swagger 1.2 response types ---

type swaggerIndex struct {
	APIs []struct {
		Path string `json:"path"` // e.g. "/courses.json"
	} `json:"apis"`
}

type swaggerResource struct {
	ResourcePath string `json:"resourcePath"` // e.g. "/courses"
	APIs         []struct {
		Path       string `json:"path"` // e.g. "/v1/courses/{course_id}"
		Operations []struct {
			Method   string `json:"method"` // "GET", "POST", etc.
			Nickname string `json:"nickname"`
		} `json:"operations"`
	} `json:"apis"`
	Models map[string]struct {
		ID          string `json:"id"`
		Description string `json:"description"`
		Properties  map[string]struct {
			Type        string `json:"type"`
			Description string `json:"description"`
		} `json:"properties"`
	} `json:"models"`
}

// swaggerParamRe replaces Swagger {param} placeholders with the :param
// convention used by the existing manifest and normalizer.
var swaggerParamRe = regexp.MustCompile(`\{[^}]+\}`)

func main() {
	syncFlag := flag.Bool("sync", false, "fetch official Canvas Swagger and write testdata/spec/canvas_endpoints.json")
	coverFlag := flag.Bool("coverage", false, "print coverage gap report and write /tmp/canvas_coverage_gap.md")
	hostFlag := flag.String("host", "", "Canvas host to fetch spec from (default: $CANVAS_SPEC_HOST or https://learn.canvas.net)")
	flag.Parse()

	host := resolveHost(*hostFlag)

	switch {
	case *syncFlag:
		if err := runSync(host); err != nil {
			fmt.Fprintln(os.Stderr, "speccheck -sync:", err)
			os.Exit(1)
		}
	case *coverFlag:
		if err := runCoverage(); err != nil {
			fmt.Fprintln(os.Stderr, "speccheck -coverage:", err)
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}

// resolveHost picks the Canvas host in priority order: flag > env > default.
func resolveHost(flagVal string) string {
	if flagVal != "" {
		return strings.TrimRight(flagVal, "/")
	}
	if env := os.Getenv("CANVAS_SPEC_HOST"); env != "" {
		return strings.TrimRight(env, "/")
	}
	return defaultHost
}

// httpGet performs a GET with browser-like User-Agent and a 30 s timeout.
// Uses the on-disk cache under .canvas-spec/ to skip already-fetched files.
func httpGet(url, cacheKey string) ([]byte, error) {
	// cacheKey is always filepath.Base(path) — a bare filename with no separators.
	// Double-applying Base here ensures path traversal is impossible even if the
	// caller supplies a path with directory components.
	cachePath := filepath.Join(cacheDir, filepath.Base(cacheKey))
	if data, err := os.ReadFile(cachePath); err == nil { // #nosec G304 -- path confined to fixed .canvas-spec/ dir; cacheKey is filepath.Base output
		return data, nil
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Write cache (best-effort; failures are non-fatal).
	if mkErr := os.MkdirAll(cacheDir, 0o750); mkErr == nil { //nolint:errcheck
		_ = os.WriteFile(cachePath, body, 0o600) //nolint:errcheck
	}
	return body, nil
}

// runSync fetches the official Canvas Swagger 1.2 spec and writes the manifests.
func runSync(host string) error {
	fmt.Printf("speccheck: fetching index from %s/doc/api/api-docs.json\n", host)

	indexData, err := httpGet(host+"/doc/api/api-docs.json", "api-docs.json")
	if err != nil {
		return fmt.Errorf("fetching index: %w", err)
	}

	var idx swaggerIndex
	if err := json.Unmarshal(indexData, &idx); err != nil {
		return fmt.Errorf("parsing index: %w", err)
	}
	if len(idx.APIs) == 0 {
		return fmt.Errorf("index listed 0 resources — unexpected response from %s", host)
	}

	listed := len(idx.APIs)
	fmt.Printf("speccheck: index lists %d resource files; fetching...\n", listed)

	// Fetch all resource files concurrently (bounded concurrency to be polite).
	type result struct {
		name string
		res  *swaggerResource
		err  error
	}
	results := make(chan result, listed)
	sem := make(chan struct{}, 8) // max 8 in-flight

	var wg sync.WaitGroup
	for _, api := range idx.APIs {
		// api.path is e.g. "/courses.json"
		resourceName := strings.TrimPrefix(api.Path, "/")                           // "courses.json"
		resourceName = strings.TrimSuffix(resourceName, filepath.Ext(resourceName)) // "courses"
		url := host + "/doc/api" + api.Path
		cacheKey := filepath.Base(api.Path)

		wg.Add(1)
		go func(name, fetchURL, ck string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			data, fetchErr := httpGet(fetchURL, ck)
			if fetchErr != nil {
				results <- result{name: name, err: fetchErr}
				return
			}
			var sr swaggerResource
			if parseErr := json.Unmarshal(data, &sr); parseErr != nil {
				results <- result{name: name, err: parseErr}
				return
			}
			results <- result{name: name, res: &sr}
		}(resourceName, url, cacheKey)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	seen := map[string]bool{}
	var endpoints []Endpoint
	models := map[string]*ModelDef{}
	var fetchErrors []string

	for r := range results {
		if r.err != nil {
			fetchErrors = append(fetchErrors, fmt.Sprintf("%s: %v", r.name, r.err))
			fmt.Fprintf(os.Stderr, "  WARNING: skipping %s: %v\n", r.name, r.err)
			continue
		}
		extractEndpoints(r.res, seen, &endpoints)
		extractModels(r.res, r.name, models)
	}

	fetched := listed - len(fetchErrors)
	pctFetched := float64(fetched) / float64(listed)
	if pctFetched < minResourceFraction {
		return fmt.Errorf(
			"only %d/%d resources fetched (%.0f%% < required %.0f%%); refusing to write truncated manifest.\nErrors:\n  %s",
			fetched, listed, pctFetched*100, minResourceFraction*100,
			strings.Join(fetchErrors, "\n  "),
		)
	}
	if len(fetchErrors) > 0 {
		fmt.Fprintf(os.Stderr, "WARNING: %d/%d resources had fetch errors (skipped)\n", len(fetchErrors), listed)
	}

	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Path != endpoints[j].Path {
			return endpoints[i].Path < endpoints[j].Path
		}
		return endpoints[i].Method < endpoints[j].Method
	})

	if err := writeManifest(host, endpoints); err != nil {
		return err
	}
	if err := writeModels(host, models); err != nil {
		return err
	}
	fmt.Printf("speccheck: %d endpoints, %d models -> %s\n", len(endpoints), len(models), manifestPath)
	return nil
}

// extractEndpoints walks the swagger resource and appends discovered endpoints
// to the slice, deduplicating via seen.
func extractEndpoints(res *swaggerResource, seen map[string]bool, endpoints *[]Endpoint) {
	for _, api := range res.APIs {
		// Swagger path is /v1/courses/{id} — prepend /api to match CLI paths.
		rawPath := "/api" + api.Path
		// Convert Swagger {param} braces to :param colon-style used by the normalizer.
		colonPath := swaggerParamRe.ReplaceAllStringFunc(rawPath, func(m string) string {
			// Strip braces: {course_id} -> :course_id
			inner := m[1 : len(m)-1]
			return ":" + inner
		})
		// Strip trailing slash and query strings.
		if i := strings.IndexByte(colonPath, '?'); i >= 0 {
			colonPath = colonPath[:i]
		}
		colonPath = strings.TrimRight(colonPath, "/")

		for _, op := range api.Operations {
			method := strings.ToUpper(op.Method)
			key := method + "|" + colonPath
			if seen[key] {
				continue
			}
			seen[key] = true
			*endpoints = append(*endpoints, Endpoint{Method: method, Path: colonPath})
		}
	}
}

// extractModels walks the swagger resource models and merges them into the map.
func extractModels(res *swaggerResource, resourceName string, models map[string]*ModelDef) {
	for modelName, model := range res.Models {
		if _, dup := models[modelName]; dup {
			continue // first definition wins (models may appear in multiple resources)
		}
		fields := make(map[string]ModelField, len(model.Properties))
		for propName, prop := range model.Properties {
			fields[propName] = ModelField{
				Type:        prop.Type,
				Description: prop.Description,
			}
		}
		models[modelName] = &ModelDef{
			Resource: resourceName,
			Fields:   fields,
		}
	}
}

// writeManifest serialises the endpoint list to testdata/spec/canvas_endpoints.json.
func writeManifest(host string, endpoints []Endpoint) error {
	man := Manifest{
		GeneratedNote: "Auto-generated by tools/speccheck -sync from official Canvas Swagger 1.2. Do not edit manually. Re-run with: make spec-sync",
		Source:        host + "/doc/api",
		Endpoints:     endpoints,
	}
	return writeJSON(manifestPath, man)
}

// writeModels serialises the model map to testdata/spec/canvas_models.json.
func writeModels(host string, models map[string]*ModelDef) error {
	// Sort model names for a stable, reviewable file.
	sorted := make(map[string]*ModelDef, len(models))
	for k, v := range models {
		// Sort fields within each model.
		if v != nil && len(v.Fields) > 0 {
			sorted[k] = v
		}
	}
	mm := ModelsManifest{
		GeneratedNote: "Auto-generated by tools/speccheck -sync from official Canvas Swagger 1.2. Do not edit manually.",
		Source:        host + "/doc/api",
		Models:        sorted,
	}
	return writeJSON(modelsPath, mm)
}

// writeJSON marshals v to JSON and writes it atomically.
func writeJSON(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	// #nosec G306 -- committed spec manifest/models, world-readable by design
	return os.WriteFile(path, append(out, '\n'), 0o644)
}

// runCoverage reports which documented endpoints the CLI has not implemented.
// This is network-free: it reads only the committed manifest.
func runCoverage() error {
	man, err := loadManifest()
	if err != nil {
		return err
	}

	cliEndpoints, err := harvestCLIEndpoints(apiSrcDir)
	if err != nil {
		return fmt.Errorf("harvesting CLI endpoints from %s: %w", apiSrcDir, err)
	}

	// Build the CLI's (method, normalized-path) set. Coverage is method-granular:
	// a documented endpoint counts as implemented only if the CLI issues a request
	// with the SAME HTTP verb on a matching path.
	cliSet := map[string]bool{}
	for _, e := range cliEndpoints {
		cliSet[e.Method+" "+normalizePath(e.Path)] = true
	}
	// implemented reports whether the CLI covers a documented (method, path),
	// allowing the multi-context aliases the service layer builds via templates.
	implemented := func(method, norm string) bool {
		if cliSet[method+" "+norm] {
			return true
		}
		for _, alias := range collapseContextPath(norm) {
			if cliSet[method+" "+alias] {
				return true
			}
		}
		return false
	}

	type resourceGap struct {
		resource    string
		implemented int
		total       int
		missing     []string
	}
	byResource := map[string]*resourceGap{}
	for _, ep := range man.Endpoints {
		seg := firstSegment(ep.Path)
		if _, ok := byResource[seg]; !ok {
			byResource[seg] = &resourceGap{resource: seg}
		}
		rg := byResource[seg]
		rg.total++
		if implemented(ep.Method, normalizePath(ep.Path)) {
			rg.implemented++
		} else {
			rg.missing = append(rg.missing, fmt.Sprintf("%s %s", ep.Method, ep.Path))
		}
	}

	resources := make([]string, 0, len(byResource))
	for r := range byResource {
		resources = append(resources, r)
	}
	sort.Strings(resources)

	totalImpl, totalDoc := 0, 0
	for _, rg := range byResource {
		totalImpl += rg.implemented
		totalDoc += rg.total
	}
	pct := 0
	if totalDoc > 0 {
		pct = totalImpl * 100 / totalDoc
	}

	summary := fmt.Sprintf("CLI implements %d of %d official Canvas API endpoint patterns (%d%%)", totalImpl, totalDoc, pct)
	fmt.Println(summary)
	fmt.Println()

	type gapEntry struct {
		resource string
		gap      int
	}
	var gaps []gapEntry
	for _, r := range resources {
		rg := byResource[r]
		if len(rg.missing) > 0 {
			gaps = append(gaps, gapEntry{r, len(rg.missing)})
		}
	}
	sort.Slice(gaps, func(i, j int) bool { return gaps[i].gap > gaps[j].gap })

	fmt.Println("Top resource areas with unimplemented official Canvas endpoints:")
	for _, g := range gaps {
		rg := byResource[g.resource]
		fmt.Printf("  %-40s %d missing / %d total\n", g.resource, g.gap, rg.total)
	}

	var sb strings.Builder
	sb.WriteString("# Canvas CLI Official Spec Coverage Gap Report\n\n")
	fmt.Fprintf(&sb, "**%s**\n\n", summary)
	sb.WriteString("Source: official Canvas Swagger 1.2 (refresh with `make spec-sync`)\n\n")
	sb.WriteString("| Resource | Implemented | Documented | Gap |\n")
	sb.WriteString("|----------|------------|------------|-----|\n")
	for _, r := range resources {
		rg := byResource[r]
		fmt.Fprintf(&sb, "| %s | %d | %d | %d |\n", r, rg.implemented, rg.total, len(rg.missing))
	}
	sb.WriteString("\n## Missing Endpoints (official Canvas spec but not implemented by CLI)\n\n")
	for _, r := range resources {
		rg := byResource[r]
		if len(rg.missing) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "### %s\n\n", r)
		for _, m := range rg.missing {
			fmt.Fprintf(&sb, "- `%s`\n", m)
		}
		sb.WriteString("\n")
	}

	// #nosec G306 -- coverage report, world-readable by design
	if err := os.WriteFile(coverageOut, []byte(sb.String()), 0o644); err != nil {
		return err
	}
	fmt.Printf("\nFull report written to %s\n", coverageOut)
	return nil
}

// loadManifest reads testdata/spec/canvas_endpoints.json.
func loadManifest() (*Manifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w (run `make spec-sync` first)", manifestPath, err)
	}
	var man Manifest
	if err := json.Unmarshal(data, &man); err != nil {
		return nil, err
	}
	if len(man.Endpoints) == 0 {
		return nil, fmt.Errorf("manifest has no endpoints — run `make spec-sync`")
	}
	return &man, nil
}

// firstSegment returns the first path segment after /api/v1/ (or /api/lti/).
func firstSegment(path string) string {
	for _, prefix := range []string{"/api/v1/", "/api/lti/"} {
		if after, ok := strings.CutPrefix(path, prefix); ok {
			seg, _, _ := strings.Cut(after, "/")
			return seg
		}
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}
