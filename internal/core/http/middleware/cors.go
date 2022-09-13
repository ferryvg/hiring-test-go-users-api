package middleware

import "github.com/valyala/fasthttp"

func CorsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		origin := string(ctx.Request.Header.Peek(fasthttp.HeaderOrigin))

		if origin != "" {
			ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, origin)
			ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, "true")
		}

		if string(ctx.Request.Header.Method()) == fasthttp.MethodOptions {
			return
		}

		next(ctx)
	}
}
