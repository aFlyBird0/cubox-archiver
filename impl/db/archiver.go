package db

import (
	"gorm.io/gorm"

	"github.com/aFlyBird0/cubox-archiver/core"
	"github.com/aFlyBird0/cubox-archiver/core/cubox"
)

type Archiver struct {
	db *gorm.DB
}

type Item struct {
	//gorm.Model
	cubox.Item
}

func (item *Item) TableName() string {
	// todo: set table name by config file
	return "cubox_archiver"
}

// todo: generate different db client by config file
func NewArchiver(db *gorm.DB) *Archiver {
	return &Archiver{db: db}
}

func (a *Archiver) Operate(item *cubox.Item) error {
	return a.db.Create(&Item{Item: *item}).Error
}

func (a *Archiver) ExistingKeys() (map[string]struct{}, error) {
	items := make([]Item, 0)
	if err := a.db.Select("user_search_engine_id").Find(&items).Error; err != nil {
		return nil, err
	}
	keys := make(map[string]struct{})
	for _, item := range items {
		keys[item.UserSearchEngineID] = struct{}{}
	}
	return keys, nil
}

var _ core.Archiver = (*Archiver)(nil)
