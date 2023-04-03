## Protobuf 协议目录结构设计指引

这里直接以具备代表性的微服务背景来进行说明。

假设有两个微服务：user、 product，推荐以如下proto目录结构设计：
```shell
proto_src
├── common
│   └── common.proto
├── product.svc.proto
└── user.svc.proto
```

几个要点：
- 推荐一个微服务使用一个proto文件来管理，在命名上进行服务的标识。而不是一个微服务一个子目录来管理，那样会产生很多子目录，增加维护的复杂度。
- 为了避免`common`部分的协议引用了业务协议，可以设计一个子目录`common/`来存放公共协议。

### 关于protoc命令
从`build_mac.sh`中可以看到protoc命令使用的`--go_out=paths=source_relative:./proto_dst`这样的相对路径来生成代码，这样的目的是在`proto_dst/`下
生成与`proto_src/`相同目录结构的代码文件。

proto文件中的`option go_package=...`能够保证在proto文件之间存在import时，生成的代码中的import路径也是正确的。


## 运行示例
当前所在的`./proto_example`目录包含了运行脚本所需的全部源（protoc编译器、proto协议文件、protoc提供的常用协议规范实现`protoc-v3.5.1/`）
```shell
sh build_mac.sh
```