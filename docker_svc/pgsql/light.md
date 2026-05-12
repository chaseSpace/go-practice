## 轻量化部署 pgsql (docker on linux)

启动 PostgreSQL:

```shell
POSTGRES_PASSWORD=123
docker run -d --name pgsql \
      -p 5432:5432 \
      -v ~/docker/pgsql:/var/lib/postgresql \
      -v /etc/localtime:/etc/localtime \
      -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
       postgres:18
```

m1 mac需要指定平台拉取镜像：

```shell
docker run -d --name pgsql \
      -p 5432:5432 \
      -v ~/docker/pgsql:/var/lib/postgresql \
      -v /etc/localtime:/etc/localtime \
      -e POSTGRES_PASSWORD='123' \
      --platform linux/amd64 \
       postgres:18
```

PostgreSQL 18+ 镜像默认将数据放在 `/var/lib/postgresql/18/docker`，所以建议把宿主机目录挂载到 `/var/lib/postgresql`，不要再直接挂载 `/var/lib/postgresql/data`。

如果之前已经用旧挂载方式启动过，可能会报 `/var/lib/postgresql/data (unused mount/volume)`。本地测试环境可以换一个新的宿主机目录，或先备份旧目录后重新初始化；生产数据不要直接切镜像大版本，需要按 PostgreSQL 大版本升级流程处理。

其他常用命令：

```shell
# 在宿主机尝试连接
psql -h 127.0.0.1 -p 5432 -U postgres

# 或直接进入容器
docker exec -it pgsql psql -U postgres

# 删除容器
docker stop pgsql && docker rm pgsql
```

## 设置 pgsql 远程登录

官方 PostgreSQL 镜像默认监听容器内所有地址。通过 `-p 5432:5432` 暴露端口后，可以使用宿主机 IP 和上面设置的 postgres 密码远程连接。

如果需要限制来源，请调整防火墙或自定义挂载 `postgresql.conf`、`pg_hba.conf`。

## 测试表

```
create database test;
\c test;

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    gender BIT(1) NOT NULL,
    money NUMERIC(2) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- drop table users;
```
