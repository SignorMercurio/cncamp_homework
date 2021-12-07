package httpserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SignorMercurio/cncamp_homework/logger"
	"go.uber.org/zap"
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
	l := fmt.Sprintf("Time: %s\nMethod: %s\nRequestURI: %s\nIP: %s\nStatus: %d", time.Now().Format(time.UnixDate), r.Method, r.RequestURI, r.RemoteAddr, http.StatusOK)
	checkStrEqual("logs", w.Body.String(), l, t)
}

func TestHealthCheckWithMiddlewares(t *testing.T) {
	r, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	srv := metricsMiddleware(loggingMiddleware(http.HandlerFunc(healthCheck)))
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
	// Init logger
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	NewServer(":8000")

	buf, err := ioutil.ReadFile("httpserver.log")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(buf), `"msg":"Start Listening...","address":":8000"`) {
		t.Error("Failed to start a new server")
	}
}
