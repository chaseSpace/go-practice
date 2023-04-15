## 功能说明
开启云端审核后，在App后台可以主动调用RESTAPI接口，送审音视图文等相关内容，其中图文同步返回机审结果，音视频通过异步回调的形式返回机审结果。

## 接口调用说明
### 请求 URL 示例
```
https://xxxxxx/v4/im_msg_audit/content_moderation?sdkappid=88888888&identifier=admin&usersig=xxx&random=99999999&contenttype=json
```
### 请求参数说明

下表仅列出调用本接口时涉及修改的参数及其说明，更多参数详情请参考 [REST API 简介](https://cloud.tencent.com/document/product/269/1519)。

| 参数               | 说明                                 |
| ------------------ | ------------------------------------ |
| xxxxxx | SDKAppID 所在国家/地区对应的专属域名：<br><li>中国：`console.tim.qq.com`</li><li>新加坡：`adminapisgp.im.qcloud.com`</li><li>首尔： `adminapikr.im.qcloud.com`</li><li>法兰克福：`adminapiger.im.qcloud.com`</li><li>孟买：`adminapiind.im.qcloud.com`</li><li>硅谷：`adminapiusa.im.qcloud.com`</li>|
| v4/recentcontact/mark_contact  | 请求接口                             |
| sdkappid           | 创建应用时即时通信 IM 控制台分配的 SDKAppID |
| identifier         | 必须为 App 管理员帐号，更多详情请参见 [App 管理员](https://cloud.tencent.com/document/product/269/31999#app-.E7.AE.A1.E7.90.86.E5.91.98)                |
| usersig            | App 管理员帐号生成的签名，具体操作请参见 [生成 UserSig](https://cloud.tencent.com/document/product/269/32688)    |
| random             | 请输入随机的32位无符号整数，取值范围0 - 4294967295                 |
|contenttype|请求格式固定值为`json`|

### 最高调用频率

200次/秒。

### 请求包示例
```
{
    "AuditName":"C2C",
    "ContentType":"Text",
    "Content":"违规词汇"
}

```

### 请求包字段说明

| 字段 | 类型|属性| 说明 |
|---------|---------|----|---------|
| AuditName|String| 必填|表明送审策略，取值:C2C/Group/UserInfo/GroupInfo/GroupMemberInfo/RelationChain|
| ContentType|String| 必填|送审类型，取值：Text/Image/Audio/Video。|
| Content|String| 必填|送审内容，最大限制8KB，当审核文件时，填对应URL。其中图片审核最大不超过5MB。|

### 应答包体示例

```
{
    "ActionStatus": "OK",
    "ErrorCode": 0,
    "ErrorInfo": "",
    "RequestId": "91fa78f3-18c8-4b20-9c56-5845df18f634",
    "Result": "Block",
    "Score": 100,
    "Label": "Polity",
    "Keywords": [
        "违规词汇"
    ]
}
```


### 应答包字段说明

| 字段|类型 |说明 |
|---------|---------|---------|
| ActionStatus| String | 请求处理的结果，OK 表示处理成功，FAIL 表示失败  |
| ErrorCode| Integer | 错误码，0表示成功，非0表示失败|
| ErrorInfo| String | 错误信息  |
| RequestId| String | 审核标示，音视频异步审核通过RequestId从回调获取审核结果 https://cloud.tencent.com/document/product/269/78633#.E6.AD.A5.E9.AA.A44.EF.BC.9A.E5.AE.A1.E6.A0.B8.E7.BB.93.E6.9E.9C.E9.85.8D.E7.BD.AE  |
| Result| String | 图文审核建议，Pass/Review/Block |
| Score| Integer | 图文审核恶意值，0-100，恶意程度与Score成正比。 |
| Label| String | 送审内容命中的标签 Normal/Polity/Porn/Illegal/Abuse/Terror/Ad/Sexy/Composite |


[](id:ErrorCode)
## 错误码说明
除非发生网络错误（例如502错误），否则该接口的 HTTP 返回码均为200。实际的错误码、错误信息是通过应答包体中的 ResultCode、ResultInfo、ErrorCode 以及 ErrorInfo 来表示的。
公共错误码（60000到79999）请参见 [错误码](https://cloud.tencent.com/document/product/269/1671)。
本 API 私有错误码如下：

| 错误码 | 描述                                                         |
| ------ | ------------------------------------------------------------ |
| 60003  | 请求参数无效。 |
| 60020  | 未开启云端审核服务。                |
| 60022  | 请求内部错误，请联系我们。  |
| 93000  | 送审内容超过了最大限制8KB。|