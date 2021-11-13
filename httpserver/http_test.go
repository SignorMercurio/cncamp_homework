package httpserver

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func checkStatusCode(got, expect int, t *testing.T) {
	if got != expect {
		t.Errorf("Wrong status code, expect %d, got %d", expect, got)
	}
}

func checkStrEqual(key, got, expect string, t *testing.T) {
	if got != expect {
		t.Errorf("Wrong %s, expect %s, got %s", key, expect, got)
	}
}

func TestInjectResquestHeaders(t *testing.T) {
	r, err := http.NewRequest("GET", "/header", nil)
	key, value := "Custom-Header", "custom_value"
	r.Header.Set(key, value)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	srv := http.HandlerFunc(injectResquestHeaders)
	srv.ServeHTTP(w, r)

	checkStatusCode(w.Code, http.StatusOK, t)
	checkStrEqual("header", w.Header().Get("X-"+key), value, t)
}

func TestGetEnv(t *testing.T) {
	r, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	srv := http.HandlerFunc(getEnv)
	srv.ServeHTTP(w, r)

	checkStatusCode(w.Code, http.StatusOK, t)
	if w.Header().Get("X-Version") == "" {
		t.Error("Invalid version info")
	}
}

func TestWriteLog(t *testing.T) {
	r, err := http.NewRequest("GET", "/log", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	srv := http.HandlerFunc(writeLog)
	srv.ServeHTTP(w, r)

	checkStatusCode(w.Code, http.StatusOK, t)
	l := fmt.Sprintf("Time: %s; IP: %s; Status: %d\n", time.Now().Format(time.UnixDate), r.RemoteAddr, http.StatusOK)
	checkStrEqual("logs", w.Body.String(), l, t)
}

func TestHealthCheck(t *testing.T) {
	r, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	srv := http.HandlerFunc(healthCheck)
	srv.ServeHTTP(w, r)

	checkStatusCode(w.Code, http.StatusOK, t)
}

func TestWelcome(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	srv := http.HandlerFunc(welcome)
	srv.ServeHTTP(w, r)

	checkStatusCode(w.Code, http.StatusOK, t)
	checkStrEqual("response body", w.Body.String(), welcome_msg, t)
}

func TestCatchAll(t *testing.T) {
	r, err := http.NewRequest("GET", "/nowhere", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	srv := http.HandlerFunc(catchAll)
	srv.ServeHTTP(w, r)

	checkStatusCode(w.Code, http.StatusNotFound, t)
	checkStrEqual("response body", w.Body.String(), not_found_msg, t)
}

func TestNewServer(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	NewServer(":8000")
	if !strings.Contains(buf.String(), "Listening on :8000...") {
		t.Error("Failed to start a new server")
	}
}
