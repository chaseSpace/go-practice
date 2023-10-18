## k8s使用教程

<!-- TOC -->
  * [k8s使用教程](#k8s使用教程)
    * [0. 安装docker](#0-安装docker)
    * [1. 创建程序和使用docker打包镜像](#1-创建程序和使用docker打包镜像)
    * [2. 推送到docker仓库](#2-推送到docker仓库)
    * [3. 了解并创建Pod](#3-了解并创建pod)
      * [3.1 创建nginx pod](#31-创建nginx-pod)
      * [3.2 创建pod](#32-创建pod)
      * [3.3 查看nginx-pod状态](#33-查看nginx-pod状态)
      * [3.4 与pod交互](#34-与pod交互)
      * [3.5 Pod 与 Container 的不同](#35-pod-与-container-的不同)
      * [3.6 创建go程序的pod](#36-创建go程序的pod)
      * [3.7 pod有哪些状态](#37-pod有哪些状态)
    * [4. 使用Deployment](#4-使用deployment)
      * [4.1 部署deployment：](#41-部署deployment)
      * [4.2 修改deployment](#42-修改deployment)
      * [4.3 更新deployment](#43-更新deployment)
      * [4.4 回滚部署](#44-回滚部署)
      * [4.4 滚动更新（Rolling Update）](#44-滚动更新rolling-update)
      * [4.5 deployment的扩缩容](#45-deployment的扩缩容)
      * [4.6 k8s的镜像管理](#46-k8s的镜像管理)
    * [参考资料](#参考资料)
<!-- TOC -->

**环境准备**：

```
10.0.2.2 k8s-master  
10.0.2.3 k8s-node1
```

### 0. 安装docker
[Centos安装docker](https://www.runoob.com/docker/centos-docker-install.html)

```shell
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
docker version

systemctl start docker
systemctl enable docker
```

### 1. 创建程序和使用docker打包镜像

1. 编写一个简单的[main.go](./minikube/main.go)
2. 编写[Dockerfile](Dockerfile)

打包镜像（替换leigg为你的docker账户名）

```shell
docker build . -t leigg/hellok8s:v1
```

这里有个小问题，（修改代码后）重新构建镜像若使用同样的镜像名会导致旧的镜像的名称和tag变成`<none>`，可通过下面的命令来一键删除：

```shell
docker image prune -f
# docker system prune # 删除
```

测试运行：

```shell
docker run --rm -p 3000:3000 leigg/hellok8s:v1
```

### 2. 推送到docker仓库

先登录

```shell
$ docker login  # 然后输入自己的docker账户和密码，没有先去官网注册
```

推送

```shell
docker push leigg/hellok8s:v1
```

### 3. 了解并创建Pod

Pod 是 Kubernetes 最小的可部署单元，**通常包含一个或多个容器**。
它们可以容纳紧密耦合的容器，例如运行在同一主机上的应用程序和其辅助进程。但是，在生产环境中，通常使用其他资源来更好地管理和扩展服务。

Pod是 Kubernetes 中创建和管理的、最小的可部署的计算单元。

#### 3.1 创建nginx pod

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

#### 3.2 创建pod

运行第一条k8s命令创建pod：

```shell
kubectl apply -f nginx.yaml
```

#### 3.3 查看nginx-pod状态

```shell
kubectl get po nginx-pod
```

查看全部pods：`kubectl get pods`

#### 3.4 与pod交互

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

#### 3.5 Pod 与 Container 的不同

在刚刚创建的资源里，在最内层是我们的服务 nginx，运行在 container 容器当中， container (容器) 的**本质是进程**，而 pod 是管理这一组进程的资源。

所以 pod 可以管理多个 container，在某些场景例如服务之间需要文件交换(日志收集)，本地网络通信需求(使用 localhost 或者 Socket 文件进行本地通信)，
在这些场景中使用 pod 管理多个 container 就非常的推荐。而这，也是 k8s 如何处理服务之间复杂关系的第一个例子。

**Pod定义**  
Pod 是 Kubernetes 最小的可部署单元，通常包含一个或多个容器。它们可以容纳紧密耦合的容器，例如运行在同一主机上的应用程序和其辅助进程。但是，在生产环境中，通常使用其他资源来更好地管理和扩展服务。

#### 3.6 创建go程序的pod

定义[pod.yaml](./pod.yaml)

启动pod：

```shell
$ kk apply -f pod.yaml
# 几秒后
$ kk get pods
NAME      READY   STATUS    RESTARTS   AGE
go-http   1/1     Running   0          17s
```

临时开启端口转发（在master节点）：

```shell
# 绑定pod端口3000到 master节点的3000端口
kubectl port-forward go-http 3000:3000
```
现在pod提供的http服务可以在master节点上可用。

打开另一个会话测试：
```shell
$ curl http://localhost:3000
[v1] Hello, Kubernetes!#
```

#### 3.7 pod有哪些状态

- Pending（挂起）： Pod 正在等待调度。
- ContainerCreating（容器创建中）： Pod 已经被调度，但其中的容器尚未完全创建和启动。
- Running（运行中）： Pod 中的容器已经在运行。
- Succeeded（已成功）： 所有容器都成功终止，任务或工作完成，特指那些批处理任务而不是常驻容器。
- Failed（已失败）： 至少一个容器以非零退出码终止。
- Unknown（未知）： 无法获取 Pod 的状态。

### 4. 使用Deployment
通常，Pod不会被（通过pod.yaml）直接创建和管理，而是由更高级别的控制器，如Deployment，来创建和管理。
这是因为Deployment提供了更强大的应用程序管理功能。

- **应用管理**：Deployment是Kubernetes中的一个控制器，用于管理应用程序的部署和更新。它允许你定义应用程序的期望状态，然后确保集群中的副本数符合这个状态。

- **自愈能力**：Deployment可以自动修复故障，如果Pod失败，它将启动新的Pod来替代。这有助于确保应用程序的高可用性。

- **滚动更新**：Deployment支持滚动更新，允许你逐步将新版本的应用程序部署到集群中，而不会导致中断。

- **副本管理**：Deployment负责管理Pod的副本，可以指定应用程序需要的副本数量，Deployment将根据需求来自动调整。

- **声明性配置**：Deployment的配置是声明性的，你只需定义所需的状态，而不是详细指定如何实现它。Kubernetes会根据你的声明来管理应用程序的状态。

先创建一个[deployment文件](./deployment.yaml)， 用来编排多个pod。

#### 4.1 部署deployment：
```shell
$ kk apply -f deployment.yaml
deployment.apps/hellok8s-go-http created

# 查看启动的pod
$ kk get deployments                
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
hellok8s-go-http   2/2     2            2           3m
```
还可以查看pod运行的node：
```shell
# 这里的IP是pod ip，属于部署k8s集群时规划的pod网段
# NODE就是集群中的node名称
$ kk get pod -o wide
NAME                                READY   STATUS    RESTARTS   AGE   IP           NODE        NOMINATED NODE   READINESS GATES
hellok8s-go-http-55cfd74847-5jw7f   1/1     Running   0          68s   20.2.36.75   k8s-node1   <none>           <none>
hellok8s-go-http-55cfd74847-zlf49   1/1     Running   0          68s   20.2.36.74   k8s-node1   <none>           <none>
```

**删除pod会自动重启一个，确保可用的pod数量与`replicas`保持一致，不再演示**。

#### 4.2 修改deployment
通过vi修改内容中的replicas=3，再次部署，开始之前，我们使用下面的命令来观察pod数量变化
```shell
$ kubectl get pods --watch
NAME                                   READY   STATUS    RESTARTS   AGE
hellok8s-go-http-58cb496c84-cft9j   1/1     Running   0          4m7s


# 在另一个CLI执行 kk apply ...

hellok8s-go-http-58cb496c84-sdrt2   0/1     Pending   0          0s
hellok8s-go-http-58cb496c84-sdrt2   0/1     Pending   0          0s
hellok8s-go-http-58cb496c84-pjkp9   0/1     Pending   0          0s
hellok8s-go-http-58cb496c84-pjkp9   0/1     Pending   0          0s
hellok8s-go-http-58cb496c84-sdrt2   0/1     ContainerCreating   0          0s
hellok8s-go-http-58cb496c84-pjkp9   0/1     ContainerCreating   0          0s
hellok8s-go-http-58cb496c84-pjkp9   1/1     Running             0          1s
hellok8s-go-http-58cb496c84-sdrt2   1/1     Running             0          1s
```

#### 4.3 更新deployment
这一步通过修改main.go来模拟实际项目中的服务更新，修改后的文件是[main2.go](./main2.go)。

重新构建镜像：
```shell
docker build . -t leigg/hellok8s:v2
```

再次push镜像到仓库：
```shell
docker push leigg/hellok8s:v2
```
然后更新deployment：
```shell
$ kubectl set image deployment/hellok8s-go-http hellok8s=leigg/hellok8s:v2

$ 查看更新过程
$ kubectl rollout status deployment/hellok8s-go-http
Waiting for deployment "hellok8s-go-http" rollout to finish: 2 out of 3 new replicas have been updated...
Waiting for deployment "hellok8s-go-http" rollout to finish: 2 out of 3 new replicas have been updated...
Waiting for deployment "hellok8s-go-http" rollout to finish: 2 out of 3 new replicas have been updated...
Waiting for deployment "hellok8s-go-http" rollout to finish: 1 old replicas are pending termination...
Waiting for deployment "hellok8s-go-http" rollout to finish: 1 old replicas are pending termination...
deployment "hellok8s-go-http" successfully rolled  # OK

# 也可以直接查看pod信息，会观察到pod正在更新（这是一个启动新pod，删除旧pod的过程，最终会维持到所配置的replicas数量）
$ kk get pods
NAMESPACE     NAME                                       READY   STATUS              RESTARTS      AGE
default       go-http                                    1/1     Running             0             14m
default       hellok8s-go-http-55cfd74847-5jw7f          1/1     Terminating         0             27m
default       hellok8s-go-http-55cfd74847-z29dl          1/1     Running             0             23m
default       hellok8s-go-http-55cfd74847-zlf49          1/1     Running             0             27m
default       hellok8s-go-http-668c7f75bd-m56pm          0/1     ContainerCreating   0             0s
default       hellok8s-go-http-668c7f75bd-qlrk5          1/1     Running             0             14s

# 绑定其中一个pod来测试
$ kk port-forward hellok8s-go-http-668c7f75bd-m56pm 3000:3000
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
```
在另一个会话窗口执行
```shell
$ curl http://localhost:3000
[v2] Hello, Kubernetes!
```

这里演示的更新是容器更新，修改deployment.yaml的其他配置也属于更新。

#### 4.4 回滚部署
如果新的镜像无法正常启动，则旧的pod不会被删除，但需要回滚，使deployment回到正常状态。

按照下面的步骤进行：

1. 修改main.go，将最后监听端口那行先注释，添加一行：panic("something went wrong")
2. 构建镜像: docker build . -t leigg/hellok8s:v2_problem
3. push镜像：docker push leigg/hellok8s:v2_problem
4. 更新deployment使用的镜像：kubectl set image deployment/hellok8s-go-http hellok8s=leigg/hellok8s:v2_problem
5. 观察：kubectl rollout status deployment/hellok8s-go-http   （会停滞，按 Ctrl-C 停止观察）
6. 观察pod：kubectl get pods

```shell
$ kk get pods
NAME                                READY   STATUS             RESTARTS     AGE
go-http                             1/1     Running            0            36m
hellok8s-go-http-55cfd74847-fv2kp   1/1     Running            0            17m
hellok8s-go-http-55cfd74847-l78pb   1/1     Running            0            17m
hellok8s-go-http-55cfd74847-qtb59   1/1     Running            0            17m
hellok8s-go-http-7c9d684dd-msj2c    0/1     CrashLoopBackOff   1 (4s ago)   6s

# CrashLoopBackOff状态表示重启次数过多，过一会儿再试，这表示pod内的容器无法正常启动，或者启动就立即退出了

# 查看每个副本集每次更新的pod情况（包含副本数量、上线时间、使用的镜像tag）
# DESIRED-预期数量，CURRENT-当前数量，READY-可用数量
# -l 进行标签筛选
$ kubectl get rs -l app=hellok8s -o wide
NAME                          DESIRED   CURRENT   READY   AGE   CONTAINERS   IMAGES                      SELECTOR
hellok8s-go-http-55cfd74847   0         0         0       76s   hellok8s     leigg/hellok8s:v1           app=hellok8s,pod-template-hash=55cfd74847
hellok8s-go-http-668c7f75bd   3         3         3       55s   hellok8s     leigg/hellok8s:v2           app=hellok8s,pod-template-hash=668c7f75bd
hellok8s-go-http-7c9d684dd    1         1         0       11s   hellok8s     leigg/hellok8s:v2_problem   app=hellok8s,pod-template-hash=7c9d684dd
```

现在进行回滚：
```shell
# 先查看deployment更新记录
$ kk rollout history deployment/hellok8s-go-http               
deployment.apps/hellok8s-go-http 
REVISION  CHANGE-CAUSE
1         <none>
2         <none>
3         <none>

# 现在回到revision 2，可以先查看它具体信息（主要看用的哪个镜像tag）
$ kk rollout history deployment/hellok8s-go-http --revision=2
deployment.apps/hellok8s-go-http with revision #2
Pod Template:
  Labels:	app=hellok8s
	pod-template-hash=668c7f75bd
  Containers:
   hellok8s:
    Image:	leigg/hellok8s:v2
    Port:	<none>
    Host Port:	<none>
    Environment:	<none>
    Mounts:	<none>
  Volumes:	<none>

# 确认后，回滚（到上个版本）
$ kubectl rollout undo deployment/hellok8s-go-http  #到指定版本 --to-revision=2          
deployment.apps/hellok8s-go-http rolled back

# 检查副本集状态（所处的版本）
$ kk get rs -l app=hellok8s -o wide                                
hellok8s-go-http-55cfd74847   0         0         0       9m31s   hellok8s     leigg/hellok8s:v1           app=hellok8s,pod-template-hash=55cfd74847
hellok8s-go-http-668c7f75bd   3         3         3       9m10s   hellok8s     leigg/hellok8s:v2           app=hellok8s,pod-template-hash=668c7f75bd
hellok8s-go-http-7c9d684dd    0         0         0       8m26s   hellok8s     leigg/hellok8s:v2_problem   app=hellok8s,pod-template-hash=7c9d684dd

# 恢复正常
$ kk get deployments hellok8s-go-http
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
hellok8s-go-http   3/3     3            3           7m42s
```

#### 4.5 滚动更新（Rolling Update）
k8s 1.15版本起支持滚动更新，即先创建新的pod，创建成功后再删除旧的pod，确保更新过程无感知，大大降低对业务影响。

在 deployment 的资源定义中, spec.strategy.type 有两种选择:

- RollingUpdate: 逐渐增加新版本的 pod，逐渐减少旧版本的 pod。（常用）
- Recreate: 在新版本的 pod 增加前，先将所有旧版本 pod 删除（针对那些不能多进程部署的服务）

另外，还可以通过以下字段来控制升级 pod 的速率：
- maxSurge: 最大峰值，用来指定可以创建的超出期望 Pod 个数的 Pod 数量。
- maxUnavailable: 最大不可用，用来指定更新过程中不可用的 Pod 的个数上限。

如果不设置，deployment会有默认的配置：
```shell
$ kk describe -f deployment.yaml
Name:                   hellok8s-go-http
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

为了明确地指定deployment的更新方式，我们需要在yaml中配置：
```shell
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-go-http
spec:
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: 3
省略其他熟悉的配置项。。。
```
这样，我们通过`k apply`命令时会以滚动更新方式进行。
>从`maxSurge: 1`可以看出更新时最多会出现4个pod，从`maxUnavailable: 1`可以看出最少会有2个pod正常运行。

注意：无论是通过`kubectl set image ...`还是`kubectl rollout restart deployment xxx`方式更新deployment都会遵循配置进行滚动更新。

#### 4.6 deployment的扩缩容
```shell
# 指定副本数量
$ kubectl scale deployment/hellok8s-go-http --replicas=10
deployment.apps/hellok8s-go-http scaled

# 观察到副本集版本并没有变化，而是数量发生变化
$ kubectl get rs -l app=hellok8s -o wide                 
NAME                          DESIRED   CURRENT   READY   AGE   CONTAINERS   IMAGES                      SELECTOR
hellok8s-go-http-55cfd74847   0         0         0       33m   hellok8s     leigg/hellok8s:v1           app=hellok8s,pod-template-hash=55cfd74847
hellok8s-go-http-668c7f75bd   10        10        10      33m   hellok8s     leigg/hellok8s:v2           app=hellok8s,pod-template-hash=668c7f75bd
hellok8s-go-http-7c9d684dd    0         0         0       32m   hellok8s     leigg/hellok8s:v2_problem   app=hellok8s,pod-template-hash=7c9d684dd
```

#### 4.7 k8s的镜像管理
TODO


### 参考资料

- [总教程](https://github.com/guangzhengli/k8s-tutorials/blob/main/docs/pre.md)
- [Docker教程](https://yeasy.gitbook.io/docker_practice/)
- [kubectl全部命令-官方](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [国内Centos机器安装clash代理](../../pure_doc/use_clash_linux.md)
