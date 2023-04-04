package cubox

// cuboxItemRaw 网页接口获得的数据结构
type cuboxItemRaw struct {
	UserSearchEngineID string      `json:"userSearchEngineID"` // Cubox 内部 ID
	Title              string      `json:"title"`              // 标题
	Description        string      `json:"description"`        // 描述
	TargetURL          string      `json:"targetURL"`          // 原始链接
	ResourceURL        string      `json:"resourceURL"`
	HomeURL            string      `json:"homeURL"`
	ArchiveName        string      `json:"archiveName"` // 不明，感觉是 Cubox 内部用的
	Content            string      `json:"content"`
	ArticleName        string      `json:"articleName"` // 不明
	ArticleWordCount   int         `json:"articleWordCount"`
	Byline             string      `json:"byline"`
	Cover              string      `json:"cover"`      // 封面图片链接（注意网页请求返回的是相对路径，cover/2022081418282457259/98737.jpg，这是绝对：https://image.cubox.pro/cover/2022081418282457259/98737.jpg）
	ArticleURL         string      `json:"articleURL"` // 不明
	LittleIcon         string      `json:"littleIcon"` // 小图标
	Archiving          bool        `json:"archiving"`
	StarTarget         bool        `json:"starTarget"`
	HasMark            bool        `json:"hasMark"`
	IsRead             interface{} `json:"isRead"`
	MarkCount          int         `json:"markCount"`
	Tags               []struct {
		TagID      string      `json:"tagID"`
		Name       string      `json:"name"`
		Rank       interface{} `json:"rank"` // 类型不明
		UpdateTime string      `json:"updateTime"`
		ParentId   string      `json:"parentId"`
	} `json:"tags"`
	AllTags            []interface{} `json:"allTags"`
	Marks              []interface{} `json:"marks"`
	GroupId            string        `json:"groupId"`
	GroupName          string        `json:"groupName"`
	CreateTime         string        `json:"createTime"` // 2022-09-04T15:15:13:151+08:00
	UpdateTime         string        `json:"updateTime"`
	Status             string        `json:"status"`
	Finished           bool          `json:"finished"`
	InBlackOrWhiteList bool          `json:"inBlackOrWhiteList"`
	Type               int           `json:"type"` // 类型，已知 0 是链接文章，1 是文本，2 是速记，3 是图片
}
