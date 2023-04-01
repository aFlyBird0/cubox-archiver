package util

import "github.com/parnurzeal/gorequest"

func SetGoRequestHeader(request *gorequest.SuperAgent, auth, cookie string) *gorequest.SuperAgent {
	request.Set("authority", "cubox.pro").
		Set("accept", "application/json, text/plain, */*").
		Set("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7").
		Set("authorization", auth).
		Set("cache-control", "no-cache").
		Set("cookie", cookie).
		Set("dnt", "1").
		Set("pragma", "no-cache").
		Set("referer", "https://cubox.pro/my/inbox").
		Set("sec-ch-ua", "\"Chromium\";v=\"111\", \"Not(A:Brand\";v=\"8\"").
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "\"macOS\"").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-origin").
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	return request
}
