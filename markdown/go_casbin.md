## Go 使用 casbin

Casbin是一个强大的、高效的开源访问控制框架，其权限管理机制支持多种访问控制模型，ACL/RBAC/ABAC。

### casbin 特性

- 支持自定义请求的格式，默认的请求格式为{subject, object, action}
- 具有访问控制模型model和策略policy两个核心概念
- 支持RBAC中的多层角色继承，不止主体可以有角色，资源也可以具有角色
- 支持超级用户，如 root 或 Administrator，超级用户可以不受授权策略的约束访问任意资源
- 支持多种内置的操作符，如 keyMatch，方便对路径式的资源进行管理，如 /foo/bar 可以映射到 /foo*

文档：

- https://github.com/casbin/casbin
- https://casbin.org/docs/zh-CN/overview

### 原理

在 Casbin 中, 访问控制模型被抽象为基于 PERM (Policy, Effect, Request, Matcher) 的一个文件。 因此，切换或升级项目的授权机制与修改配置一样简单。

#### PERM模型

PERM(Policy, Effect, Request, Matchers)模型很简单, 但是反映了权限的本质 – 访问控制。

- Policy: 定义权限的规则
- Effect: 定义组合了多个 Policy 之后的结果, allow/deny
- Request: 访问请求, 也就是谁想操作什么
- Matcher: 判断 Request 是否满足 Policy

#### 核心配置 - model file

casbin 是基于 PERM 的, 所以 model file 中主要就是定义 PERM 4 个部分。

1. **请求定义**

```
[request_definition]
r = sub, obj, act  // 表示casbin接受的一个请求由 sub, obj, act 3 个部分组成
```

2. **策略定义**

```
[policy_definition]
p = sub, obj, act  // 定义一个策略叫p，匹配请求中的sub, obj, act
p2 = sub, act // 定义一个策略叫p2，匹配请求中的sub, act，表示 sub 所有的资源都能执行 act
p3 = sub, obj, act, eft // 还可以支持 指定 eft 为 deny 或 allow的策略
```

这个策略定义最终会应用到外部导入的多条 policy rule 上。

3. **策略效果**

```
[policy_effect]
e = some(where (p.eft == allow)) // 表示有任意一条 policy rule 满足, 则最终结果为 allow
```

4. **匹配规则**

```
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act // 最常见的匹配模型，定义了 request 和 policy 匹配的方式, p.eft 是 allow 还是 deny, 就是基于此来决定的
```

5. **角色定义**

仅当你的系统中的存在角色继承关系时（如RBAC模型），才需要定义角色模型。

```
[role_definition]
g = _, _
g2 = _, _
g3 = _, _, _
```

g, g2, g3 表示不同的 RBAC 体系,一般选择其一， `_, _` 表示用户和角色，`_, _, _` 表示用户, 角色, 域。

#### 如何工作

casbin服务首先加载model和policy rules，然后根据所提供的请求信息进行判断：

- 通过 matchers 将req 与 具体的policy rule 进行匹配，将匹配结果输入 policy_effect
- matchers可以匹配N个policy rule，并将这些rule的 effect，输入给 policy_effect，上面的`e = some(where (p.eft == allow))`
  表示有任何一条 policy rule 的eft为allow（若未匹配中rule，则eft是deny），则最终结果为 allow。
    - casbin支持多种策略效果，参考[policy_effect](https://casbin.org/zh/docs/syntax-for-models#策略效果)。

### 支持多种模型

- ACL
- 带超级用户的ACL
- RBAC
- 带资源角色的RBAC：用户和资源同时可以拥有角色（或组）
- 带有域/租户的RBAC：用户可以为不同的域/租户拥有不同的角色集
- ABAC (基于属性的访问控制)：可以使用类似"resource.Owner"的语法糖来获取资源的属性。
- RESTful：支持像"/res/*"，"/res/:id"这样的路径，以及像"GET"，"POST"，"PUT"，"DELETE"这样的HTTP方法。
- 拒绝优先：同时支持允许和拒绝授权，其中拒绝优先于允许。
- 优先级：策略规则可以设置优先级，类似于防火墙规则。

参考页面：https://casbin.org/zh/docs/supported-models

### ACL 模型

Casbin中最基本、最简单的model是ACL。ACL中的model CONF为:

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
# 若是带有root用户的ACL，则在上面的match rule 后追加：|| r.sub == "root"
```

`policy_rule.csv`:

```
p, alice, data1, read  // 第一条规则：alice 可以对 data1 进行 read 操作
p, bob, data2, write // 第二条规则：bob 可以对 data2 进行 write 操作
```

实际项目中的 policy rule 可能会有很多，通常存储在数据库中，并且不由casbin管理。

### RBAC模型

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _ // 用户和角色

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

`policy_rule.csv`:

```csv
p, alice, data1, read
p, bob, data2, write
p, data2_admin, data2, read
p, data2_admin, data2, write

g, alice, data2_admin
```

### model 加载

由于model只是一串文本，因此可以直接通过字符串加载，也可以通过文件加载。

```shell
# csv是策略文件，实践中一般是来自数据库
e := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
```

在测试时，也可通过代码加载：

```shell
// Initialize the model from Go code.
m := casbin.NewModel()
m.AddDef("r", "r", "sub, obj, act")
m.AddDef("p", "p", "sub, obj, act")
m.AddDef("g", "g", "_, _")
m.AddDef("e", "e", "some(where (p.eft == allow))")
m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

// Load the policy rules from the .CSV file adapter.
// 使用自己的 adapter 替换。
a := persist.NewFileAdapter("examples/rbac_policy.csv")

// 创建一个 enforcer。
e := casbin.NewEnforcer(m, a)
```

#### 从字符串加载

```go
// Initialize the model from a string.
text := `
    [request_definition]
    r = sub, obj, act
    
    [policy_definition]
    p = sub, obj, act
    
    [role_definition]
    g = _, _
    
    [policy_effect]
    e = some(where (p.eft == allow))
    
    [matchers]
    m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
`
m := NewModel(text)

// Load the policy rules from the .CSV file adapter.
// Replace it with your adapter to avoid files.
a := persist.NewFileAdapter("examples/rbac_policy.csv")

// Create the enforcer.
e := casbin.NewEnforcer(m, a)
```

### match函数

你可以在匹配器中指定函数，使其更强大。你可以使用内置函数或指定你自己的函数。具体可以参考[官方文档](https://casbin.org/zh/docs/function)。

keyMatch是Casbin库中的一种匹配函数，用于匹配请求对象（r.obj）是否符合策略中的对象（p.obj）。

keyMatch规则支持两种通配符：

- *：匹配任意长度的字符串，但不包括路径分隔符/。
- **：匹配任意长度的字符串，包括路径分隔符/。

通过使用keyMatch规则，可以进行Restful 风格的路径匹配。它主要是在model中的`matchers`块中使用：

```shell
[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
```

加载以下策略规则：

```shell
p, alice, /alice_data/*, GET
p, cathy, /cathy_data, (GET)|(POST)
```

那么就有`alice`可以`GET`路径`/alice_data/r2`；`cathy`可以`GET`或`POST`路径`/cathy_data`。

除了 `keyMatch`，Casbin 还提供了其他匹配函数，如：keyMatch2/keyMatch3... 在上面的官方文档中查询详情。

## 参考

- https://www.hanhandato.top/archives/golang-casbin