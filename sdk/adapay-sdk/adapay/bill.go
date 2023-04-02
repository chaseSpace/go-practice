package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type billInterface interface {
	Download(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Bill struct {
	*Adapay
}

func (b *Bill) Download(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + BILL_DOWNLOAD
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, b.HandleConfig(multiMerchConfigId...))
}
