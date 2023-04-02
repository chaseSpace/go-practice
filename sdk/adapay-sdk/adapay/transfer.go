package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type TransferInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
	List(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Transfer struct {
	*Adapay
}

func (t *Transfer) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + TRANSFER_CREATE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, t.HandleConfig(multiMerchConfigId...))
}

func (t *Transfer) List(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + TRANSFER_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, t.HandleConfig(multiMerchConfigId...))
}
