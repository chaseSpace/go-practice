# Plan 1：创建 Go 与 Eino 学习项目

## 本节只解决什么

确认项目能编译、能导入 Eino，并知道如何运行一个指定测试。完成本节不需要 API Key，也不调用任何外部模型。

## 理论基础

Go 的模块由 `go.mod` 描述：模块路径、Go 版本和直接依赖都记录在这里，校验和记录在 `go.sum`。本项目把 Eino 核心库和 OpenAI 适配器列为依赖；
后面的练习会使用核心库的 `schema`、Prompt、Tool 和 ADK，也会在 Plan 2 构造 OpenAI ChatModel。

先把“能导入框架”与“能调用模型”分开。前者是本地编译问题，后者还涉及密钥、网络、模型权限和费用。分开后，绝大多数学习测试可稳定离线运行。

`schema.Message` 是后续练习最常用的数据类型。它不是模型本身，而是一条带角色和内容的对话记录。Plan 1 只验证它能被创建。

## 动手步骤

在仓库根目录执行：

```bash
go test ./...
go test ./testfiles -run TestPlan1ProjectSetup -v
```

打开 `testfiles/plan1_test.go`。测试创建 `schema.UserMessage("你好，Eino")`，并断言角色为 `schema.User`。这说明 Go 能解析依赖，也说明你已经接触到后续所有模型调用都会使用的消息类型。

## 容易混淆的点

- `go test` 成功不代表 API Key 有效；它只说明本地代码和离线测试正常。
- 不要把 `go get` 当成每次运行的必要动作。依赖已写入 `go.mod` 后，Go 会按模块文件解析。
- `main.go` 以后可以作为手工运行入口；本教程首先使用测试，让每一步都有自动断言。

## 本节验收

`TestPlan1ProjectSetup` 通过，且你能用 `-run` 只运行它。下一节才连接真实模型。
