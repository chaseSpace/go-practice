package main

import (
	_ "webapp/internal/packed"

	_ "webapp/internal/boot"
	_ "webapp/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"webapp/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
