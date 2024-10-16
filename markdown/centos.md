# CentOS

## 更新源

```shell
cd /etc/yum.repos.d/
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
yum clean all
yum makecache 
```

来自 https://blog.csdn.net/wade3015/article/details/94494929