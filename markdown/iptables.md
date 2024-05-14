## iptables

### 1. 四表五链

四表（Tables）：

- filter 表：用于数据包的过滤和防火墙功能，它是默认的表。
- nat 表：用于网络地址转换（NAT）和端口转发等地址操作。
- mangle 表：用于修改数据包的特定字段，例如 TTL（Time to Live）等。
- raw 表：用于在数据包进入网络协议栈之前进行处理，主要用于连接跟踪和数据包标记等。

五链（Chains）：

- INPUT 链：用于处理目标主机接收的数据包，包括来自本地进程或服务的数据包以及通过转发到达的数据包。
- OUTPUT 链：用于处理目标主机发出的数据包，包括由本地进程或服务生成的数据包以及通过转发发送的数据包。
- FORWARD 链：用于处理从一个网络接口进入防火墙并通过另一个网络接口转发的数据包，用于实现网络流量的转发。
- PREROUTING 链：用于在数据包到达目标主机之前进行处理，可以进行 DNAT（目标网络地址转换）等操作。
- POSTROUTING 链：用于在数据包离开目标主机之前进行处理，可以进行 SNAT（源网络地址转换）等操作。

流程图:
<img src="./images/iptables.png" width="900" height="600">

### 2. 查看链规则

```shell
iptables -L  #　默认filter表
iptables -t nat -L 
iptables -t mangle -L
iptables -t raw -L
```