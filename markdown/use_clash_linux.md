## 安装 clash linux 版

目的，为了国内 Centos 主机能够畅通世界互联网。

准备：
- 一个可使用 clash 连接的订阅链接
- 一台国内 Centos 服务器
- 暴露外网 9090 端口，之后可关闭

```shell
cd ~
wget https://github.com/Dreamacro/clash/releases/download/v1.14.0/clash-linux-amd64-v1.14.0.gz # 拉取clash压缩包，也可直接GitHub下载后上传文件
gunzip clash-linux-amd64-v1.14.0.gz # 解压文件
mv clash-linux-amd64-v1.14.0 /usr/local/bin/clash # 更改文件名
chmod u+x /usr/local/bin/clash # 赋权


mkdir -p ~/.config/clash/
# -O指定具体路径
wget -O ~/.config/clash/config.yml __sub_link
# -P指定位置
wget -P ~/.config/clash/ https://github.com/Dreamacro/maxmind-geoip/releases/download/20230312/Country.mmdb # 拉取Country.mmdb文件，也可直接GitHub下载后上传文件


vi ~/.config/clash/config.yml  # 修改9090端口对应的密码
nohup ./clash &

# 其他配置
export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=socks5://127.0.0.1:7891
ping k8s.io

# unset http_proxy https_proxy all_proxy


# 设置开机启动
# 创建 systemd 脚本，脚本文件路径为 /etc/systemd/system/clash.service，内容如下：
[Unit]
Description=clash daemon

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/clash
Restart=on-failure

[Install]
WantedBy=multi-user.target

# --------------------- 111 ---------------------
systemctl daemon-reload
#systemctl start clash
systemctl enable clash

# 看日志
journalctl -u clash.service -f

# 设置env
vi /etc/profile.d/proxy.sh # 写入下面内容
# export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=socks5://127.0.0.1:7891
source /etc/profile
```

然后服务器开放外网 9090，访问：[razord 控制面板](http://clash.razord.top)，输入服务器 IP 和刚才的密码，enjoy~