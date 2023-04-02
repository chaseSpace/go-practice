package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
)

type fastPayInterface interface {
	CardBind(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
	CardBindConfirm(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
	CardBindlist(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
	Confirm(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
	SmsCode(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type FastPay struct {
	*Adapay
}

func (f *FastPay) CardBind(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := PAGE_BASE_URL + FAST_CARD_APPLY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, f.HandleConfig(multiMerchConfigId...))
}

func (f *FastPay) CardBindConfirm(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := PAGE_BASE_URL + FAST_CARD_CONFIRM
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, f.HandleConfig(multiMerchConfigId...))
}

func (f *FastPay) CardBindlist(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := PAGE_BASE_URL + FAST_CARD_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, f.HandleConfig(multiMerchConfigId...))
}

func (f *FastPay) Confirm(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + FAST_PAY_CONFIRM
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, f.HandleConfig(multiMerchConfigId...))
}

func (f *FastPay) SmsCode(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + FAST_PAY_SMS_CODE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, f.HandleConfig(multiMerchConfigId...))
}
