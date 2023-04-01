package util

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
)

func CreateNotionDatabase(token, databaseName, parentPageID string) (databaseID string, err error) {
	client := notionapi.NewClient(notionapi.Token(token))
	var parent notionapi.Parent
	// page id 不为空，创建在指定页面下
	if parentPageID != "" {
		parent = notionapi.Parent{
			PageID: notionapi.PageID(parentPageID),
			Type:   notionapi.ParentTypePageID,
		}
	} else {
		// 否则创建在工作区下
		parent = notionapi.Parent{
			Workspace: true,
			Type:      notionapi.ParentTypeWorkspace,
			PageID:    notionapi.PageID(""),
		}
	}
	req := notionapi.DatabaseCreateRequest{
		Parent: parent,
		Title: []notionapi.RichText{
			{Text: &notionapi.Text{Content: databaseName}},
		},
		Properties: notionapi.PropertyConfigs{
			"Name":        notionapi.TitlePropertyConfig{Type: notionapi.PropertyConfigTypeTitle},
			"CuboxID":     notionapi.RichTextPropertyConfig{Type: notionapi.PropertyConfigTypeRichText},
			"Description": notionapi.RichTextPropertyConfig{Type: notionapi.PropertyConfigTypeRichText},
			"Tags": notionapi.MultiSelectPropertyConfig{
				Type: notionapi.PropertyConfigTypeMultiSelect,
				MultiSelect: notionapi.Select{
					Options: []notionapi.Option{},
				},
			},
			"CreateTime": notionapi.DatePropertyConfig{Type: notionapi.PropertyConfigTypeDate},
			"UpdateTime": notionapi.DatePropertyConfig{Type: notionapi.PropertyConfigTypeDate},
			"Type": notionapi.SelectPropertyConfig{
				Type: notionapi.PropertyConfigTypeSelect,
				Select: notionapi.Select{
					Options: []notionapi.Option{},
				},
			},
			"URL":        notionapi.URLPropertyConfig{Type: notionapi.PropertyConfigTypeURL},
			"Cover":      notionapi.FilesPropertyConfig{Type: notionapi.PropertyConfigTypeFiles},
			"LittleIcon": notionapi.FilesPropertyConfig{Type: notionapi.PropertyConfigTypeFiles},
			"GroupName": notionapi.SelectPropertyConfig{
				Type: notionapi.PropertyConfigTypeSelect,
				Select: notionapi.Select{
					Options: []notionapi.Option{},
				},
			},
		},
	}
	db, err := client.Database.Create(context.TODO(), &req)
	if err != nil {
		return "", fmt.Errorf("create database failed: %w", err)
	}
	return string(db.ID), nil
}
