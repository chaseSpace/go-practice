package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/chaseSpace/lumberjack/v2"
)

func TestXXX(t *testing.T) {
	g := func() (int, error) {
		return 1, fmt.Errorf("111")
	}

	g2 := func() (err error) {
		v, err := g()
		_ = v
		return
	}
	println(111, g2)
}

func Test_lumberjack(t *testing.T) {
	// 创建日志目录
	if err := os.MkdirAll("./logs", 0755); err != nil {
		log.Fatal(err)
	}

	// 创建 lumberjack logger 配置
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 输出日志文件路径
		MaxSize:    1,                // 每个日志文件最大 10MB
		MaxBackups: 5,                // 保留 5 个备份
		MaxAge:     30,               // 保留 30 天
		Compress:   false,            // 启用压缩
	}

	// 创建多写入目标：同时写入文件和标准输出
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)

	// 创建日志实例
	fileLogger := log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// 模拟一些日志输出
	for i := 0; i < 100; i++ {
		fileLogger.Printf("处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求处理请求 %d", i)
		println(i)
	}

	// 关闭 logger
	defer lumberjackLogger.Close()
}
