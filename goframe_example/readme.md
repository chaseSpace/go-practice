## GoFrame 框架入门

> [!TIP]
> 点击网页右侧的目录按钮查看大纲，支持跳转。

这是一款功能相当大而全的 Go 基础开发框架，包含 web 开发以及微服务开发所需的各种组件/工具，基本不需要再去寻找或引入其他单个仓库。GoFrame
官方文档大而全，
第一次看可能会阅读困难，本文档将带你快速入门 GoFrame，部分主题提供官文链接，待你用到时再去了解。

显著特点：

- 丰富的开发组件
- 简单易用，详细的文档
- 跟踪和错误堆栈功能
- ORM 组件
- 工程设计规范
- CLI 工具
- 支持 protobuf
- OpenTelemetry、OpenAPIV3 等

资料索引：

- Github：https://github.com/gogf/gf
- 快速开始：https://goframe.org/pages/viewpage.action?pageId=57183742

### 准备工作

GoFrame 官方缩写`gf`，最新版本 v2.7.0，于 2024-4-8 发布，要求 Go >= 1.18.0。开始前，先创建一个 Go 项目（略），然后安装它：

```shell
# 安装主库
go get -u -v github.com/gogf/gf/v2

# 安装gf工具
go install github.com/gogf/gf/cmd/gf/v2@latest

# 获取帮助
gf gen -h
```

工具基本使用：

```shell
gf -v  # 会自动检测当前项目使用的GoFrame版本（解析go.mod）
gf up -cf  # 在升级框架代码库的同时，升级gf工具和修复本地代码在升级中的不兼容变更
gf init webapp # 创建一个项目，使用gf推荐的目录结构
gf init mymono -m # 创建一个微服务大仓

# 代码生成*  v2.5+
```

**Windows安装Make工具**

gf
为项目提供了Makefile来管理开发过程中的各种脚本命令，比单独使用gf命令要更高效。Windows下载[make安装程序](../bin/win/make-3.81.exe)，
其他系统自行搜索。

