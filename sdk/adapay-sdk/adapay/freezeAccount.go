package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type FreezeAccountInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
	List(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type FreezeAccount struct {
	*Adapay
}

func (f *FreezeAccount) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + FREEZE_ACCOUNT_FREEZE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, f.HandleConfig(multiMerchConfigId...))
}

func (f *FreezeAccount) List(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + FREEZE_ACCOUNT_FREEZE_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, f.HandleConfig(multiMerchConfigId...))
}
