### 轻量化部署 redis

```shell
docker run --name redis5 -d -p 6379:6379 redis:5.0 --requirepass "123"
```

docker exec -it redis5 redis-cli -a 123