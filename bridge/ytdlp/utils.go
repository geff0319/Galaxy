package ytdlp

import (
	"os"
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
