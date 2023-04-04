package core

import "github.com/aFlyBird0/cubox-archiver/core/cubox"

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

func NewDeduplicateWithKeysInitiator(initiator KeysInitiator) (*Deduplicate, error) {
	keys, err := initiator.ExistingKeys()
	if err != nil {
		return nil, err
	}

	return NewDeduplicateWithKeys(keys), nil
}

func (d *Deduplicate) Remain(item *cubox.Item) bool {
	key := item.UserSearchEngineID
	if _, ok := d.keys[key]; ok {
		return false
	}
	d.keys[key] = struct{}{}

	return true
}
