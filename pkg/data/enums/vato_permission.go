package enums

type VatoActionType int

const (
	VatoActionClient              VatoActionType = 1
	VatoActionDriver              VatoActionType = 2
	VatoActionTrip                VatoActionType = 3
	VatoActionFinancialInfo       VatoActionType = 4
	VatoActionFinancialManagement VatoActionType = 5
	VatoActionAnalytics           VatoActionType = 6
	VatoActionRoot                VatoActionType = 7
	VatoActionBot                 VatoActionType = 8
	VatoActionDocument            VatoActionType = 9
	VatoActionVehicle             VatoActionType = 10
	VatoActionFare                VatoActionType = 47
	VatoActionCustomerManagement  VatoActionType = 49
	VatoActionReview              VatoActionType = 48
	VatoActionOfficeManagement    VatoActionType = 50
	VatoActionManifest            VatoActionType = 51
	VatoActionVehicleModify       VatoActionType = 53
	VatoActionWithdrawOrder       VatoActionType = 55
	VatoActionTopUpOrder          VatoActionType = 54
	VatoActionTripTransaction     VatoActionType = 56
	VatoActionPromotion           VatoActionType = 57
	VatoActionNotification        VatoActionType = 58
	VatoActionExporter            VatoActionType = 59
	VatoActionImporter            VatoActionType = 60
	VatoActionBankInfo            VatoActionType = 61
	VatoActionReferral            VatoActionType = 62
	VatoActionFareReference       VatoActionType = 64
	VatoActionFareAdjustment      VatoActionType = 65
	VatoActionSystemConfig        VatoActionType = 66
	VatoActionTripCoordinator     VatoActionType = 68
	VatoActionUserManagement      VatoActionType = 69
	VatoActionTaxiAnalytics       VatoActionType = 70
	VatoActionTaxiManagement      VatoActionType = 71
	VatoActionBusLineManagement   VatoActionType = 72
	VatoActionBusLineSupport      VatoActionType = 73
	VatoActionFood                VatoActionType = 74
	VatoActionTripCollector       VatoActionType = 75
	VatoActionAppConfigurations   VatoActionType = 76
)

func (a VatoActionType) GetValue() int {
	return int(a)
}

type VatoPermissionType int

const (
	VatoPermissionListing VatoPermissionType = 1
	VatoPermissionDetail  VatoPermissionType = 2
	VatoPermissionUpdate  VatoPermissionType = 4
	VatoPermissionCreate  VatoPermissionType = 8
	VatoPermissionApprove VatoPermissionType = 32
	VatoPermissionDelete  VatoPermissionType = 64
)

func (p VatoPermissionType) GetValue() int {
	return int(p)
}
