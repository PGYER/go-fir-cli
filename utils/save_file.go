package utils

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/PGYER/go-fir-cli/api"
)

func SaveAnswer(apiAppInfo *api.ApiAppInfo) error {
	// fileName := "fir-cli-answer.json"

	data := make(map[string]string)
	now := time.Now()

	data["short"] = apiAppInfo.Short
	data["name"] = apiAppInfo.Name
	data["download_url"] = apiAppInfo.DownloadUrl
	data["app_id"] = apiAppInfo.Id
	data["release_id"] = apiAppInfo.MasterReleaseId
	data["time"] = now.Format("2006-01-02 15:03:04")

	json, _ := json.Marshal(data)
	// 将 data 转化为 json 后存储到本地answer.json 中
	ioutil.WriteFile("go-fir-cli-answer.json", json, 0644)

	return nil
}
