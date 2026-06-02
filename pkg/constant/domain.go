package constant

import "github.com/google/uuid"

// Domain constants for schema categorization
const (
	VATO                   = "vato"
	TanQuangDung           = "tqd"
	VATO_TENANT_ID         = "00000000-0000-7000-0000-000000000001"
	TanQuangDung_TENANT_ID = "00000000-0000-7000-0000-000000000002"
)

var DefaultTenantID = uuid.MustParse(VATO_TENANT_ID)
var VatoTenantID = uuid.MustParse(VATO_TENANT_ID)
var TanQuangDungTenantID = uuid.MustParse(TanQuangDung_TENANT_ID)

var mapDomainTenant = map[string]uuid.UUID{
	VATO:         VatoTenantID,
	TanQuangDung: TanQuangDungTenantID,
}

func GetDomain(tenantId uuid.UUID) string {
	for k, v := range mapDomainTenant {
		if v == tenantId {
			return k
		}
	}
	return ""
}
