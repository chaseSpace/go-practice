package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type adapayToolsInterface interface {
	DownloadBill(bill_date string, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	UserIdentity(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	VerifySign(signData string, originalData string, multiMerchConfigId ...string) error
}

type AdapayTools struct {
	*Adapay
}

func (b *AdapayTools) DownloadBill(bill_date string, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + BILL_DOWNLOAD

	reqParam := make(map[string]interface{})
	reqParam["bill_date"] = bill_date

	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, b.HandleConfig(multiMerchConfigId...))
}

func (u *AdapayTools) UserIdentity(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + USER_IDENTITY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, u.HandleConfig(multiMerchConfigId...))
}

func (v *AdapayTools) VerifySign(signData string, originalData string, multiMerchConfigId ...string) error {

	return adapayCore2.RsaSignVerify(signData, originalData, v.HandleConfig(multiMerchConfigId...))
}
