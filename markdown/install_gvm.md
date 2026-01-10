## 使用 GVM 安装go

```shell
bash <<(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source ~/.gvm/scripts/gvm

# 依赖
yum -y install binutils bison gcc make

# 1.24.11  1.25.5  1.23.12

gvm install go1.23.12 --binary
gvm list
gvm use go1.23.12

# 安装dlv
go install github.com/go-delve/delve/cmd/dlv@v1.23.12
```