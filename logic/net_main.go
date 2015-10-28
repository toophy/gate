package logic

import (
	"fmt"
	"github.com/toophy/gate/help"
)

func Main_go() {
	RegMsgProc()

	go help.GetApp().Listen("main_listen", "tcp", ":8001", OnListenRet)
}

func OnListenRet(typ string, name string, id int, info string) bool {
	name_fix := name
	if len(name_fix) == 0 {
		name_fix = fmt.Sprintf("Conn[%d]", id)
	}

	switch typ {
	case "listen failed":
		help.GetApp().LogFatal("%s : Listen failed[%s]", name_fix, info)

	case "listen ok":
		help.GetApp().LogInfo("%s : Listen ok.", name_fix)

	case "accept failed":
		help.GetApp().LogFatal(info)
		return false

	case "accept ok":
		help.GetApp().LogDebug("%s : Accept ok", name_fix)

	case "connect failed":
		help.GetApp().LogError("%s : Connect failed[%s]", name_fix, info)

	case "connect ok":
		help.GetApp().LogDebug("%s : Connect ok", name_fix)

	case "read failed":
		help.GetApp().LogError("%s : Connect read[%s]", name_fix, info)

	case "pre close":
		help.GetApp().LogDebug("%s : Connect pre close", name_fix)

	case "close failed":
		help.GetApp().LogError("%s : Connect close failed[%s]", name_fix, info)

	case "close ok":
		help.GetApp().LogDebug("%s : Connect close ok.", name_fix)
	}

	return true
}
