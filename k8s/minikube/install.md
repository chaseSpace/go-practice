## Centos上使用minikube安装k8s

<!-- TOC -->

* [Centos上使用minikube安装k8s](#centos上使用minikube安装k8s)
    * [参考链接](#参考链接)
    * [0. 安装最新docker](#0-安装最新docker)
    * [1. 安装启动minikube](#1-安装启动minikube)
    * [2. 创建程序和使用docker打包镜像](#2-创建程序和使用docker打包镜像)
    * [3. 推送到docker仓库](#3-推送到docker仓库)
    * [4. 了解并创建Pod](#4-了解并创建pod)
        * [4.1 创建nginx pod](#41-创建nginx-pod)
        * [4.2 安装kubectl](#42-安装kubectl)
        * [4.3 创建pod](#43-创建pod)
        * [4.4 查看nginx-pod状态](#44-查看nginx-pod状态)
        * [4.5 与pod交互](#45-与pod交互)
        * [4.6 Pod 与 Container 的不同](#46-pod-与-container-的不同)
        * [4.7 创建go程序的pod](#47-创建go程序的pod)
        * [4.8 pod有哪些状态](#48-pod有哪些状态)
    * [5. Deployment](#5-deployment)

<!-- TOC -->
目录：

环境：

```
- 两台机，相同配置
    - OS: Centos v7.9
    - Mem: 4c8g
    - Disk: 100g
```

#### 参考链接

- [总教程](https://github.com/guangzhengli/k8s-tutorials/blob/main/docs/pre.md)
- [Docker教程](https://yeasy.gitbook.io/docker_practice/)
- [kubectl全部命令-官方](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)

### 0. 安装最新docker

- [Centos升级/安装docker](https://www.cnblogs.com/wdliu/p/10194332.html)
    - 注意换国内源

### 1. 安装启动minikube

安装

```shell
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

启动（minikube要求较新的docker版本）

```shell
# --force允许在root执行
# --image-mirror-country='cn' 是加速minikube自身资源的下载安装
minikube start --force --image-mirror-country=cn
```
如果多次在`Pulling base image...`未知出错，需要清理缓存：`minikube delete --all --purge`，再重新下载。


其他启动参数：
```
--driver=*** 从1.5.0版本开始，Minikube缺省使用系统优选的驱动来创建Kubernetes本地环境，比如您已经安装过Docker环境，minikube 将使用 docker 驱动
--cpus=2: 为minikube虚拟机分配CPU核数
--memory=2048mb: 为minikube虚拟机分配内存数
--registry-mirror=*** 为了提升拉取Docker Hub镜像的稳定性，可以为 Docker daemon 配置镜像加速，参考阿里云镜像服务
--kubernetes-version=***: minikube 虚拟机将使用的 kubernetes 版本
```

查看启动状态：

```shell
$ minikube status
minikube
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```

minikube命令速查：

```shell
minikube 命令速查

minikube stop 不会删除任何数据，只是停止 VM 和 k8s 集群。

minikube delete 删除所有 minikube 启动后的数据。

minikube ip 查看集群和 docker enginer 运行的 IP 地址。

minikube pause 暂停当前的资源和 k8s 集群

minikube status 查看当前集群状态
```

### 2. 创建程序和使用docker打包镜像

1. 编写一个简单的[main.go](./main.go)
2. 编写[Dockerfile](./Dockerfile)
    - 坑：因为运行go程序和编译的不是一个镜像？所以在编译程序时需要关闭CGO，否则启动main时会提示main文件找不到的问题。

打包镜像（替换leigg为你的docker账户名）

```shell
docker build . -t leigg/hellok8s:v1
```

这里有个小问题，（修改代码后）重新构建镜像若使用同样的镜像名会导致旧的镜像的名称和tag变成`<none>`，可通过下面的命令来一键删除：

```shell
docker rmi $(docker images -f "dangling=true" -q)
```

测试运行：

```shell
docker run --rm -p 3000:3000 leigg/hellok8s:v1
```

### 3. 推送到docker仓库

先登录

```shell
$ docker login  # 然后输入自己的docker账户和密码，没有先去官网注册
```

推送

```shell
docker push leigg/hellok8s:v1
```

### 4. 了解并创建Pod

Pod 是 Kubernetes 最小的可部署单元，**通常包含一个或多个容器**。
它们可以容纳紧密耦合的容器，例如运行在同一主机上的应用程序和其辅助进程。但是，在生产环境中，通常使用其他资源来更好地管理和扩展服务。

Pod是 Kubernetes 中创建和管理的、最小的可部署的计算单元。

#### 4.1 创建nginx pod

```yaml
# nginx.yaml
apiVersion: v1
kind: Pod  # 资源类型=pod
metadata:
  name: nginx-pod  # 需要唯一
spec:
  containers: # pod内的容器组
    - name: nginx-container
      image: nginx  # 镜像默认来源 DockerHub
```

#### 4.2 安装kubectl

由于minikube下载kubectl命令太慢，所以笔者自行下载kubectl。

先导入源

```shell
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
```

再配置和安装最新版本的k8s组件

```shell
setenforce 0
yum install -y kubelet kubeadm kubectl
systemctl enable kubelet && systemctl start kubelet
```

> 安装指定版本  
> `yum install -y kubelet-<version> kubectl-<version> kubeadm-<version>`

#### 4.3 创建pod

运行第一条k8s命令创建pod：

```shell
kubectl apply -f nginx.yaml
```

#### 4.4 查看nginx-pod状态

```shell
kubectl get po nginx-pod
```

查看全部pods：`kubectl get pods`

#### 4.5 与pod交互

添加端口转发，然后就可以在宿主机访问nginx-pod

```shell
# 宿主机4000映射到pod的80端口
# 这条命令是阻塞的，仅用来调试pod服务是否正常运行
kubectl port-forward nginx-pod 4000:80

# 打开另一个控制台
curl http://127.0.0.1:4000
```

其他命令：

```shell
kubectl delete pod nginx-pod # 删除pod
kubectl delete -f nginx.yaml  # 删除配置文件内的全部资源
 
kubectl exec -it nginx-pod -- /bin/bash   # 进入pod shell

# 支持 --tail LINES_NUM
kubectl logs -f nginx-pod  # 查看日志（stdout/stderr）
```

#### 4.6 Pod 与 Container 的不同

在刚刚创建的资源里，在最内层是我们的服务 nginx，运行在 container 容器当中， container (容器) 的本质是进程，而 pod 是管理这一组进程的资源。

所以 pod 可以管理多个 container，在某些场景例如服务之间需要文件交换(日志收集)，本地网络通信需求(使用 localhost 或者 Socket 文件进行本地通信)，
在这些场景中使用 pod 管理多个 container 就非常的推荐。而这，也是 k8s 如何处理服务之间复杂关系的第一个例子。

**Pod定义**  
Pod 是 Kubernetes 最小的可部署单元，通常包含一个或多个容器。它们可以容纳紧密耦合的容器，例如运行在同一主机上的应用程序和其辅助进程。但是，在生产环境中，通常使用其他资源来更好地管理和扩展服务。

#### 4.7 创建go程序的pod

定义pod.yaml:

```yaml
# go-http.yaml
apiVersion: v1
kind: Pod
metadata:
  name: go_http
spec:
  containers:
    - name: go_http-container
      image: leigg/hellok8s:v1
```

启动pod：

```shell
$ k apply -f go-http.yaml
➜  install_k8s_all k get pods
NAME      READY   STATUS              RESTARTS   AGE
go-http   0/1     ContainerCreating   0          16s
➜  install_k8s_all k get pods
NAME      READY   STATUS    RESTARTS   AGE
go-http   1/1     Running   0          17s
```

开启端口转发：

```shell
kubectl port-forward go-http 3000:3000
```

#### 4.8 pod有哪些状态

- Pending（挂起）： Pod 正在等待调度。
- ContainerCreating（容器创建中）： Pod 已经被调度，但其中的容器尚未完全创建和启动。
- Running（运行中）： Pod 中的容器已经在运行。
- Succeeded（已成功）： 所有容器都成功终止，任务或工作完成。
- Failed（已失败）： 至少一个容器以非零退出码终止。
- Unknown（未知）： 无法获取 Pod 的状态。

### 5. 了解Deployment