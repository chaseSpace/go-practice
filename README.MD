## Go-Practice

### 关于vendor
2023在使用Go mod的同时仍然使用**vendor**是为了保证变更环境时依然能够使用包源失效的pkg。


### 拉取仓库后

```shell
go mod tidy
go mod vendor  将下载的包保存到 /vendor 才能正常引用
```