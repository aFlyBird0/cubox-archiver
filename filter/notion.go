package filter

import (
	"context"

	"github.com/jomei/notionapi"
)

func NewNotionDuplicateFilter(token, dbID string) (Filter, error) {
	client := notionapi.NewClient(notionapi.Token(token))
	keys, err := initKeysFromNotion(client, dbID)
	if err != nil {
		return nil, err
	}
	return NewDeduplicateWithKeys(keys), nil
}

// todo 考虑把这块逻辑变成一个接口，这样就不需要为每一个数据持久化源写单独的过滤器了
// 大概就是，原来写 Notion 去重器需要写个 NewNotionDuplicateFilter
// 现在只要这样 NewDeduplicateWithKeys(NewNotionKeysGetter(client, dbID))
func initKeysFromNotion(client *notionapi.Client, dbID string) (map[string]struct{}, error) {
	keys := make(map[string]struct{})

	cursor := notionapi.Cursor("")
	for {
		query := notionapi.DatabaseQueryRequest{
			PageSize:    5,
			StartCursor: cursor,
		}

		res, err := client.Database.Query(context.Background(), notionapi.DatabaseID(dbID), &query)
		if err != nil {
			return nil, err
		}
		pages := res.Results
		for _, page := range pages {
			IDProp := page.Properties["CuboxID"].(*notionapi.RichTextProperty)
			keys[IDProp.RichText[0].PlainText] = struct{}{}
		}
		if res.HasMore {
			cursor = res.NextCursor
		} else {
			break
		}
	}

	return keys, nil
}
