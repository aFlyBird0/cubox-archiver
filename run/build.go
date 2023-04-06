package run

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/aFlyBird0/cubox-archiver/config"
	"github.com/aFlyBird0/cubox-archiver/core"
	cuboxType "github.com/aFlyBird0/cubox-archiver/core/cubox"
	cuboxImpl "github.com/aFlyBird0/cubox-archiver/impl/cubox"
	"github.com/aFlyBird0/cubox-archiver/impl/notion"
)

const defaultNotionDatabaseName = "Cubox归档"

func Build(conf config.Config) (
	cuboxType.Source, []core.Archiver) {

	if err := handleNotionDatabaseCreate(conf); err != nil {
		logrus.Fatalln(err)
	}

	if err := checkOnlyOneDelete(conf); err != nil {
		logrus.Fatalln(err)
	}

	// 从cubox获取数据源
	getter := cuboxImpl.NewArchivedCuboxSource(conf.Cubox.Auth, conf.Cubox.Cookie)

	var archivers []core.Archiver
	for _, archiverConf := range conf.Archivers {
		if archiverConf.Enable == false {
			logrus.Infof("archiver %s is disabled", archiverConf.Type)
			continue
		}
		archiver, err := NewArchiverFromConfig(archiverConf)
		if err != nil {
			logrus.Fatalln(err)
		}
		// 如果配置了删除 cubox，则在 archiver 之后添加一个删除 cubox 的操作
		if archiverConf.DeleteCuboxAfterSave {
			logrus.Infof("归档器 %s 配置了保存成功后删除 cubox", archiverConf.Type)
			archiver = archiverWithDelete(archiver, conf.Cubox.Auth, conf.Cubox.Cookie)
		}
		archivers = append(archivers, archiver)
	}

	return getter, archivers
}

func handleNotionDatabaseCreate(conf config.Config) error {
	notionIndex := -1

	var (
		token        string
		databaseID   string
		databaseName string
		pageID       string
		err          error
	)
	for i, archiver := range conf.Archivers {
		// 如果没有启用，跳过
		if !archiver.Enable {
			continue
		}
		// 如果是 notion
		if archiver.Type == "notion" {
			notionIndex = i
			databaseID = archiver.Options["databaseID"]
			databaseName = archiver.Options["databaseName"]
			pageID = archiver.Options["pageID"]
			token = archiver.Options["token"]
			if token == "" {
				logrus.Fatalln("请填写 Notion Token")
			}
			if databaseID == "" && pageID == "" {
				logrus.Fatalln("请填写数据库ID或者页面ID")
			}
			// 如果有数据库ID，说明不需要自动新建数据库
			if databaseID != "" {
				return nil
			}
			break
		}
	}

	// 如果没有 notion，直接返回
	if notionIndex == -1 {
		return nil
	}

	if databaseName == "" {
		databaseName = defaultNotionDatabaseName
	}
	// 如果没有数据库ID，说明需要自动新建数据库
	databaseID, err = notion.CreateNotionDatabase(token, databaseName, pageID)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.Info("-----------------------")
	logrus.Info("生成的数据库ID为：", databaseID)
	// 其实也可以在创建好数据库后马上把数据同步过去
	// 但可能有的用户在后面第二次执行的时候会忘记填写数据库ID
	// 导致本次都是新建的数据库，而不是同步到之前的数据库
	// 干脆强制把两步分开
	logrus.Info("请将该ID填写到配置文件的DatabaseID中，重新运行程序")
	os.Exit(0)
	return nil
}

func checkOnlyOneDelete(conf config.Config) error {
	var deleteCount int
	for _, archiver := range conf.Archivers {
		if archiver.Enable && archiver.DeleteCuboxAfterSave {
			deleteCount++
		}
	}
	if deleteCount > 1 {
		return errors.New("只能有一个archiver设置为deleteCuboxAfterSave")
	}
	return nil
}

func NewArchiverFromConfig(conf config.Archiver) (core.Archiver, error) {
	if conf.Options == nil {
		conf.Options = make(map[string]string)
	}
	if conf.Enable == false {
		return nil, fmt.Errorf("archiver %s is disabled", conf.Type)
	}
	switch conf.Type {
	case "notion":
		token := conf.Options["token"]
		databaseID := conf.Options["databaseID"]
		return notion.NewArchiver(token, databaseID)
	default:
		return nil, fmt.Errorf("未被支持的archiver类型: %s", conf.Type)
	}
}
