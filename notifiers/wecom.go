package notifiers

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type WeComNotifier struct {
	Key string
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	PicUrl      string `json:"picurl"`
}

func (w *WeComNotifier) Notify(message string) error {

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
				

	}`)
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", w.Key)
	resp, err := resty.New().R().SetBody(jsonStr).SetHeader("Content-Type", "application/json").Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("请求失败 %s, %s", resp.Status(), string(resp.Body()))
	}

}
