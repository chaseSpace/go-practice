## AlmaLinux-10

### 初始化

```shell
sudo dnf install -y wget curl git vim unzip zip telnet net-tools 
dnf install -y binutils bison gcc make
```

### 网络相关

```shell
ip a
ip link show
ip route

curl ifconfig.me

#查看所有ip
hostname -I

# networkManager
dnf install NetworkManager

systemctl enable NetworkManager
systemctl start NetworkManager

nmcli device status

netstat -antup
```

### 系统相关

```shell
cat /etc/os-release
cat /etc/redhat-release
uname -a
```

### 其他软件

```shell
# 安装docker&docker-compose
dnf install -y dnf-plugins-core device-mapper-persistent-data lvm2 curl

# 安装docker、docker-compose
dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
curl -L https://github.com/docker/compose/releases/download/v2.36.2/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

docker --version
docker-compose --version
```
