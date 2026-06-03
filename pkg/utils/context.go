package utils

import (
	"context"

	"github.com/google/uuid"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/guid"
)

func ExtractRequestID(ctx context.Context) uuid.UUID {
	if val, ok := ctx.Value(constant.CTX_KEY_REQUEST_ID).(string); ok {
		return uuid.MustParse(val)
	}
	return guid.New()
}

func ExtractLocale(ctx context.Context) string {
	var lang string
	if val, ok := ctx.Value(constant.CTX_KEY_LANGUAGE).(string); ok {
		lang = val
	}
	if lang == "" {
		lang = constant.DEFAULT_LOCALE
	}
	return lang
}

func ExtractHttpUrlPathReq(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_HTTP_URL_PATH_REQ).(string); ok {
		return val
	}
	return ""
}

func ExtractHttpMethodReq(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_HTTP_METHOD_REQ).(string); ok {
		return val
	}
	return ""
}

func ExtractIPAddress(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_CLIENT_IP).(string); ok {
		return val
	}
	return ""
}

func ExtractChannelReq(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_CHANNEL).(string); ok {
		return val
	}
	return ""
}

// func ExtractUserAgent(ctx context.Context) string {
// 	if val, ok := ctx.Value(constant.CTX_KEY_USER_AGENT).(string); ok {
// 		return val
// 	}
// 	return ""
// }
