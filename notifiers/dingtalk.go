package notifiers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// https://open.dingtalk.com/document/group/custom-robot-access

type DingTalkNotifier struct {
	Key         string
	SecretToken string
}

func (d *DingTalkNotifier) Notify(message string) error {
	var jsonStr string

	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", d.Key)

	if d.SecretToken != "" {
		timestamp := time.Now().UnixMilli()
		signature, err := sign(d.SecretToken, timestamp)
		if err != nil {
			return err
		}
		url = fmt.Sprintf("%s&timestamp=%d&sign=%s", url, timestamp, signature)
	}

	jsonStr = fmt.Sprintf(`{
		"msgtype": "text",
		"text": {
			"content": "%s"
		}
	}`, message)

	resp, err := resty.New().R().SetBody(jsonStr).SetHeader("Content-Type", "application/json").Post(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("请求失败 %s, %s", resp.Status(), string(resp.Body()))
	}

}

func sign(secret string, t int64) (string, error) {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data), nil
}
