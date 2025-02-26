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
yum install -y git wget
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