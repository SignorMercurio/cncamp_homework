package httpserver

import (
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

func NewServer(addr string) error {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/reqHdr":
			injectResquestHeaders(ctx)
		case "/getEnv":
			getEnv(ctx)
		case "/log":
			writeLog(ctx)
		case "/healthz":
			healthCheck(ctx)
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
	log.Printf("IP: %s; Status: %d\n", ctx.RemoteAddr().String(), ctx.Response.StatusCode())
	ctx.Response.SetBody([]byte(`This access attempt has been logged.`))
}

func healthCheck(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBody([]byte("Awesome"))
}
