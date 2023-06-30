package crontask

import (
	"github.com/robfig/cron/v3"
)

// doc: https://pkg.go.dev/github.com/robfig/cron
var cronMain *cron.Cron

func Start() {
	cronMain = cron.New()
	_, err := cronMain.AddFunc("0 0 0 * * *", computeAccostFlow) // 0点计算搭讪流量
	if err != nil {
		panic(err)
	}
}

func Stop() {
	cronMain.Stop()
}