- [工具完整介绍](https://goframe.org/pages/viewpage.action?pageId=1114260&src=contextnavpagetreemode)

### 代码结构分层

不管是单服务还是微服务，都是五层架构，api/controller/logic/dao/model。

- api: 请求和响应定义。
- controller：业务接口定义，入参解析验证后传递给 service，以及对出参的维护。
  - 可以直接调用 dao 层实现业务逻辑，当认为逻辑可能会被多个接口复用时，再将逻辑下放到 service 层。
  - 可以调用一个或多个 service 方法来得到结果。
- logic（service）：业务逻辑实现，会通过 gf 工具生成 service 包。
  - logic 层的实现可以调用其他 service 包来完成逻辑。
  - service 包是根据 logic 层定义的方法生成对应的接口（interface）。
  - 若涉及多表访问（事务），则通过 tx 传参调用 dao 层不同方法。
- dao：数据库访问。
  - gf 工具为 dao 层生成一个个以表名命名的 go 文件（用到的表名在 config 中配置），开发者在每个文件中定义表相关的 dao 方法。
  - 方法实现应尽量简洁，让 service 层来多次调用不同 dao 方法得到结果，以避免对 dao 层的频繁修改。
- model：数据库表实体，但也可以是接口需要的其他结构体。
  - 包含`do`和`entity`两个 gf 工具维护的子目录，其中存放与实际数据表一致的代码结构体。
  - 在 model 层则存放人工定义的业务需要的模型结构。
  - model 层的实体可以在上面几层之间共享（除了 api 层），而非每一层定义不同的结构体。

请求入口为 controller，然后内部再依次调用 logic/dao/model。

### 工程目录设计

gf 规范了工程目录结构。

```shell
/
├── api # 对外提供服务的输入/输出数据结构定义。考虑到版本管理需要，往往以api/xxx/v1...存在。
├── hack # 存放项目开发工具、脚本等内容。例如，CLI工具的配置，各种shell/bat脚本等文件。
├── internal # 业务逻辑存放目录。通过Golang internal特性对外部隐藏可见性。（如api层不能调用internal下面的model包，但反之可以）
│   ├── cmd # 命令行管理目录。可以管理维护多个命令行。
│   ├── consts  # 项目所有常量定义。
│   ├── controller # 接收/解析用户输入参数的入口/接口层。
│   ├── dao
│   |   └── internal # 存放gf工具生成的dao对象，仅允许dao包内调用
│   ├── logic # 业务逻辑封装管理，特定的业务逻辑实现和封装。往往是项目中最复杂的部分。
│   ├── model # 数据结构管理模块，管理 api/dao/service 层所有的出入参数
│   |   ├── do # 自动生成。用于dao数据操作中业务模型与实例模型转换，由工具维护
│   │   └── entity # 自动生成。是模型与数据集合的一对一关系，由工具维护
│   └── service # 用于业务模块解耦的接口定义层。具体的接口实现在logic中进行注入。
├── manifest # 包含程序编译、部署、运行、配置的文件
│   ├── config # 配置文件存放目录。
│   ├── docker # Docker镜像相关依赖文件，脚本文件等等。
│   ├── deploy # 部署相关的文件。默认提供了Kubernetes集群化部署的Yaml模板，通过 kustomize 管理。
│   └── protobuf # GRPC协议时使用的protobuf协议定义文件，协议文件编译后生成go文件到api目录。
├── resource # 静态资源文件。这些文件往往可以通过 资源打包/镜像编译 的形式注入到发布文件中。
├── utility
├── go.mod
└── main.go # 程序入口文件。
```

#### 请求流转

- cmd：负责引导程序启动，显著的工作是初始化逻辑、注册路由对象、启动 server 监听、阻塞运行程序直至 server 退出。
- controller：接收 Req 请求对象后做一些业务逻辑校验，可以直接在 controller 中实现业务逻辑，或者调用一个或多个 service
  实现业务逻辑，将执行结果封装为约定的 Res 数据结构对象返回。
- model：管理了所有的业务模型，service 资源的 Input/Output 输入输出数据结构都由 model 层来维护。
- service：接口层，用于解耦业务模块，service 没有具体的业务逻辑实现，具体的业务实现是依靠 logic 层注入的。
- logic：业务逻辑实现，需要通过调用 dao 来实现数据的操作，调用 dao 时需要传递 do 数据结构对象，用于传递查询条件、输入数据。dao
  执行完毕后通过 Entity 数据模型将数据结果返回给 service 层。
- dao：通过框架的 ORM 抽象层组件与底层真实的数据库交互。

#### dao 层大改进

![](../img/gf_dao.png)

其中的`dao.DoctorUser`对象是`gf gen dao`自动生成的单表访问对象，工具还会生成具体的表结构体，位于`model/entity/`
，无需人工维护，保障了代码数据结构体与表模型的强一致性。

#### 关于 model

model 分为两种，一种是**数据模型**，是数据库表中的结构；另一种是**业务模型**，具体又分为**接口输入/输出模型** 与
**业务输入/输出模型**。

**数据模型**

它由 gf 工具自动生成到`internal/model/entity`。

**接口出入模型**

定义在 api 接口层中，供工程项目所有的层级调用，例如 controller, logic,
model 中均可以调用 api 层的输入输出模型。在 GoFrame 框架规范中，这部分输出输出模型名称以`XxxReq`和`XxxRes`格式命名。

**业务输入/输出模型**

用于服务内部模块/组件之间的方法调用交互，特别是 controller->service 或者 service->service 之间的调用。这部分模型定义在
model 模型层中,
如`internal/model/some_biz.go`，这部分输入输出模型名称通常以`XxxInput`和`XxxOutput`格式命名。

### 微服务大仓管理

官方建议当微服务数量少于 50 个时（匹配大部分项目），建议采用单仓库管理所有代码。代码结构设计上也是利用 go 的**internal 特性
**
来避免不同服务之间的代码引用，**服务之间仅允许 api 包的引用**。
参考[微服务大仓管理模式](https://goframe.org/pages/viewpage.action?pageId=87246764)。

### 接口化与泛型设计

GoFrame 框架的核心组件均采用了接口化设计，比如
gcfg/gcache/gredis/gsession/gdb，它们都是以接口形式提供服务，并且每个方法返回的对象都是一个泛型`gvar.Var`
，可以随意转换为其他类型，对于一些自定义类型，可以通过泛型对象的`Scan`方法转换为具体类型。

### 隐式与显式初始化

GoFrame 框架的很多模块都采用了隐式初始化。但对于业务代码，应该避免使用隐式查询，以免无法直观的了解项目的启动顺序。正确示例：

```shell
# go伪代码

function initBasis() {
  initConf()
  initLogger()
  initDB()
  initRedis()
  initKafka()
  ...
}

initBasis()
startServer()
```

### Context 传递

参考[Context: 业务流程共享变量](https://goframe.org/pages/viewpage.action?pageId=3672552)。

### 开始开发

#### 创建项目

```shell
gf init webapp
cd webapp
```

#### 定义 api

在`/api/模块/版本/定义文件.go`中定义具体的请求响应体，例如默认的`HelloReq`和`HelloRes`，
比如按照此范式定义，才能生成控制器代码。

请求/响应体可以通过 tag
定义字段注释，校验规则。具体参考[数据校验-校验规则](https://goframe.org/pages/viewpage.action?pageId=1114367)。

#### 生成控制器代码

```shell
# -merge是将 hello.go 中的多个 req&res 合并存放到 controller/ 的单个文件中，否则是一个req对应一个文件，不推荐。
gf gen ctrl -merge # or `make ctrl`，但make不支持参数输入，只能修改一下makefile，示例项目已修改
```

生成 `api/hello/hello.go`
的控制器接口代码，以及对应的控制器实现`controller/hello/hello.go| hello_new.go| hello_v1_hello.go`。
**其中的`hello.go`是一个空文件**，用于存放 controller 内部使用的变量、常量、数据结构定义，或者 init 方法定义等，官方建议用不到也留着先。。

注意：api目录中定义的 `req&res` 必须符合命名规范，否则生成的控制器代码会报红！

> [!IMPORTANT]
> 有 bug！概率性出现生成的 controller 代码没有 import 对应 api 结构体，需要手动处理。导致报红！v2.7.0

> [!IMPORTANT]
> v2.7.0版本的gf工具，对于api目录下已经注释掉的请求/响应结构体，也会生成控制器代码，已告知官方，回复说在 v2.7.1 修复。

#### 生成 dao 代码

`gen dao`命令是 CLI 中最频繁使用的命令。但是这个命令的参数很多，一般通过`hack/config.yml`管理。参考：

```yaml
gfcli:
  gen:
    dao:
      # link定义数据库类型和连接信息
      - link: "mysql:root:123@tcp(127.0.0.1:3306)/test"
        tables: "users" # 指定当前数据库中需要执行代码生成的数据表。如果为空，表示数据库的所有表都会生成!!!
        jsonCase: "CamelLower" # 指定model中生成的数据实体对象中json标签名称规则，还支持 Snake 等
        group: default
        stdTime: true  # 当数据表字段类型为时间类型时，代码生成的属性类型使用标准库的time.Time而不是框架的*gtime.Time类型。
        withTime: true # 为每个自动生成的代码文件增加生成时间注释
        gJsonSupport: true # 当数据表字段类型为JSON类型时，代码生成的属性类型使用*gjson.Json类型。
        importPrefix: "" # 用于指定生成Go文件的import路径前缀
        descriptionTag: false # 用于指定是否为数据模型结构体属性增加desription的标签，内容为对应的数据表字段注释。一般不需要开，因为有注释了。
        typeMapping: # 自定义类型映射，从版本v2.5开始支持。
          decimal:
            type: float64

        #      - link: "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
        #        path: "./my-app" # 生成dao和model文件的存储目录地址
        #        prefix: "primary_" # 生成的文件名、对象名的前缀
        #        tables: "user, userDetail"
        #        group: order
        #        noJsonTag: true # 生成的数据模型中，字段不带有json标签
        #        noModelComment: true # 不要注释
        #        clear: true # !!!自动删除数据库中不存在对应数据表的本地dao/do/entity代码文件。

      # sqlite需要自行编译带sqlite驱动的gf，下载库代码后修改路径文件（gf\cmd\gf\internal\cmd\cmd_gen_dao.go）的import包，取消注释即可。sqlite驱动依赖了gcc
#      - link: "sqlite:./file.db"
#        group: order
#        tablesEx: "order" # 指定当前数据库中需要排除代码生成的数据表
```

生成代码：

```shell
model/do/users.go  # 不关心
model/entity/users.go # 自动生成的 users 表结构体，可以被 controller/logic/dao 层引用

dao/users.go   # 用于crud users表的DAO对象
dao/internal/user.go
```

参考 [数据规范-gen dao](https://goframe.org/pages/viewpage.action?pageId=3673173)。

#### 生成 service 代码

service 代码是一层接口抽象，从 logic 抽象而来。操作顺序是先写 logic 代码，然后生成 service 代码，然后 controller 和其他
service 可以调用它。

定义 logic 层方法需要满足范式`logic/xxx/*.go`，即只能有一个二级目录，否则 gf 无法解析。

步骤：

- 定义`internal/logic/menu/user.go`
  - 可能需要定义每个方法的`*Input`和`*Output`结构体。
- 生成 service 方法：`gf gen service` 或 `make service`
  - 生成文件：`service/menu.go`
- 在 logic 代码中编写 init 函数将方法实现注册到 service，参考[logic/menu/user.go](webapp/internal/logic/menu/user.go)
  - 生成的`logic/logic.go`会自动引入`logic/menu`包，确保 service 接口的注册，而 logic 包又会被 main 文件引入。

然后就可以在 controller 通过 service 包来调用 logic 代码了，例如：`service.User().GetUser(ctx, &menu.GetUserInput{})`
，当然，其他 service 包也可根据需要调用。

参考[数据规范-gen service](https://goframe.org/pages/viewpage.action?pageId=49770772)。

#### 生成 enums

gf 工具扫描 api 目录中的请求/响应结构体中使用的全部枚举，然后生成 swagger 文档所需的 go 文件，在项目启动后的 swagger
文档中可以看到字段枚举值。

步骤如下：

- 定义一个 api 请求&响应，其中包含 const 枚举，参考[hello.go](webapp/api/hello/v1/hello.go)中的`Hello3Req`；
- 然后执行`gf gen enums` 或 `make enums`，生成`internal/boot/enums.go`；
- 首次生成 boot 目录时，将其添加到 main.go 中初始化即可。

然后`gf run main.go`启动项目，访问输出中的 swagger 地址，在 web 端可以看到接口`Hello3Req`中的字段 PodState 列出了枚举值。

#### 生成 pb

如果你的项目用不到 protobuf，忽略此章节。

识别 proto 文件，生成对应的 pb Go 文件。这个功能建议不要用 gf 提供的`gen pb`
命令，因为它没有集成 protoc 执行文件，需要自行下载，而这就会产生两者版本的兼容性问题（`gf gen pb`命令执行报错）。

步骤如下：

- 下载 protoc 文件，并安装插件：protoc-gen-go, protoc-gen-go-grpc；
  - 可直接在[这个页面](https://github.com/chaseSpace/go-microsvc-template/tree/main/tool_win/protoc_v24)下载全部。
- 定义`manifest/protobuf/svc_a/a.proto`文件；
- 使用 protoc 命令生成代码到指定位置；

protoc 命令如下：

```shell
mkdir -p ./api/pb
OUTPUT_DIR=./api/pb
SRC=./manifest/protobuf/
protoc -I ./manifest/protobuf \
      --go_out=$OUTPUT_DIR \
      --go_opt=paths=source_relative \
      $SRC/*/*.proto $SRC/*.proto
```

pb 文件中引入了 protobuf 库，记得拉取：`go get -u google.golang.org/protobuf`
。笔者修改了 Makefile 文件，用以上命令替换了原来的`gen pb`脚本，参考 [makefile 引用的 hack.mk](webapp/hack/hack.mk)
中的`pb`部分。

> [!IMPORTANT]
> 还剩下一个问题，由于是使用原生 protoc 命令生成代码，所以控制器的代码需要自己写。

#### 生成 pbentity

之前的 dao 生成只是将数据库中的表模型生成 go 代码结构体，即`model/entity`
，它主要是给 Go 程序使用，但在异构服务或前后端使用 pb 协作的场景中，则需要通过 proto 文件来共享数据模型。这时候
go 代码的 entity 就没法使用了，所以 pbentity 就派上用场。

步骤：

- 在`hack/config.yaml`中加入`gfcli.gen.pbentity`配置，其中包含数据库连接信息、生成目标位置等；
  - 确保数据库表存在，能连接。
- 执行：`gf gen pbentity` 或 `make pbentity`，然后可以观察到目标位置`manifest/protobuf/pbentity`包含新的 proto 文件，其中包含表模型。

#### 运行

```shell
gf run main.go # 或make all
```

使用**gf**工具启动的好处是会自动监控项目中所有 go 文件变化，然后进行自动编译重启。同时可以在 `hack/config.yaml`
中加入`gfcli.run`配置，其中可以包含二进制存放位置、go 编译参数等。

#### Build

编译 Go 程序。

- 在`hack/config.yaml`中加入`gfcli.build`配置；
- 执行`gf build` 或 `make build`；

#### 打包镜像

`gf docker`命令默认使用`manifest/docker/Dockerfile`将已经编译好的二进制打包到 alpine
镜像。相关参数通过 `hack/config.yaml`
中的`gfcli.docker`配置。

`make image`是对该命令的封装，自动上传是`make image.push`。

### 资源管理

这个功能不是常用的，只有在项目运行用到非 go 静态文件时（比如 ip2region 文件等）才可能用到，它只是一种将静态文件打包到 go
执行文件中的方法。

参考[这个页面](https://goframe.org/pages/viewpage.action?pageId=1114671)。

### 日常开发文档索引

查看[核心组件](https://goframe.org/pages/viewpage.action?pageId=1114409)页面。

### 其他工具

GoFrame 提供了极其丰富的日常工具使用，包括但不限于数据结构、定时器、锁、上下文、日志管理、时间转换、缓存、类型转换、编码和加密等。

需要时在[这个页面](https://goframe.org/pages/viewpage.action?pageId=1114411)搜索。

### 微服务开发

[这个页面](https://goframe.org/pages/viewpage.action?pageId=77852968)。

### Web 开发

Demo 项目是一个 HTTP 服务，在[cmd.go](webapp/internal/cmd/cmd.go)中启动了一个 HTTP 服务器，其中使用了 gf 提供的 HTTP 框架。
由`ghttp`模块实现。实现了丰富完善的相关组件，例如：Router、Cookie、Session、路由注册、配置管理、模板引擎、缓存控制等等，
支持热重启、热更新、多域名、多端口、多实例、HTTPS、Rewrite、PProf 等等特性。

参考[Web 框架文档](https://goframe.org/pages/viewpage.action?pageId=1114405&src=contextnavpagetreemode)。

### 其他

GoFrame 还提供了 TCP/UDP、Websocket 协议组件、可观测性，根据实际需求探索使用。