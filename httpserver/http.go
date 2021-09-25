package httpserver

import (
	"fmt"
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

func NewServer(addr string) error {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/header":
			injectResquestHeaders(ctx)
		case "/version":
			getEnv(ctx)
		case "/log":
			writeLog(ctx)
		case "/healthz":
			healthCheck(ctx)
		case "/":
			welcome(ctx)
		default:
			ctx.Error("Path not found", fasthttp.StatusNotFound)
		}
	}
	log.Printf("Listening on %s...", addr)

	return fasthttp.ListenAndServe(addr, requestHandler)
}

func injectResquestHeaders(ctx *fasthttp.RequestCtx) {
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.Set("X-"+string(key), string(value))
	})
	ctx.Response.SetBody([]byte(`Check the Request Headers in Responses Headers (those with the "X" prefix)`))
}

func getEnv(ctx *fasthttp.RequestCtx) {
	v := os.Getenv("VERSION")
	if v == "" {
		v = "0.0.0"
	}
	ctx.Response.Header.Set("X-Version", v)
	ctx.Response.SetBody([]byte(`Check the version in Response Headers`))
}

func writeLog(ctx *fasthttp.RequestCtx) {
	l := fmt.Sprintf("IP: %s; Status: %d\n", ctx.RemoteAddr().String(), ctx.Response.StatusCode())
	log.Println(l)
	ctx.Response.SetBody([]byte(l))
}

func healthCheck(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBody([]byte("Awesome"))
}

func welcome(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBody([]byte(`
	Welcome to cncamp_homework HTTP Server:
	
	- Access /header to find your Request Headers in the Response Headers
	- Access /version to get the VERSION environment variable
	- Access /log to write logs in the server
	- Access /healthz for a health check
	`))
}
