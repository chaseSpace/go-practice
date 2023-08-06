package gutil

import (
	"encoding/json"
	"fmt"
	"util/common"
	"utils/config"
)

func GenConnId() string {
	part1 := common.RandomBytes2(1, 4)
	part2 := common.RandomBytes2(0, 4)
	return "conn-" + string(part1) + "-" + string(part2)
}

var Conf *config.GatewayCfg

func MustInit() {
	Conf = config.Cfg.Gateway
	if Conf == nil {
		panic("conf is nil!")
	}

	b, _ := json.MarshalIndent(Conf, "", " ")
	fmt.Println("Gateway Conf:\n", string(b))
}
