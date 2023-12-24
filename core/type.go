package core

import (
	"fmt"
	"time"
)

const (
	Link CuboxContentType = iota
	Text
	ShortHand // 速记
	Picture

	Unknown = 99
)

var ContentTypeSet = []CuboxContentType{Link, Text, ShortHand, Picture}

type (
	Item struct {
		UserSearchEngineID string `json:"userSearchEngineID"` // Cubox 内部 ID
		Title              string `json:"title"`              // 标题
		Description        string `json:"description"`        // 描述
		TargetURL          string `json:"targetURL"`          // 原始链接
		//ResourceURL        interface{} `json:"resourceURL"`
		//HomeURL            string      `json:"homeURL"`
		//ArchiveName string `json:"archiveName"` // 不明，感觉是 Cubox 内部用的
		Content string `json:"content"`
		//ArticleName string `json:"articleName"` // 不明
		//ArticleWordCount   int         `json:"articleWordCount"`
		//Byline             string      `json:"byline"`
		Cover string `json:"cover"` // 封面图片链接
		//ArticleURL         interface{} `json:"articleURL"` // 不明
		LittleIcon string `json:"littleIcon"` // 小图标，一般是网站icon一类
		//Archiving          bool        `json:"archiving"`
		//StarTarget         bool        `json:"starTarget"`
		//HasMark            bool        `json:"hasMark"`
		//IsRead             interface{} `json:"isRead"`
		//MarkCount          int         `json:"markCount"`
		Tags []Tag `json:"tags"`
		//AllTags            []interface{}    `json:"allTags"`
		//Marks              []interface{}    `json:"marks"`
		GroupId    string    `json:"groupId"`
		GroupName  string    `json:"groupName"`
		CreateTime time.Time `json:"createTime"`
		UpdateTime time.Time `json:"updateTime"`
		Status     string    `json:"status"`
		//Finished           bool             `json:"finished"`
		//InBlackOrWhiteList bool             `json:"inBlackOrWhiteList"`
		Type CuboxContentType `json:"type"` // 类型，已知 0 是链接文章，1 是文本
	}

	CuboxContentType int
)

func (c CuboxContentType) String() string {
	switch c {
	case Link:
		return "链接"
	case Text:
		return "文本"
	case ShortHand:
		return "速记"
	case Picture:
		return "图片"
	default:
		return fmt.Sprintf("Unknown-%d", c)
	}
}

func TypeFromString(s string) CuboxContentType {
	switch s {
	case "链接":
		return Link
	case "文本":
		return Text
	case "速记":
		return ShortHand
	case "图片":
		return Picture
	default:
		return Unknown
	}
}

type Tag struct {
	TagID      string      `json:"tagID"`
	Name       string      `json:"name"`
	Rank       interface{} `json:"rank"` // 类型不明
	UpdateTime time.Time   `json:"updateTime"`
	ParentId   string      `json:"parentId"`
}
