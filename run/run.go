package run

import (
	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"

	"github.com/sirupsen/logrus"

	"github.com/aFlyBird0/cubox-archiver/config"
	"github.com/aFlyBird0/cubox-archiver/core"
	"github.com/aFlyBird0/cubox-archiver/core/cubox"
)

func Run(cfg config.Config) {
	// 根据配置，构建处理链
	getter, archivers := Build(cfg)

	logrus.Info("数据源获取器，归档器构建完成，开始处理数据")

	// 执行
	Process(getter, archivers)

	logrus.Info("所有处理完成，程序结束")
}

// Process 处理数据获取和归档
// todo 支持多种 cubox 数据源/多个数据源合并
func Process(source cubox.Source, archivers []core.Archiver) {
	// 从cubox获取数据源
	cuboxChan := make(chan *cubox.Item, 100)
	//stop := make(chan struct{})
	go source.List(cuboxChan)

	// 将数据源转换成go-streams的格式
	cuboxChanAny := make(chan any, 100)
	go transformChan(cuboxChan, cuboxChanAny)
	chanSource := ext.NewChanSource(cuboxChanAny)

	// 针对每个 archiver，fan out 出一个流，并行处理
	archiveFlows := flow.FanOut(chanSource, len(archivers))

	// done 用来控制所有的 archiver 都处理完毕
	done := make(chan struct{})
	for i := range archiveFlows {
		i := i
		deduplicate, err := core.NewDeduplicateWithKeysInitiator(archivers[i])
		if err != nil {
			logrus.Fatalln(err)
		}
		archiveFlows[i].
			Via(flow.NewFilter(deduplicate.Remain, 10)).
			To(core.NewOperatorSink(archivers[i], done))
	}

	for i := 0; i < len(archivers); i++ {
		<-done
	}
}

// 转换一下专利chan的类型，go-streams的输入是chan any
func transformChan(in <-chan *cubox.Item, out chan<- any) {
	for e := range in {
		out <- e
	}
	close(out)
}
