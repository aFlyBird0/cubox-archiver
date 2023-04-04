package stdout

import (
	"github.com/sirupsen/logrus"

	"github.com/aFlyBird0/cubox-archiver/core"
	"github.com/aFlyBird0/cubox-archiver/core/cubox"
)

type StdoutOperator struct {
}

func NewStdoutOperator() *StdoutOperator {
	return &StdoutOperator{}
}

func (s StdoutOperator) Operate(item *cubox.Item) error {
	logrus.Info(item.Title)
	return nil
}

var _ core.Operator = &StdoutOperator{}
