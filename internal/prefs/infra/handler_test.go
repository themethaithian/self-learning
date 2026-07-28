package infra

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/themethaithian/self-learning/internal/platform/httpserver"
	"github.com/themethaithian/self-learning/internal/prefs/domain"
)

var _ httpserver.Registrar = (*Handler)(nil)

type fakeService struct {
	value  string
	ok     bool
	getErr error

	setErr    error
	setCalled bool
	setRaw    *string
}

func (f *fakeService) FocusTrack(context.Context) (string, bool, error) {
	return f.value, f.ok, f.getErr
}

func (f *fakeService) SetFocusTrack(_ context.Context, raw *string) error {
	f.setCalled = true
	f.setRaw = raw
	return f.setErr
}

func newTestLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return slog.New(slog.NewTextHandler(&buf, nil)), &buf
}

func doRequest(h *Handler, method string, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	h.Register(mux)

	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, "/api/v1/prefs/focus-track", reader)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHandlerGetFocusTrack_Unset(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodGet, "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want application/json; charset=utf-8", ct)
	}
	if got, want := rec.Body.String(), `{"track":null}`+"\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerGetFocusTrack_Set(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{value: "ai-systems", ok: true}, logger)

	rec := doRequest(h, http.MethodGet, "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got, want := rec.Body.String(), `{"track":"ai-systems"}`+"\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerGetFocusTrack_ServiceErrorLeaksNoInternals(t *testing.T) {
	logger, logBuf := newTestLogger()
	h := NewHandler(&fakeService{getErr: errors.New("mysql: connection refused on 10.0.0.5")}, logger)

	rec := doRequest(h, http.MethodGet, "")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	body := rec.Body.String()
	for _, leak := range []string{"mysql", "10.0.0.5", "connection refused"} {
		if strings.Contains(body, leak) {
			t.Fatalf("body = %q leaks internal error detail %q", body, leak)
		}
	}
	if !strings.Contains(logBuf.String(), "connection refused") {
		t.Fatalf("log output = %q, want it to contain the real error", logBuf.String())
	}
}

func TestHandlerPutFocusTrack_Success(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, `{"track":"ai-systems"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !svc.setCalled {
		t.Fatal("expected SetFocusTrack to be called")
	}
	if svc.setRaw == nil || *svc.setRaw != "ai-systems" {
		t.Fatalf("SetFocusTrack called with %v, want \"ai-systems\"", svc.setRaw)
	}
	if got, want := rec.Body.String(), `{"track":"ai-systems"}`+"\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

// TestHandlerPutFocusTrack_Idempotent pins the ticket's DoD: PUTting the
// same value twice must both return 200, never an error the second time.
func TestHandlerPutFocusTrack_Idempotent(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	for i := 0; i < 2; i++ {
		rec := doRequest(h, http.MethodPut, `{"track":"ai-systems"}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("call %d: status = %d, want %d", i, rec.Code, http.StatusOK)
		}
	}
}

func TestHandlerPutFocusTrack_Clear(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, `{"track":null}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !svc.setCalled {
		t.Fatal("expected SetFocusTrack to be called")
	}
	if svc.setRaw != nil {
		t.Fatalf("SetFocusTrack called with %v, want nil", svc.setRaw)
	}
	if got, want := rec.Body.String(), `{"track":null}`+"\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerPutFocusTrack_InvalidFormat(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{setErr: fmt.Errorf("prefs: set focus track: %w", domain.ErrInvalidFocusTrack)}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, `{"track":"Not A Valid Slug!"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestHandlerPutFocusTrack_InvalidJSON(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodPut, `not json`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestHandlerPutFocusTrack_UnknownField guards against a typo'd field name
// (e.g. "trac" instead of "track") decoding as a no-op Track == nil and
// silently clearing the preference instead of failing loudly.
func TestHandlerPutFocusTrack_UnknownField(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, `{"trac":"ai-systems"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if svc.setCalled {
		t.Fatal("SetFocusTrack must not be called for a body with an unknown field")
	}
}

// TestHandlerPutFocusTrack_TrailingJSON guards dec.More(): a plain
// json.Decoder.Decode call silently stops after the first JSON value, so a
// body carrying two concatenated values would otherwise decode (and act on)
// only the first one instead of failing.
func TestHandlerPutFocusTrack_TrailingJSON(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, `{"track":"a"}{"track":"b"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if svc.setCalled {
		t.Fatal("SetFocusTrack must not be called for a body with trailing JSON")
	}
}

// TestHandlerPutFocusTrack_EmptyBodyClears pins intended PUT/replace
// semantics: a body with no "track" field at all decodes the same as an
// explicit {"track":null} — clearing the pref — which is correct REST
// replace behavior, not an oversight.
func TestHandlerPutFocusTrack_EmptyBodyClears(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, `{}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !svc.setCalled {
		t.Fatal("expected SetFocusTrack to be called")
	}
	if svc.setRaw != nil {
		t.Fatalf("SetFocusTrack called with %v, want nil", svc.setRaw)
	}
}

func TestHandlerPutFocusTrack_ServiceErrorLeaksNoInternals(t *testing.T) {
	logger, logBuf := newTestLogger()
	svc := &fakeService{setErr: errors.New("mysql: connection refused on 10.0.0.5")}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, `{"track":"ai-systems"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	body := rec.Body.String()
	for _, leak := range []string{"mysql", "10.0.0.5", "connection refused"} {
		if strings.Contains(body, leak) {
			t.Fatalf("body = %q leaks internal error detail %q", body, leak)
		}
	}
	if !strings.Contains(logBuf.String(), "connection refused") {
		t.Fatalf("log output = %q, want it to contain the real error", logBuf.String())
	}
}

func TestHandlerRejectsNonGETNonPUT(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodPost, "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
