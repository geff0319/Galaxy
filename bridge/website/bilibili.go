package website

import "C"
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type progressWriter struct {
	totalBytes   int64
	currentBytes int64
	writer       io.Writer
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	if err != nil {
		return n, err
	}
	pw.currentBytes += int64(n)
	return n, nil
}

type Bilibili struct {
	sourceUrl string
}
type CidResponse struct {
	Data struct {
		Bvid  string `json:"bvid"`
		Cid   int64  `json:"cid"`
		Pic   string `json:"pic"`
		Title string `json:"title"`
	} `json:"data"`
}
type VideoInfoResponse struct {
	Data struct {
		AcceptDescription []string `json:"accept_description"`
		Dash              struct {
			Video []struct {
				BaseUrl string `json:"base_url"`
				Codecs  string `json:"codecs"`
				Width   int    `json:"width"`
				Height  int    `json:"height"`
			} `json:"video"`
			Audio []struct {
				BaseUrl string `json:"base_url"`
			} `json:"audio"`
		} `json:"dash"`
	} `json:"data"`
}
type BiliMetadata struct {
	ctx           context.Context
	Cl            context.CancelFunc `json:"-"`
	DoneChan      chan struct{}      `json:"-"`
	SavedFilePath string
	Cr            CidResponse
	Vir           VideoInfoResponse
	WriteFn       func(string, float32) `json:"-"`
	pWriter       *progressWriter
}

func NewBlibili(url string) *Bilibili {
	return &Bilibili{
		sourceUrl: url,
	}
}

func (b *Bilibili) Compile() (string, bool) {
	return "", false
}

func (b *Bilibili) AppCompile() (string, bool) {
	// 定义正则表达式来匹配 HTTP 链接
	re := regexp.MustCompile(`https?://[^\s]+`)

	// 查找匹配的链接
	Links := re.FindAllString(b.sourceUrl, -1)

	// 输出所有匹配的链接
	if len(Links) != 0 {
		req, err := http.NewRequest(http.MethodGet, Links[0], nil)
		if err != nil {
			gefflog.Err("解析b站分享链接失败：" + err.Error())
			return "", false
		}
		req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
		res, err := client.Do(req)
		if err != nil {
			gefflog.Err("解析b站分享链接失败：" + err.Error())
			return "", false
		}
		defer res.Body.Close()
		if res.Request.URL != nil {
			return res.Request.URL.String(), true
		}
		return "", false
	}
	return "", false
}

var client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		//禁止复用连接，防止同一个连接长时间大流量被限速
		DisableKeepAlives: true,
	},
}

func getCid(bv string) ([]byte, error) {
	url := "https://api.bilibili.com/x/web-interface/view?bvid=" + bv
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println(string(body))
	return body, nil
}
func getStream(bv, ck string, cid int64) (*VideoInfoResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/player/wbi/playurl?fnver=0&fnval=3216&fourk=1&qn=127&bvid=%s&cid=%d", bv, cid)
	fmt.Println(url)
	var err error
	url, err = sign(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	if ck != "" {
		req.AddCookie(&http.Cookie{Name: "SESSDATA", Value: ck})
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var vif VideoInfoResponse
	err = json.Unmarshal(body, &vif)
	if err != nil {
		return nil, err
	}

	return &vif, nil
}

func GetBilibiliInfo(url, ck string) (*BiliMetadata, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return nil, err
	}
	segments := strings.Split(u.Path, "/")
	if len(segments) < 2 {
		return nil, errors.New("获取bv失败")
	}
	bvid := segments[2]
	var cr CidResponse
	res, err := getCid(bvid)
	if err != nil {
		return nil, errors.New("getCid err: " + err.Error())
	}
	err = json.Unmarshal(res, &cr)
	if err != nil {
		return nil, err
	}
	infoResp, err := getStream(cr.Data.Bvid, ck, cr.Data.Cid)
	if err != nil {
		return nil, errors.New("getStream err: " + err.Error())
	}

	md := BiliMetadata{
		Cr:  cr,
		Vir: *infoResp,
	}

	return &md, nil
}

func (bmd *BiliMetadata) Download(basePath string) error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	bmd.ctx = ctx
	bmd.Cl = cancel
	bmd.DoneChan = make(chan struct{})

	err = bmd.DA()
	if err != nil {
		return err
	}
	err = bmd.DV()
	if err != nil {
		return err
	}
	err = bmd.Merge(basePath)
	if err != nil {
		return err
	}
	return nil
}

