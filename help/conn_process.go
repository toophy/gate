package help

import ()

// func (t *ScsKernel) AddConn(client *ClientConn) {
// 	LockScs()
// 	if t.OpenTcpSvr {
// 		t.Conns[t.ConnLast] = client
// 		t.ConnLast++
// 	}
// 	UnlockScs()
// }

// func (t *ScsKernel) DelConn(id int) {
// 	LockScs()
// 	if t.OpenTcpSvr {
// 		delete(t.Conns, id)
// 	}
// 	UnlockScs()
// }

// func (t *ScsKernel) ConnProc(client *ClientConn) {

// 	GetLog().Trace("Connection from: ", client.Address)

// 	for {
// 		client.Stream.Seek(0)
// 		err := client.Msg.ReadData(client.Conn)

// 		if err == nil {

// 			client.Stream.Seek(MaxHeader)
// 			msg_code := client.Stream.ReadU2()

// 			fc, ok := t.MsgProc[msg_code]
// 			if ok {
// 				fc(t, client)
// 			}
// 		} else {
// 			GetLog().Trace(err.Error())
// 			break
// 		}
// 	}

// 	LockScs()
// 	if t.OpenTcpSvr {
// 		if client.IsMapSvr() {
// 			t.MapStatus = 0
// 		} else if client.IsWorldSvr() {
// 			t.WorldStatus = 0
// 		} else if client.IsWatchSvr() {
// 			t.WatchStatus = 0
// 		} else if client.IsDbSvr() {
// 			t.TdbStatus = 0
// 		} else if client.IsLoginSvr() {
// 			t.LoginStatus = 0
// 		}
// 	}
// 	UnlockScs()

// 	err := client.Conn.Close()
// 	GetLog().Trace("Closed connection:", client.Address)
// 	if err != nil {
// 		GetLog().Trace("ERROR: " + "Close:" + " " + err.Error())
// 	}

// 	master := client.Master
// 	master.DelConn(client.Id)
// }
