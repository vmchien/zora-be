package utils

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"go.opentelemetry.io/otel/trace"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/encode"
)

type RequestInfo struct {
	TenantID  string
	RequestID string
	UserID    string
	DeviceID  string
	Channel   string
	Version   string
	Language  string
	Phone     string
	Email     string
	ClientIP  string
	FullName  string
	UserToken string
	TokenType string
}

func GetRequestInfoMD(ctx context.Context) *RequestInfo {
	if md, ok := metadata.FromServerContext(ctx); ok {
		fullNameDecoded, _ := encode.DecodeBase64(md.Get(constant.CTX_KEY_MD_USER_FULLNAME))
		return &RequestInfo{
			TenantID:  md.Get(constant.CTX_KEY_MD_TENANT_ID),
			RequestID: md.Get(constant.CTX_KEY_MD_REQUEST_ID),
			UserID:    TryGetUserID(md.Get(constant.CTX_KEY_MD_USER_ID)),
			DeviceID:  md.Get(constant.CTX_KEY_MD_DEVICE_ID),
			Channel:   md.Get(constant.CTX_KEY_MD_CHANNEL),
			Version:   md.Get(constant.CTX_KEY_MD_VERSION),
			Language:  md.Get(constant.CTX_KEY_MD_LANGUAGE),
			Phone:     md.Get(constant.CTX_KEY_MD_USER_PHONE),
			Email:     md.Get(constant.CTX_KEY_MD_USER_EMAIL),
			ClientIP:  md.Get(constant.CTX_KEY_MD_CLIENT_IP),
			FullName:  fullNameDecoded,
			UserToken: md.Get(constant.CTX_KEY_MD_USER_TOKEN),
			TokenType: md.Get(constant.CTX_KEY_MD_TOKENT_TYPE),
		}
	}
	return nil
}

func GetRequestInfo(ctx context.Context) *RequestInfo {
	return &RequestInfo{
		UserID:   TryGetUserIDFromContext(ctx),
		FullName: TryGetUserFullNameFromContext(ctx),
		Email:    TryGetUseEmailFromContext(ctx),
		Phone:    TryGetPhoneNumberFromContext(ctx),
	}
}

func TryGetUserIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_USER_ID).(string); ok {
		return val
	}
	return constant.AnonymousUserID
}

func TryGetClientIPFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_CLIENT_IP).(string); ok {
		return val
	}
	return ""
}

func TryGetUserAgentFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_USER_AGENT).(string); ok {
		return val
	}
	return ""
}

func TryGetPhoneNumberFromContext(ctx context.Context) string {
	if phone, ok := ctx.Value(constant.CTX_KEY_USER_PHONE).(string); ok {
		return phone
	}
	return ""
}

func TryGetUserFullNameFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_USER_FULLNAME).(string); ok {
		return val
	}
	return constant.AnonymousUserName
}

func TryGetUseEmailFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_USER_EMAIL).(string); ok {
		return val
	}
	return ""
}

func TryGetUserTokenFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(constant.CTX_KEY_USER_TOKEN).(string); ok {
		return val
	}
	return ""
}

func TryGetUserID(id string) string {
	if strings.Trim(id, " ") == "" {
		return constant.AnonymousUserID
	}
	return id
}

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(constant.CTX_KEY_LOG_TRACE_ID).(string); ok && traceID != "" {
		return traceID
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return spanCtx.TraceID().String()
	}
	return ""
}
