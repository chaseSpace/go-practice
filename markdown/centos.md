# CentOS

## 更新源

仅国内机器需要。

```shell
cd /etc/yum.repos.d/
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
yum clean all
yum makecache
```

## 常用

```shell
yum install -y git wget make
```

```shell
# Install Docker
yum install -y yum-utils device-mapper-persistent-data lvm2

yum-config-manager --add-repo http://download.docker.com/linux/centos/docker-ce.repo #中央仓库
# yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo    #阿里仓库

yum list docker-ce --showduplicates | sort -r
yum -y install docker-ce-26.1.4-1.el7

systemctl enable docker && systemctl start docker

yum erase docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine \
                  docker-ce
```

来自 https://blog.csdn.net/wade3015/article/details/94494929

## 安装Nginx

不建议也不需要docker安装nginx。

```shell
sudo yum install -y epel-release
sudo yum install -y nginx

#安装成功后，默认的网站目录为： /usr/share/nginx/html
#默认的配置文件为：/etc/nginx/nginx.conf
#自定义配置文件目录为: /etc/nginx/conf.d/

systemctl start nginx && systemctl enable nginx
```

### 使用acme.sh申请免费证书

https://github.com/acmesh-official/acme.sh/wiki/说明

## 安装chrome

```shell
wget https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm

yum install libX11 libXcursor libXdamage libXext libXcomposite libXi libXrandr gtk3 libappindicator-gtk3 xdg-utils \
libXScrnSaver liberation-fonts \
alsa-lib liberation-fonts libgbm xdg-utils

rpm -ivh google-chrome-stable_current_x86_64.rpm

google-chrome --version


# 卸载
yum remove google-chrome-stable

# cli调试 --repl是进入交互模式
google-chrome --headless --no-sandbox --enable-unsafe-swiftshader --remote-debugging-port=9222 --disable-gpu --repl 
https://www.google.com 
```