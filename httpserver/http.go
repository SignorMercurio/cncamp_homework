package httpserver

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SignorMercurio/cncamp_homework/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	welcome_msg = `
	Welcome to cncamp_homework HTTP Server:

	- Access /header to find your Request Headers in the Response Headers
	- Access /version to get the VERSION environment variable
	- Access /log to write logs in the server
	- Access /healthz for a health check
	`
	not_found_msg = "Path not found"
)

func NewServer(addr string) *http.Server {
	metrics.Register()
	r := mux.NewRouter()

	r.HandleFunc("/header", injectResquestHeaders)
	r.HandleFunc("/version", getEnv)
	r.HandleFunc("/log", writeLog)
	r.HandleFunc("/healthz", healthCheck)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", welcome)
	r.PathPrefix("/").HandlerFunc(catchAll)
	r.Use(metricsMiddleware)
	r.Use(loggingMiddleware)

	zap.S().Infow("Start Listening...", "address", addr)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := metrics.NewTimer()
		defer timer.ObserveTotal()
		delay := randInt(10, 2000)
		time.Sleep(time.Millisecond * time.Duration(delay))
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.S().Infow(
			"Logged client access",
			"Time",
			time.Now().Format(time.UnixDate),
			"Method",
			r.Method,
			"RequestURI",
			r.RequestURI,
			"IP",
			r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
	})
}

func injectResquestHeaders(w http.ResponseWriter, r *http.Request) {
	for key, value := range r.Header {
		w.Header().Set("X-"+key, strings.Join(value, "; "))
	}
	w.Write([]byte(`Check the Request Headers in Responses Headers (those with the "X" prefix)`))
}

func getEnv(w http.ResponseWriter, r *http.Request) {
	v := os.Getenv("VERSION")
	zap.S().Debugw("Read VERSION from env", "version", v)
	if v == "" {
		v = "0.0.0"
	}
	w.Header().Set("X-Version", v)
	w.Write([]byte(`Check the version in Response Headers`))
}

func writeLog(w http.ResponseWriter, r *http.Request) {
	l := fmt.Sprintf("Time: %s\nMethod: %s\nRequestURI: %s\nIP: %s\nStatus: %d", time.Now().Format(time.UnixDate), r.Method, r.RequestURI, r.RemoteAddr, http.StatusOK)
	zap.S().Debugw("Composed the current log", "log", l)
	w.Write([]byte(l))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Awesome"))
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(welcome_msg))
}

func catchAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(not_found_msg))
}
