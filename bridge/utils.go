package bridge

import (
	"github.com/lxn/win"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func GetPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	path = filepath.Join(Env.BasePath, path)
	path = filepath.Clean(path)
	return path
}

func GetProxy(_proxy string) func(*http.Request) (*url.URL, error) {
	proxy := http.ProxyFromEnvironment

	if _proxy != "" {
		if !strings.HasPrefix(_proxy, "http") {
			_proxy = "http://" + _proxy
		}
		proxyUrl, err := url.Parse(_proxy)
		if err == nil {
			proxy = http.ProxyURL(proxyUrl)
		}
	}

	return proxy
}

func ConvertByte2String(byte []byte) string {
	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
	return string(decodeBytes)
}

func (a *App) AbsolutePath(path string) FlagResult {
	log.Printf("AbsolutePath: %s", path)

	if filepath.IsAbs(path) {
		return FlagResult{true, path}
	}

	path = GetPath(path)

	return FlagResult{true, path}
}

func (a *App) HideToolWindow() {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(AppTitle))
	// 获取当前窗口的扩展风格
	currentExStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	// 添加 WS_EX_TOOLWINDOW 样式
	newExStyle := (currentExStyle | win.WS_EX_TOOLWINDOW) &^ win.WS_EX_APPWINDOW
	//win.SetWindowLong(hwnd, win.GWL_EXSTYLE, win.GetWindowLong(hwnd, win.GWL_EXSTYLE)|win.WS_EX_LAYERED)
	// 应用新的扩展风格
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, newExStyle)
}
func (a *App) ShowToolWindow() {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(AppTitle))
	// 获取当前窗口的扩展风格
	currentExStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	// 添加 WS_EX_TOOLWINDOW 样式
	newExStyle := currentExStyle &^ win.WS_EX_TOOLWINDOW
	// 应用新的扩展风格
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, newExStyle)
}
