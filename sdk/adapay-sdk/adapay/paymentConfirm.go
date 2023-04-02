package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"strings"
)

type paymentConfirmInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type PaymentConfirm struct {
	*Adapay
}

func (p *PaymentConfirm) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + PAYMENT_CONFIRM
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, p.HandleConfig(multiMerchConfigId...))
}

func (p *PaymentConfirm) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + strings.Replace(PAYMENT_QUERY_CONFIRM, "{payment_confirm_id}", adapayCore2.ToString(reqParam["payment_confirm_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, p.HandleConfig(multiMerchConfigId...))
}

func (p *PaymentConfirm) QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + PAYMENT_QUERY_CONFIRM_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, p.HandleConfig(multiMerchConfigId...))
}
