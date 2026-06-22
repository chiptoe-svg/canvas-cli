package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBrandService_GetVariables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/brand_variables" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		bv := BrandVariables{IcBrandPrimary: "#E66000"}
		if err := json.NewEncoder(w).Encode(bv); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBrandService(client)

	bv, err := svc.GetVariables(context.Background())
	if err != nil {
		t.Fatalf("GetVariables failed: %v", err)
	}
	if bv.IcBrandPrimary != "#E66000" {
		t.Errorf("unexpected primary color: %s", bv.IcBrandPrimary)
	}
}

func TestBrandService_GetVariablesForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/3/brand_variables" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		bv := BrandVariables{IcBrandNavBgd: "#0770A3"}
		if err := json.NewEncoder(w).Encode(bv); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBrandService(client)

	bv, err := svc.GetVariablesForAccount(context.Background(), 3)
	if err != nil {
		t.Fatalf("GetVariablesForAccount failed: %v", err)
	}
	if bv.IcBrandNavBgd != "#0770A3" {
		t.Errorf("unexpected nav color: %s", bv.IcBrandNavBgd)
	}
}

func TestBrandService_GetVariablesForCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/7/brand_variables" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		bv := BrandVariables{IcBrandPrimary: "#AABBCC"}
		if err := json.NewEncoder(w).Encode(bv); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewBrandService(client)

	bv, err := svc.GetVariablesForCourse(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetVariablesForCourse failed: %v", err)
	}
	if bv.IcBrandPrimary != "#AABBCC" {
		t.Errorf("unexpected color: %s", bv.IcBrandPrimary)
	}
}
