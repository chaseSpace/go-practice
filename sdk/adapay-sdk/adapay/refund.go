package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"strings"
)

type refundInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Refund struct {
	*Adapay
}

func (r *Refund) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + strings.Replace(REFUND_CREATE, "{payment_id}", adapayCore2.ToString(reqParam["payment_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, r.HandleConfig(multiMerchConfigId...))
}

func (r *Refund) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + REFUND_QUERY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, r.HandleConfig(multiMerchConfigId...))
}
