package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"strings"
)

type memberInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Update(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type Member struct {
	*Adapay
}

func (m *Member) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + MEMBER_CREATE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, m.HandleConfig(multiMerchConfigId...))
}

func (m *Member) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + strings.Replace(MEMBER_QUERY, "{member_id}", adapayCore2.ToString(reqParam["member_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, m.HandleConfig(multiMerchConfigId...))
}

func (m *Member) Update(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + MEMBER_UPDATE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, m.HandleConfig(multiMerchConfigId...))
}

func (m *Member) QueryList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + MEMBER_QUERY_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, m.HandleConfig(multiMerchConfigId...))
}
