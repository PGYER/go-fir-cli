package notifiers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"betaqr.com/go_fir_cli/api"
	"github.com/go-resty/resty/v2"
)

// https://open.dingtalk.com/document/group/custom-robot-access

type DingTalkNotifier struct {
	Key         string
	SecretToken string
}

func (d *DingTalkNotifier) BuildAppPubishedMessage(apiAppInfo *api.ApiAppInfo, CustomMsg, DownloadUrl string) string {
	jsonStr := fmt.Sprintf(`{
		"msgtype": 'markdown',
		"markdown": {
		  "title": "%s uploaded",
		  "text": "#### %s(%s)\n\n>uploaded at #{%s}\n\nurl: %s\n\n%s\n\n ![](https://api.appmeta.cn/welcome/qrcode?url=%s)"
		}
	  }`, apiAppInfo.Name, apiAppInfo.Name, apiAppInfo.Type, time.Now(), DownloadUrl, CustomMsg, url.PathEscape(DownloadUrl))

	return jsonStr
}

func (d *DingTalkNotifier) Notify(jsonStr string) error {

	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", d.Key)

	if d.SecretToken != "" {
		timestamp := time.Now().UnixMilli()
		signature, err := sign(d.SecretToken, timestamp)
		if err != nil {
			return err
		}
		url = fmt.Sprintf("%s&timestamp=%d&sign=%s", url, timestamp, signature)
	}

	resp, err := resty.New().R().SetBody(jsonStr).SetHeader("Content-Type", "application/json").Post(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("请求失败 %s, %s", resp.Status(), string(resp.Body()))
	}
	return nil

}

func sign(secret string, t int64) (string, error) {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data), nil
}
