// main.go
package main

import (
	"github.com/toophy/gate/config"
	"github.com/toophy/gate/help"
)

// Gogame framework version.
const (
	VERSION = "0.0.2"
)

func main() {
	help.GetApp().Start(config.LogDir, config.ProfFile)

	// 主协程
	go main_go()

	// 等待结束
	help.GetApp().WaitExit()
}

func main_go() {
	RegMsgProc()

	go help.GetApp().Listen("tcp", ":8001")
}

func RegMsgProc() {
	help.GetApp().RegMsgFunc(1, on_c2g_login)
}

func on_c2g_login(c *help.ClientConn) {
	if c.Id > 0 {
		name := c.Stream.ReadStr()
		println(name)
	}
}
