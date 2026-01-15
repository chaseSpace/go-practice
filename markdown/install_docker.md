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

> **如何让虚拟机使用宿主机代理**：将宿主机的以太网卡共享给 VMNET8，虚拟机无需设置任何proxy。
> - 经过测试，即使clash允许lan连接，虚拟机也无法访问宿主机的clash端口，应该是VMware限制。

### 写入源

```shell
sudo tee /etc/docker/daemon.json <<-'EOF'
{
    "registry-mirrors": [
    "https://docker.m.daocloud.io",
    "https://docker.imgdb.de",
    "https://docker-0.unsee.tech",
    "https://docker.hlmirror.com",
    "https://cjie.eu.org"
    ]
}
EOF

systemctl restart docker
```

### 安装docker-compose

```shell
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
chmod +x /usr/local/bin/docker-compose && \
docker-compose --version
```