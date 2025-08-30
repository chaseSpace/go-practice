## Etcd

etcd 是一个 分布式、高可用、强一致性的键值存储系统，由 CoreOS（现为 Red Hat 旗下）开发，使用 Go 语言实现。它是 Kubernetes、Cloud
Foundry、OpenStack 等系统的核心组件，用于 服务发现、配置管理、分布式锁 等场景。

官方文档：https://etcd.io/docs/v3.6/quickstart/

### Docker 启动单节点

```shell
# 单节点最简命令（v3.5.21）,一般使用次新主版本（最新稳定版本是 3.6）
# Note: 2379是client连接端口，2380是组集群时节点之间的连接端口
docker run -d --name etcd \
  -p 2379:2379 \
  gcr.io/etcd-development/etcd:v3.5.21 \
  /usr/local/bin/etcd \
  --name node1 \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://127.0.0.1:2379
```

**验证**

```shell
# 看集群成员
docker exec -it etcd etcdctl --endpoints=http://localhost:2379 member list

# 看集群健康
docker exec -it etcd etcdctl --endpoints=127.0.0.1:2379 endpoint health
```