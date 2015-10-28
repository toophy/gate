package help

import (
	"bytes"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// 消息函数类型
type MsgFunc func(*ClientConn)
type ConnRetFunc func(string, string, int, string) bool

type AppBase struct {
	baseGoNumStart int
	baseGoNumEnd   int
	Listener       map[string]*ListenConn // 本地侦听端口
	RemoteSvr      map[string]*ClientConn // 远程服务连接
	Conns          map[int]*ClientConn    // 连接池
	ConnLast       int                    // 最后连接Id
	MsgProc        []MsgFunc              // 消息处理函数注册表
	MsgProcCount   int                    // 消息函数数量
	log_Buffer     []byte                 // 线程日志缓冲
	log_BufferLen  int                    // 线程日志缓冲长度
	log_TimeString string                 // 时间格式(精确到秒2015.08.13 16:33:00)
	log_Header     [LogMaxLevel]string    // 各级别日志头
	log_FileBuff   bytes.Buffer           // 日志总缓冲, Tid_world才会使用
	log_FileHandle *os.File               // 日志文件, Tid_world才会使用
}

// 程序控制核心
var app *AppBase

func GetApp() *AppBase {
	if app == nil {
		app = &AppBase{}
		app.init()
	}
	return app
}

// App初始化
func (this *AppBase) init() {
	this.Listener = make(map[string]*ListenConn, 10)
	this.RemoteSvr = make(map[string]*ClientConn, 10)
	this.Conns = make(map[int]*ClientConn, 1000)
	this.MsgProc = make([]MsgFunc, 8000)

	this.ConnLast = 1

	// 日志初始化
	this.log_Buffer = make([]byte, LogBuffMax)
	this.log_BufferLen = 0

	this.log_TimeString = time.Now().Format("15:04:05")
	this.MakeLogHeader()

	this.log_FileBuff.Grow(LogBuffSize)

	if !IsExist(LogFileName) {
		os.Create(LogFileName)
	}
	file, err := os.OpenFile(LogFileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	this.log_FileHandle = file
	this.log_FileHandle.Seek(0, 2)
	// 第一条日志
	this.LogDebug("\n          %s服务器启动\n", AppName)

	go this.Flush_log()
}

// 程序开启
func (this *AppBase) Start() {

	runtime.GOMAXPROCS(1)

	// 检查log目录
	if !IsExist(LogDir) {
		os.MkdirAll(LogDir, os.ModeDir)
	}

	// 创建pprof文件
	f, err := os.Create(LogDir + "/" + ProfFile)
	if err != nil {
		this.LogWarn(err.Error())
	}
	pprof.StartCPUProfile(f)
	this.baseGoNumStart = runtime.NumGoroutine()
}

// 等待协程结束
func (this *AppBase) WaitExit() {

	this.baseGoNumEnd = this.baseGoNumStart
	if runtime.NumGoroutine() > this.baseGoNumStart {
		this.baseGoNumEnd = this.baseGoNumStart + 1
	}

	for {
		<-time.Tick(2 * time.Second)
		if runtime.NumGoroutine() == this.baseGoNumEnd {
			pprof.StopCPUProfile()
			this.LogInfo("bye bye.")
			break
		}
	}

	// 关闭日志文件
	if this.log_FileHandle != nil {
		this.log_FileHandle.Close()
	}
}
