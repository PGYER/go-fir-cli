package notifiers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/PGYER/go-fir-cli/api"
	"github.com/go-resty/resty/v2"
)

type LarkNotifier struct {
	Key         string
	SecretToken string
}

// https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN

func (l *LarkNotifier) BuildAppPubishedMessage(apiAppInfo *api.ApiAppInfo, CustomMsg, DownloadUrl string) string {
	partialJSON := fmt.Sprintf(`
		"msg_type": "post",
		"content": {
			"post": {
				"zh_cn": {
					"title": "%s uploaded",
					"content": [
						[
							{
								"tag": "text",
								"text": "%s(%s) uploaded at %s"
							}
						],
						[
							{
								"tag": "a",
								"text": "Download",
								"href": "%s"
							}
						],
						[
							{
								"tag": "text",
								"text": "%s"
							}
						]
					]
				}
			}
		}`, apiAppInfo.Name, apiAppInfo.Name, apiAppInfo.Type, time.Now(), DownloadUrl, CustomMsg)

	return partialJSON
}

func (l *LarkNotifier) Notify(partialJsonStr string) error {
	var jsonStr string
	if l.SecretToken != "" {
		timestamp := time.Now().Unix()
		signature, err := GenSign(l.SecretToken, timestamp)
		if err != nil {
			return err
		}

		jsonStr = fmt.Sprintf(`{
			"timestamp": %v,
			"sign": "%s",
			%s
		}`, timestamp, signature, partialJsonStr)

	} else {
		jsonStr = fmt.Sprintf(`{
			%s
		} `, partialJsonStr)
	}

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", l.Key)
	resp, err := resty.New().R().SetBody(jsonStr).SetHeader("Content-Type", "application/json").Post(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("请求失败 %s, %s", resp.Status(), string(resp.Body()))
	}
	return nil
}

func GenSign(secret string, timestamp int64) (string, error) {
	//Encode timestamp + key with SHA256, and then with Base64
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
