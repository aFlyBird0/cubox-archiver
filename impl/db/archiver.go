package db

import (
	"gorm.io/gorm"

	"github.com/aFlyBird0/cubox-archiver/core"
)

const defaultTableName = "cubox_archiver"

type Archiver struct {
	db        *gorm.DB
	tableName string
}

type Item struct {
	//gorm.Model
	core.Item
}

// todo: generate different db client by config file
func NewArchiver(db *gorm.DB, tableName string) *Archiver {
	if tableName == "" {
		tableName = defaultTableName
	}

	return &Archiver{db: db, tableName: tableName}
}

func (a *Archiver) Operate(item *core.Item) error {
	return a.dbWithTableName().Create(&Item{Item: *item}).Error
}

func (a *Archiver) ExistingKeys() (map[string]struct{}, error) {
	items := make([]Item, 0)
	if err := a.dbWithTableName().Select("user_search_engine_id").Find(&items).Error; err != nil {
		return nil, err
	}
	keys := make(map[string]struct{})
	for _, item := range items {
		keys[item.UserSearchEngineID] = struct{}{}
	}
	return keys, nil
}

func (a *Archiver) dbWithTableName() *gorm.DB {
	return a.db.Table(a.tableName)
}

var _ core.Archiver = (*Archiver)(nil)
