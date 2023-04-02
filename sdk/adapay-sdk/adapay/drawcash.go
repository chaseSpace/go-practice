package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type drawcashInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Drawcash struct {
	*Adapay
}

func (c *Drawcash) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + CREATE_CASHS
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, c.HandleConfig(multiMerchConfigId...))
}

func (c *Drawcash) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + QUERY_CASHS_STAT
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, c.HandleConfig(multiMerchConfigId...))
}
