## 启动Centos

```shell
# 启动并直接进入容器内
docker run  -it --rm --name centos --network host -v ~/docker/centos/data:/data centos
``` 