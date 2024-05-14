## 使用 Iptables(未跑通)

[科普：关于 SNAT 与 MASQUERADE 的区别](https://blog.51cto.com/dengaosky/2129665)

### 开启内核转发
```shell
#/etc/sysctl.conf文件中添加net.ipv4.ip_forward=1，表示开启内核流量转发
sysctl -p  # 生效
```
```shell
sudo yum install iptables-services iptables-devel
systemctl enable iptables.service 
systemctl start iptables.service
```

配置流量转发：
```shell
# 修改指定端口进入的流量的目的地址
iptables -t nat -A PREROUTING -p tcp --dport 9050 -j DNAT --to-destination 175.178.200.251:9050
# 对出去的流量的源IP使用当前主机网卡的IP
iptables -t nat -A POSTROUTING -j MASQUERADE  

# 设置默认转发策略
iptables -P FORWARD ACCEPT
# 允许通过已建立连接和相关连接的数据包
iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT

# 查看nat列表
iptables -t nat -L -n

service iptables save && service iptables restart

# 删除nat配置
iptables -t nat -D PREROUTING/POSTROUTING/INPUT/OUTPUT（其中之一） int（第几条）
```

清空 iptables：
```shell
# 清空全部
iptables -F

# 清空指定链
iptables -F INPUT
iptables -F FORWARD
iptables -F OUTPUT

# 只清空nat列表
iptables -t nat -F
```


## 使用 Firewalld（更简单，但不支持域名）
```shell
firewall-cmd --zone=public --permanent --add-masquerade # 开启nat
firewall-cmd --add-port=9050/tcp --permanent #　放开端口
#开启TCP流量转发
firewall-cmd --add-forward-port=port=9050:proto=tcp:toaddr=175.178.200.251:toport=9050 --permanent
# 重启生效
firewall-cmd --reload

# 删除指定的转发规则
forward_ports=(
9050
8000  # long conn
80 443 # nginx
9999 # uauth
9998 # uauth
10000 # friend
)
to_host='175.178.200.251'
for port in "${forward_ports[@]}"; do
    firewall-cmd --remove-port=$port/tcp --permanent
    firewall-cmd --remove-forward-port=port=$port:proto=tcp:toaddr=$to_host:toport=$port --permanent
    iptables -t nat -D PREROUTING -p tcp --dport $port -j DNAT --to-destination $to_host:$port # nat规则需要这样删除才生效
done
firewall-cmd --reload
```

其他常用：
```shell
firewall-cmd --list-all
firewall-cmd --state

systemctl stop firewalld
systemctl restart firewalld

# 删除规则
# <zone_name>：要操作的区域（例如 public、internal、dmz 等）。
# <rule_type>：要删除的规则类型，如 service、port、rich-rule 等。
# <rule>：要删除的规则。
firewall-cmd --zone=<zone_name> --remove-<rule_type>=<rule>

# 举例
firewall-cmd --zone=public --remove-port=80/tcp
```

## 使用 socat（使用简单且支持域名）
**严重问题**：启动进程一段时间后，由于未知原因会挂掉，未查出原因！
```shell
yum install -y socat

nohup socat TCP4-LISTEN:8000,reuseaddr,fork TCP4:testapi.aklivechat.com:8000>> socat.log 2>&1 &
```


### 扩展：使用 NC 命令测试端口
```shell
# 安装nc
wget http://vault.centos.org/6.6/os/x86_64/Packages/nc-1.84-22.el6.x86_64.rpm
rpm -iUv nc-1.84-22.el6.x86_64.rpm


# 监听端口
nc -l 8888

# 其他主机
nc [ip] 8888
```

快速启动 http 服务器：
```shell
python2 -m SimpleHTTPServer 8888
```