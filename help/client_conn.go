package help

import (
	"net"
)

// 客户端连接: GMS,TDB
type ClientConn struct {
	Type    string
	Address string
	Id      int
	Sid     int
	Msg     Ty_net_msg
	Stream  Ty_msg_stream
	Conn    *net.TCPConn
}

func (c *ClientConn) InitClient(id int, con *net.TCPConn) {
	c.Type = "Null"
	c.Address = con.RemoteAddr().String()
	c.Id = id
	c.Conn = con
	c.Msg.InitNetMsg()
	c.Stream.InitMsgStream(&c.Msg)
}

func (c *ClientConn) IsNull() bool {
	return c.Type == "Null"
}
