package ytdlp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"
)

const template = `download:
{
	"eta":%(progress.eta)s, 
	"percentage":"%(progress._percent_str)s",
	"speed":%(progress.speed)s
}`

const (
	StatusPending = iota
	StatusDownloading
	StatusCompleted
	StatusErrored
)

// Process descriptor
type Process struct {
	Id       string           `json:"id"`
	Url      string           `json:"url"`
	Params   []string         `json:"params"`
	Info     DownloadInfo     `json:"info"`
	Progress DownloadProgress `json:"progress"`
	Output   DownloadOutput   `json:"output"`
	proc     *os.Process
	//Logger   *slog.Logger
}

// Starts spawns/forks a new ytdlp process and parse its stdout.
// The process is spawned to outputting a custom progress text that
// Resembles a JSON Object in order to Unmarshal it later.
// This approach is anyhow not perfect: quotes are not escaped properly.
// Each process is not identified by its PID but by a UUIDv4
func (p *Process) Start() {
	// escape bash variable escaping and command piping, you'll never know
	// what they might come with...
	p.Params = slices.DeleteFunc(p.Params, func(e string) bool {
		match, _ := regexp.MatchString(`(\$\{)|(\&\&)`, e)
		return match
	})

	p.Params = slices.DeleteFunc(p.Params, func(e string) bool {
		return e == ""
	})

	//out := DownloadOutput{
	//	Path:     filepath.Join(bridge.Env.BasePath, "video"),
	//	Filename: "%(title)s.%(ext)s",
	//}
	//
	//if p.Output.Path != "" {
	//	out.Path = p.Output.Path
	//}
	//
	//if p.Output.Filename != "" {
	//	out.Filename = p.Output.Filename
	//}
	//
	//buildFilename(&p.Output)

	//TODO: it spawn another one ytdlp process, too slow.
	//go p.GetFileName(&out)
	p.Output.SavedFilePath = filepath.Join(YdpConfig.DownloadPath, p.Info.FileName)
	gefflog.Info(p.Info.FileName)
	baseParams := []string{
		strings.Split(p.Url, "?list")[0], //no playlist
		"--newline",
		"--no-colors",
		"--no-playlist",
		"--progress-template",
		strings.NewReplacer("\n", "", "\t", "", " ", "").Replace(template),
	}

	// if user asked to manually override the output path...
	if !(slices.Contains(p.Params, "-P") || slices.Contains(p.Params, "--paths")) {
		p.Params = append(p.Params, "-o")
		p.Params = append(p.Params, p.Output.SavedFilePath)
	}
	if strings.Contains(p.Url, "x.com") {
		if !IsFileExist(YdpConfig.BasePath + "/data/yt-dlp/cookies.txt") {
			gefflog.Err("下载X视频,cookies.txt不存在")
		}
		baseParams = append(baseParams, "--cookies", YdpConfig.BasePath+"/data/yt-dlp/cookies.txt")
	}

	params := append(baseParams, p.Params...)
	gefflog.Info(params)

	// ----------------- main block ----------------- //
	if !IsYtDlpExist() {
		YdpConfig.Mq.eventBus.Publish("notify", "error", "启动任务失败:ytdlp程序不存在,请下载")
		p.Progress.Status = StatusErrored
		return
	}
	cmd := exec.Command(YdpConfig.YtDlpPath, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	r, err := cmd.StdoutPipe()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to connect to stdout: err=%s", err.Error()))
		YdpConfig.Mq.eventBus.Publish("notify", "error", "启动任务失败")
		p.Progress.Status = StatusErrored
		return
	}

	if err := cmd.Start(); err != nil {
		gefflog.Err(fmt.Sprintf("failed to start ytdlp process: err=%s", err.Error()))
		YdpConfig.Mq.eventBus.Publish("notify", "error", "启动任务失败")
		p.Progress.Status = StatusErrored
		return
	}

	p.proc = cmd.Process

	// --------------- progress block --------------- //
	var (
		sourceChan = make(chan []byte)
		doneChan   = make(chan struct{})
	)

	// spawn a goroutine that does the dirty job of parsing the stdout
	// filling the channel with as many stdout line as ytdlp produces (producer)
	go func() {
		scan := bufio.NewScanner(r)

		defer func() {
			r.Close()
			p.Complete()

			doneChan <- struct{}{}

			close(sourceChan)
			close(doneChan)
		}()

		for scan.Scan() {
			sourceChan <- scan.Bytes()
		}
	}()

	// Slows down the unmarshal operation to every 500ms
	go func() {
		Sample(time.Millisecond*500, sourceChan, doneChan, func(event []byte) {
			var progress ProgressTemplate

			if err := json.Unmarshal(event, &progress); err != nil {
				return
			}

			p.Progress = DownloadProgress{
				Status:     StatusDownloading,
				Percentage: progress.Percentage,
				Speed:      progress.Speed,
				ETA:        progress.Eta,
			}
			//gefflog.Info(fmt.Sprintf("progress: id=%s, url=%s, percentage=%s", p.getShortId(), p.Url, progress.Percentage))

		})
	}()

	// ------------- end progress block ------------- //
	cmd.Wait()
}

// Keep process in the memoryDB but marks it as complete
// Convention: All completed processes has progress -1
// and speed 0 bps.
func (p *Process) Complete() {
	p.Progress = DownloadProgress{
		Status:     StatusCompleted,
		Percentage: "-1",
		Speed:      0,
		ETA:        0,
	}
	gefflog.Info(fmt.Sprintf("finished: id=%s, url=%s", p.getShortId(), p.Url))

}

