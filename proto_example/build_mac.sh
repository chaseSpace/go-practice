#!/bin/bash
# 安装protoc
# protoc是C++实现，所以不能使用go get安装
# https://github.com/google/protobuf/releases 指定平台、版本下载

# 两个插件安装 （默认最新）
# go env -w GOPROXY=https://proxy.golang.org,direct   (不能用http://goproxy.io)
# go get -d google.golang.org/protobuf/cmd/protoc-gen-go@latest   # 指定版本 @v1.30.0，失败则多次几次
# go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc

<<comment
# protoc-gen-go 和 protoc-gen-go-grpc 插件之间的区别
  早期：早期版本的protoc-gen-go由托管在 github.com/golang/protobuf，当时的版本支持同时生成pb msg的结构体、序列化方法、service interfaces + pb msg的grpc代码
      那时还是使用 --go_out=plugins=grpc 的flag。
      所谓的gRPC代码指的是：gRPC服务所需的server/client stub代码，这些代码允许你快速创建grpc server/client，以快速尝鲜grpc
  后来：较新的版本托管在 google.golang.org/protobuf（github镜像仓库：https://github.com/protocolbuffers/protobuf-go，不能使用go get下载）
      此仓库托管的插件版本重新划分了生成代码的职责，proto-gen-go 不再生成grpc代码（仅含msg的结构体、序列化方法、service interfaces）；
      proto-gen-go使用的flag也随之更新：--go_out=.
      而grpc部分代码转由 google.golang.org/grpc/cmd/protoc-gen-go-grpc 生成,
      proto-gen-go-grpc使用的flag是：--go-grpc_out=.

      二者可以同时工作：protoc -I. --go-grpc_out=. --go_out=. *.proto
comment


../bin/mac/protoc --version
../bin/mac/protoc-gen-go --version
../bin/mac/protoc-gen-go-grpc --version

export PATH=$PATH:./bin
echo '生成pb...'


# 注意这里不需要指定输出位置，因为全路径已经在所有proto文件中定义好（也是推荐的做法）
# --go_out=plugins=grpc:.是旧版的protoc flag，已不再支持
# 重点：使用source_relative将基于设置的相对路径生成不带有嵌套目录结构的代码，并且在pb内存在引用的同时保证生成的代码也是正确import路径（按照proto文件中的go_package那样）

../bin/mac/protoc -I=./proto_src -I=./protoc-v3.5.1 --go-grpc_out=paths=source_relative:./proto_dst --go_out=paths=source_relative:./proto_dst $(find ./proto_src/ -name '*.proto')

<<comment
除此之外，官方支持如下方式为特定文件指定 import path,  --go_opt=M${PROTO_FILE}=${GO_IMPORT_PATH}
protoc --proto_path=src \
  --go_opt=Mprotos/buzz.proto=example.com/project/protos/fizz \
  --go_opt=Mprotos/bar.proto=example.com/project/protos/foo \
  protos/buzz.proto protos/bar.proto
comment