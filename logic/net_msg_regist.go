package logic

import (
	"github.com/toophy/gate/help"
)

func RegMsgProc() {
	help.GetApp().RegMsgFunc(1, on_c2g_login)
}
