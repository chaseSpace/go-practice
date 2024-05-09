package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	"go.uber.org/zap"
	"os"
)

type JsonOutputsForLogger struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

// LoggingJsonHandler is a example handler for logging JSON format content.
var LoggingJsonHandler glog.Handler = func(ctx context.Context, in *glog.HandlerInput) {
	jsonForLogger := JsonOutputsForLogger{
		Time:    in.TimeFormat,
		Level:   gstr.Trim(in.LevelFormat, "[]"),
		Content: gstr.Trim(in.Content),
	}
	jsonBytes, err := json.Marshal(jsonForLogger)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		return
	}
	in.Buffer.Write(jsonBytes)
	in.Buffer.WriteString("\n")
	in.Next(ctx)
}

func main() {
	fmt.Printf("%+.2f", 1.)
	type User struct {
		Name string
	}

	var xx = []*User{{Name: "x1"}, {Name: "x2"}}

	logger, _ := zap.NewProduction()
	defer logger.Sync() // zap底层有缓冲。在任何情况下执行 defer logger.Sync() 是一个很好的习惯
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// 字段是松散类型，不是强类型
		"user", xx,
	)

	fmt.Printf("%+v", xx)

	glog.SetDefaultHandler(LoggingJsonHandler)

	ctx := context.Background()
	g.Log().Debug(ctx, "Debugging...")
	glog.Warning(ctx, xx)
}
