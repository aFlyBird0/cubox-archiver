package stdout

import (
	"github.com/sirupsen/logrus"

	"github.com/aFlyBird0/cubox-archiver/core"
)

type StdoutOperator struct {
}

func NewStdoutOperator() *StdoutOperator {
	return &StdoutOperator{}
}

func (s StdoutOperator) Operate(item *core.Item) error {
	logrus.Info(item.Title)
	return nil
}

var _ core.Operator = &StdoutOperator{}
