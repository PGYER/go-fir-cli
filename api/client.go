package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PGYER/go-fir-cli/constants"
	"github.com/go-resty/resty/v2"
)

// Login 用 api token 换取 email, 顺便把 token 缓存到 f.ApiToken.
func (f *FirApi) Login(token string) error {
	resp, err := resty.New().R().
		SetQueryParam("api_token", token).
		SetHeader("Content-Type", "application/json").
		Get(domain + "/user")
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New("登录失败, 请检查 api token 是否正确")
	}

	var userInfo UserInfo
	if err := json.Unmarshal(resp.Body(), &userInfo); err != nil {
		return fmt.Errorf("解析用户信息失败: %w", err)
	}
	f.Email = userInfo.Email
	f.ApiToken = token
	return nil
}

// FetchAppInfo 获取 app 当前的元信息 (id / short / download_domain 等).
// 依赖 UploadPrepare 设置的 appPrepareUploadData.
func (f *FirApi) FetchAppInfo() (*resty.Response, error) {
	if f.appPrepareUploadData == nil {
		return nil, errors.New("UploadPrepare 尚未调用, 无法获取 app id")
	}
	resp, err := resty.New().R().
		SetHeader("User-Agent", constants.USER_AGENT).
		SetHeader("Content-Type", "application/json").
		SetQueryParam("api_token", f.ApiToken).
		Get(domain + "/apps/" + f.appPrepareUploadData.Id)
	if err != nil {
		return nil, fmt.Errorf("获取 app 最新内容失败: %w", err)
	}
	return resp, nil
}

// forcePinRelease 对应 Ruby --force_pin_history: 上传完成后把这个 release 固定在下载页面.
func (f *FirApi) forcePinRelease(appId, releaseId string) error {
	endpoint := fmt.Sprintf("%s/apps/%s/releases/%s/force_set_history", domain, appId, releaseId)
	resp, err := resty.New().R().
		SetHeader("User-Agent", constants.USER_AGENT).
		SetQueryParam("api_token", f.ApiToken).
		Post(endpoint)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return errors.New(resp.Status() + ": " + string(resp.Body()))
	}
	return nil
}
