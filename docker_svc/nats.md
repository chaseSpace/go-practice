## NATS

MqNats NATS 是一个轻量级、高性能的消息系统，核心设计理念是简单、可靠、可扩展。
它支持发布/订阅、请求/响应、队列等多种消息模式，适用于微服务、IoT、边缘计算等场景
NATS VS Kafka，如果要持久化，NATS性能不如Kafka，否则远高于Kafka

- Github：https://github.com/nats-io/nats-server

### Docker 启动

```shell
# -js 启用 JetStream 持久化功能
docker run -p 4222:4222 --name nats -tid nats:latest -js
```

### JetStream

JetStream 是 NATS 的持久化消息扩展，提供了消息持久化、重试机制和确认策略等功能。与默认的 Core NATS（非持久化）相比，JetStream
支持以下特性：

- 消息持久化：消息会被存储，即使服务重启也不会丢失。
- 消息重试：消费失败时，可根据 ACK 策略进行重试。
- ACK 机制：消费者需确认消息处理状态。

JetStream 支持三种 ACK 策略：

- Explicit（显式确认）：默认策略。每条消息必须通过 msg.Ack()、msg.Nak() 或 msg.Term() 显式确认，否则会超时重投。
- All（批量确认）：收到一批消息后一次性确认整批。
- None（无需确认）：服务器认为消息发送即成功，不等待确认，不保证投递。

**💡 对比 Core NATS**：Core NATS 默认不持久化消息，不支持重试和 ACK，每条消息最多投递一次。