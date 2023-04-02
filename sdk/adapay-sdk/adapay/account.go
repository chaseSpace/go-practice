package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type accountInterface interface {
	Payment(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Account struct {
	*Adapay
}

func (w *Account) Payment(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := PAGE_BASE_URL + WALLET_PAY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, w.HandleConfig(multiMerchConfigId...))
}
