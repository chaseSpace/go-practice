package adapay

import (
	"errors"
	"fmt"
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"strings"
)

type settleAccountInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	DeleteAccount(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Detail(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Update(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Balance(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Commission(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	CommissionList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type SettleAccount struct {
	*Adapay
}

func (s *SettleAccount) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + SETTLE_ACCOUNT_CREATE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, s.HandleConfig(multiMerchConfigId...))
}

func (s *SettleAccount) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {

	if reqParam["settle_account_id"] == nil || reqParam["settle_account_id"] == "" {
		return make(map[string]interface{}), nil, errors.New(fmt.Sprintf("请求参数错误 ==> %s", "settle_account_id"))
	}
	reqUrl := BASE_URL + strings.Replace(SETTLE_ACCOUNT_QUERY, "{settle_account_id}", adapayCore2.ToString(reqParam["settle_account_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, s.HandleConfig(multiMerchConfigId...))
}

func (s *SettleAccount) DeleteAccount(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + SETTLE_ACCOUNT_DELETE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, s.HandleConfig(multiMerchConfigId...))
}

func (s *SettleAccount) Detail(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + SETTLE_ACCOUNT_DETAIL_QUERY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, s.HandleConfig(multiMerchConfigId...))
}

func (s *SettleAccount) Update(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + SETTLE_ACCOUNT_MODIFY
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, s.HandleConfig(multiMerchConfigId...))
}

func (s *SettleAccount) Balance(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + SETTLE_ACCOUNT_BALANCE
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, s.HandleConfig(multiMerchConfigId...))
}

func (s *SettleAccount) Commission(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + SETTLE_ACCOUNT_COMMISSIONS
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.POST, reqParam, s.HandleConfig(multiMerchConfigId...))
}

func (s *SettleAccount) CommissionList(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + SETTLE_ACCOUNT_COMMISSIONS_LIST
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, s.HandleConfig(multiMerchConfigId...))
}
