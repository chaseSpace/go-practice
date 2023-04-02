package adapayCore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func ToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	default:
		r, _ := json.Marshal(v)
		return string(r)
	}
}


func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}

func ReadMerchConfig(configPath string) (*MerchSysConfig, error) {

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("ReadMerchConfig Read Config File Failed ---> path: %s", configPath))
	}


	tmpMsc := MerchSysConfig{}
	jsonError := json.Unmarshal(data, &tmpMsc)
	if jsonError != nil {
		return nil, errors.New(
			fmt.Sprintf("ReadMerchConfig Read Config File Failed ---> path: %s", configPath))
	}


	if tmpMsc.ApiKeyLive == "" || tmpMsc.ApiKeyTest == "" ||
		tmpMsc.RspPubKey == "" || tmpMsc.RspPriKey == "" {
		return nil, errors.New(
			fmt.Sprintf("ReadMerchConfig Read Config File Failed ---> path: %s", configPath))
	}

	return &tmpMsc, nil
}
