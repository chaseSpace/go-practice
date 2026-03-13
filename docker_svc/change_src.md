## Change docker src

```shell
mkdir /etc/docker
echo '''{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.m.daocloud.io",
    "https://dockerproxy.com",
    "https://docker.mirrors.sjtug.sjtu.edu.cn"
  ]
}
''' > /etc/docker/daemon.json

sudo systemctl daemon-reload
sudo systemctl restart docker
```