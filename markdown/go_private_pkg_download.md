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
// - 如果密码中包含被URL保留的字符，例如`!@#$%^&+-`等，参照下面步骤进行转义
# 原始密码
password='1234!@#$'
# url转义（需要已安装python）
password=$(echo -n $password | python -c "import sys, urllib.parse; print(urllib.parse.quote(sys.stdin.read()))")
git config --global url."http://lilei:${password}@coding.xxx.cn".insteadOf "http://coding.xxx.cn"

# 若环境不支持python命令，请使用以下网站进行在线转换，将**密码**进行转换后再粘贴到命令中
# https://www.urlencoder.org/

// 删除配置
#git config --global --unset url."http://lilei:123@coding.xxx.cn".insteadOf

// 检查配置
git config --global -l


// 执行拉取（也是更新命令）
go get coding.xxx.cn/lilei/something
```