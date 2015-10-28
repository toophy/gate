// main.go
package main

import (
	"github.com/toophy/gate/app"
)

// Gogame framework version.
const (
	VERSION = "0.0.2"
)

func main() {
	if app.GetApp().Start(100, app.Evt_lay1_time) {
		// 主协程
		go app.Main_go()
		// 等待结束
		app.GetApp().WaitExit()
	}
}
