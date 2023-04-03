package adapay

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONData 回调结构
type adapayCallback struct {
	CreatedTime string `json:"created_time"` // 20200103105148
	Expend      Expend `json:"expend"`
	ID          string `json:"id"`       // 支付商订单号
	OrderNo     string `json:"order_no"` // 我方订单号
	PayAmt      string `json:"pay_amt"`  // 0.01
	PayChannel  string `json:"pay_channel"`
	Status      string `json:"status"` // succeeded
}
type Expend struct {
	BankType     string `json:"bank_type"`
	OpenID       string `json:"open_id"`
	SubOpenID    string `json:"sub_open_id"`
	BuyerLogonID string `json:"buyer_logon_id"`
}

// 服务器为POST回调，默认超时时间为5秒，超时后会重试3次；不支持HTTP重定向；服务器对应答不是200~300之间的错误，会默认重试3次；
// 异步通知服务器对HTTPS不认证验签和ALLOW_ALL_HOSTNAME_VERIFIER；如商户自定义通知端口，请使用8000-9005内端口；
// URL上请勿附带参数；异步回调请求编码集为：UTF-8。 异步回调参数为：Event对象，以key=value方式POST发送数据
func runCallbackServ() {
	http.HandleFunc("/adapay", func(w http.ResponseWriter, r *http.Request) {
		var err error
		data := r.FormValue("data")
		err = adapayCli.AdapayTools().VerifySign(data, r.FormValue("sign"))
		if err != nil {
			return
		}
		callbackV := new(adapayCallback)
		err = json.Unmarshal([]byte(data), callbackV)
		if err != nil {
			return
		}

		// log

		_ = callbackV.OrderNo
		_ = callbackV.ID
	})

	log.Println(http.ListenAndServe(":8080", nil))

}
