# Install Node.js

```shell
# centos 7 不要安装 v18+, 肯定失败！
yum install -y gcc openssl-devel gcc-c++
wget https://nodejs.org/download/release/v16.20.0/node-v16.20.0-linux-x64.tar.xz

tar xf node-v16.20.0-linux-x64.tar.xz -C /usr/local/
mv /usr/local/node-v16.20.0-linux-x64 /usr/local/node
ln -s /usr/local/node/bin/node /usr/local/bin/node
ln -s /usr/local/node/bin/npm /usr/local/bin/npm


npm config set registry https://registry.npmmirror.com
npm config set bin-links false # 避免共享文件夹导致安装失败
```

其他：

```shell
npm install -g cnpm
```

## NVM 安装

https://juejin.cn/post/70brew install nvm00652162950758431

```shell
# 不要随意安装较新版本，系统可能不支持！
nvm install v16.20.0
```

## windows 安装 v18

https://nodejs.org/en/download/package-manager

FNM是win10的包管理器，默认安装：C:\Users{your user name}\AppData\Roaming\fnm，安装完node后需要设置env。
Win + R 输入：`%APPDATA%/fnm/node-versions`，找到node bin目录，然后把bin目录添加到path中。