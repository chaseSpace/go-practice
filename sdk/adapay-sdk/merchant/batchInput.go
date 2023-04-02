package merchant

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type batchInputInterface interface {
	MerConf(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	QueryMerConf(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	MerResidentModify(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type BatchInput struct {
	*Merchant
}

func (b *BatchInput) MerConf(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + MER_CONF
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, b.HandleConfig(multiMerchConfigId...))
}

func (b *BatchInput) QueryMerConf(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + QUERY_MER_CONF
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, b.HandleConfig(multiMerchConfigId...))
}

func (b *BatchInput) MerResidentModify(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + MER_RESIDENT_MODIFY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, b.HandleConfig(multiMerchConfigId...))
}
