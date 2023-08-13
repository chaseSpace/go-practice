## Centos上使用minikube安装k8s

<!-- TOC -->
  * [Centos上使用minikube安装k8s](#centos上使用minikube安装k8s)
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
    * [5. 了解Deployment](#5-了解deployment)
      * [5.1 部署deployment：](#51-部署deployment)
      * [5.2 修改deployment](#52-修改deployment)
      * [5.3 使用新的镜像更新pod](#53-使用新的镜像更新pod)
      * [5.4 滚动更新（Rolling Update）](#54-滚动更新rolling-update)
      * [5.5 minikube的镜像管理](#55-minikube的镜像管理)
      * [5.6 deployment的回滚](#56-deployment的回滚)
<!-- TOC -->

**环境准备**：

```
- 两台机，相同配置
    - OS: Centos v7.9
    - Mem: 4c8g
    - Disk: 100g
```

**参考资料：**

- [总教程](https://github.com/guangzhengli/k8s-tutorials/blob/main/docs/pre.md)
- [Docker教程](https://yeasy.gitbook.io/docker_practice/)
- [kubectl全部命令-官方](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [国内Centos机器安装clash代理](../../pure_doc/use_clash_linux.md)

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

>如备有代理，可参考前面**参考资料**中的文档连接代理后再直接下载kubectl

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

在刚刚创建的资源里，在最内层是我们的服务 nginx，运行在 container 容器当中， container (容器) 的**本质是进程**，而 pod 是管理这一组进程的资源。

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
- Succeeded（已成功）： 所有容器都成功终止，任务或工作完成，特指那些批处理任务而不是常驻容器。
- Failed（已失败）： 至少一个容器以非零退出码终止。
- Unknown（未知）： 无法获取 Pod 的状态。

### 5. 了解Deployment
先创建一个[deployment文件](./deployment.yaml)， 用来编排多个pod。

#### 5.1 部署deployment：
```shell
root@VM-0-13-centos ~/install_k8s » k apply -f deployment.yaml
deployment.apps/hellok8s-deployment created

# 查看启动的pod
root@VM-0-13-centos ~/install_k8s » k get pods
NAME                                   READY   STATUS    RESTARTS   AGE
hellok8s-deployment-784d5f676d-zcnr6   1/1     Running   0          17s
```

**删除pod会自动重启一个**。

#### 5.2 修改deployment
通过vi修改内容中的replicas=3，再次部署，开始之前，我们使用下面的命令来观察pod数量变化
```shell
root@VM-0-13-centos ~/install_k8s » kubectl get pods --watch
NAME                                   READY   STATUS    RESTARTS   AGE
hellok8s-deployment-58cb496c84-cft9j   1/1     Running   0          4m7s


# 在另一个CLI执行 k apply ...

hellok8s-deployment-58cb496c84-sdrt2   0/1     Pending   0          0s
hellok8s-deployment-58cb496c84-sdrt2   0/1     Pending   0          0s
hellok8s-deployment-58cb496c84-pjkp9   0/1     Pending   0          0s
hellok8s-deployment-58cb496c84-pjkp9   0/1     Pending   0          0s
hellok8s-deployment-58cb496c84-sdrt2   0/1     ContainerCreating   0          0s
hellok8s-deployment-58cb496c84-pjkp9   0/1     ContainerCreating   0          0s
hellok8s-deployment-58cb496c84-pjkp9   1/1     Running             0          1s
hellok8s-deployment-58cb496c84-sdrt2   1/1     Running             0          1s
```

#### 5.3 使用新的镜像更新pod
这一步通过修改main.go来模拟实际项目中的服务更新，修改后的文件是[main2.go](./main2.go)。

需要再次push镜像到仓库：
```shell
docker push leigg/hellok8s:v2
```
然后重新部署并测试：
```shell
root@VM-0-13-centos ~/install_k8s » k apply -f deployment.yaml
deployment.apps/hellok8s-deployment configured

root@VM-0-13-centos ~/install_k8s » k port-forward hellok8s-deployment-c7fdf4bc9-wh46w 3000:3000
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
Handling connection for 3000
```
在另一个CLI窗口执行
```shell
root@VM-0-13-centos ~ » curl http://localhost:3000
[v2] Hello, Kubernetes!
```

#### 5.4 滚动更新（Rolling Update）
上一步骤的更新方式比较粗暴，因为它是在新的镜像拉取后立即同时更新全部旧pod，这会导致
服务短暂不可用。如果新的镜像有问题，**这会导致更新失败，服务宕机**。

所以我们要使用更安全的滚动更新

>不过在笔者使用的`v1.27`版本中，通过`k apply`同样是滚动更新了。

在 deployment 的资源定义中, spec.strategy.type 有两种选择:

- RollingUpdate: 逐渐增加新版本的 pod，逐渐减少旧版本的 pod。（常用）
- Recreate: 在新版本的 pod 增加前，先将所有旧版本 pod 删除（针对那些不能多进程部署的服务）

另外，还可以通过以下字段来控制升级 pod 的速率：
- maxSurge: 最大峰值，用来指定可以创建的超出期望 Pod 个数的 Pod 数量。
- maxUnavailable: 最大不可用，用来指定更新过程中不可用的 Pod 的个数上限。

如果不设置，deployment会有默认的配置：
```shell
root@VM-0-13-centos ~/install_k8s » k describe -f deployment.yaml
Name:                   hellok8s-deployment
Namespace:              default
CreationTimestamp:      Sun, 13 Aug 2023 21:09:33 +0800
Labels:                 <none>
Annotations:            deployment.kubernetes.io/revision: 1
Selector:               app=aaa,app1=hellok8s
Replicas:               3 desired | 3 updated | 3 total | 3 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge # <------ 看这
省略。。。
```
**所以，在使用滚动更新时，k8s会始终保持服务可用，在新的pod未完全正常启动前，不会停止旧的pod。**


#### 5.5 minikube的镜像管理
当我们启动pod时，引用的镜像会从远程拉取到本地（而不是`docker images`），存入minikube自身的本地镜像库中管理，它和docker images是不同的东西。
```shell
# alias m='minikube'
root@VM-0-13-centos ~/install_k8s » m image -h
管理 images

Available Commands:
  build         在 minikube 中构建一个容器镜像
  load          将镜像加载到 minikube 中
  ls            列出镜像
  pull          拉取镜像
  push          推送镜像
  rm            移除一个或多个镜像
  save          从 minikube 中保存一个镜像
  tag           为镜像打标签

Use "minikube <command> --help" for more information about a given command.
root@VM-0-13-centos ~/install_k8s » m image ls
registry.k8s.io/pause:3.9
registry.k8s.io/kube-scheduler:v1.27.3
registry.k8s.io/kube-proxy:v1.27.3
registry.k8s.io/kube-controller-manager:v1.27.3
registry.k8s.io/kube-apiserver:v1.27.3
registry.k8s.io/etcd:3.5.7-0
registry.k8s.io/coredns/coredns:v1.10.1
gcr.io/k8s-minikube/storage-provisioner:v5
docker.io/leigg/hellok8s:v2   <----------------
docker.io/leigg/hellok8s:v1   <----------------
```
也就是说，`docker rmi`删除的镜像是不会影响minikube的镜像库的。即使通过`m image rm`删除了本地的一个minikube管理的镜像，
再启动deployment，也可以启动的，因为minikube会去远程镜像库Pull，除非远程仓库也删除了这个镜像。
重新启动后，可通过`m image ls`再次看到被删除的镜像又出现了。


#### 5.6 deployment的回滚
首次部署deployment后，通过`k rollout history`命令看到其第一次部署记录：
```shell
root@VM-0-13-centos ~/install_k8s » k rollout history -f deployment.yaml
deployment.apps/hellok8s-deployment
REVISION  CHANGE-CAUSE
1         <none>
```
因为只有一次记录，所以无法执行回滚命令`k rollout undo`：
```shell
root@VM-0-13-centos ~/install_k8s » k rollout undo -f deployment.yaml
error: no rollout history found for deployment "hellok8s-deployment"
```
现在我们修改`deployment.yaml`，使用v2镜像，然后再次部署，现在查看其部署记录：
```shell
root@VM-0-13-centos ~/install_k8s » k rollout history -f deployment.yaml
deployment.apps/hellok8s-deployment
REVISION  CHANGE-CAUSE
1         <none>
2         <none>
```
顺便查看deployment使用的镜像：
```shell
root@VM-0-13-centos ~/install_k8s » k describe -f deployment.yaml
Name:                   hellok8s-deployment
Namespace:              default
CreationTimestamp:      Sun, 13 Aug 2023 21:22:44 +0800
Labels:                 <none>
Annotations:            deployment.kubernetes.io/revision: 2
Selector:               app=aaa,app1=hellok8s
Replicas:               3 desired | 3 updated | 3 total | 3 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=aaa
           app1=hellok8s
  Containers:
   hellok8s-container:
    Image:        leigg/hellok8s:v2  # <--------------
```

现在可以进行回滚：
```shell
root@VM-0-13-centos ~/install_k8s » k rollout undo -f deployment.yaml
deployment.apps/hellok8s-deployment rolled back
```
然后通过上面的命令再次验证deployment使用的镜像即可。现在再看一下部署记录：
```shell
root@VM-0-13-centos ~/install_k8s » k rollout history -f deployment.yaml
deployment.apps/hellok8s-deployment
REVISION  CHANGE-CAUSE
2         <none>
3         <none>
```
可以看到 1 消失了，多了个 3。这不是因为最多保存2条，而是因为3和1是相同的镜像，只显示1条记录。
下面通过部署`v3`镜像来验证这一点。

执行下面的步骤：
- 修改`main.go`，在接口返回`v3`字样，保存
- 重新build v3镜像，并且push到docker远程仓库
- 需改`deployment.yaml`引用v3镜像，然后部署

验证：
```shell
root@VM-0-13-centos ~/install_k8s » k rollout history -f deployment.yaml
deployment.apps/hellok8s-deployment
REVISION  CHANGE-CAUSE
2         <none>
3         <none>
4         <none>
```

>注意：无论什么原因导致这次更换镜像的部署失败了，都不影响revision号的递增。


上面演示的是回滚到上个版本，但可以回滚到指定revision版本：
```shell
k rollout undo -f deployment.yaml --to-revision=2
```
如果我们回滚到2，那么同样的道理，revision 2消失，增加revision 5。

另外，在回滚前，我们可能需要查看这个revision的各项配置信息（容器、镜像、端口、挂载），可以查看：
```shell
root@VM-0-13-centos ~/install_k8s » k rollout history -f deployment.yaml --revision=2                                                                                             1 ↵
deployment.apps/hellok8s-deployment with revision #2
Pod Template:
  Labels:	
    app=aaa
	app1=hellok8s
	pod-template-hash=66695888cf
  Containers:
   hellok8s-container:
    Image:	leigg/hellok8s:v2
    Port:	<none>
    Host Port:	<none>
    Environment:	<none>
    Mounts:	<none>
  Volumes:	<none>
```