## Install MySQL

```shell
yum install -y mysql mysql-devel
```

注意：若本机docker安装了mysql，在宿主机使用mysql cli连接时需要制定`-h 127.0.0.1`，而不是默认的 `-h localhost`
，原因参考[csdn][0]。


[0]: https://blog.csdn.net/Aria_Miazzy/article/details/92803246