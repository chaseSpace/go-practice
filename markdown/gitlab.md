## Gitlab Runner维护文档

[官方介绍](https://docs.gitlab.cn/runner/)

Runner是CI/CD的具体执行者，也就是构建机。通常它是与Gitlab服务器交互，通过Gitlab服务器的API来获取任务，执行任务，并上报执行结果。

然后Runner一般是与需要部署的环境是处于同一网络或可以连接到目标环境的（才能完成部署步骤）。Runner的维护主要有以下几点：

- 配置runner（参考下面的管理命令）
- 清理工作目录

默认工作目录（定期清理）：

```shell
# fetch 模式下 基本不用清理
ls /home/gitlab-runner/builds
rm -rf /home/gitlab-runner/builds/*
```

**runner管理命令**

```shell
gitlab-runner register # 可多次注册为不同类型的runner
gitlab-runner verify # 检查所有runner状态
gitlab-runner list # 列出所有runner，token可用来删除runner

gitlab-runner start     # 启动
gitlab-runner stop      # 停止

# 删除runner
gitlab-runner unregister --url http://coding.xxx.cn/ --token xxxxxx
# 删除的另一种方式
gitlab-runner verify --delete -t xxxxxx -u http://coding.xxxx.cn/
```

### 其他步骤

Runner通常使用特定的`gitlab-runner`用户来执行任务，为了让runner能够执行特定命令，需要添加必要的环境变量添加到 `/.env`中，
设置gitlab-runner用户的`.bashrc`文件：

```shell
# root操作
echo 'source /.env' >> /home/gitlab-runner/.bashrc
chown gitlab-runner:gitlab-runner /home/gitlab-runner/.bashrc
```

手动登陆gitlab-runner用户，验证各命令是否正常执行；

去runner机器创建ssh公钥，然后添加到gitlab中，确认git命令正常执行。

> 或者使用官方推荐的方法（更复杂）：https://docs.gitlab.cn/jh/ci/ssh_keys/

### 其他帮助

- [注册Runner](https://docs.gitlab.cn/runner/register/)