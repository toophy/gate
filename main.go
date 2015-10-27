// main.go
package main

import (
	"fmt"
	"github.com/toophy/pangu/help"
	"github.com/toophy/pangu/thread"
	"os"
	"runtime"
	"runtime/pprof"
)

// Gogame framework version.
const (
	VERSION = "0.0.2"
)

func main() {

	runtime.GOMAXPROCS(1)

	// 检查log目录
	if !help.IsExist(thread.LogDir) {
		os.MkdirAll(thread.LogDir, os.ModeDir)
	}

	// 创建pprof文件
	f, err := os.Create(thread.LogDir + "/" + thread.ProfFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// 主协程

	// 等待结束
}
