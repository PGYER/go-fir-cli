package notifiers

import (
	"fmt"
	"net/url"
	"time"

	"betaqr.com/go_fir_cli/api"
	"github.com/go-resty/resty/v2"
)

type WeComNotifier struct {
	Key string
}

func (w *WeComNotifier) BuildAppPubishedMessage(apiAppInfo *api.ApiAppInfo, CustomMsg, DownloadUrl string) string {
	jsonStr := fmt.Sprintf(`{
		"msgtype": "news",
		"news": {
			"articles": [
				{
					"title": "%s",
					"description": "%s (%s) uploaded at %s",
					"url": "%s",
					"picurl": "https://api.appmeta.cn/welcome/qrcode?url=%s"
				}]
		}
	}`, apiAppInfo.Name, apiAppInfo.Name, apiAppInfo.Type, time.Now(), DownloadUrl, url.PathEscape(DownloadUrl))
	return jsonStr
}

func (w *WeComNotifier) Notify(jsonStr string) error {

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", w.Key)
	resp, err := resty.New().R().SetBody(jsonStr).SetHeader("Content-Type", "application/json").Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("请求失败 %s, %s", resp.Status(), string(resp.Body()))
	}
	return nil
}
