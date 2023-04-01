package filter

import "github.com/aFlyBird0/cubox-archiver/cubox"

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

type Deduplicate struct {
	keys map[string]struct{}
}

// NewDeduplicate 保留不重复的
func NewDeduplicate() *Deduplicate {
	return &Deduplicate{
		keys: make(map[string]struct{}),
	}
}

// NewDeduplicateWithKeys 保留不重复的，并传入已有的元素列表
func NewDeduplicateWithKeys(keys map[string]struct{}) *Deduplicate {
	if keys == nil {
		keys = make(map[string]struct{})
	}

	return &Deduplicate{
		keys: keys,
	}
}

func (d *Deduplicate) Remain(item *cubox.Item) bool {
	key := item.UserSearchEngineID
	if _, ok := d.keys[key]; ok {
		return false
	}
	d.keys[key] = struct{}{}

	return true
}
