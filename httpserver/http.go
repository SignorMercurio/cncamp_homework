package httpserver

import (
	"fmt"
	"log"
	"net/http"
)

func NewServer(addr string) error {
	http.HandleFunc("/", http.HandlerFunc(defaultHandler))

	log.Printf("Listening on %s...", addr)

	return http.ListenAndServe(addr, nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Here are your HTTP Request Headers:\n\n")
	printHeaders(w, r.Header)

	fmt.Fprint(w, "Here are your HTTP Request Headers:\n\n")
	printHeaders(w, r.Response.Header)
}

func printHeaders(w http.ResponseWriter, h http.Header) {
	for name, headers := range h {
		for _, header := range headers {
			fmt.Fprintf(w, "%s: %s\n", name, header)
		}
	}
}
