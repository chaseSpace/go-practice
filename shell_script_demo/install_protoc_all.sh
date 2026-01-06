#!/bin/zsh
set -e

echo "====== （必须使用zsh执行）开始安装 protoc 和 protoc-gen-go ======"

# 版本信息
PROTOC_VERSION="3.12.4"
PROTOC_GEN_GO_VERSION="1.36.6"

PROTOC_FILE="protoc-${PROTOC_VERSION}-linux-x86_64.zip"
PROTOC_GEN_GO_FILE="protoc-gen-go.v${PROTOC_GEN_GO_VERSION}.linux.amd64.tar.gz"

# 下载地址(linux)
PROTOC_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/$PROTOC_FILE"
PROTOC_GEN_GO_URL="https://github.com/protocolbuffers/protobuf-go/releases/download/v${PROTOC_GEN_GO_VERSION}/$PROTOC_GEN_GO_FILE"

# 目录
WORK_DIR="$HOME/Downloads"
PROTOC_INSTALL_DIR="/usr/local/protoc"
BIN_DIR="/usr/local/bin"

# 创建工作目录
mkdir -p ${WORK_DIR}
cd ${WORK_DIR}

# ------------------------------------------------
# 1. 安装 protoc
# ------------------------------------------------
echo ">>> 下载 protoc ${PROTOC_VERSION}"
if [ ! -f $PROTOC_FILE ]; then
  wget -q ${PROTOC_URL}
else
  echo "$PROTOC_FILE 已存在"
fi

UNZIP_DIR=protoc-${PROTOC_VERSION}-linux-x86_64

echo ">>> 解压 protoc"
unzip -oq protoc-${PROTOC_VERSION}-linux-x86_64.zip -d $UNZIP_DIR

echo ">>> 安装 protoc 到 ${PROTOC_INSTALL_DIR}"
sudo rm -rf ${PROTOC_INSTALL_DIR}
sudo mv $UNZIP_DIR ${PROTOC_INSTALL_DIR}

# ------------------------------------------------
# 2. 安装 protoc-gen-go
# ------------------------------------------------
echo ">>> 下载 protoc-gen-go ${PROTOC_GEN_GO_VERSION}"
if [ ! -f $PROTOC_GEN_GO_FILE ]; then
  wget -q ${PROTOC_GEN_GO_URL}
else
  echo "$PROTOC_GEN_GO_FILE 已存在"
fi

echo ">>> 解压 protoc-gen-go"
tar -xzf protoc-gen-go.v${PROTOC_GEN_GO_VERSION}.linux.amd64.tar.gz

echo ">>> 安装 protoc-gen-go 到 ${BIN_DIR}"
sudo mv protoc-gen-go ${BIN_DIR}/
sudo chmod +x ${BIN_DIR}/protoc-gen-go

# ------------------------------------------------
# 3. 配置 PATH（只追加一次）
# ------------------------------------------------
echo ">>> 配置 PATH"

if ! grep -q "/usr/local/protoc/bin" ~/.zshrc; then
  echo "export PATH=$PATH:/usr/local/protoc/bin" >> ~/.zshrc
fi

if ! grep -q "/usr/local/bin" ~/.zshrc; then
  echo "export PATH=$PATH:/usr/local/bin" >> ~/.zshrc
fi

source ~/.zshrc

# ------------------------------------------------
# 4. 验证安装
# ------------------------------------------------
echo "====== 验证 ======"
protoc --version
protoc-gen-go --version

echo "====== 安装完成 🎉 ======"
