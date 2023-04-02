package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type walletInterface interface {
	Login(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Wallet struct {
	*Adapay
}

func (w *Wallet) Login(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := PAGE_BASE_URL + WALLET_LOGIN
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, w.HandleConfig(multiMerchConfigId...))
}
