# kubeadm搭建k8s集群

## 1. 准备资源

- 至少一台linux机器，本文使用Centos7
- 2G+内存
- 2Cpu+
- 所有机器内网互通
    - 通过子网ping
- 不同的hostname，mac地址，product_uuid
    - 检查mac地址：`ip link` or `ifconfig -a`
    - 检查product_uuid：`sudo cat /sys/class/dmi/id/product_uuid`

- 禁用swap内存，让kubelet正常工作
    - centos：sudo swapoff -a （临时）
- 内网开放k8s需要的端口：
    - Control plane
        - tcp：6443 all
        - tcp：2379-2380 kube-apiserver, etcd
        - tcp：10250 Self
        - tcp：10259 Self
        - tcp：10257 Self
    - Worker node
        - tcp：10250 Self, Control plane
        - tcp：30000-32767 All

    - 测试端口开放：`nc 127.0.0.1 6443`

## 2. 安装容器runtime

k8s使用 Container Runtime Interface（CRI）来连接你选择的runtime。

### 2.1 Linux支持的CRI的端点

| Runtime                           | Path to Unix domain socket                 |
|-----------------------------------|--------------------------------------------|
| containerd                        | unix:///var/run/containerd/containerd.sock |
| CRI-O                             | unix:///var/run/crio/crio.sock             |
| Docker Engine (using cri-dockerd) | unix:///var/run/cri-dockerd.sock           |


## 3. 安装 kubeadm, kubelet and kubectl

```shell
# 设置阿里云为源
$ cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
       http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

# ubuntu
apt-get update && apt-get install -y apt-transport-https

curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 

cat <<EOF > /etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF

apt-get update
# 2023-10-15 已经出了1.28
apt-get install -y kubelet=1.25.14-00 kubeadm=1.25.14-00 kubectl=1.25.14-00
# 查看软件仓库包含哪些版本 apt-cache madison kubelet
# 删除 apt-get remove  -y kubelet kubeadm kubectl

# 检查版本
kubelet --version
kubeadm version -o json
kubectl version -o json

# centos 继续
sudo setenforce 0
sudo sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

# centos 安装各组件
sudo yum install -y kubelet-1.25.14 kubeadm-1.25.14 kubectl-1.25.14 --disableexcludes=kubernetes

# 开机启动，且立即启动
sudo systemctl enable --now kubelet
```

## 4. 为kubelet和runtime配置相同的cgroup driver

Container runtimes推荐使用`systemd`作为kubeadm的driver，而不是kubelet默认的`cgroupfs`driver。

从k8s v1.22起，kubeadm默认使用`systemd`作为cgroupDriver。

https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/configure-cgroup-driver/

所以使用高于v1.22的版本，这步就不用配置。

## 5. 使用kubeadmin创建集群

```shell

```