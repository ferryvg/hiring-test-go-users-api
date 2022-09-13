package http

import "github.com/valyala/fasthttp"

// Interface implemented by endpoint middlewares
type Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

// Build endpoint from provided handler and middlewares
func BuildHandler(handler fasthttp.RequestHandler, middlewares ...Middleware) fasthttp.RequestHandler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
