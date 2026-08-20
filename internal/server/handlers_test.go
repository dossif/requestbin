package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dossif/requestbin/internal/config"
)

func newTestHandler() handler {
	return handler{Service: config.Config{Log: slog.New(slog.NewTextHandler(io.Discard, nil))}}
}

func TestHandlerRequestStatusDefault(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.handlerRequestStatus(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHandlerRequestStatusExplicit(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/status/{status}", h.handlerRequestStatus)
	req := httptest.NewRequest(http.MethodGet, "/status/201", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 201 {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
}

func TestHandlerRequestStatusInvalid(t *testing.T) {
	cases := []string{"abc", "100", "999"}
	for _, status := range cases {
		t.Run(status, func(t *testing.T) {
			h := newTestHandler()
			mux := http.NewServeMux()
			mux.HandleFunc("/status/{status}", h.handlerRequestStatus)
			req := httptest.NewRequest(http.MethodGet, "/status/"+status, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != 599 {
				t.Fatalf("status = %d, want 599", rec.Code)
			}
		})
	}
}

func TestHandlerRequestStatusBody(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/status/{status}", h.handlerRequestStatus)
	req := httptest.NewRequest(http.MethodGet, "/status/201?foo=bar&foo=baz", nil)
	req.Header.Set("Cookie", "session=abc123")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var got struct {
		Request struct {
			Path    string              `json:"Path"`
			Pattern string              `json:"Pattern"`
			Query   map[string][]string `json:"Query"`
			Cookies map[string]string   `json:"Cookies"`
		} `json:"Request"`
		Response struct {
			Status int `json:"Status"`
		} `json:"Response"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if got.Response.Status != 201 {
		t.Errorf("Response.Status = %d, want 201", got.Response.Status)
	}
	if got.Request.Path != "/status/201" {
		t.Errorf("Request.Path = %q, want /status/201", got.Request.Path)
	}
	if got.Request.Pattern != "/status/{status}" {
		t.Errorf("Request.Pattern = %q, want /status/{status}", got.Request.Pattern)
	}
	if want := []string{"bar", "baz"}; len(got.Request.Query["foo"]) != len(want) {
		t.Errorf("Request.Query[foo] = %v, want %v", got.Request.Query["foo"], want)
	}
	if got.Request.Cookies["session"] != "abc123" {
		t.Errorf("Request.Cookies[session] = %q, want abc123", got.Request.Cookies["session"])
	}
}

func TestHandlerFavicon(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	rec := httptest.NewRecorder()
	h.handlerFavicon(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/x-icon" {
		t.Errorf("Content-Type = %q, want image/x-icon", ct)
	}
	if rec.Body.Len() == 0 {
		t.Error("favicon body is empty")
	}
}
