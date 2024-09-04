package ytdlp

import (
	"github.com/ge-fei-fan/gefflog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

// Used to unmarshall ytdlp progress
type ProgressTemplate struct {
	Percentage string  `json:"percentage"`
	Speed      float32 `json:"speed"`
	Size       string  `json:"size"`
	Eta        float32 `json:"eta"`
}

// Defines where and how the download needs to be saved
type DownloadOutput struct {
	Path          string
	Filename      string
	SavedFilePath string `json:"savedFilePath"`
}

// Progress for the Running call
type DownloadProgress struct {
	Status     int     `json:"process_status"`
	Percentage string  `json:"percentage"`
	Speed      float32 `json:"speed"`
	ETA        float32 `json:"eta"`
}

// Used to deser the ytdlp -J output
type DownloadInfo struct {
	Id          string    `json:"Id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Thumbnail   string    `json:"thumbnail"`
	Resolution  string    `json:"resolution"`
	Size        int64     `json:"filesize_approx"`
	VCodec      string    `json:"vcodec"`
	ACodec      string    `json:"acodec"`
	Extension   string    `json:"ext"`
	OriginalURL string    `json:"original_url"`
	FileName    string    `json:"filename"`
	CreatedAt   time.Time `json:"created_at"`
	//RequestedDownloads struct {
	//	Filename string `json:"filename"`
	//} `json:"requested_downloads"`
}

// Used to deser the formats in the -J output
type DownloadFormats struct {
	Formats   []Format `json:"formats"`
	Best      Format   `json:"best"`
	Thumbnail string   `json:"thumbnail"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
}

// A skimmed ytdlp format node
type Format struct {
	Format_id   string  `json:"format_id"`
	Format_note string  `json:"format_note"`
	FPS         float32 `json:"fps"`
	Resolution  string  `json:"resolution"`
	VCodec      string  `json:"vcodec"`
	ACodec      string  `json:"acodec"`
	Size        float32 `json:"filesize_approx"`
}

// struct representing the response sent to the client
// as JSON-RPC result field
type ProcessResponse struct {
	Id       string           `json:"id"`
	Url      string           `json:"url"`
	Progress DownloadProgress `json:"progress"`
	Info     DownloadInfo     `json:"info"`
	Output   DownloadOutput   `json:"output"`
	Params   []string         `json:"params"`
}

// struct representing the current status of the memoryDB
// used for serializaton/persistence reasons
type Session struct {
	Processes []ProcessResponse `json:"processes"`
}

// struct representing the intent to stop a specific process
type AbortRequest struct {
	Id string `json:"id"`
}

// struct representing the intent to start a download
type DownloadRequest struct {
	Id     string
	URL    string   `json:"url"`
	Path   string   `json:"path"`
	Rename string   `json:"rename"`
	Params []string `json:"params"`
}

// struct representing request of creating a netscape cookies file
type SetCookiesRequest struct {
	Cookies string `json:"cookies"`
}

// represents a user defined collection of ytdlp arguments
type CustomTemplate struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type YtDlpConfig struct {
	BasePath     string
	DownloadPath string        `yaml:"downloadPath"` //视频保存路径
	YtDlpPath    string        `yaml:"ytDlpPath"`    //下载程序路径
	QueueSize    int           `yaml:"queueSize"`
	Mdb          *MemoryDB     `yaml:"-"`
	Mq           *MessageQueue `yaml:"-"`
}

var YdpConfig YtDlpConfig

func InitYtDlpConfig(basePath string) {
	var mdb MemoryDB
	YdpConfig = YtDlpConfig{
		BasePath:     basePath,
		DownloadPath: basePath + "/data/yt-dlp-download/",
		YtDlpPath:    basePath + "/data/yt-dlp/yt-dlp.exe",
		QueueSize:    8,
		Mdb:          &mdb,
	}
	b, err := os.ReadFile(basePath + "/data/ytdlp.yaml")
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(basePath+"/data/ytdlp.yaml"), os.ModePerm)
		content, err := yaml.Marshal(YdpConfig)
		if err != nil {
			gefflog.Err("生成ytdlp.yaml配置失败")
		}
		if err = os.WriteFile(basePath+"/data/ytdlp.yaml", content, 0644); err != nil {
			gefflog.Err("生成ytdlp.yaml配置失败")
		}
	} else {
		if err := yaml.Unmarshal(b, &YdpConfig); err != nil {
			gefflog.Err("获取ytdlp.yaml配置失败")
		}
	}

	mq, err := NewMessageQueue()
	if err != nil {
		panic(err)
	}
	YdpConfig.Mq = mq
	YdpConfig.Mq.SetupConsumers()

	go YdpConfig.Mdb.Restore(basePath, mq)

}
