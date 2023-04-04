package notion

import (
	"context"

	"github.com/jomei/notionapi"
)

func (o *Archiver) ExistingKeys() (map[string]struct{}, error) {
	keys := make(map[string]struct{})

	cursor := notionapi.Cursor("")
	for {
		query := notionapi.DatabaseQueryRequest{
			PageSize:    5,
			StartCursor: cursor,
		}

		res, err := o.client.Database.Query(context.TODO(), notionapi.DatabaseID(o.databaseID), &query)
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
