package core

import (
	"sync"
)

type Deduplicate struct {
	mu   sync.RWMutex
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

func (d *Deduplicate) Remain(item *Item) bool {
	key := item.UserSearchEngineID
	d.mu.RLock()
	if _, ok := d.keys[key]; ok {
		d.mu.RUnlock()
		return false
	}
	d.mu.RUnlock()

	d.mu.Lock()
	d.keys[key] = struct{}{}
	d.mu.Unlock()

	return true
}
