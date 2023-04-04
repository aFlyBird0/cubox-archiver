package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/aFlyBird0/cubox-archiver/config"
	"github.com/aFlyBird0/cubox-archiver/run"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("运行失败: %v", err)
	}
}

var (
	cfg     config.Config
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "github.com/aFlyBird0/cubox-archiver",
	Short: "转存Cubox稍后读内容，支持同步到多端并自动删除Cubox内容（目前仅支持 Notion）",
	Long:  `转存Cubox稍后读内容，支持同步到多端（目前仅支持 Notion）`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("请使用 --help 查看帮助")
	},
}

var fromFileCmd = &cobra.Command{
	Use:   "from-file",
	Short: "转存Cubox稍后读内容，从文件读取配置",
	Long:  `转存Cubox稍后读内容，从文件读取配置`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(cfgFile)
		if err != nil {
			logrus.Fatalf("Error reading Config file: %v", err)
		}

		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			logrus.Fatalf("Error parsing Config file: %v", err)
		}

		logrus.Info("配置读取成功")
		//logrus.Infof("%+v\n", cfg)

		run.Run(cfg)
	},
}

var (
	cuboxAuth   string
	cuboxCookie string

	notionEnabled      bool
	notionToken        string
	notionPageID       string
	notionDatabaseID   string
	notionDatabaseName string

	deleteCuboxAfterSaveToNotion bool
)

var fromFlagCmd = &cobra.Command{
	Use:   "from-flag",
	Short: "转存Cubox稍后读内容，从命令行读取配置",
	Long:  `转存Cubox稍后读内容，从命令行读取配置`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Config{
			Cubox: config.CuboxConfig{
				Auth:   cuboxAuth,
				Cookie: cuboxCookie,
			},
		}
		if notionEnabled {
			cfg.Archivers = append(cfg.Archivers, config.Archiver{
				Type:   "notion",
				Enable: notionEnabled,
				Options: map[string]string{
					"token":        notionToken,
					"pageID":       notionPageID,
					"databaseID":   notionDatabaseID,
					"databaseName": notionDatabaseName,
				},
				DeleteCuboxAfterSave: deleteCuboxAfterSaveToNotion,
			})
		}

		logrus.Info("配置读取成功")
		//logrus.Infof("%+v\n", cfg)
		run.Run(cfg)
	},
}

func init() {
	rootCmd.AddCommand(fromFileCmd)
	rootCmd.AddCommand(fromFlagCmd)

	fromFileCmd.Flags().StringVarP(&cfgFile, "file", "f", "config.yaml", "配置文件路径")

	fromFlagCmd.Flags().StringVarP(&cuboxAuth, "cubox-auth", "a", "", "Cubox Auth")
	fromFlagCmd.Flags().StringVarP(&cuboxCookie, "cubox-cookie", "c", "", "Cubox Cookie")

	fromFlagCmd.Flags().BoolVarP(&notionEnabled, "archiver-notion-enable", "e", false, "是否启用Notion")
	fromFlagCmd.Flags().StringVarP(&notionToken, "notion-token", "t", "", "Notion Token")
	fromFlagCmd.Flags().StringVarP(&notionPageID, "notion-page-id", "p", "", "Notion Page ID，将在该Page下自动创建Database")
	fromFlagCmd.Flags().StringVarP(&notionDatabaseID, "notion-database-id", "d", "", "Notion Database ID")
	fromFlagCmd.Flags().StringVarP(&notionDatabaseName, "notion-database-name", "n", "Cubox归档", "Notion Database Name")
	fromFlagCmd.Flags().BoolVarP(&deleteCuboxAfterSaveToNotion, "delete-cubox-after-save-to-notion", "s", false, "是否在保存到Notion后删除Cubox")
}
