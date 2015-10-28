package help

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

const (
	LogDebugLevel = 0                // 日志等级 : 调试信息
	LogInfoLevel  = 1                // 日志等级 : 普通信息
	LogWarnLevel  = 2                // 日志等级 : 警告信息
	LogErrorLevel = 3                // 日志等级 : 错误信息
	LogFatalLevel = 4                // 日志等级 : 致命信息
	LogMaxLevel   = 5                // 日志最大等级
	LogLimitLevel = LogInfoLevel     // 显示这个等级之上的日志(控制台)
	LogBuffMax    = 20 * 1024 * 1024 // 日志缓冲
)

const (
	AppName     = "gate"
	LogBuffSize = 10 * 1024 * 1024
	LogDir      = "../log"
	ProfFile    = AppName + "_prof.log"
	LogFileName = LogDir + "/" + AppName + ".log"
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

func (this *AppBase) AddConn(c *ClientConn) {
	this.Conns[this.ConnLast] = c
	this.ConnLast++
}

func (this *AppBase) DelConn(id int) {
	if _, ok := this.Conns[id]; ok {
		if len(this.Conns[id].Name) > 0 {
			delete(this.RemoteSvr, this.Conns[id].Name)
		}
		delete(this.Conns, id)
	}
}

func (this *AppBase) GetConnById(id int) *ClientConn {
	if v, ok := this.Conns[id]; ok {
		return v
	}
	return nil
}

func (this *AppBase) GetConnByName(name string) *ClientConn {
	if v, ok := this.RemoteSvr[name]; ok {
		return v
	}
	return nil
}

func (this *AppBase) RegMsgFunc(id int, f MsgFunc) {
	this.MsgProc[id] = f

	if id > this.MsgProcCount {
		this.MsgProcCount = id
	}
}

func (this *AppBase) Listen(name, net_type, address string, onRet ConnRetFunc) {
	if len(address) == 0 || len(address) == 0 || len(net_type) == 0 {
		onRet("listen failed", name, 0, "listen failed")
		return
	}

	// 打开本地TCP侦听
	serverAddr, err := net.ResolveTCPAddr(net_type, address)

	if err != nil {
		onRet("listen failed", name, 0, "Listen Start : port failed: '"+address+"' "+err.Error())
		return
	}

	listener, err := net.ListenTCP(net_type, serverAddr)
	if err != nil {
		onRet("listen failed", name, 0, "TcpSerer ListenTCP: "+err.Error())
		return
	}

	ln := new(ListenConn)
	ln.InitListen(name, net_type, address, listener)
	this.Listener[name] = ln

	onRet("listen ok", name, 0, "")

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if !onRet("accept failed", name, 0, "TcpSerer Accept: "+err.Error()) {
				break
			}
			continue
		}
		c := new(ClientConn)
		c.InitClient(this.ConnLast, conn)
		this.AddConn(c)
		onRet("accept ok", "", c.Id, "")

		go this.ConnProc(c, onRet)
	}
}

func (this *AppBase) Connect(name, net_type, address string, onRet ConnRetFunc) {
	if len(address) == 0 || len(net_type) == 0 || len(name) == 0 {
		onRet("connect failed", name, 0, "listen failed")
		return
	}

	// 打开本地TCP侦听
	remoteAddr, err := net.ResolveTCPAddr(net_type, address)

	if err != nil {
		onRet("connect failed", name, 0, "Connect Start : port failed: '"+address+"' "+err.Error())
		return
	}

	conn, err := net.DialTCP(net_type, nil, remoteAddr)
	if err != nil {
		onRet("connect failed", name, 0, "Connect dialtcp failed: '"+address+"' "+err.Error())
	} else {
		c := new(ClientConn)
		c.InitClient(this.ConnLast, conn)
		c.Name = name
		this.RemoteSvr[name] = c
		this.AddConn(c)

		onRet("connect ok", name, c.Id, "")
		go this.ConnProc(c, onRet)
	}
}

func (this *AppBase) ConnProc(c *ClientConn, onRet ConnRetFunc) {

	for {
		c.Stream.Seek(0)
		err := c.Msg.ReadData(c.Conn)

		if err == nil {

			c.Stream.Seek(MaxHeader)
			msg_code := c.Stream.ReadU2()

			if msg_code >= 0 && msg_code <= this.MsgProcCount && this.MsgProc[msg_code] != nil {
				this.MsgProc[msg_code](c)
			}

		} else {
			onRet("read failed", c.Name, c.Id, err.Error())
			break
		}
	}

	onRet("pre close", c.Name, c.Id, "")

	err := c.Conn.Close()
	if err != nil {
		onRet("close failed", c.Name, c.Id, err.Error())
	} else {
		onRet("close ok", c.Name, c.Id, "")
	}

	GetApp().DelConn(c.Id)
}

// 线程日志 : 生成日志头
func (this *AppBase) MakeLogHeader() {
	this.log_Header[LogDebugLevel] = this.log_TimeString + " [D] "
	this.log_Header[LogInfoLevel] = this.log_TimeString + " [I] "
	this.log_Header[LogWarnLevel] = this.log_TimeString + " [W] "
	this.log_Header[LogErrorLevel] = this.log_TimeString + " [E] "
	this.log_Header[LogFatalLevel] = this.log_TimeString + " [F] "
}

// 线程日志 : 调试[D]级别日志
func (this *AppBase) LogDebug(f string, v ...interface{}) {
	this.LogBase(LogDebugLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 信息[I]级别日志
func (this *AppBase) LogInfo(f string, v ...interface{}) {
	this.LogBase(LogInfoLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 警告[W]级别日志
func (this *AppBase) LogWarn(f string, v ...interface{}) {
	this.LogBase(LogWarnLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 错误[E]级别日志
func (this *AppBase) LogError(f string, v ...interface{}) {
	this.LogBase(LogErrorLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 致命[F]级别日志
func (this *AppBase) LogFatal(f string, v ...interface{}) {
	this.LogBase(LogFatalLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 手动分级日志
func (this *AppBase) LogBase(level int, info string) {
	if level >= LogDebugLevel && level < LogMaxLevel {
		s := this.log_Header[level] + info
		s = strings.Replace(s, "\n", "\n"+this.log_Header[level], -1) + "\n"

		this.Add_log(s)

		if level >= LogLimitLevel {
			fmt.Print(s)
		}
	} else {
		fmt.Println("LogBase : level failed : ", level)
	}
}

// 增加日志到缓冲
func (this *AppBase) Add_log(d string) {
	this.log_FileBuff.WriteString(d)
}

// go协程 刷新缓冲日志到文件,每5s写盘一次
func (this *AppBase) Flush_log() {

	for {
		if this.log_FileBuff.Len() > 0 {
			this.log_FileHandle.Write(this.log_FileBuff.Bytes())
			this.log_FileBuff.Reset()
		}

		<-time.Tick(5 * time.Second)
	}
}
