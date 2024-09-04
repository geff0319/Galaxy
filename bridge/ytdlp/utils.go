package ytdlp

import (
	"regexp"
)

func AppYoutubeCompile(str string) (string, bool) {
	// 定义 YouTube 短链接的正则表达式
	youtubeShortURLRegex := `https://youtu\.be/[\w-]+`

	// 编译正则表达式
	youtubeRe := regexp.MustCompile(youtubeShortURLRegex)

	// 查找所有匹配的链接
	youtubeLinks := youtubeRe.FindAllString(str, -1)

	//// 输出结果
	if len(youtubeLinks) != 0 {
		return youtubeLinks[0], true
	}
	return "", false
}

func AppXCompile(str string) (string, bool) {
	xURLpattern := `^https://x\.com.*`

	// 编译正则表达式
	xRe := regexp.MustCompile(xURLpattern)
	xLinks := xRe.FindAllString(str, -1)
	if len(xLinks) != 0 {
		return xLinks[0], true
	}
	return "", false
}
