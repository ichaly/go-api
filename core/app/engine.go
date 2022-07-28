package app

import (
	"github.com/dosco/graphjin/core"
)

func NewEngine(c *Config) (g *core.GraphJin, e error) {
	g, e = core.NewGraphJin(&c.Core, c.GetDB())
	return
}
