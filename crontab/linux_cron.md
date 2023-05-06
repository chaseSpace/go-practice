## Linux crontab

### 命令用法
crontab (选项) file

- -e：编辑某个用户的crontab文件内容。如果不指定用户，则表示编辑当前用户的crontab文件；  
- -l：显示某个用户的crontab文件内容，如果不指定用户，则表示显示当前用户的crontab文件内容；
- -r：从/var/spool/cron目录中删除某个用户的crontab文件，如果不指定用户，则默认删除当前用户的crontab文件；
- -u<用户名称>：指定要设定计时器的用户名称，例如，“-u haha”表示设定haha用户的crontab服务，此参数一般由root用户来运行；
- -i 在删除用户的crontab文件时给确认提示。
- -s (selinux context)

### 时间语法
推荐使用crontab语法解释网站：https://crontab.guru/

### 服务管理
在 /etc/rc.d/rc.local 中添加 service crond start 这一行。

```shell
# 启动服务
/sbin/service crond start 
# 关闭服务
/sbin/service crond stop 
# 重启服务
/sbin/service crond restart 
# 重新载入配置
/sbin/service crond reload
```

### 定时任务分类
Linux下的任务调度分为两类，系统任务调度和用户任务调度。

**系统任务调度**：系统周期性所要执行的工作，比如写缓存数据到硬盘、日志清理等。
在`/etc/crontab`文件，这个就是系统任务调度的配置文件。

**用户任务调度**：用户定期要执行的工作，比如用户数据备份、定时邮件提醒等。
用户可以使用 crontab 工具来定制自己的计划任务。
在crontab 文件都被保存在`/var/spool/cron`目录中。其文件名与用户名一致

### 常用命令
```shell
crontab -l    # 列出所有用户的任务（不用sudo则列出当前用户的任务）
crontab –e    # 编辑当前用户的任务列表
crontab -r    # 删除当前用户的全部定时任务，若加上sudo就是root用户的，建议先通过 -l查看存在哪些任务
crontab –i    # 在删除用户的crontab文件时给确认提示
```
如果报错：`No such file or directory`，可以重启服务`sudo service crond restart`，再执行命令。

增删改任务后无需重启crontab。

快速添加任务：
```shell
echo "0 0 * * * root python -c 'import random; import time; time.sleep(random.random() * 3600)' && certbot renew" | sudo tee -a /var/spool/cron/root > /dev/null
```

### 查看执行情况

#### 1. 查看执行记录
这个路径下记录了**哪一天执行了什么任务**，但不含执行结果
```shell
ls /var/log/cron*
/var/log/cron  /var/log/cron-20221211  /var/log/cron-20221218  /var/log/cron-20221225  /var/log/cron-20230101
```

#### 2. 查看执行结果
执行结果发送到用户邮箱，请确定你的定时任务属于哪个用户。
```shell
less /var/spool/mail/用户名
```

### 限制用户使用cron
- /etc/cron.allow：将可以使用 crontab 的帐号写入其中，若不在这个文件内的使用者则不可使用 crontab
- /etc/cron.deny：将不可以使用 crontab 的帐号写入其中，若未记录到这个文件当中的使用者，就可以使用 crontab

以优先顺序来说， /etc/cron.allow 比 /etc/cron.deny 要优先， 而判断上面，这两个文件只选择一个来限制而已，因此，建议你只要保留一个即可， 免得影响自己在配置上面的判断！一般来说，系统默认是保留 /etc/cron.deny ，你可以将不想让他运行 crontab 的那个使用者写入 /etc/cron.deny 当中，一个帐号一行！


### 备份和恢复定时任务
```shell
# backup
crontab -l > $HOME/mycron

# recover
sudo cp $HOME/mycron /var/spool/cron/用户名
```