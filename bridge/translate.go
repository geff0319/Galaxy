package bridge

import (
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"gopkg.in/yaml.v3"
	"os"
)

func NewClient() (*tmt.Client, error) {
	b, err := os.ReadFile(Env.BasePath + "/data/user.yaml")
	if err != nil {
		return nil, errors.New("翻译密钥未配置")
	}
	yaml.Unmarshal(b, &Config)
	credential := common.NewCredential(Config.Translate.TencentTanslateSecretId, Config.Translate.TencentTanslateSecretKey)
	cpf := profile.NewClientProfile()
	client, err := tmt.NewClient(credential, regions.Guangzhou, cpf)
	return client, err
}

func (a *App) TencentTextTranslate(sourceText, sourceLang, targetLang string) FlagResult {
	c, err := NewClient()
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	var projectId int64 = 0
	untranslatedText := ""
	request := tmt.NewTextTranslateRequest()
	request.SourceText = &sourceText
	request.Source = &sourceLang
	request.Target = &targetLang
	request.ProjectId = &projectId
	request.UntranslatedText = &untranslatedText

	response, err := c.TextTranslate(request)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	fmt.Printf("%s", *response.Response.TargetText)
	return FlagResult{true, *response.Response.TargetText}
}
