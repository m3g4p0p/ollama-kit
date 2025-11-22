package console

import "m3g4p0p/agents/util"

func WithStyle(style string) util.Option[Console] {
	return func(c *Console) {
		c.style = style
	}
}
