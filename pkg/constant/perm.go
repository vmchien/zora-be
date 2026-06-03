package constant

import "time"

const (
	ROUTE_OVERRIDE_TENANT_KEY_PREFIX_CACHE = "cfg:route-overrides:"
	ROUTE_OVERRIDE_TENANT_EXPIRES_DURATION = 15 * time.Minute
)

const (
	RBAC_ROLE_TENANT_KEY_PREFIX_CACHE = "data:rbac-roles:"
	RBAC_ROLE_TENANT_EXPIRES_DURATION = 24 * 60 * time.Minute
)

const (
	RBAC_USER_ROLE_TENANT_KEY_CACHE        = "data:rbac-user-roles:"
	RBAC_USER_ROLE_TENANT_EXPIRES_DURATION = DEFAULT_JWT_DURATION
)

const (
	PBAC_POLICY_TENANT_KEY_PREFIX_CACHE = "data:pbac-policies:"
	PBAC_POLICY_TENANT_EXPIRES_DURATION = 24 * 60 * time.Minute
)

const (
	// Cache key prefix for holiday ticket rules per tenant (per year)
	HOLIDAY_TICKET_RULE_TENANT_KEY_PREFIX_CACHE = "vato-buslines:holiday-ticket-rule:"
	HOLIDAY_TICKET_RULE_TENANT_EXPIRES_DURATION = 365 * 24 * time.Hour
)

const (
	OLD_SYSTEM_BOOKING_CACHE_KEY_PREFIX_CACHE = "pmv-booking-id:"
	OLD_SYSTEM_BOOKING_CACHE_EXPIRES_DURATION = 60 * time.Minute
)

const (
	SEAT_HOLDING_CACHE_KEY_PREFIX_CACHE = "PartnerTrip:"
)
