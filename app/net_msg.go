package app

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	MaxDataLen     = 5080
	MaxSendDataLen = 4000
	MaxHeader      = 2
)

// 网络消息体
type Ty_net_msg struct {
	Data []byte
	Len  int
}

func (t *Ty_net_msg) InitNetMsg() {
	t.Len = 0
	t.Data = make([]byte, MaxDataLen)
}

func (t *Ty_net_msg) PrintData() {
	fmt.Println(t.Data[:t.Len+2])
}

func (t *Ty_net_msg) ReadData(conn *net.TCPConn) error {

	t.Len = 0
	length, err := io.ReadFull(conn, t.Data[:2])
	if length != MaxHeader {
		GetApp().LogWarn("Packet header : %d != %d", length, MaxHeader)
		return err
	}
	if err != nil {
		return err
	}

	body_len := int(t.Data[1]) + (int(t.Data[0]) << 8)

	if body_len > (MaxDataLen - 2) {
		err = errors.New("Body too much")
		return err
	}

	t.Len = body_len + 2
	return t.ReadBody(conn)
}

func (t *Ty_net_msg) ReadBody(conn *net.TCPConn) error {

	length, err := io.ReadFull(conn, t.Data[2:t.Len])
	if length != (t.Len - 2) {
		GetApp().LogWarn("Packet length : %d != %d ", length, t.Len-2)
		return err
	}
	if err != nil {
		return err
	}
	// 注意 : 可以解密

	return nil
}

func (t *Ty_net_msg) Send(conn *net.TCPConn) error {
	if t.Len > MaxHeader && t.Len < MaxSendDataLen {

		t.Data[0] = byte((t.Len & 0xFF00) >> 8)
		t.Data[1] = byte(t.Len & 0xFF)

		_, err := conn.Write(t.Data[:MaxHeader+t.Len])
		if err != nil {
			GetApp().LogWarn(err.Error())
			return err
		}
	}

	return nil
}
