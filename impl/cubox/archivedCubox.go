package cubox

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"

	"github.com/aFlyBird0/cubox-archiver/core"
	"github.com/aFlyBird0/cubox-archiver/util"
)

func NewArchivedCuboxSource(auth, cookie string) *ArchivedCuboxSource {
	return &ArchivedCuboxSource{auth: auth, cookie: cookie}
}

type ArchivedCuboxSource struct {
	auth   string
	cookie string
}

var _ = core.Source(&ArchivedCuboxSource{})

func (client *ArchivedCuboxSource) List(cuboxChan chan *core.Item) {
	const archiving = true
	// 先请求第一页试试
	items, pageCount, totalCounts := client.requestCubox(archiving, 1, "")
	logrus.Infof("Cubox共有%d条记录\n", totalCounts)
	logrus.Info("已经获取到第1页的Cubox记录, 共", len(items), "条")
	for _, item := range items {
		cuboxChan <- item
	}

	// 并发处理后面的页
	if pageCount >= 2 && len(items) >= 1 {
		lastID := items[len(items)-1].UserSearchEngineID
		client.handleNextPages(cuboxChan, archiving, pageCount, lastID)
	} else {
		close(cuboxChan)
		logrus.Info("已经获取完所有的Cubox记录")
	}
}

func (client *ArchivedCuboxSource) requestCubox(archiving bool, page int, lastBookmarkId string) (res []*core.Item, pageCount, totalCount int) {
	const url = "https://cubox.pro/c/api/v2/search_engine/my"

	dataResp := cuboxItemResponse{}
	request := gorequest.New().Get(url)
	if lastBookmarkId != "" {
		request = request.Param("lastBookmarkId", lastBookmarkId)
	}
	request = request.
		Param("asc", "false").
		Param("page", strconv.FormatInt(int64(page), 10)).
		Param("filters", "").
		Param("archiving", strconv.FormatBool(archiving))
	request = SetGoRequestHeader(request, client.auth, client.cookie)
	httpResp, body, errs := request.EndStruct(&dataResp)
	if len(errs) > 0 {
		logrus.Fatalln(fmt.Sprintf("failed to request cubox content, err: %v", errs))
	}
	if httpResp == nil || httpResp.StatusCode != http.StatusOK {
		logrus.Fatalln(fmt.Sprintf("failed to request cubox content, http code: %v, body: %v", httpResp.StatusCode, body))
	}

	res = make([]*core.Item, 0, len(dataResp.Data))

	for _, itemRaw := range dataResp.Data {
		if itemRaw != nil {
			res = append(res, client.convertCuboxItem(itemRaw))
		}
	}

	return res, dataResp.PageCount, dataResp.TotalCounts
}

func (client *ArchivedCuboxSource) handleNextPages(cuboxChan chan<- *core.Item, archiving bool, pageCount int, lastBookmarkId string) {
	for page := 2; page <= pageCount; page += 1 {
		time.Sleep(time.Second * 1)
		items, _, _ := client.requestCubox(archiving, page, lastBookmarkId)
		logrus.Infof("已经获取到第%d页的Cubox记录, 共%d条\n", page, len(items))
		for _, item := range items {
			cuboxChan <- item
		}
		lastBookmarkId = items[len(items)-1].UserSearchEngineID
	}
	close(cuboxChan)
	logrus.Info("已经获取完所有的Cubox记录")
}