func (bmd *BiliMetadata) speed() {
	var lastBytes int64

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			speed := bmd.pWriter.currentBytes - lastBytes
			lastBytes = bmd.pWriter.currentBytes
			percentage := float64(bmd.pWriter.currentBytes) / float64(bmd.pWriter.totalBytes) * 100
			bmd.WriteFn(strconv.Itoa(int(percentage)), float32(speed))
		case <-bmd.DoneChan:
			ticker.Stop()
			gefflog.Info("speed 退出")
			//bmd.WriteFn("0",0)
			return
		}
	}
}
func (bmd *BiliMetadata) DV() error {

	req, err := http.NewRequestWithContext(bmd.ctx, http.MethodGet, bmd.Vir.Data.Dash.Video[0].BaseUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://www.bilibili.com")
	resp, err := client.Do(req)
	if err != nil {
		gefflog.Err("client request err: " + err.Error())
		return err
	}
	var file *os.File
	file, err = os.OpenFile(bmd.SavedFilePath+".video", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		gefflog.Err("OpenFile video err: " + err.Error())
		return err
	}
	defer file.Close()
	_ = file.Truncate(0)
	pw := progressWriter{
		totalBytes:   resp.ContentLength,
		currentBytes: int64(0),
		writer:       file,
	}
	fmt.Println(resp.ContentLength)
	bmd.pWriter = &pw
	go bmd.speed()
	//defer func() {
	//	bmd.DoneChan <- struct{}{}
	//	close(bmd.DoneChan)
	//}()
	_, err = io.Copy(bmd.pWriter, resp.Body)
	if err != nil {
		gefflog.Err("CopyFile video err: " + err.Error())
		return err
	}
	return nil
}

func (bmd *BiliMetadata) DA() error {
	req, err := http.NewRequestWithContext(bmd.ctx, http.MethodGet, bmd.Vir.Data.Dash.Audio[0].BaseUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://www.bilibili.com")
	resp, err := client.Do(req)
	if err != nil {
		gefflog.Err("client request err: " + err.Error())
		return err
	}
	var file *os.File
	file, err = os.OpenFile(bmd.SavedFilePath+".audio", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		gefflog.Err("OpenFile audio err: " + err.Error())
		return err
	}
	defer file.Close()
	_ = file.Truncate(0)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		gefflog.Err("CopyFile audio err: " + err.Error())
		return err
	}
	return nil
}

func (bmd *BiliMetadata) Merge(basePath string) error {

	video := bmd.SavedFilePath + ".audio"
	audio := bmd.SavedFilePath + ".video"
	ffPath := basePath + "/data/yt-dlp/ffmpeg.exe"
	cmd := exec.Command(ffPath, "-y", "-i", video, "-i", audio, "-c", "copy", bmd.SavedFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		gefflog.Err("Merge error: " + err.Error())
		return errors.New("Merge error: " + err.Error())
	}
	err = os.Remove(video)
	if err != nil {
		gefflog.Err("Remove video error: " + err.Error())
		return errors.New(video + "Remove video error: " + err.Error())
	}
	err = os.Remove(audio)
	if err != nil {
		gefflog.Err("Remove audio error: " + err.Error())
		return errors.New(audio + "Remove audio error: " + err.Error())
	}
	return nil
}
