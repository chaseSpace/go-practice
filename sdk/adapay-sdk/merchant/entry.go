package merchant

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type entryInterface interface {
	BatchEntrys(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	QueryEntry(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Entry struct {
	*Merchant
}

func (e *Entry) BatchEntrys(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + BATCH_ENTRYS
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, e.HandleConfig(multiMerchConfigId...))
}

func (e *Entry) QueryEntry(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + QUERY_ENTRY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, e.HandleConfig(multiMerchConfigId...))
}
