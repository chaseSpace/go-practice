package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"strings"
)

type corpMemberInterface interface {
	Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)

	Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error)
}

type CorpMember struct {
	*Adapay
}

func (c *CorpMember) Create(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + CORP_MEMBERS_CREATE

	return adapayCore2.UploadAdaPay(reqUrl, reqParam, "attach_file", c.HandleConfig(multiMerchConfigId...))
}

func (c *CorpMember) Query(reqParam map[string]interface{}, multiMerchConfigId ...string) (map[string]interface{}, *adapayCore2.ApiError, error) {
	reqUrl := BASE_URL + strings.Replace(CORP_MEMBERS_QUERY, "{member_id}", adapayCore2.ToString(reqParam["member_id"]), -1)
	return adapayCore2.RequestAdaPay(reqUrl, adapayCore2.GET, reqParam, c.HandleConfig(multiMerchConfigId...))
}
