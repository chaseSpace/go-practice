package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type UnFreezeAccountInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
	List(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type UnFreezeAccount struct {
	*Adapay
}

func (f *UnFreezeAccount) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + UnFREEZE_ACCOUNT_FREEZE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, f.HandleConfig(multiMerchConfigId...))
}

func (f *UnFreezeAccount) List(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + UnFREEZE_ACCOUNT_FREEZE_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, f.HandleConfig(multiMerchConfigId...))
}
