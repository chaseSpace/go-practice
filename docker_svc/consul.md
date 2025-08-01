## Consul docker部署

### 开发模式

在开发环境下，可以使用docker启动一个单节点consul服务。

```
docker run -d --name=consul -p 8500:8500 hashicorp/consul agent -dev -ui -client=0.0.0.0
```

在宿主机访问Consul UI：http://localhost:8500