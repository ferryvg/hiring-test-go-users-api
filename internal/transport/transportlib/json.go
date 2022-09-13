package transportlib

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func ResponseJsonStatus(ctx *fasthttp.RequestCtx, msg string, statusCode int) error {
	data := map[string]interface{}{
		"code":    statusCode,
		"message": msg,
	}

	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	ctx.Response.Reset()
	ctx.SetStatusCode(statusCode)
	ctx.SetContentTypeBytes([]byte("application/json"))
	ctx.SetBody(payload)

	return nil
}
