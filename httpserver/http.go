package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(addr string) error {
	r := mux.NewRouter()

	r.HandleFunc("/header", injectResquestHeaders)
	r.HandleFunc("/version", getEnv)
	r.HandleFunc("/log", writeLog)
	r.HandleFunc("/healthz", healthCheck)
	r.HandleFunc("/", welcome)
	r.PathPrefix("/").HandlerFunc(catchAll)

	log.Printf("Listening on %s...", addr)

	return http.ListenAndServe(addr, r)
}

func injectResquestHeaders(w http.ResponseWriter, r *http.Request) {
	for key, value := range r.Header {
		w.Header().Set("X-"+key, strings.Join(value, "; "))
	}
	w.Write([]byte(`Check the Request Headers in Responses Headers (those with the "X" prefix)`))
}

func getEnv(w http.ResponseWriter, r *http.Request) {
	v := os.Getenv("VERSION")
	if v == "" {
		v = "0.0.0"
	}
	w.Header().Set("X-Version", v)
	w.Write([]byte(`Check the version in Response Headers`))
}

func writeLog(w http.ResponseWriter, r *http.Request) {
	l := fmt.Sprintf("Time: %s; IP: %s; Status: %d\n", time.Now().Format(time.UnixDate), r.RemoteAddr, http.StatusOK)
	log.Println(l)
	w.Write([]byte(l))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Awesome"))
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
	Welcome to cncamp_homework HTTP Server:

	- Access /header to find your Request Headers in the Response Headers
	- Access /version to get the VERSION environment variable
	- Access /log to write logs in the server
	- Access /healthz for a health check
	`))
}

func catchAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Path not found"))
}
