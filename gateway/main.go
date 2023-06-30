package main

import (
	"service/gateway2/gutil"
	"service/gateway2/tcp_impl"
	"utils/config"
)

// 规范且优雅的TCP长连接服务器实现
// 包含规范化的 router, conn_manager, connection, request, datapack等模块实现，具体查看 iface目录下的接口定义
// 每层定义接口的原因其一是规范且标准化，更重要的其二是让每个模块的功能更清晰，不至于随意给每个模块增减方法，每个模块的用途在编写前都已设计明确
func main() {
	config.InitDBConfig()
	config.InitDedicateConf(&config.Cfg.Gateway, "gateway.json")
	gutil.MustInit()

	s := tcp_impl.NewServer("winter")

	// 所有router写在这里，一目了然
	s.AddRouter(0, new(tcp_impl.PingRouter))

	s.Serve()
}
