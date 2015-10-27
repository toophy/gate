// main.go
package main

import (
	"github.com/toophy/gate/config"
	"github.com/toophy/gate/help"
	"github.com/toophy/gate/logic"
)

// Gogame framework version.
const (
	VERSION = "0.0.2"
)

func main() {
	help.GetApp().Start(config.LogDir, config.ProfFile)

	// 主协程
	go logic.Main_go()

	// 等待结束
	help.GetApp().WaitExit()
}
