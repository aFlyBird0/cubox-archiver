package config

type Config struct {
	Cubox CuboxConfig `yaml:"cubox"`
	//Notion               NotionConfig `yaml:"notion"`
	//DeleteCuboxAfterSave bool         `yaml:"deleteCuboxAfterSave"` // 是否在保存后删除Cubox
	Archivers []Archiver `yaml:"archivers"`
}

type Archiver struct {
	Type                 string            `yaml:"type"`
	Enable               bool              `yaml:"enable"`
	DeleteCuboxAfterSave bool              `yaml:"deleteCuboxAfterSave"` // 是否在保存后删除Cubox
	Options              map[string]string `yaml:"options"`
}

type CuboxConfig struct {
	Auth   string // Cubox Authorization Header	`
	Cookie string
}

type NotionConfig struct {
	Token        string `yaml:"token"`        // https://developers.notion.com/docs
	PageID       string `yaml:"pageID"`       // 如果传入了PageID, 将在该Page下自动创建Database
	DatabaseID   string `yaml:"databaseID"`   // 指定要使用的DatabaseID（可以手动创建，也可以使用PageID自动创建，再配置该项）
	DatabaseName string `yaml:"databaseName"` // 在Page下创建Database时使用，用于指定Database的名称。
}
