package ytdlp

import (
	"os"
	"regexp"
)

func IsYtDlpExist() bool {
	_, err := os.Stat(YdpConfig.YtDlpPath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func IsDirExists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if f.IsDir() {
		return true
	}
	return false
}

// 过滤 Windows 不支持的文件名字符
func sanitizeFileName(name string) string {
	// 定义一个正则表达式模式来匹配 Windows 不支持的字符
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	// 用下划线替换这些字符
	sanitized := re.ReplaceAllString(name, "_")
	return sanitized
}
