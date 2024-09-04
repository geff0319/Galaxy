package bridge

import (
	"encoding/json"
	"fmt"
	"galaxy/bridge/ytdlp"
	"github.com/ge-fei-fan/gefflog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
	"os"
)

// 获取视频清晰度和名称
func (a *App) GetVideoMeta(url string) FlagResultWithData {
	p := &ytdlp.Process{
		Url:    url,
		Params: []string{},
		Output: ytdlp.DownloadOutput{},
	}
	p.SetPending()
	if err := p.SetMetadata(); err != nil {
		gefflog.Err("解析视频失败: " + err.Error())
		return FlagResultWithData{false, fmt.Sprintf("解析视频失败"), nil}
	}
	return FlagResultWithData{true, "解析视频成功", p}
}

func (a *App) DownloadYoutubeByKey(p string, retry bool) FlagResult {

	var process ytdlp.Process
	if err := json.Unmarshal([]byte(p), &process); err != nil {
		return FlagResult{false, "下载视频失败,解析任务信息出错"}
	}
	if !retry {
		if process.Info.Id == "" {
			return FlagResult{false, "下载视频失败,任务信息为空"}
		}
		if ytdlp.YdpConfig.Mdb.IsProcessExist(&process) {
			return FlagResult{false, "下载视频失败,任务已存在"}
		} else {
			ytdlp.YdpConfig.Mdb.Set(&process)
			ytdlp.YdpConfig.Mq.PublishByTopic("process:downloading", &process)
			return FlagResult{true, process.Id}
		}
	} else {
		dstPrecess, err := ytdlp.YdpConfig.Mdb.Get(process.Id)
		if err != nil {
			return FlagResult{false, "任务id获取失败"}
		}
		if dstPrecess.Info.Id == "" {
			gefflog.Info("重试解析: " + dstPrecess.Id)
			ytdlp.YdpConfig.Mq.Publish(dstPrecess)
		} else {
			gefflog.Info("重试下载: " + dstPrecess.Id)
			dstPrecess.SetPending()
			ytdlp.YdpConfig.Mq.PublishByTopic("process:downloading", dstPrecess)
		}
		return FlagResult{true, dstPrecess.Id}
	}

}

func (a *App) DownloadYoutube(url string, params []string) FlagResult {
	p := &ytdlp.Process{
		Url:    url,
		Params: params,
		Output: ytdlp.DownloadOutput{},
	}
	id := ytdlp.YdpConfig.Mdb.Set(p)
	ytdlp.YdpConfig.Mq.Publish(p)
	return FlagResult{Flag: true, Data: id}
}

// 缓存的数据持久化到文件
func (a *App) Persist() FlagResult {
	err := ytdlp.YdpConfig.Mdb.Persist(Env.BasePath)
	if err != nil {
		gefflog.Err("视频信息保存失败")
		return FlagResult{Flag: false, Data: err.Error()}
	}
	return FlagResult{Flag: true, Data: "视频信息保存成功"}
}

func (a *App) All() FlagResultWithData {
	res := ytdlp.YdpConfig.Mdb.All()
	return FlagResultWithData{
		Flag: true,
		Msg:  "操作成功",
		Data: res,
	}
}

func (a *App) UpdateYtDlpConfig() FlagResult {
	b, err := os.ReadFile(Env.BasePath + "/data/ytdlp.yaml")
	if err != nil {
		gefflog.Err("更新配置失败: " + err.Error())
		return FlagResult{false, "更新配置失败"}
	}
	if err := yaml.Unmarshal(b, &ytdlp.YdpConfig); err != nil {
		gefflog.Err("更新配置失败: " + err.Error())
		return FlagResult{false, "更新配置失败"}
	}
	return FlagResult{true, "更新配置成功"}
}

func (a *App) Delete(id string) FlagResult {
	gefflog.Info("删除任务 id:" + id)
	p, err := ytdlp.YdpConfig.Mdb.Get(id)
	if err != nil {
		gefflog.Err(fmt.Sprintf("删除任务失败 id: %s, err: %s", id, err.Error()))
		return FlagResult{false, "删除任务失败:"}
	}
	if p == nil {
		return FlagResult{false, "删除任务失败: nil process"}
	}
	if p.Progress.Status == ytdlp.StatusPending || p.Progress.Status == ytdlp.StatusDownloading {
		if err := p.Kill(); err != nil {
			gefflog.Err(fmt.Sprintf("删除任务失败 id: %s, err: %s", id, err.Error()))
			return FlagResult{false, "删除任务失败"}
		}
		ytdlp.YdpConfig.Mdb.Delete(p.Id)
	} else {
		ytdlp.YdpConfig.Mdb.Delete(p.Id)
	}

	return FlagResult{true, "删除成功"}
}

func MqNotifyConsumer(a *App) {
	err := ytdlp.YdpConfig.Mq.SetConsumer("notify", func(level, data string) {
		runtime.EventsEmit(a.Ctx, "notify", level, data)
	})
	if err != nil {
		gefflog.Err(err.Error())
	}
}
