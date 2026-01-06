## Go-Practice

### 关于vendor

2026 在使用Go mod的同时仍然使用**vendor**是为了保证变更项目环境时依然能够从容应对使用包源失效的情况。


### 拉取仓库后

```shell
sh build.sh
```

### 关于代码

- 每个目录都是不同的主题，目录中会包含 `xxx_test.go` 来执行目录下的示例代码。
- 根目录的`main.go`作用不大。

### CI

ci相关文件

```
├── .githooks
│   ├── pre-commit
│   └── pre-push
├── .golangci.yml
├── golangci-lint.sh
└── init.sh
```

拉取项目后，先执行`init.sh`配置环境，然后可以使用。

### 文档

- [iptables 详解](markdown/iptables.md)
- [实施 Git 提交消息规范](markdown/make_commit_rule.md)
- [openssl 使用示例](markdown/openssl_example.md)
- [Linux 实现端口转发](./markdown/port_forward.md)
- [Linux上使用Clash](markdown/use_clash_linux.md)
- [使用 letsEncrypt](markdown/use_letsEncrypt.md)
- [使用 pm2 管理服务](markdown/use_pm2.md)
- [GoFrame 框架入门](goframe_example/readme.md)