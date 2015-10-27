package logic

import (
	"github.com/toophy/gate/help"
)

func on_c2g_login(c *help.ClientConn) {
	if c.Id > 0 {
		name := c.Stream.ReadStr()
		println(name)
	}
}
