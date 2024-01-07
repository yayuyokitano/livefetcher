package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CommitMutating(t *testing.T, r *http.Request, execute HTTPImplementer, token string, affectedRows int64) {
	t.Helper()
	w := httptest.NewRecorder()
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	HandleMethod(execute, w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Status code is %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() == "" {
		t.Errorf("Body is empty")
	}
	if w.Body.String() != fmt.Sprintf("%d", affectedRows) {
		t.Errorf("Body is %s, want %s", w.Body.String(), fmt.Sprintf("%d", affectedRows))
	}
}
