package adapay

import (
	"fmt"
	"go-practice/sdk/adapay-sdk/adapay"
	adapayCore "go-practice/sdk/adapay-sdk/adapay-core"
	"log"
	"time"
)

/*
注意：adapay-sdk并不是一个go mod项目，所以不要将其复制到vendor目录中（若你的项目使用vendor）
避免被go mod tidy 或 go mod vendor 命令识别不到而删除
*/
var adapayCli *adapay.Adapay
var adapayDefaultConf *adapayCore.MerchSysConfig
var err error

const timeLayout = "20060102150405"

func Init() {
	// json文件以商户拼音为名字，便于代码使用
	//dir := "./adapay_conf/"
	dir := "../_ignore/" // 与 ./adapay_conf/下的内容一致
	adapayCli, err = adapay.InitMultiMerchSysConfigs(dir)

	if err != nil {
		log.Fatalln("初始化配置失败：", err)
	}
	configFileName := "adapay"
	adapayDefaultConf = adapayCli.MultiMerchSysConfigs[configFileName]
	if adapayDefaultConf == nil {
		panic("上海汇付配置初始化失败，请确认json配置文件名！")
	}
	fmt.Printf("初始化 adapay 成功, %#v\n", adapayDefaultConf)
	// 之后  可从 adapayCli.MultiMerchSysConfigs["adapay"] 获取json配置
}

func RunExample() {
	Init()

	createPaymentParams := make(map[string]interface{})
	createPaymentParams["order_no"] = time.Now().UnixNano() // 字母、数字、下划线的组合
	createPaymentParams["app_id"] = adapayDefaultConf.AppId
	createPaymentParams["pay_channel"] = "alipay_wap"
	createPaymentParams["pay_amt"] = "0.01" // 元，保留2位小数
	createPaymentParams["goods_title"] = "测试标题"
	createPaymentParams["goods_desc"] = "测试描述" // <=42 char
	createPaymentParams["currency"] = "cny"    // cny:ＲＭＢ
	createPaymentParams["time_expire"] = time.Now().Add(time.Hour * 2).Format(timeLayout)
	createPaymentParams["notify_url"] = ""

	// “adapay”是配置文件的名字，否则会找不到配置
	// 创建 Payment对象
	data, apiError, err := adapayCli.Payment().Create(createPaymentParams, "adapay")
	if err != nil || apiError != nil { // 网络或本应用异常
		fmt.Println(err)
		fmt.Println(apiError)
		return
	}

	redirectUrl := ""
	if expend := data["expend"]; expend != nil {
		if _expend, _ := expend.(map[string]interface{}); _expend != nil {
			redirectUrl, _ = _expend["pay_info"].(string)
		}
	}
	fmt.Println(111, redirectUrl)
	fmt.Println(data)
	runCallbackServ()
}
