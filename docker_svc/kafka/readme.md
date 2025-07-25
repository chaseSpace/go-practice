## kafka 4.x 快速部署

4.0 版本开始废弃 zookeeper，使用 KRaft 替代。

使用 [docker-compose.yml](docker-compose.yml) 进行部署。


```shell
docker-compose up -d

# 查看容器状态
docker-compose ps
```

```shell
# 创建测试 topic
docker exec -it kafka4 kafka-topics.sh --create \
  --topic demo --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

# 生产消息
docker exec -it kafka4 kafka-console-producer.sh \
  --bootstrap-server localhost:9092 --topic demo

# 消费消息（另开终端）
docker exec -it kafka4 kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 --topic demo --from-beginning
```