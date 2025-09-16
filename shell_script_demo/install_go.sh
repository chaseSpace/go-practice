#!/usr/bin/env bash


#注意
#  1. 若已安装go，请从PATH中删除已有go版本的ENV，否则脚本执行后仍然会识别旧版
#  2. 请不要在goland的terminal中执行此脚本，因为此处打开的shell中的PATH会自动设置当前goland配置的go版本，导致脚本执行完后，本来已经安装成功，但因PATH设置不正确导致未能识别新版本go。

version=go1.21.12

target=$version.linux-amd64.tar.gz

GOPATH='/usr/local/gopath'
GO_INSTALL='/usr/local'
BASHRC='/root/.bashrc'
ZSHRC='/root/.zshrc'

if [ `uname` == "Darwin" ]; then
  echo "This is macOS"
  GOPATH=$HOME/gopath
  GO_INSTALL=$HOME
  BASHRC=$HOME/.bash_profile
  target=$version.darwin-amd64.tar.gz
fi


PATH=$(echo $PATH | sed 's/ /\\ /g')  # 将PATH中的空格转义，否则source会报错！

# mac遇到 wget 无法建立SSL连接，请修改代理为直连。
now=$(pwd) && cd $GO_INSTALL && \
wget https://studygolang.com/dl/golang/$target && \
tar xzf $target && \
rm -rf $target && \
mkdir -p $GOPATH && mkdir -p $GOPATH/bin

# Set environment variables in shell config files
# Check if ZSHRC exists, otherwise use BASHRC
if [ -f "$ZSHRC" ]; then
  {
    echo "export PATH=$PATH:$GO_INSTALL/go/bin:$GOPATH/bin"
    echo "export GOPATH=$GOPATH"
    echo "export GOPROXY=https://goproxy.cn"
    echo "export GO111module=on"
    echo "export GOROOT=$GO_INSTALL/go"
  } >> $ZSHRC && source $ZSHRC
else
  {
    echo "export PATH=$PATH:$GO_INSTALL/go/bin:$GOPATH/bin"
    echo "export GOPATH=$GOPATH"
    echo "export GOPROXY=https://goproxy.cn"
    echo "export GO111module=on"
    echo "export GOROOT=$GO_INSTALL/go"
  } >> $BASHRC && source $BASHRC
fi

cd $now && go version