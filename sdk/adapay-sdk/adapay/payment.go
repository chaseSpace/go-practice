package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"strings"
)

type paymentInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Close(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Payment struct {
	*Adapay
}

func (p *Payment) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + PAYMENT_CREATE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, p.HandleConfig(multiMerchConfigId...))
}

func (p *Payment) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + strings.Replace(PAYMENT_QUERY, "{payment_id}", adapayCore2.ToString(reqParam["payment_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, p.HandleConfig(multiMerchConfigId...))
}

func (p *Payment) QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + PAYMENT_LIST_QUERY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, p.HandleConfig(multiMerchConfigId...))
}

func (p *Payment) Close(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + strings.Replace(PAYMENT_CLOSE, "{payment_id}", adapayCore2.ToString(reqParam["payment_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, p.HandleConfig(multiMerchConfigId...))
}
