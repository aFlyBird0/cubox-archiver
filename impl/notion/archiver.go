package notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

type (
	Archiver struct {
		databaseID string
		client     *notionapi.Client
	}
)

func NewArchiver(token, databaseID string) (*Archiver, error) {
	if token == "" {
		return nil, fmt.Errorf("notion token 不能为空")
	}
	if databaseID == "" {
		return nil, fmt.Errorf("notion databaseID 不能为空")
	}
	client := notionapi.NewClient(notionapi.Token(token), notionapi.WithRetry(3))

	if _, err := client.User.Me(context.TODO()); err != nil {
		return nil, fmt.Errorf("notion token 无效或网络异常(尝试获取用户信息失败): %w", err)
	}
	logrus.Info("notion 归档器初始化成功")
	return &Archiver{
		databaseID: databaseID,
		client:     client,
	}, nil
}

var (
	NewOperator      = NewArchiver
	NewKeysInitiator = NewArchiver
)
