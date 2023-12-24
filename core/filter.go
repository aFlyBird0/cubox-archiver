package core

type Filter interface {
	Remain(item *Item) bool
}

// KeysInitiator 能够获取已有的数据的Keys
type KeysInitiator interface {
	ExistingKeys() (map[string]struct{}, error) // 获取已有的数据的Keys
}

// NewFilterAllRemain 保留所有
func NewFilterAllRemain() *AllRemain {
	return &AllRemain{}
}

type AllRemain struct {
}

func (f *AllRemain) Remain(item *Item) bool {
	return true
}

type FirstN struct {
	n int
}

// NewFilterFirstN 保留前 n 个
func NewFilterFirstN(n int) *FirstN {
	return &FirstN{n: n}
}

func (f *FirstN) Remain(item *Item) bool {
	if f.n <= 0 {
		return false
	}

	f.n--
	return true
}
