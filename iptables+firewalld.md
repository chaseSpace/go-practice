## Iptables 使用

### 开启内核转发
```shell
#/etc/sysctl.conf文件中添加net.ipv4.ip_forward=1，表示开启内核流量转发
sysctl -p  # 生效
```
```shell
sudo yum install iptables-services iptables-devel
systemctl enable iptables.service && sudo systemctl start iptables.service
```


配置流量转发：
```shell
iptables -t nat -A PREROUTING -p tcp --dport 9050 -j DNAT --to-destination 175.178.200.251:9050
iptables -t nat -A POSTROUTING -j MASQUERADE

iptables -P FORWARD ACCEPT
iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT

# 查看nat列表
iptables -t nat -L -n

service iptables save && service iptables restart

# 删除nat配置
iptables -t nat -D PREROUTING/POSTROUTING/INPUT/OUTPUT（其中之一） int（第几条）
```

清空iptables：
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


### 使用Firewalld（更简单）
```shell
firewall-cmd --zone=public --permanent --add-masquerade # 开启nat
firewall-cmd --add-port=9050/tcp --permanent #　放开端口
#开启TCP流量转发
firewall-cmd --add-forward-port=port=9050:proto=tcp:toaddr=175.178.200.251:toport=9050 --permanent

# 重启生效
firewall-cmd --reload

# 删除指定的转发规则
firewall-cmd --remove-forward-port=port=9050:proto=tcp:toaddr=175.178.200.251:toport=9050 --permanent 
```

其他常用：
```shell
firewall-cmd --list-all

# 删除规则
# <zone_name>：要操作的区域（例如 public、internal、dmz 等）。
# <rule_type>：要删除的规则类型，如 service、port、rich-rule 等。
# <rule>：要删除的规则。
firewall-cmd --zone=<zone_name> --remove-<rule_type>=<rule>

# 举例
sudo firewall-cmd --zone=public --remove-port=80/tcp
```

### 扩展：使用NC命令测试端口
```shell
# 安装nc
wget http://vault.centos.org/6.6/os/x86_64/Packages/nc-1.84-22.el6.x86_64.rpm
rpm -iUv nc-1.84-22.el6.x86_64.rpm


# 监听端口
nc -l 8888

# 其他主机
nc [ip] 8888
```
