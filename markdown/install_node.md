# Install Node.js

Centos 7 安装 node v18（使用glibc兼容的版本）：

```shell
nvm unload # 卸载nvm

# https://unofficial-builds.nodejs.org/download/release/$VER/    18.20.7是V18最后一个版本

VER=v20.17.0
wget https://unofficial-builds.nodejs.org/download/release/$VER/node-$VER-linux-x64-glibc-217.tar.gz
mkdir -p node-$VER && tar -xzf node-$VER-linux-x64-glibc-217.tar.gz -C node-$VER --strip-components 1
ln -sf ~/node-$VER/bin/* /usr/local/bin/
```

其他：

```shell
npm install -g cnpm
```

## NVM 安装

https://juejin.cn/post/7394823316584972325

### Linux 安装

```shell
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
# 或者
wget -qO- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
```

根据项目要求按照指定版本：

```shell
nvm install v16.20.0

nvm install v18.20.0

npm config set registry https://registry.npmmirror.com
npm config get registry

{
  echo "export NVM_DIR="$HOME/.nvm""
  echo "[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh""  # This loads nvm
  echo "[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion""  # This loads nvm bash_completion
} >> ~/.zshrc
```

## windows 安装 v18

https://nodejs.org/en/download/package-manager

FNM是win10的包管理器，默认安装：C:\Users{your user name}\AppData\Roaming\fnm，安装完node后需要设置env。
Win + R 输入：`%APPDATA%/fnm/node-versions`，找到node bin目录，然后把bin目录添加到path中。

## windows 使用nvm

https://github.com/coreybutler/nvm-windows/releases

## windows 配置yarn代理

```shell
yarn config delete proxy
yarn config set https-proxy http://127.0.0.1:7890
```