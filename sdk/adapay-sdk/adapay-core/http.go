package adapayCore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sort"
	"strings"
)

func DeleteEmptyValue(src map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})

	for key, value := range src {
		if key != "" && value != nil && value != "" {
			resultMap[key] = value
		}
	}
	return resultMap
}

func FormatSignSrcText(method string, paramMap map[string]interface{}) (string, error) {
	validParamMap := DeleteEmptyValue(paramMap)

	if strings.EqualFold(method, "GET") {
		keys := make([]string, 0)
		for k := range validParamMap {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		tmpList := make([]string, 0)
		for i := range keys {
			switch tmpValue := validParamMap[keys[i]]; tmpValue.(type) {
			case string:
				tmpList = append(tmpList, fmt.Sprintf("%s=%s", keys[i], tmpValue))
				continue
			case interface{}:
				rs, err := json.Marshal(tmpValue)
				if err != nil || rs == nil || string(rs) == "" {
					continue
				}
				tmpList = append(tmpList, fmt.Sprintf("%s=%s", keys[i], string(rs)))
				continue
			}
		}
		return strings.Join(tmpList, "&"), nil

	} else if strings.EqualFold(method, "POST") {
		postResult, postErr := json.Marshal(validParamMap)
		return string(postResult), postErr

	} else {
		Println("FormatSignSrcText error ... msg: Incorrect method")
		return "", errors.New("Unknow Method: \"" + method + "\"")
	}
}

func FilterApiError(respBodyString string) (*ApiError, error) {

	apiErrorResult := ApiError{}
	jsonError := json.Unmarshal([]byte(respBodyString), &apiErrorResult)
	if jsonError != nil {
		return nil, errors.New(fmt.Sprintf("异常数据 ==> %s", respBodyString))
	}

	if apiErrorResult.ErrorCode != "" {
		return &apiErrorResult, nil
	}

	return nil, nil
}

func HandleResponse(resp *http.Response, msc *MerchSysConfig) (map[string]interface{}, *ApiError, error) {
	resultData := make(map[string]interface{})
	if resp == nil {
		return resultData, nil, errors.New(fmt.Sprintf("网络异常"))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && (body == nil || string(body) == "") {
		return resultData, nil, errors.New(fmt.Sprintf("Adapay 异常应答 ==> status_code<%d>, error_message<%s>",
			resp.StatusCode, err.Error()))
	}

	respDataMap := make(map[string]string)
	tmpJsonErr1 := json.Unmarshal(body, &respDataMap)
	if tmpJsonErr1 != nil {
		return resultData, nil, errors.New(fmt.Sprintf("解析 Adapay 应答数据异常 ==> %s", string(body)))
	}

	err = RsaSignVerify(respDataMap["signature"], respDataMap["data"], msc)
	if err != nil {
		return resultData, nil, errors.New("check signature error !")
	}

	apiErr, err := FilterApiError(respDataMap["data"])
	if err != nil {
		return resultData, nil, errors.New(fmt.Sprintf("解析 Adapay 应答数据异常 ==> %s", err.Error()))
	}
	if apiErr != nil {
		return resultData, apiErr, nil
	}

	jsonErr := json.Unmarshal([]byte(respDataMap["data"]), &resultData)
	if jsonErr != nil {
		return nil, nil, errors.New(fmt.Sprintf("解析返回数据异常 ==> %s", resultData))
	}

	return resultData, nil, nil
}

func DoGetReq(url string, params map[string]interface{}, msc *MerchSysConfig) (*http.Response, error) {

	params = DeleteEmptyValue(params)
	paramText, _ := FormatSignSrcText("GET", params)

	request, _ := http.NewRequest("GET", url+"?"+paramText, nil)
	request.Header.Add("authorization", msc.ApiKeyLive)
	request.Header.Add("sdk_version", "go_v"+Version)
	sign, _ := RsaSign(url+paramText, msc)
	request.Header.Add("signature", sign)

	return http.DefaultClient.Do(request)
}

func DoPostReq(url string, params map[string]interface{}, msc *MerchSysConfig) (*http.Response, error) {

	params = DeleteEmptyValue(params)
	paramText, _ := FormatSignSrcText("POST", params)

	request, _ := http.NewRequest("POST", url, strings.NewReader(paramText))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("authorization", msc.ApiKeyLive)
	request.Header.Add("sdk_version", "go_v"+Version)
	sign, err := RsaSign(url+paramText, msc)
	_ = err
	//log.Println(111, err.Error())
	request.Header.Add("signature", sign)

	return http.DefaultClient.Do(request)
}

func DoUploadFile(url string, params map[string]interface{}, fileParamName string, msc *MerchSysConfig) (*http.Response, error) {

	params = DeleteEmptyValue(params)

	var filepath string
	if str, ok := params[fileParamName].(string); ok {
		filepath = str
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fileParamName, filepath)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	delete(params, fileParamName)

	for key, val := range params {
		if str, ok := val.(string); ok {
			_ = writer.WriteField(key, str)
		}
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", url, body)

	paramText, _ := FormatSignSrcText("GET", params)

	request.Header.Set("content-Type", writer.FormDataContentType())
	request.Header.Add("authorization", msc.ApiKeyLive)
	request.Header.Add("sdk_version", "go_v"+Version)
	sign, _ := RsaSign(url+paramText, msc)
	request.Header.Add("signature", sign)

	return http.DefaultClient.Do(request)
}

func GetMerchSysConfig(mscId string) MerchSysConfig {
	if value, exist := GMultiMerchSysConfigs[mscId]; exist {
		return value
	} else {
		return GDefaultMerchSysConfig
	}
}

type RequestMethod string

const (
	POST   RequestMethod = "POST"
	GET    RequestMethod = "GET"
	UPLOAD RequestMethod = "UPLOAD"
)

func RequestAdaPay(reqUrl string, requestMethod RequestMethod, reqParam map[string]interface{}, msc *MerchSysConfig) (map[string]interface{}, *ApiError, error) {
	data := make(map[string]interface{})

	var resp *http.Response
	var err error

	switch requestMethod {
	case POST:
		resp, err = DoPostReq(reqUrl, reqParam, msc)
		break
	case GET:
		resp, err = DoGetReq(reqUrl, reqParam, msc)
		break
	}

	if err != nil {
		return data, nil, errors.New(fmt.Sprintf("网络异常 ==> %s", err.Error()))
	}

	data, apiErr, err := HandleResponse(resp, msc)
	if apiErr != nil || err != nil {
		return data, apiErr, err
	}

	return data, nil, nil
}

func UploadAdaPay(reqUrl string, reqParam map[string]interface{}, fileParamName string, msc *MerchSysConfig) (map[string]interface{}, *ApiError, error) {
	data := make(map[string]interface{})

	var resp *http.Response
	var err error

	resp, err = DoUploadFile(reqUrl, reqParam, fileParamName, msc)

	if err != nil {
		return data, nil, errors.New(fmt.Sprintf("网络异常 ==> %s", err.Error()))
	}

	data, apiErr, err := HandleResponse(resp, msc)
	if apiErr != nil || err != nil {
		return data, apiErr, err
	}

	return data, nil, nil
}
