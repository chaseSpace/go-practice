## Prometheus的安装使用

### 学习

[Prometheus book 中文](https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/promql/prometheus-metrics-types)

[Prometheus 实战](https://songjiayang.gitbooks.io/prometheus/content/)

学习链接
- [什么是指标类型?](https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/promql/prometheus-metrics-types)
- [初识查询语言PromQL](https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/promql/prometheus-query-language)
- [理解四大指标类型以及对应PromQL](https://pandaychen.github.io/2020/04/11/PROMETHEUS-METRICS-INTRO/)

### 安装

>参考：[掘金文章—使用Prometheus搭建监控系统](https://juejin.cn/post/7095578954660053006)

从 https://prometheus.io/download/ 下载安装包

```shell
tar xvfz prometheus-2.35.0.linux-amd64.tar.gz
cd prometheus-2.35.0.linux-amd64

# 阻塞启动  观察日志是否正常  需要排错时加 --log.level=debug
./prometheus --config.file=prometheus.yml --web.listen-address=:11000


# 若正常，ctrl+c 重新后台启动
nohup ./prometheus --config.file=prometheus.yml --web.listen-address=:11000 &
```

通过HTTP API删除prometheus数据：
```shell
# 删除具有标签kubernetes_name="redis"的时间序列指标
curl -X POST -g 'http://localhost:9090/api/v1/admin/tsdb/delete_series?match[]={kubernetes_name="redis"}'

# 删除所有的数据
curl -X POST -g 'http://localhost:9090/api/v1/admin/tsdb/delete_series?match[]={__name__=~".+"}'

# 彻底清除
curl -X PUT　http://localhost:9090/api/v1/admin/tsdb/clean_tombstones

```
不过需要注意的是上面的 API 调用并不会立即删除数据，实际数据任然还存在磁盘上，会在后面进行数据清理。
要确定何时删除旧数据，可以使用--storage.tsdb.retention参数进行配置（默认情况下，Prometheus 会将数据保留15天）。

访问 http://ip:11000 进入仪表台。

### 安装grafana（作为prometheus的前端大屏展示）
[官网安装](https://grafana.com/grafana/download?edition=enterprise&platform=linux)

```shell
# Centos 安装
sudo yum install -y https://dl.grafana.com/enterprise/release/grafana-enterprise-9.5.1-1.x86_64.rpm

systemctl start grafana-server
systemctl enable grafana-server # 开机启动

tail -f /var/log/grafana/grafana.log # 查实时日志

# 配置文件：/etc/grafana/grafana.ini
# bin：/usr/share/grafana/bin/grafana
# logs：/var/log/grafana
# data: /var/lib/grafana
# plugins: /var/lib/grafana/plugins
# 功能配置目录: /etc/grafana/provisioning 用于自动化配置和部署，例如数据源、面板、警报等
```
访问grafana前端：http://ip:3000，登录用户名密码admin/admin

>初次修改/etc/grafana/grafana.ini时，记得删掉该行开头的分号，才会生效。

### 配置dashboard
参考：https://www.cnblogs.com/linuxk/p/12030478.html

### node_exporter
`xxx_exporter` 一般是用来直接安装在数据采集边缘端的导出客户端，而`node_exporter`就是用来采集linux服务器各项指标的导出客户端。

在安装prometheus后，一般首先用它来进行测试。具体安装使用命令如下：

```shell
# 从 https://prometheus.io/download/#node_exporter 下载安装包        
tar xvf node_exporter-1.5.0.linux-amd64.tar.gz 
cd node_exporter-1.5.0.linux-amd64/
./node_exporter --web.listen-address=:11002
nohup ./node_exporter --web.listen-address=:11002 &
```
由于`node_exporter`不需要存储数据，也就不需要通过docker运行。

### Pushgateway
https://github.com/prometheus/pushgateway

#### 作用
Pushgateway是一种将指标数据推送到Prometheus的中间代理。它的作用是允许那些非常短暂或不稳定的任务（如批处理作业和临时服务）向Prometheus报告指标，并使得这些指标可以被Prometheus轻松地收集和查询。

#### 使用场景
- 瞬时任务：需要收集指标的任务非常短暂，无法使用Prometheus的pull模型进行监控。
- 临时服务：需要在运行时收集指标的服务是动态生成的，无法预先注册到Prometheus中。
- 批处理作业：需要对一次性的批处理作业（如ETL、数据清洗等）进行指标监控。
- 有较多客户端需要收集指标：使用pushgateway可以避免在prometheus配置大量的job

#### 工作流程
在没有pushgateway时，prometheus是通过`prometheus.yml`中的job配置去主动收集每个客户端的端点暴露的数据。有了pushgateway后，
只需要在`prometheus.yml`中添加一项pushgateway的job，然后对于需要收集数据的客户端，也无需在`prometheus.yml`中添加job，而是直接将数据推送到pushgateway。

数据流向：
- 多个客户端的数据汇聚到pushgateway进行存储；
- prometheus定时通过pushgateway暴露的端点进行数据收集；


#### 缺点：
- 存在单点故障：pushgateway故障后，影响所连接客户端的数据上报；
- 增加了存储负担：原本数据就存在prometheus，现在还需在pushgateway进行存储，即需要定时对pushgateway进行旧数据清理；

#### 安装使用
下面演示的docker安装使用pushgateway的步骤。

测试：
```shell
docker run -it --rm -p 9222:9091 \
  -v /service/pushgateway:/pushgateway \
  --name pp prom/pushgateway \
  --persistence.file="/pushgateway/data.file" \
  --persistence.interval="3s"

# 其中`--persistence.file`是设置pushgateway持久化文件（默认不持久化，使用内存）
# `--persistence.interval=`是持久化时间间隔，默认`5m`，可以改为`15s`或`1m` (不宜过于频繁)
```
查看pushgateway状态：` curl -X GET http://pushgateway.example.org:9091/api/v1/status | jq`

推送数据进行测试：`echo "some_metric 3.14" | curl --data-binary @- http://127.0.0.1:9222/metrics/job/some_job`

检查数据是否被收到：`curl -X GET http://127.0.0.1:9222/api/v1/metrics | jq`

删除某个job下的某个instance中的全部指标数据：`  curl -X DELETE http://pushgateway.example.org:9091/metrics/job/some_job/instance/some_instance`

删除某个job下的全部指标数据：`  curl -X DELETE http://127.0.0.1:9222/metrics/job/some_job`

#### 生产部署
参考本仓库下的 [docker-compose.yml](../docker_svc/prometheus-grafana/docker-compose.yml)