## MongoDB Docker部署

```shell
dir=/Users/lei/docker/mongo
mkdir -p $dir/{data,logs}

docker run -d --name mongo --restart=always \
  -p 27017:27017 \
  -v $dir/data:/data/db \
  -v $dir/logs:/var/log/mongodb \
  -e TZ=Asia/Shanghai \
  -e MONGO_INITDB_ROOT_USERNAME=root \
  -e MONGO_INITDB_ROOT_PASSWORD=123 \
  mongo:latest
```

### 基本操作

```shell
use admin

db.createUser(
    {
        user: "user1",
        pwd: "123456",
        roles: [{ role: "root", db: "admin" }]
    }
);


db.updateUser("user1", { roles: [{ role: "readWrite", db: "admin" }] })
db.changeUserPassword("user1", "123456")
db.dropUser("user1")

show users
```