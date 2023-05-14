## docker-compose管理prometheus+grafana

首先安装docker-compose，略。

### 1、安装prometheus
```shell
mkdir -p /service/prometheus/data
mkdir /service/grafana

chmod 777 -R /service

# 复制 prometheus.yml 到 /service/prometheus，略
# 复制 docker-compose.yml 到/service/prometheus，略
# 复制 grafana.ini 到/service/grafana，略

cd /service/prometheus

# 直接启动  先观察日志是否正常，然后访问 http://IP:11000  (本目录下的配置文件修改了它的默认9090端口为11000)
docker-compose up 

# prometheus的WEB UI中点击 Status——Targets，确认 prometheus.yml中配置的两个target的state是否up
# 若是，则表示prometheus成功监控到自己的指标数据，否则需要排错

# 确认无误，然后退出进程，重新以后台模式启动几个容器

docker-compose up -d

# 常用命令：
# docker-compose restart <容器名>  # 不加容器名就是全部dc管理的容器

# docker-compose rm -f   # 删除所有docker-compose管理的容器
# docker-compose port prometheus 11001   查看容器映射宿主机的端口
# docker-compose ps -q prometheus     查看容器ID
```

### 2、在宿主机启动node_exporter
https://prometheus.io/download/#node_exporter 下载压缩包
```shell
# 解压
tar xvfz node_exporter-1.3.1.linux-amd64.tar.gz

# 运行
cd ./node_exporter-1.3.1.linux-amd64

# 默认端口9100 可不改
nohup ./node_exporter --web.listen-address=:11002 &

# 开放端口11002
# 检查是否启动成功
访问http://ip:11002/metrics 
```