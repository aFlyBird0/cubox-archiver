package notion

import (
	"context"
	"fmt"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"

	"github.com/aFlyBird0/cubox-archiver/core/cubox"
	"github.com/aFlyBird0/cubox-archiver/util"
)

func (o *Archiver) Operate(item *cubox.Item) error {
	_, err := o.createNewPage(item)
	if err != nil {
		return fmt.Errorf("create new page: %w", err)
	}
	logrus.Info("成功将【", item.Title, "】同步到Notion")

	return nil
}

func (o *Archiver) createNewPage(item *cubox.Item) (*notionapi.Page, error) {
	properties := map[string]notionapi.Property{
		// todo 把这里的属性名换成常量，在创建数据库的时候也要用到
		"Name": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{Content: item.Title, Link: &notionapi.Link{Url: item.TargetURL}},
				},
			},
		},
		"CuboxID": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{Content: item.UserSearchEngineID},
				},
			},
		},
		"Description": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{Content: item.Description},
				},
			},
		},
		"Tags":       convertTags(item.Tags),
		"CreateTime": convertTime(item.CreateTime),
		"UpdateTime": convertTime(item.UpdateTime),
		"Type":       convertType(item.Type),
	}
	if item.TargetURL != "" {
		properties["URL"] = notionapi.URLProperty{
			URL: item.TargetURL,
		}
	}
	if item.Cover != "" {
		properties["Cover"] = notionapi.FilesProperty{
			Files: []notionapi.File{
				{
					Name:     util.Trunc(item.Cover, 100),
					External: &notionapi.FileObject{URL: item.Cover},
					Type:     notionapi.FileTypeExternal,
				},
			},
		}
	}
	if item.LittleIcon != "" {
		properties["LittleIcon"] = notionapi.FilesProperty{
			Files: []notionapi.File{
				{
					Name:     util.Trunc(item.LittleIcon, 100),
					External: &notionapi.FileObject{URL: item.LittleIcon},
					Type:     notionapi.FileTypeExternal,
				},
			},
		}
	}
	if item.GroupName != "" {
		properties["GroupName"] = notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: item.GroupName,
			},
		}
	}

	children := make([]notionapi.Block, 0)
	if item.Content != "" {
		contents := strings.Split(item.Content, "\n")
		richTextes := make([]notionapi.RichText, 0, len(contents))
		for _, content := range contents {
			if content == "" {
				continue
			}
			richTextes = append(richTextes, notionapi.RichText{
				Text: &notionapi.Text{Content: content},
			})
		}
		children = append(children, notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: richTextes,
			}},
		)
	}

	pageReq := notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(o.databaseID),
		},
		Properties: properties,
	}
	if len(children) > 0 {
		pageReq.Children = children
	}

	page, err := o.client.Page.Create(context.TODO(), &pageReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create notion page: %w", err)
	}

	return page, nil
}
