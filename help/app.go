package help

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// 消息函数类型
type MsgFunc func(*ClientConn)

type AppBase struct {
	baseGoNumStart int
	baseGoNumEnd   int
	Address        string              // 地址
	Listener       *net.TCPListener    // 本地侦听端口
	Conns          map[int]*ClientConn // 连接池
	ConnLast       int                 // 最后连接Id
	MsgProc        []MsgFunc           // 消息处理函数注册表
	MsgProcCount   int                 // 消息函数数量
}

// 程序控制核心
var app *AppBase

func GetApp() *AppBase {
	if app == nil {
		app = &AppBase{}
		app.Conns = make(map[int]*ClientConn, 1000)
		app.MsgProc = make([]MsgFunc, 8000)
	}
	return app
}

// 程序开启
func (this *AppBase) Start(logDir, profFile string) {

	runtime.GOMAXPROCS(1)

	// 检查log目录
	if !IsExist(logDir) {
		os.MkdirAll(logDir, os.ModeDir)
	}

	// 创建pprof文件
	f, err := os.Create(logDir + "/" + profFile)
	if err != nil {
		fmt.Println(err.Error())
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
			fmt.Println("bye bye.")
			break
		} else {
			// fmt.Println("mimi", runtime.NumGoroutine())
		}
	}
}

func (this *AppBase) AddConn(c *ClientConn) {
	this.Conns[this.ConnLast] = c
	this.ConnLast++
}

func (this *AppBase) DelConn(id int) {
	delete(this.Conns, id)
}

func (this *AppBase) RegMsgFunc(id int, f MsgFunc) {
	this.MsgProc[id] = f
}

func (this *AppBase) Listen(net_type, address string) {
	if len(this.Address) > 0 || len(address) == 0 || len(net_type) == 0 {
		fmt.Println("listen failed")
		return
	}

	this.Address = address

	// 打开本地TCP侦听
	serverAddr, err := net.ResolveTCPAddr("tcp", this.Address)

	if err != nil {
		// GetLog().Error("ScsKernel Start : port failed: '%s' %s", this.Address, err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		// GetLog().Error("TcpSerer ListenTCP: %s", err.Error())
		return
	}

	for {
		this.Listener = listener
		conn, err := listener.AcceptTCP()
		if err != nil {
			// handle error
			continue
		}
		c := new(ClientConn)
		this.ConnLast++
		c.InitClient(this.ConnLast, conn)

		go this.ConnProc(c)
	}
}

func (this *AppBase) ConnProc(c *ClientConn) {

	for {
		c.Stream.Seek(0)
		err := c.Msg.ReadData(c.Conn)

		if err == nil {

			c.Stream.Seek(MaxHeader)
			msg_code := c.Stream.ReadU2()

			if msg_code >= 0 && msg_code < this.MsgProcCount {
				this.MsgProc[msg_code](c)
			}

		} else {
			fmt.Println(err.Error())
			break
		}
	}

	err := c.Conn.Close()
	fmt.Println("Closed connection:", c.Address)
	if err != nil {
		fmt.Println("ERROR: " + "Close:" + " " + err.Error())
	}

	GetApp().DelConn(c.Id)
}
