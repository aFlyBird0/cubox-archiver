package core

import (
	"os"

	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"

	"github.com/aFlyBird0/cubox-archiver/config"
	"github.com/aFlyBird0/cubox-archiver/cubox"
	"github.com/aFlyBird0/cubox-archiver/filter"
	"github.com/aFlyBird0/cubox-archiver/operator"
	cuboxOp "github.com/aFlyBird0/cubox-archiver/operator/cubox"
	"github.com/aFlyBird0/cubox-archiver/operator/notion"
	"github.com/aFlyBird0/cubox-archiver/source"
	"github.com/aFlyBird0/cubox-archiver/util"

	"github.com/sirupsen/logrus"
)

const defaultNotionDatabaseName = "Cubox归档"

func Run(cfg config.Config) {
	// 根据配置，构建处理链
	getter, filters, operators := Build(cfg)

	// 执行
	Process(getter, filters, operators)
}

func Build(conf config.Config) (
	cubox.Source, []filter.Filter, []operator.Operator) {

	if conf.Notion.DatabaseID == "" && conf.Notion.PageID == "" {
		logrus.Fatalln("请填写数据库ID或者页面ID")
	}

	// 如果没有数据库ID，说明需要自动新建数据库
	if conf.Notion.DatabaseID == "" {
		var err error
		if conf.Notion.DatabaseName == "" {
			conf.Notion.DatabaseName = defaultNotionDatabaseName
		}
		conf.Notion.DatabaseID, err = util.CreateNotionDatabase(conf.Notion.Token, conf.Notion.DatabaseName, conf.Notion.PageID)
		if err != nil {
			logrus.Fatalln(err)
		}
		logrus.Info("-----------------------")
		logrus.Info("生成的数据库ID为：", conf.Notion.DatabaseID)
		// 其实也可以在创建好数据库后马上把数据同步过去
		// 但可能有的用户在后面第二次执行的时候会忘记填写数据库ID
		// 导致本次都是新建的数据库，而不是同步到之前的数据库
		// 干脆强制把两步分开
		logrus.Info("请将该ID填写到配置文件的DatabaseID中，重新运行程序")
		os.Exit(0)
	}

	// 从cubox获取数据源
	getter := source.NewArchivedCuboxSource(conf.Cubox.Auth, conf.Cubox.Cookie)

	// 拼接过滤器
	var filters []filter.Filter
	//filters = append(filters, filter.NewFilterAllRemain())
	//filters = append(filters, filter.NewFilterFirstN(30))
	notionFilter, err := filter.NewNotionDuplicateFilter(conf.Notion.Token, conf.Notion.DatabaseID)
	if err != nil {
		logrus.Fatalln(err)
	}
	filters = append(filters, notionFilter)

	// 拼接操作器
	var operators []operator.Operator
	//operators = append(operators, stdout.NewStdoutOperator())
	//operators = append(operators, csv.NewCsvOperator("test.csv"))
	notionOperator, err := notion.NewNotionOperator(conf.Notion.Token, conf.Notion.DatabaseID)
	if err != nil {
		logrus.Fatalln(err)
	}
	// 是否在同步完之后删除cubox中的数据
	if conf.DeleteCuboxAfterSave {
		deleteOperator := cuboxOp.NewDeleteCuboxOperator(conf.Cubox.Auth, conf.Cubox.Cookie)
		operators = append(operators, operator.NewOperatorChain(notionOperator, deleteOperator))
	} else {
		operators = append(operators, notionOperator)
	}

	return getter, filters, operators
}

func Process(getter cubox.Source, filters []filter.Filter, operators []operator.Operator) {
	// 从cubox获取数据源
	cuboxChan := make(chan *cubox.Item, 100)
	//stop := make(chan struct{})
	go getter.List(cuboxChan)

	// 将数据源转换成go-streams的格式
	cuboxChanAny := make(chan any, 100)
	go transformChan(cuboxChan, cuboxChanAny)
	chanSource := ext.NewChanSource(cuboxChanAny)

	// 用 pass through 把 source 转换成 go-streams 的 Flow
	filterFlow := chanSource.Via(flow.NewPassThrough())

	// 否则用 filters，并转换成 go-streams 的 Flow
	for _, f := range filters {
		filterFlow = filterFlow.Via(flow.NewFilter(f.Remain, 10))
	}

	// 针对每个 operator，fan out 出一个流，并行处理
	flows := flow.FanOut(filterFlow, len(operators))

	// done 用来控制所有的 operator 都处理完毕
	done := make(chan struct{})
	for i := range flows {
		i := i
		// 为每个分支流应用一个 operator
		go flows[i].To(operator.NewOperatorSink(operators[i], done))
	}

	//<-stop

	for i := 0; i < len(operators); i++ {
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
