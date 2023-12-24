package notion

import (
	"time"

	"github.com/jomei/notionapi"

	"github.com/aFlyBird0/cubox-archiver/core"
)

func convertTags(tags []core.Tag) notionapi.MultiSelectProperty {
	options := make([]notionapi.Option, len(tags))
	for i, tag := range tags {
		options[i] = notionapi.Option{
			Name: tag.Name,
		}
	}
	return notionapi.MultiSelectProperty{
		MultiSelect: options,
	}
}

func convertTime(t time.Time) notionapi.DateProperty {
	notionTime := notionapi.Date(t)
	return notionapi.DateProperty{
		Date: &notionapi.DateObject{
			Start: &notionTime,
		},
	}
}

func convertType(t core.CuboxContentType) notionapi.SelectProperty {
	return notionapi.SelectProperty{
		Select: notionapi.Option{
			Name: t.String(),
		},
	}
}
