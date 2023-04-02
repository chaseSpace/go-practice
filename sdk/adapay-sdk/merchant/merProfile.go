package merchant

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type merProfileInterface interface {
	MerProfilePicture(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	MerProfileForAudit(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type MerProfile struct {
	*Merchant
}

func (e *MerProfile) MerProfilePicture(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + MER_PROFILE_PICTURE

	return adapayCore2.UploadAdaPay(reqUrl, reqParam, "file", e.HandleConfig(multiMerchConfigId...))
}

func (e *MerProfile) MerProfileForAudit(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + MER_PROFILE_FOR_AUDIT
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, e.HandleConfig(multiMerchConfigId...))
}
