package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

// TestCheckAndUpdate_MissingChecksums verifies that the updater aborts when
// checksums.txt is absent from the release rather than installing unverified
// binaries (fail-closed behaviour).
func TestCheckAndUpdate_MissingChecksums(t *testing.T) {
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}

	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	assetName := fmt.Sprintf("canvas-cli_%s_%s%s", runtime.GOOS, archName, ext)

	// Build a minimal tar.gz archive containing a fake binary.
	binaryName := BinaryName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	var archiveBuf bytes.Buffer
	gzWriter := gzip.NewWriter(&archiveBuf)
	tarWriter := tar.NewWriter(gzWriter)
	fakeContent := []byte("fake-binary-content")
	_ = tarWriter.WriteHeader(&tar.Header{
		Name:     binaryName,
		Mode:     0755,
		Size:     int64(len(fakeContent)),
		Typeflag: tar.TypeReg,
	})
	tarWriter.Write(fakeContent)
	tarWriter.Close()
	gzWriter.Close()

	archiveBytes := archiveBuf.Bytes()

	release := Release{
		TagName: "v99.0.0",
		Assets: []Asset{
			{
				Name: assetName,
				// BrowserDownloadURL will be filled in below
			},
			// Deliberately NO checksums.txt asset
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/releases/latest"):
			// Patch the download URL now that we know the server address.
			release.Assets[0].BrowserDownloadURL = "http://" + r.Host + "/archive"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(release)
		case r.URL.Path == "/archive":
			w.WriteHeader(http.StatusOK)
			w.Write(archiveBytes)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	updater := NewUpdater("1.0.0")
	updater.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	ctx := context.Background()
	result := updater.CheckAndUpdate(ctx)

	if result.Updated {
		t.Fatal("updater installed binary despite missing checksums.txt")
	}
	if result.Error == nil {
		t.Fatal("expected an error for missing checksums.txt, got nil")
	}
	if !strings.Contains(result.Error.Error(), "checksums.txt") {
		t.Errorf("error message should mention 'checksums.txt', got: %v", result.Error)
	}
}

// TestCheckAndUpdate_WithChecksums verifies the happy path: when checksums.txt
// is present and correct the updater succeeds (or at least reaches the apply
// step, which may fail in this test because there is no real executable path).
func TestCheckAndUpdate_WithChecksums(t *testing.T) {
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}

	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	assetName := fmt.Sprintf("canvas-cli_%s_%s%s", runtime.GOOS, archName, ext)

	binaryName := BinaryName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	var archiveBuf bytes.Buffer
	gzWriter := gzip.NewWriter(&archiveBuf)
	tarWriter := tar.NewWriter(gzWriter)
	fakeContent := []byte("fake-binary-content")
	_ = tarWriter.WriteHeader(&tar.Header{
		Name:     binaryName,
		Mode:     0755,
		Size:     int64(len(fakeContent)),
		Typeflag: tar.TypeReg,
	})
	tarWriter.Write(fakeContent)
	tarWriter.Close()
	gzWriter.Close()

	archiveBytes := archiveBuf.Bytes()
	hash := sha256.Sum256(archiveBytes)
	checksumsContent := fmt.Sprintf("%s  %s\n", hex.EncodeToString(hash[:]), assetName)

	release := Release{
		TagName: "v99.0.0",
		Assets:  []Asset{
			// URLs filled at request time.
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/releases/latest"):
			release.Assets = []Asset{
				{Name: assetName, BrowserDownloadURL: "http://" + r.Host + "/archive"},
				{Name: "checksums.txt", BrowserDownloadURL: "http://" + r.Host + "/checksums"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(release)
		case r.URL.Path == "/archive":
			w.WriteHeader(http.StatusOK)
			w.Write(archiveBytes)
		case r.URL.Path == "/checksums":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, checksumsContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	updater := NewUpdater("1.0.0")
	updater.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	ctx := context.Background()
	result := updater.CheckAndUpdate(ctx)

	// The test environment likely lacks a writable executable at os.Executable(),
	// so we expect either success or an "apply" error, but NOT a checksum error.
	if result.Error != nil && strings.Contains(result.Error.Error(), "checksum") {
		t.Errorf("unexpected checksum error when checksums are valid: %v", result.Error)
	}
}
