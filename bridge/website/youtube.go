package website

import "regexp"

type Youtube struct {
	sourceUrl string
}

func NewYoutube(url string) *Youtube {
	return &Youtube{
		sourceUrl: url,
	}
}

func (y *Youtube) Compile() (string, bool) {
	return "", false
}

func (y *Youtube) AppCompile() (string, bool) {
	// 定义 YouTube 短链接的正则表达式
	youtubeShortURLRegex := `https://youtu\.be/[\w-]+`

	// 编译正则表达式
	youtubeRe := regexp.MustCompile(youtubeShortURLRegex)

	// 查找所有匹配的链接
	youtubeLinks := youtubeRe.FindAllString(y.sourceUrl, -1)

	//// 输出结果
	if len(youtubeLinks) != 0 {
		return youtubeLinks[0], true
	}
	return "", false
}
