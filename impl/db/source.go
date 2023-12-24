package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/aFlyBird0/cubox-archiver/core"
)

type Source struct {
	db        *gorm.DB
	tableName string
}

type SourceItem struct {
	//gorm.Model
	core.Item
}

func NewSource(db *gorm.DB, tableName string) *Source {
	if tableName == "" {
		tableName = defaultTableName
	}

	return &Source{db: db, tableName: tableName}
}

func (s *Source) List(items chan *core.Item) {
	defer close(items)
	sourceItems := make([]SourceItem, 0)
	// read from db per 1000 items
	for offset := 0; ; offset += 1000 {
		if err := s.dbWithTableName().Limit(1000).Offset(offset).Find(&sourceItems).Error; err != nil {
			logrus.Errorf("read from db: %v", err)
			break
		}
		if len(sourceItems) == 0 {
			break
		}
		for _, sourceItem := range sourceItems {
			items <- &sourceItem.Item
		}
	}
}

func (s *Source) dbWithTableName() *gorm.DB {
	return s.db.Table(s.tableName)
}

var _ core.Source = (*Source)(nil)
