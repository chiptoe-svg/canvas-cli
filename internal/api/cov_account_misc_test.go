package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCovAccountMiscExtra_AllMethods exercises every AccountMiscExtraService
// method against a permissive mock server for execution coverage of the
// account miscellaneous-admin surface.
func TestCovAccountMiscExtra_AllMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/accounts/1/rubrics/2/used_locations" {
			_, _ = w.Write([]byte(`[{"id":1}]`))
			return
		}
		_, _ = w.Write([]byte(`{"id":1,"name":"x"}`))
	}))
	defer server.Close()

	s := NewAccountMiscExtraService(newTestClient(t, server.URL))
	ctx := context.Background()
	body := map[string]interface{}{"name": "x"}

	_ = s.AccountRestoreUser(ctx, 1, 2)
	_, _ = s.AccountBulkUpdateUsers(ctx, 1, body)
	_, _ = s.AccountBulkEnrollment(ctx, 1, body)
	_ = s.AccountDeleteSubAccount(ctx, 1, 2)
	_, _ = s.AccountGetEnrollment(ctx, 1, 2)
	_, _ = s.AccountGetCourse(ctx, 1, 2)
	_, _ = s.AccountSelfRegistration(ctx, 1, body)
	_, _ = s.AccountGetRolesPermissions(ctx, 1)
	_, _ = s.AccountCreateSharedBrandConfig(ctx, 1, body)
	_, _ = s.AccountUpdateSharedBrandConfig(ctx, 1, 2, body)
	_, _ = s.AccountGetRubricUsedLocations(ctx, 1, 2)
	_, _ = s.AccountUploadRubric(ctx, 1, body)
	_, _ = s.AccountGetRubricUpload(ctx, 1, 2)
	_, _ = s.AccountCreateFolder(ctx, 1, body)
}
