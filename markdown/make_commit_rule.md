# 实施 Git 提交规范

实施 Git 提交规范具有多方面的好处，特别是对于团队合作和项目管理。以下是一些关键的好处：

- 提高代码可维护性
- 增强团队协作
- 简化代码审查和合并
- 改进版本控制
- 支持自动化工具
- 提升项目质量
- 便于追踪和审计

可以使用一些三方工具来帮助实施，其中包含仓库插件和 IDE 插件，团队成员习惯后效率将提高！

## 1. commit message 格式

`<type>(<scope>): <subject>`，**注意首字母小写，末尾不含句号或点**。

- type(必须)：用于说明 git commit 的类别，只允许使用下面的标识。
    - feat：新功能（feature）。
    - fix/wip: fix 是一次提交就完全修复了 bug，wip 表示一次修复分多次提交（一般是因为修改过多，最后一次使用 fix）。
    - docs：文档（documentation）。
    - style：格式（不影响代码运行的变动）。
    - refactor：重构（即不是新增功能，也不是修改 bug 的代码变动）。
    - perf：优化相关，比如提升性能、体验。
    - test：增加测试。
    - chore：构建过程或辅助工具的变动，对代码库的非功能性更改。
    - revert：回滚到上一个版本。
    - merge：代码合并。
    - ci：自动化流程配置修改。
    - build：构建脚本修改，可以包含在`ci`中。
    - sync：同步主线或分支的 Bug。
- scope(可选)：用于说明 commit 影响的范围，比如数据层、控制层、视图层等。
    - 例如在 Angular，可以是 location，browser，若有多个，也可用星号代替。
- subject(必须)：commit 目的的简短描述，不超过 50 个字符。
    - 中文项目使用中文。

示例（chatGPT 提供）：

```plain
feat(auth): 添加用户注册功能
fix(login): 修复登录页面无法响应的问题
perf(database): 优化数据库查询性能
docs(installation): 更新安装指南
refactor(auth): 重构用户认证逻辑
test(auth): 添加用户注册页面的端到端测试
build(dependencies): 更新依赖库版本
style(ui): 修正缩进和空格
revert: feat(auth): 添加用户注册功能
- This reverts commit 667ecc1654a317a13331b17617d973392f415f02.
chore(editor): 更新编辑器配置
```

## 2. 使用三方工具强制实施

若没有，则规范形同虚设。实施者必须先让团队成员认识到规范的优点、能带来的收益。

第三方工具将对提交信息进行规范检查，不符合规范，不能提交。还可以选择使用 webhook 通过发送钉钉/飞书群警告的形式进行监控，督促大家按照规范执行代码提交。
进一步，还可以加入大代码量提交监控和删除文件监控，减少研发的代码误操作。

### 2.1 软拦截

指的是通过项目根目录下的`.git`目录中添加 git hooks 文件，来在客户端本地提交时做检查和**拦截**。
客户端 hook 有 pre-commit、prepare-commit-msg、commit-msg、post-commit 等，服务端有 pre-receive、post-receive、update，具体参考[自定义 Git - Git 钩子][0]。

每个 git hook 是一个 shell 脚本，在对应时机执行，如果 exit
code 非 0，则终止下一步。所以原则上在客户端使用 hook 也能达到规范目的。参考本项目的[.githooks](../.githooks)。

### 2.2 硬拦截

在 git 服务端执行才是真正做到强制实施，因为开发者可以临时修改客户端 hook 以允许不规范的 Commit。主要使用`pre-receive`这个 hook 来拦截。

## 3. 插件推荐

**帮助填写 CommitMsg**

- Goland 插件：CommitMessage
- VSCode 插件：Commitize

**CommitMsg 检查**

- [Commitize](https://commitizen-tools.github.io/commitizen/)
- [Commitlint](https://commitlint.js.org/guides/getting-started.html)

这些插件都可以在 Client 或 Server 端安装，Server 端安装后可以将命令配置到 ServerHook 中。

## 其他 Git 规范

分支使用、tag 规范、Issue 等。

## 参考

- https://zhuanlan.zhihu.com/p/182553920

[0]: https://git-scm.com/book/zh/v2/自定义-Git-Git-钩子

