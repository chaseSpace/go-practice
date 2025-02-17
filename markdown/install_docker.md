## Install Docker

### With Centos

```shell
yum install -y yum-utils device-mapper-persistent-data lvm2

yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo


yum list docker-ce --showduplicates | sort -r

yum -y install docker-ce-26.1.4

systemctl start docker
systemctl enable docker
```

### 安装docker-compose

```shell
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
chmod +x /usr/local/bin/docker-compose && \
docker-compose --version
```