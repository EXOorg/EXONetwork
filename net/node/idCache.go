package node

import (
	"sync"

	"github.com/nknorg/nkn/common"
)

type idCache struct {
	sync.RWMutex
	list map[common.Uint256]bool
}

func (c *idCache) init() {
}

func (c *idCache) add() {
}

func (c *idCache) del() {
}

func (c *idCache) ExistedID(id common.Uint256) bool {
	// TODO
	return false
}
