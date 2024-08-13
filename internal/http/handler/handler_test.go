package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_PingHandler(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.PingHandler(w, req)

	res := w.Result()

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: %d", res.StatusCode)
	}

	b, err := io.ReadAll(w.Body)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !strings.EqualFold("pong", string(b)) {
		t.Errorf("unexpected body: %s", string(b))
	}
}
