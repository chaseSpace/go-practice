## Go 私有包下载

```shell
// 配置开启gomod
go env -w GO111MODULE="on"
// 配置GoMod私有仓库
go env -w GOPRIVATE="git@coding.xxx.cn"
// 配置不加密访问
go env -w GOINSECURE="coding.xxx.cn"
// 配置不使用代理
go env -w GONOPROXY="coding.xxx.cn"
// 配置不验证包
go env -w GONOSUMDB="coding.xxx.cn"


// 在git请求URL中嵌入认证信息
#git config --global url."http://${user}:${password}@coding.xxx.cn".insteadOf "http://coding.xxx.cn"
// - 例如
git config --global url."http://lilei:123@coding.xxx.cn".insteadOf "http://coding.xxx.cn"


// 检查配置
git config --global -l


// 执行拉取（也是更新命令）
go get coding.xxx.cn/lilei/something
```