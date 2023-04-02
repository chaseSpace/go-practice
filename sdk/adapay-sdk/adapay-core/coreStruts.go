package adapayCore

type ApiError struct {

	ErrorCode string `json:"error_code"`

	ErrorMsg string `json:"error_msg"`

	ErrorType string `json:"error_type"`

	Status string `json:"status"`

	InvalidParam string `json:"invalid_param"`
}

type MerchSysConfig struct {

	ApiKeyLive string `json:"api_key_live"`

	ApiKeyTest string `json:"api_key_test"`

	RspPubKey string `json:"rsa_public_key"`

	RspPriKey string `json:"rsa_private_key"`
	AppId string `json:"app_id"`
}

func (msc *MerchSysConfig) IsEmpty() bool {
	return msc.ApiKeyLive == "" || msc.ApiKeyTest == "" || msc.RspPubKey == "" || msc.RspPriKey == ""
}


var GDefaultMerchSysConfig = MerchSysConfig{}


var GMultiMerchSysConfigs = map[string]MerchSysConfig{}
