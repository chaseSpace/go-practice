package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type checkoutInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Checkout struct {
	*Adapay
}

func (w *Checkout) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := PAGE_BASE_URL + WALLET_CHECKOUT
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, w.HandleConfig(multiMerchConfigId...))
}

func (m *Checkout) QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + QUERY_CHECKOUT_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, m.HandleConfig(multiMerchConfigId...))
}
