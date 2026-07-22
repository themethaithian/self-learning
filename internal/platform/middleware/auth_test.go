package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearerAuth(t *testing.T) {
	const token = "correct-token"
	passthrough := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := BearerAuth(token)(passthrough)

	tests := []struct {
		name       string
		authHeader string
		setHeader  bool
		wantStatus int
	}{
		{name: "no header", wantStatus: http.StatusUnauthorized},
		{name: "wrong token", authHeader: "Bearer wrong-token", setHeader: true, wantStatus: http.StatusUnauthorized},
		{name: "correct token", authHeader: "Bearer correct-token", setHeader: true, wantStatus: http.StatusOK},
		{name: "wrong scheme", authHeader: "Basic x", setHeader: true, wantStatus: http.StatusUnauthorized},
		{name: "bearer with no space or token", authHeader: "Bearer", setHeader: true, wantStatus: http.StatusUnauthorized},
		{name: "trailing space, empty token", authHeader: "Bearer ", setHeader: true, wantStatus: http.StatusUnauthorized},
		{name: "lowercase scheme is accepted per RFC 6750", authHeader: "bearer correct-token", setHeader: true, wantStatus: http.StatusOK},
		{name: "two spaces before token per RFC 7235", authHeader: "Bearer  correct-token", setHeader: true, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.setHeader {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestBearerAuthDoesNotDistinguishMissingFromWrong(t *testing.T) {
	handler := BearerAuth("correct-token")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	noHeader := httptest.NewRequest(http.MethodGet, "/", nil)
	noHeaderRec := httptest.NewRecorder()
	handler.ServeHTTP(noHeaderRec, noHeader)

	wrongToken := httptest.NewRequest(http.MethodGet, "/", nil)
	wrongToken.Header.Set("Authorization", "Bearer wrong-token")
	wrongTokenRec := httptest.NewRecorder()
	handler.ServeHTTP(wrongTokenRec, wrongToken)

	noHeaderBody, _ := io.ReadAll(noHeaderRec.Result().Body)
	wrongTokenBody, _ := io.ReadAll(wrongTokenRec.Result().Body)

	if string(noHeaderBody) != string(wrongTokenBody) {
		t.Fatalf("response bodies differ: missing=%q wrong=%q; must be identical so callers can't tell which failed", noHeaderBody, wrongTokenBody)
	}
}

func TestBearerAuthPanicsOnEmptyToken(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("BearerAuth(\"\") did not panic; an empty token must fail at construction, not authenticate every request")
		}
	}()
	BearerAuth("")
}
