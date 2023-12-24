package core

import (
	"github.com/reugn/go-streams"
	"github.com/sirupsen/logrus"
)

// Operator 能对 Cubox Item 进行处理
type Operator interface {
	Operate(item *Item) error // 处理单个数据
}

type OperatorChain struct {
	operators []Operator
}

// NewOperatorChain 处理链，将多个 Operator 串联起来，任一 Operator 失败直接中断
func NewOperatorChain(operators ...Operator) *OperatorChain {
	return &OperatorChain{operators: operators}
}

func (c *OperatorChain) Operate(item *Item) error {
	for _, operator := range c.operators {
		err := operator.Operate(item)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ Operator = &OperatorChain{}

type OperatorSink struct {
	in       chan any
	operator Operator
}

// NewOperatorSink 将 Operator 包装成 go-streams 支持的 Sink
func NewOperatorSink(operator Operator, done chan<- struct{}) streams.Sink {
	o := OperatorSink{in: make(chan any), operator: operator}
	go o.init(done)

	return &o
}

// 处理主逻辑
func (o *OperatorSink) init(done chan<- struct{}) {
	for itemAny := range o.in {
		item := itemAny.(*Item)
		err := o.operator.Operate(item)
		if err != nil {
			logrus.Errorf("operate item: %v, err: %v", item, err)
		}
	}
	done <- struct{}{}
}

func (o *OperatorSink) In() chan<- any {
	return o.in
}

var _ streams.Sink = &OperatorSink{}