// Kill a process and remove it from the memory
func (p *Process) Kill() error {
	defer func() {
		p.Progress.Status = StatusCompleted
	}()
	// ytdlp uses multiple child process the parent process
	// has been spawned with setPgid = true. To properly kill
	// all subprocesses a SIGTERM need to be sent to the correct
	// process group
	if p.proc == nil {
		return errors.New("*os.Process not set")
	}
	p.proc.Kill()
	//process, err := os.FindProcess(p.proc.Pid)
	//if err != nil {
	//	return err
	//}
	//if err = process.Signal(syscall.SIGKILL); err != nil {
	//	return err
	//}

	return nil
}

// Returns the available format for this URL
// TODO: Move out from process.go
func (p *Process) GetFormatsSync() (DownloadFormats, error) {
	cmd := exec.Command(YdpConfig.YtDlpPath, p.Url, "-J")

	stdout, err := cmd.Output()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to retrieve metadata: err=%s", err.Error()))
		//p.Logger.Error("failed to retrieve metadata", slog.String("err", err.Error()))
		return DownloadFormats{}, err
	}

	info := DownloadFormats{URL: p.Url}
	best := Format{}

	var (
		wg            sync.WaitGroup
		decodingError error
	)

	wg.Add(2)

	log.Println(
		BgRed, "Metadata", Reset,
		BgBlue, "Formats", Reset,
		p.Url,
	)
	gefflog.Info(fmt.Sprintf("retrieving metadata: caller=%s, url=%s", "getFormats", p.Url))
	//p.Logger.Info(
	//	"retrieving metadata",
	//	slog.String("caller", "getFormats"),
	//	slog.String("url", p.Url),
	//)

	go func() {
		decodingError = json.Unmarshal(stdout, &info)
		wg.Done()
	}()

	go func() {
		decodingError = json.Unmarshal(stdout, &best)
		wg.Done()
	}()

	wg.Wait()

	if decodingError != nil {
		return DownloadFormats{}, err
	}

	info.Best = best

	return info, nil
}

func (p *Process) GetFileName(o *DownloadOutput) error {
	cmd := exec.Command(
		YdpConfig.YtDlpPath,
		"--print", "filename",
		"-o", fmt.Sprintf("%s/%s", o.Path, o.Filename),
		p.Url,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	p.Output.SavedFilePath = strings.Trim(string(out), "\n")
	return nil
}

func (p *Process) SetPending() {
	// Since video's title isn't available yet, fill in with the URL.
	p.Info = DownloadInfo{
		URL:       p.Url,
		Title:     p.Url,
		CreatedAt: time.Now(),
	}
	p.Progress.Status = StatusPending
}

func (p *Process) SetMetadata() error {
	//检查ytdlp程序是否存在
	if !IsYtDlpExist() {
		return errors.New("tydlp程序不存在,请下载")
	}
	//cmd := exec.Command(YdpConfig.YtDlpPath, p.Url, "-J")
	baseParams := []string{
		strings.Split(p.Url, "?list")[0],
		"--dump-json",
		"--no-warnings",
	}
	if strings.Contains(p.Url, "x.com") {
		if !IsFileExist(YdpConfig.BasePath + "/data/yt-dlp/cookies.txt") {
			return errors.New("下载X视频,cookies.txt不存在")
		}
		baseParams = append(baseParams, "--cookies", YdpConfig.BasePath+"/data/yt-dlp/cookies.txt")
	}
	params := append(baseParams, p.Params...)
	fmt.Println(p.Params)
	fmt.Println(params)
	cmd := exec.Command(YdpConfig.YtDlpPath, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to connect to stdout: id=%s, url=%s, err=%s", p.getShortId(), p.Url, err.Error()))
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to connect to stderr: id=%s, url=%s, err=%s", p.getShortId(), p.Url, err.Error()))
		return err
	}

	info := DownloadInfo{
		URL:       p.Url,
		CreatedAt: time.Now(),
	}

	if err := cmd.Start(); err != nil {
		gefflog.Err(fmt.Sprintf("failed to start cmd: id=%s, url=%s, err=%s", p.getShortId(), p.Url, err.Error()))
		return err
	}

	var bufferedStderr bytes.Buffer

	go func() {
		io.Copy(&bufferedStderr, stderr)
	}()
	gefflog.Info(fmt.Sprintf("retrieving metadata: id=%s, url=%s", p.getShortId(), p.Url))

	if err := json.NewDecoder(stdout).Decode(&info); err != nil {
		gefflog.Err(fmt.Sprintf("failed to Decode json : id=%s, url=%s, err=%s", p.getShortId(), p.Url, err.Error()))
		gefflog.Err(bufferedStderr.String())
		return err
	}
	gefflog.Info(info)
	p.Info = info
	p.Progress.Status = StatusPending

	if err := cmd.Wait(); err != nil {
		gefflog.Err(fmt.Sprintf("failed to wait cmd: id=%s, url=%s, err=%s", p.getShortId(), p.Url, err.Error()))
		return errors.New(bufferedStderr.String())
	}

	return nil
}

func (p *Process) getShortId() string { return strings.Split(p.Id, "-")[0] }

func buildFilename(o *DownloadOutput) {
	if o.Filename != "" && strings.Contains(o.Filename, ".%(ext)s") {
		o.Filename += ".%(ext)s"
	}

	o.Filename = strings.Replace(
		o.Filename,
		".%(ext)s.%(ext)s",
		".%(ext)s",
		1,
	)
}

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
