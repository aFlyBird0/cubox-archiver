package csv

import (
	"encoding/csv"
	"os"

	"github.com/aFlyBird0/cubox-archiver/cubox"
)

type CsvOperator struct {
	filename string
}

func NewCsvOperator(filename string) *CsvOperator {
	return &CsvOperator{filename: filename}
}

func (o *CsvOperator) Operate(item *cubox.Item) error {
	file, err := os.OpenFile(o.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{
		item.UserSearchEngineID,
		item.Title,
		//item.Description,
		item.TargetURL,
		//item.ArchiveName,
		//item.ArticleName,
		item.Cover,
		item.LittleIcon,
		item.GroupId,
		item.GroupName,
		item.CreateTime.Format("2006-01-02 15:04:05"),
		item.UpdateTime.Format("2006-01-02 15:04:05"),
		//item.Status,
	})

	return err
}
