package website

import "regexp"

type Twitter struct {
	sourceUrl string
}

func NewTwitter(url string) *Twitter {
	return &Twitter{
		sourceUrl: url,
	}
}
func (t *Twitter) Compile() (string, bool) {
	return "", true
}

func (t *Twitter) AppCompile() (string, bool) {
	xURLpattern := `^https://x\.com.*`

	// 编译正则表达式
	xRe := regexp.MustCompile(xURLpattern)
	xLinks := xRe.FindAllString(t.sourceUrl, -1)
	if len(xLinks) != 0 {
		return xLinks[0], true
	}
	return "", false
}
