## 部署

使用pm2工具管理进程，[常用命令](https://blog.csdn.net/weixin_42658813/article/details/127283913)

查看进程状态：
```shell
[root@VM-0-13-centos like]# pm2 ls
┌────┬────────────────────┬──────────┬──────┬───────────┬──────────┬──────────┐
│ id │ name               │ mode     │ ↺    │ status    │ cpu      │ memory   │
├────┼────────────────────┼──────────┼──────┼───────────┼──────────┼──────────┤
│ 4  │ admin              │ fork     │ 15   │ online    │ 0%       │ 10.9mb   │
│ 2  │ api                │ fork     │ 17   │ online    │ 0%       │ 11.7mb   │
└────┴────────────────────┴──────────┴──────┴───────────┴──────────┴──────────┘
```

查看进程日志，以admin为例：
- pm2默认为进程创建2个日志文件：
    - /root/.pm2/logs/admin-out.log   正常日志
    - /root/.pm2/logs/admin-err.log   错误日志
```shell
# --lines 20 查看最新20行，这会输出上述2个文件的最新20行，并且是实时打印
pm2 logs admin --lines 20
```

启动进程
```shell
# -time 让log带上时间
pm2 start admin --time
```

更新服务，以admin为例：
```shell
# 1. 备份现在的bin文件
$ cd /service/like
$ cp admin admin-old

# 2. 上传并覆盖现有的admin二进制

# 3. 重启(需要给执行权限，才能正确重启)
$ chmod +x admin && pm2 reload admin 
```

其他命令：
```shell
pm2 reload all # 热重启全部
pm2 restart all # 重启全部
pm2 delete APP_NAME  # 不删log
pm2 show APP_NAME # 查看应用详情
pm2 moni APP_NAME # 监控应用指标

pm2 startup # 开机启动
pm2 save # 冻结当前应用列表，以便在开机时快速恢复
```

配置日志分割：
```shell
$ pm2 install pm2-logrotate

# 支持下面配置
$ pm2 set pm2-logrotate:max_size 10M
$ pm2 set pm2-logrotate:retain 30
$ pm2 set pm2-logrotate:compress false
$ pm2 set pm2-logrotate:dateFormat YYYY-MM-DD_HH-mm-ss
$ pm2 set pm2-logrotate:workerInterval 30
$ pm2 set pm2-logrotate:rotateInterval 0 0 * * *
$ pm2 set pm2-logrotate:rotateModule true

# 修改配置
pm2 set pm2-logrotate:max_size 30M  # 单文件最大30M
pm2 set pm2-logrotate:retain 5 # 保留5个文件
pm2 reload all # 重启生效
```
其他：
- 进程崩溃时会自动重启，但频繁多次重启失败时会停止尝试