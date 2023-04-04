package core

import (
	"github.com/aFlyBird0/cubox-archiver/core/cubox"
)

type Filter interface {
	Remain(item *cubox.Item) bool
}

type KeysInitiator interface {
	ExistingKeys() (map[string]struct{}, error) // 获取已有的数据的Keys
}

// NewFilterAllRemain 保留所有
func NewFilterAllRemain() *AllRemain {
	return &AllRemain{}
}

type AllRemain struct {
}

func (f *AllRemain) Remain(item *cubox.Item) bool {
	return true
}

type FirstN struct {
	n int
}

// NewFilterFirstN 保留前 n 个
func NewFilterFirstN(n int) *FirstN {
	return &FirstN{n: n}
}

func (f *FirstN) Remain(item *cubox.Item) bool {
	if f.n <= 0 {
		return false
	}

	f.n--
	return true
}
