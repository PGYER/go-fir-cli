package notifiers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/PGYER/go-fir-cli/api"
	"github.com/go-resty/resty/v2"
)

type WeComNotifier struct {
	Key       string
	CustomMsg string // 自定义关键字 / 备注, 拼到消息描述末尾
	PicUrl    string // 对应 Ruby --wxwork_pic_url, 自定义消息封面图, 为空回退到默认二维码
}

func (w *WeComNotifier) BuildAppPubishedMessage(apiAppInfo *api.ApiAppInfo, DownloadUrl string) string {
	picUrl := w.PicUrl
	if picUrl == "" {
		picUrl = "https://api.appmeta.cn/welcome/qrcode?url=" + url.PathEscape(DownloadUrl)
	}
	description := fmt.Sprintf("%s (%s) uploaded at %s", apiAppInfo.Name, apiAppInfo.Type, time.Now().Format(time.RFC3339))
	if w.CustomMsg != "" {
		description = description + ". " + w.CustomMsg
	}
	jsonStr := fmt.Sprintf(`{
		"msgtype": "news",
		"news": {
			"articles": [
				{
					"title": "%s",
					"description": "%s",
					"url": "%s",
					"picurl": "%s"
				}]
		}
	}`, apiAppInfo.Name, description, DownloadUrl, picUrl)
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
