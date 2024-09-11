package website

import "github.com/ge-fei-fan/gefflog"

type Handler interface {
	// Compile 过滤网页端视频地址
	Compile() (string, bool)

	// AppCompile 过滤手机端分享视频地址
	AppCompile() (string, bool)
}

func NewHandler(source, url string) Handler {
	switch source {
	case "bilibili":
		return NewBlibili(url)
	case "youtube":
		return NewYoutube(url)
	case "twitter":
		return NewTwitter(url)
	default:
		return nil
	}
}

func PreprocessApp(source, url string) (string, bool) {
	h := NewHandler(source, url)
	if h == nil {
		gefflog.Info("获取到未知source： " + source)
		return "", false
	}
	return h.AppCompile()
}
