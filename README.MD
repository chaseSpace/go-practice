## Go-Practice

### 关于vendor
2023在使用Go mod的同时仍然使用**vendor**是为了保证变更环境时依然能够使用包源失效的pkg。


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