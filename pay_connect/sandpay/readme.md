## 杉德支付

[官方文档](https://www.yuque.com/sd_cw/xfq1vq/ut7292)

## 1. 对接产品——H5包装支付宝生活号（02020002）

特点：

- 无需重新发布商家APP
- 支持原生APP
- 支持uniapp
- 支持web框架APP
- 无需集成杉德SDK包和插件

### 1. 搞定签名和验签

- 在后台生成下单参数时需要对下单参数进行**签名**，然后将签名包含在URL参数中；
- 在接收sand回调时需要对传入参数**验签**，再解析；

sand文档说明了具体的签名方式：https://www.yuque.com/sd_cw/xfq1vq/ut7292#rl9yx ，
对于go语言它没提供sdk，需要自己实现，本文件下的代码已经实现。

### 2. 下单请求

先由后台构建含下单参数以及签名的URL传给前端，再由**前端H5**请求该URL以进入支付宝H5页。完整的URL示例如下：

```
https://sandcash.mixienet.com.cn/pay/h5/alipay?version=10&mer_no=6888806&mer_order_no=a765a2cf88af49ebbc8a0451b57424c1&create_time=20221117163619&expire_time=20221117173619&order_amt=0.5&notify_url=https%3A%2F%2Fwww.baidu.com&return_url=https%3A%2F%2Fwww.baidu.com&create_ip=127_0_0_1&goods_name=%E6%B5%8B%E8%AF%95&store_id=000000&product_code=02020002&clear_cycle=3&meta_option=%5B%7B%22s%22%3A%22Android%22,%22n%22%3A%22wxDemo%22,%22id%22%3A%22com.pay.paytypetest%22,%22sc%22%3A%22com.pay.paytypetest%22%7D%5D&accsplit_flag=NO&jump_scheme=&sign_type=RSA&sign=IYZ3k%2BrNMgbm5FNPuwf2tF4ACWTY9NU8c0gajU0wZ9Ll8%2BK0uD2GWqSuzxMPB53KRiQU%2BGRyeRw6kIT5a4I%2B00fEuDoXTSTvA1qserIp19mM5BLxtriv8h8CKhRA%2FlBPavB4eLmutBNWWFTpnB8oIOSBJjM7edGPbH0zlhG%2BkdaZGlWmUO7ribsewtdw2KnCIPxNKlHo3iA5gp4e9oAIiHm0aYE6wUqw1W2Fq%2FwHVXHc4Xaa%2Fl6teFvBoc7SUF4cR2EWoRNBVJkVHnngOoNfuF7ai%2BRfgvLJbQJ%2BWcL7O%2BVaqI379%2FAkx46piUN5lULp1ntq1Y6aHoUPATuX4V3ulw%3D%3D
```
>注意，后端只负责构建下单参数生成URL，传给前端，由前端去调用URL完成sand端下单以及支付流程。
> 
具体参数来自 [服务端统一下单参数](https://www.yuque.com/sd_cw/xfq1vq/ut7292#R50PD)

### 3. 获取商户私钥和杉德公钥
#### 商户私钥
我们的步骤是先以 [sand提供的导出私钥pfx文件的指导](https://open.sandpay.com.cn/product/detail/43962/44235/) 从证书中导出私钥文件(需要填写密码)。
然后从`xx.pfx`文件中提供商户私钥明文，导出私钥命令：
```shell
# step1: 导出pem格式的私钥（需输入导出pfx文件时输入的密码）
openssl pkcs12 -in priv.pfx -nocerts -out priv.pem -nodes
# step2：pem导出私钥明文
openssl rsa -in priv.pem -out priv.key

# 私钥明文示例：
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAxfmbrtezLQoS12qW42JAyjovL5z6Pe5acKizsYDA4Unxd9+O
5xgp5mwT0QffzJ5tJ0gK/Ys3jTGgSDwy8SB4+s/tGSWt/sotNbEp+PgEojA5tH5l
+JTQZcFhcUUisCcQCRwC2XHs5R3MIWR9boxt/a7+CjVG82FGLXQX/OnuJ0oXwVm2
F+vUBIkEgohlVXZ36yusrKQ9Yude1X7eI579EO46CAjVN9mO8/9bZpj5lzKmFZ0O
pjyhkNEJsMEjDmGM1+GcazRP67jy1YebkZibcjkvQtKYAVdsFdtWyIbUdcSZLIYw
...
NEW6IQ5JbeGsdZKYvd/H4wKBgBeoivjNc6N8QrxzFeVni/rK+POrMOVXf2AAIo98
y2LA88KFgxaKmIcRDKwg3OCdOr03ucyQWliXYn6nT+QwXCCdtOAJzhcZhqcgbBnd
0a2cDY3b5C1XnO03U+fyduGK3Myos/Y95sQePPQIBq11Nr8pvOtSoebzfuRg2GS0
tMtzAoGAYdkS7WbFc6lIhu6XMOsE5QnKOEjl0BCbUCgEdrJ7L7qO/wb2v0EpYAwm
/yZPhuAcjS7VlCWSB2+nVUV5vlaBg+xzviEDcRlh38hS8Vx3C1bkedEALxHZ1tTV
Hy243upf3gvbXj2NTvafNxmkIHdG7DAUThDsHR4rGgym7QOGsjU=
-----END RSA PRIVATE KEY-----
```

#### 杉德公钥
杉德提供其证书下载，可从证书中提取出杉德公钥：`openssl x509 -in sand.cer -inform der -pubkey`
>[杉德公钥下载](https://open.sandpay.com.cn/open/downFileById?id=900db638fffe46859490deb21bf2d294)
```shell
# 公钥明文示例
-----BEGIN CERTIFICATE-----
MIIDJjCCAg6gAwIBAgIGAViVUCVLMA0GCSqGSIb3DQEBCwUAMFQxDTALBgNVBAMT
BFNBTkQxDTALBgNVBAsTBFNBTkQxDTALBgNVBAoTBFNBTkQxCzAJBgNVBAcTAlNI
MQswCQYDVQQIEwJTSDELMAkGA1UEBhMCQ04wHhcNMTYxMTIzMDc1MDA3WhcNMjYx
MTIyMDc1MDA3WjBUMQ0wCwYDVQQDEwRzYW5kMQ0wCwYDVQQLEwRTQU5EMQ0wCwYD
VQQKEwRTQU5EMQswCQYDVQQHEwJTSDELMAkGA1UECBMCU0gxCzAJBgNVBAYTAkNO
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyIwo8Jq6XiUSY8cMrDfT
...
vwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCd6W65I4SJ2BkH8RsBnZNDpQ7fMTYt
VfQtBctptmEtbSlNz7WKrBFcjsZ6KvsMoUy0ftvaddHKu7t8UhW29C0vWGW00Ihf
dkjGYMWxrNcicX/4KJPjFWwPXw11Vc3NuSlVJrr+eh5OSLHPnpvxoKs4I+55hXS4
4Ch1x2LZ4rLsQ6vrVJz2mnygg2JEeredh74XAMAgAWGZ/Tqn4/QWpjFDggHOF8I9
eXddK5yiD+cJ3EcZDYr4LGaaG95XQfvdKNl0igAfFmGd3Sxg5MrFnFJbDsqE0HAF
crEaCK5rqJVMQdvZEWO4j6c6ZX9WCMfjcXRbZonE3b1DSXCrh1uildk2
-----END CERTIFICATE-----

```

### 4. 异步回调

[回调文档](https://open.sandpay.com.cn/product/detail/43314/43801/43805)

下面列出的是关键字段，以及示例数据
```shell
{
	"head": {
		"version": "1.0",
		"respTime": "20230321180517",
		"respCode": "000000",
		"respMsg": "成功"
	},
	"body": {
		"mid": "6888801119246",
		"orderCode": "alipayh5-sand-18703a0f312",
		"tradeNo": "alipayh5-sand-18703a0f312",
		"clearDate": "20230321",
		"totalAmount": "000000000100",
		"orderStatus": "1",
		"payTime": "20230321180517",
		"settleAmount": "000000000100",
		"buyerPayAmount": "000000000100",
		"discAmount": "000000000000",
		"txnCompleteTime": "20230321180516",
		"payOrderCode": "20230321001343020000000000044686",
		"accLogonNo": "kyt***@163.com",
		"accNo": "208800******8700",
		"midFee": "000000000001",
		"extraFee": "000000000000",
		"specialFee": "000000000000",
		"plMidFee": "000000000000",
		"bankserial": "952023032122001488701432093215",
		"externalProductCode": "00002022",
		"cardNo": "208800******8700",
		"creditFlag": "",
		"bid": "",
		"benefitAmount": "000000000000",
		"remittanceCode": "",
		"extend": ""
	}
}```

```