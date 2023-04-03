package sandpay

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smartwalle/crypto4go"
	"net/http"
	"net/url"
	"time"
)

const (
	// H5包装支付宝生活号-02020002
	api_H5Alipay = "https://sandcash.mixienet.com.cn/pay/h5/alipay"
)

var (
	ErrLoadPrivateKey = errors.New("sandpay: load private key failed")
	ErrLoadPublicKey  = errors.New("sandpay: load public key failed")
)

type Client struct {
	merchantNo    string
	privateKey    *rsa.PrivateKey
	sandPublicKey *rsa.PublicKey
}

var defaultClient *Client

func InitClient(merchantNo, priKey, sandPubKey string) (*Client, error) {
	_priKey, err := crypto4go.ParsePKCS1PrivateKey(crypto4go.FormatPKCS1PrivateKey(priKey))
	if err != nil {
		_priKey, err = crypto4go.ParsePKCS8PrivateKey(crypto4go.FormatPKCS8PrivateKey(priKey))
		if err != nil {
			return nil, ErrLoadPrivateKey
		}
	}
	_pubKey, err := crypto4go.ParsePublicKey(crypto4go.FormatPublicKey(sandPubKey))
	if err != nil {
		return nil, ErrLoadPublicKey
	}
	cli := &Client{
		merchantNo:    merchantNo,
		privateKey:    _priKey,
		sandPublicKey: _pubKey,
	}
	if defaultClient == nil {
		defaultClient = cli
	}
	return cli, nil
}

type OrderParams struct {
	m map[string]string
}

type OrderResp struct {
	rawBytes []byte
}

func BuildPayUrlForH5Alipay(orderNo string, amount float64,
	notifyUrl, returnUrl,
	createIp, goodsName string,
	createAt time.Time, hint string) (string, error) {

	if defaultClient == nil {
		panic("sandpay: must init defaultClient by call InitClient() firstly")
	}
	p := buildOrderParams(orderNo, amount, notifyUrl, returnUrl, createIp, goodsName, createAt)
	return defaultClient.BuildPayUrlForH5Alipay(p, hint)
}

// BuildPayUrlForH5Alipay - H5包装支付宝生活号-02020002
func (c *Client) BuildPayUrlForH5Alipay(p *OrderParams, hint string) (string, error) {
	// 添加此场景特定的参数
	p.m["mer_no"] = c.merchantNo
	p.m["pay_extra"] = `{}`
	p.m["accsplit_flag"] = "NO"
	p.m["product_code"] = "02020002"
	p.m["jump_scheme"] = ""

	params := url.Values{}
	for k, v := range p.m {
		params.Add(k, v)
	}
	sign, err := signWithPKCS1v15(params, c.privateKey, crypto.SHA1)
	if err != nil {
		return "", errors.New("sandpay - 签名失败: " + err.Error())
	}
	params.Add("sign", sign)

	orderUrl := fmt.Sprintf("%s?%s", api_H5Alipay, params.Encode())
	return orderUrl, nil
}

const timeLayout = "20060102150405"

// BuildOrderParams
// - amount: 元
// - notifyUrl: 回调地址，地址需向杉德报备
// - return_url: 支付后的重定向地址，支付成功和失败都会返回此页面
func buildOrderParams(orderNo string, amount float64,
	notifyUrl, returnUrl,
	createIp, goodsName string,
	createAt time.Time) *OrderParams {
	_urlParams := map[string]string{
		"version":      "10",
		"mer_order_no": orderNo, // 自定义，最小长度12位
		"create_time":  createAt.Format(timeLayout),
		"order_amt":    fmt.Sprintf("%.2f", amount), // "0.11"元，建议>=0.1
		"notify_url":   notifyUrl,
		"return_url":   returnUrl,
		"create_ip":    createIp,
		"sign_type":    "RSA",
		"store_id":     "000000",
		"expire_time":  createAt.Add(time.Minute * 30).Format(timeLayout),
		"goods_name":   goodsName,
		"clear_cycle":  "3", // 清算模式，3-D1;0-T1;1-T0;2-D0 TODO
		"meta_option":  `[{"s":"Android","n":"","id":"","sc":""},{"s":"IOS","n":"","id":"","sc":""}]`,
		//"activity_no":    "", // 详情咨询业务员
		//"benefit_amount": "", // 详情咨询业务员
		//"extend":         `{}`,  // 可选
	}
	return &OrderParams{_urlParams}
}

// JSONData 回调使用的结构体
type JSONData struct {
	Head Head `json:"head"`
	Body Body `json:"body"`
}
type Head struct {
	Version  string `json:"version"`
	RespTime string `json:"respTime"`
	RespCode string `json:"respCode"`
	RespMsg  string `json:"respMsg"`
}
type Body struct {
	Mid                 string `json:"mid"`
	OrderCode           string `json:"orderCode"`
	TradeNo             string `json:"tradeNo"`
	ClearDate           string `json:"clearDate"`
	TotalAmount         string `json:"totalAmount"`
	OrderStatus         string `json:"orderStatus"`
	PayTime             string `json:"payTime"`
	SettleAmount        string `json:"settleAmount"`
	BuyerPayAmount      string `json:"buyerPayAmount"`
	DiscAmount          string `json:"discAmount"`
	TxnCompleteTime     string `json:"txnCompleteTime"`
	PayOrderCode        string `json:"payOrderCode"`
	AccLogonNo          string `json:"accLogonNo"`
	AccNo               string `json:"accNo"`
	MidFee              string `json:"midFee"`
	ExtraFee            string `json:"extraFee"`
	SpecialFee          string `json:"specialFee"`
	PlMidFee            string `json:"plMidFee"`
	Bankserial          string `json:"bankserial"`
	ExternalProductCode string `json:"externalProductCode"`
	CardNo              string `json:"cardNo"`
	CreditFlag          string `json:"creditFlag"`
	Bid                 string `json:"bid"`
	BenefitAmount       string `json:"benefitAmount"`
	RemittanceCode      string `json:"remittanceCode"`
	Extend              string `json:"extend"`
}

func (c *Client) ParseCallbackJSON(data, sign string) (*JSONData, error) {
	b := []byte(data)
	_, err := verify(b, sign, c.sandPublicKey)
	if err != nil {
		return nil, errors.New("sandpay: verify sign failed")
	}
	j := new(JSONData)
	_ = json.Unmarshal(b, j)
	return j, nil
}

func ParseCallbackJSON(data, sign string) (*JSONData, error) {
	if defaultClient == nil {
		panic("sandpay: must init defaultClient by call InitClient()")
	}
	return defaultClient.ParseCallbackJSON(data, sign)
}

func AckCallback(w http.ResponseWriter) {
	// 回复下面的内容，则杉德不再回调此订单
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("respCode=000000"))
}
