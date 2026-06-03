package constant

// HTTP Headers -> GRPC Metadata
const (
	CTX_KEY_TENANT_ID         = "x-tenant-id"
	CTX_KEY_REQUEST_ID        = "x-request-id"
	CTX_KEY_LANGUAGE          = "x-lang"
	CTX_KEY_USER_AGENT        = "x-user-agent"
	CTX_KEY_CLIENT_IP         = "x-client-ip"
	CTX_KEY_HTTP_URL_PATH_REQ = "x-http-url-path"
	CTX_KEY_HTTP_METHOD_REQ   = "x-http-method"
	CTX_KEY_DEVICE_ID         = "x-device-id"
	CTX_KEY_USER_ID           = "x-user-id"
	CTX_KEY_CHANNEL           = "x-channel"
	CTX_KEY_VERSION           = "x-version"
	CTX_KEY_USER_PHONE        = "x-user-phone"
	CTX_KEY_USER_EMAIL        = "x-user-email"
	CTX_KEY_USER_FULLNAME     = "x-user-fullname"
	CTX_KEY_DOMAIN            = "x-domain"

	CTX_KEY_MD_TENANT_ID = "x-md-tenant-id"
	CTX_KEY_MD_LANGUAGE  = "x-md-lang"
	// CTX_KEY_MD_USER_AGENT        = "x-md-user-agent"
	CTX_KEY_MD_CLIENT_IP         = "x-md-client-ip"
	CTX_KEY_MD_REQUEST_ID        = "x-md-request-id"
	CTX_KEY_MD_HTTP_URL_PATH_REQ = "x-md-http-url-path"
	CTX_KEY_MD_HTTP_METHOD_REQ   = "x-md-http-method"
	CTX_KEY_MD_DEVICE_ID         = "x-md-device-id"
	CTX_KEY_MD_USER_ID           = "x-md-user-id"
	CTX_KEY_MD_CHANNEL           = "x-md-channel"
	CTX_KEY_MD_VERSION           = "x-md-version"
	CTX_KEY_MD_USER_PHONE        = "x-md-user-phone"
	CTX_KEY_MD_USER_EMAIL        = "x-md-user-email"
	CTX_KEY_MD_USER_FULLNAME     = "x-md-user-fullname"
	CTX_KEY_MD_DOMAIN            = "x-md-domain"
	CTX_KEY_MD_TOKENT_TYPE       = "x-md-token-type"

	// TODO: remove token
	CTX_KEY_USER_TOKEN    = "x-user-token"
	CTX_KEY_MD_USER_TOKEN = "x-md-usertoken"
	CTX_KEY_TOKENT_TYPE   = "x-token-type"

	CTX_KEY_SERIAL_NUMBER = "serial-number"
)

const (
	CTX_KEY_LOG_TRACE_ID  = "trace.id"
	CTX_KEY_LOG_SPAN_ID   = "span.id"
	CTX_KEY_LOG_DOMAIN_ID = "domain.id"
	CTX_KEY_LOGGER        = "logger"
)

const (
	CTX_KEY_MD_CLAIMS_USER   = "x-md-nova-claims-userid"
	CTX_KEY_MD_CLAIMS_TENANT = "x-md-nova-claims-tenantid"
)
const (
	CTX_KEY_MD_ERROR_BIZ_CODE       = "x-md-nova-error-biz-code"
	CTX_KEY_MD_ERROR_BIZ_MESSAGE    = "x-md-nova-error-biz-msg"
	CTX_KEY_MD_ERROR_BIZ_SPECS      = "x-md-nova-error-biz-specs"
	CTX_KEY_MD_ERROR_HTTP_STATUS    = "x-md-nova-error-http-status"
	CTX_KEY_MD_ERROR_INTERNAL_SPECS = "x-md-nova-error-internal-specs"
)