// 将 cuboxItemRaw 转换成 Item
func (client *ArchivedCuboxSource) convertCuboxItem(raw *cuboxItemRaw) (item *core.Item) {
	item = &core.Item{}
	item.UserSearchEngineID = raw.UserSearchEngineID
	item.Title = raw.Title
	item.Description = raw.Description
	item.TargetURL = raw.TargetURL
	//item.ArchiveName = raw.ArchiveName
	//item.ArticleName = raw.ArticleName
	if raw.Cover != "" {
		item.Cover = "https://image.cubox.pro/" + raw.Cover
	}
	item.LittleIcon = raw.LittleIcon
	for _, tag := range raw.Tags {
		updateTime, err := convertTime(tag.UpdateTime)
		if err != nil {
			logrus.Errorf("convert time failed, err: %v", err)
		}
		item.Tags = append(item.Tags, core.Tag{TagID: tag.TagID, Name: tag.Name, Rank: tag.Rank, UpdateTime: updateTime, ParentId: tag.ParentId})
	}
	item.GroupId = raw.GroupId
	item.GroupName = raw.GroupName
	createTime, err := convertTime(raw.CreateTime)
	if err != nil {
		logrus.Errorf("convert time failed, err: %v", err)
	}
	item.CreateTime = createTime
	updateTime, err := convertTime(raw.UpdateTime)
	if err != nil {
		logrus.Errorf("convert time failed, err: %v", err)
	}
	item.UpdateTime = updateTime
	item.Status = raw.Status
	//item.Finished = raw.Finished
	//item.InBlackOrWhiteList = raw.InBlackOrWhiteList
	item.Type = core.CuboxContentType(raw.Type)

	// 把链接全部 encode 一下（因为Notion里的链接会被自动解码，导致去重失败）
	item.TargetURL, err = util.EncodeURL(item.TargetURL)
	if err != nil {
		logrus.Errorf("encode url failed, err: %v", err)
	}
	item.Cover, err = util.EncodeURL(item.Cover)
	if err != nil {
		logrus.Errorf("encode url failed, err: %v", err)
	}
	item.LittleIcon, err = util.EncodeURL(item.LittleIcon)
	if err != nil {
		logrus.Errorf("encode url failed, err: %v", err)
	}

	// 单独处理文本和随手记类型的内容
	if item.Type == core.Text || item.Type == core.ShortHand {
		getContentURL := "https://cubox.pro/c/api/bookmark/content"
		request := gorequest.New().Get(getContentURL).Timeout(time.Second*10).
			Param("bookmarkId", item.UserSearchEngineID)
		request = SetGoRequestHeader(request, client.auth, client.cookie)

		type response struct {
			Code int `json:"code"`
			Data struct {
				Content string `json:"content"`
				Marks   any    `json:"marks"`
			}
			Message string `json:"message"`
		}
		var resp response
		httpResp, body, errs := request.EndStruct(&resp)
		if errs != nil {
			logrus.Errorf("请求cubox文本类型的内容失败, err: %v", errs)
			return item
		}
		if httpResp.StatusCode != 200 {
			logrus.Errorf("请求cubox文本类型的内容失败, status code: %d, body: %s\n", httpResp.StatusCode, body)
			return item
		}
		if resp.Code != 200 {
			logrus.Errorf("请求cubox文本类型的内容失败, code: %d, message: %s\n", resp.Code, resp.Message)
			return item
		}
		item.Content = resp.Data.Content
	}

	return item
}

// 将时间字符串转换成 time.Time
// 2022-07-31T01:26:33:603+08:00
func convertTime(raw string) (time.Time, error) {
	if raw == "" || raw == "null" {
		return time.Time{}, nil
	}

	// 把倒数第二个冒号替换成点，以转换成标准的 RFC3339 格式
	timeRFC3339 := formatTimeToRFC3339(raw)
	createTime, err := time.Parse("2006-01-02T15:04:05.999Z07:00", timeRFC3339)
	return createTime, err
}

// 把时间的倒数第二个冒号替换成点，以转换成标准的 RFC3339 格式
func formatTimeToRFC3339(s string) string {
	i := strings.LastIndex(s[:strings.LastIndex(s, ":")], ":")
	if i == -1 {
		return s
	}

	return s[:i] + "." + s[i+1:]
}

type cuboxItemResponse struct {
	Message     string          `json:"message"`
	Code        int             `json:"code"`
	Data        []*cuboxItemRaw `json:"data"`
	PageCount   int             `json:"pageCount"`
	TotalCounts int             `json:"totalCounts"`
}
