package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"strings"
)

type paymentReverseInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type PaymentReverse struct {
	*Adapay
}

func (p *PaymentReverse) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + PAYMENT_REVERSE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, p.HandleConfig(multiMerchConfigId...))
}

func (p *PaymentReverse) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + strings.Replace(PAYMENT_QUERY_REVERSE, "{reverse_id}", adapayCore2.ToString(reqParam["reverse_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, p.HandleConfig(multiMerchConfigId...))
}

func (p *PaymentReverse) QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + PAYMENT_QUERY_REVERSE_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, p.HandleConfig(multiMerchConfigId...))
}
