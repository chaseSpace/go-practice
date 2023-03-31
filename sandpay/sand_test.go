package sandpay

import (
	"log"
	"net/http"
	"time"
)

func Example() {
	var (
		merchantNo string
		privKey    string
		sandPubKey string
	)

	_, err := InitClient(merchantNo, privKey, sandPubKey)
	if err != nil {
		panic(err)
	}

	//func BuildPayUrlForH5Alipay() (string, error) {
	var (
		orderNo                                   string
		amount                                    float64
		notifyUrl, returnUrl, createIp, goodsName string
		createAt                                  time.Time
		hint                                      string
	)
	orderUrl, err := BuildPayUrlForH5Alipay(orderNo, amount, notifyUrl, returnUrl, createIp, goodsName, createAt, hint)
	if err != nil {
		log.Panicln("BuildPayUrlForH5Alipay", err)
	}

	_ = orderUrl
	/*
		orderUrl是sand下单地址加上我们的下单参数，我们的后台把这个地址传给前端H5，再由前端H5去请求该地址，以调起支付宝H5支付页，并由用户完成支付
	*/

	// 关于回调：
	// -	我们的后台服务启动一个HTTP服务器，提供一个API给sand调用，例如http://our_host:9999/sandpay_callback
	http.HandleFunc("/sandpay_callback", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400) // 先写个400，后面成功改为200
		if err := r.ParseForm(); err != nil {
			log.Println("SandpayCallback ParseForm", err)
			return
		}
		// 1. 提取数据&验签
		sign := r.FormValue("sign")
		if sign == "" {
			log.Println("SandpayCallback sign EMPTY!!!")
			return
		}
		// 2. 解析出sand给我们的支付结果，其中包含我们的之前给的单号
		dataObj, err := ParseCallbackJSON(r.FormValue("data"), sign)
		if err != nil {
			log.Println("SandpayCallback ParseCallbackJSON", err)
			return
		}

		_ = dataObj.Body.OrderCode    // 我们的单号
		_ = dataObj.Body.PayOrderCode // sand的单号 如：20230321001343010000000000144686

		// 下面逻辑
		/*
			1. 使用orderId查出我们的订单
			2. 更新我们的订单状态（成功/失败）
			3. 执行状态逻辑
			4. 若我们的逻辑执行成功，则返回ACK，否则return
		*/
		AckCallback(w)
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
