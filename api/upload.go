package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/PGYER/go-fir-cli/analysis"
	"github.com/PGYER/go-fir-cli/constants"
	"github.com/cheggaaa/pb"
	"github.com/go-resty/resty/v2"
)

// ensureAppFileInfo lazily 解析一次 ipa/apk, 之后所有上传环节复用 cache.
// 修复了重构前 UploadPrepare 调用两次 NewUploadAppService 导致 ipa 被解压两遍的 bug.
func (f *FirApi) ensureAppFileInfo(file string) error {
	if f.appFileInfo != nil {
		return nil
	}
	svc, err := analysis.NewUploadAppService(file)
	if err != nil {
		return fmt.Errorf("解析 %s 失败: %w", file, err)
	}
	info, err := svc.GetAppFileInfo()
	if err != nil {
		return fmt.Errorf("读取 %s 元信息失败: %w", file, err)
	}
	f.uploadAppService = svc
	f.appFileInfo = info
	return nil
}

// UploadPrepare 向 fir.im 申请本次上传的服务器凭证 (icon / binary 的 CDN url + token).
func (f *FirApi) UploadPrepare(file string) (AppPrepareUploadData, error) {
	var apiUploadJson AppPrepareUploadData

	if err := f.ensureAppFileInfo(file); err != nil {
		return apiUploadJson, err
	}

	firUploadPrepare := &FirUploadPrepare{
		ApiToken:       f.ApiToken,
		Type:           f.appFileInfo.Type,
		BundleId:       f.appFileInfo.BundleId,
		Fname:          path.Base(file),
		ForceUpload:    "ali",
		SkipIconUpload: f.SkipUpdateIcon,
		ManualCallback: true,
		Protocol:       "https",
	}
	jsonBytes, _ := json.Marshal(firUploadPrepare)

	resp, err := resty.New().R().
		SetBody(jsonBytes).
		SetHeader("User-Agent", constants.USER_AGENT).
		SetHeader("Content-Type", "application/json").
		Post(domain + "/apps")
	if err != nil {
		return apiUploadJson, err
	}
	if resp.StatusCode() >= 400 {
		return apiUploadJson, fmt.Errorf("请求失败 %s, %s", resp.Status(), string(resp.Body()))
	}

	if err := json.Unmarshal(resp.Body(), &apiUploadJson); err != nil {
		return apiUploadJson, fmt.Errorf("解析 prepare 响应失败: %w", err)
	}
	f.appPrepareUploadData = &apiUploadJson
	return apiUploadJson, nil
}

// Upload 是发布的入口: 登录 → prepare → 上传图标 → 上传 binary → 回调 → 拉最新 app info → (可选) pin release.
// 失败时通过 error 返回, 调用方 (main.go) 负责 os.Exit, 这样 api package 是可复用 / 可测的.
func (f *FirApi) Upload(file string) error {
	if err := f.Login(f.ApiToken); err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	uploadingInfo, err := f.UploadPrepare(file)
	if err != nil {
		return fmt.Errorf("获取上传信息失败: %w", err)
	}

	if f.SkipUpdateIcon {
		fmt.Println("跳过图标上传 (--skip-update-icon)")
	} else {
		if _, err := f.uploadAppIcon(uploadingInfo); err != nil {
			// 图标上传失败不阻塞流程, 保留旧行为 (Ruby 也是 best-effort).
			fmt.Println("图标上传出错, 继续上传 app 文件:", err)
		} else {
			fmt.Println("图标上传完毕, 开始上传 app 文件...")
		}
	}

	if _, err := f.uploadAppFile(uploadingInfo, file); err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}
	fmt.Println("上传成功")

	callbackResp, err := f.manualCallback(file, f.appFileInfo, uploadingInfo)
	if err != nil {
		return fmt.Errorf("上传成功但回调失败: %w", err)
	}
	var manualCallbackResp ManualCallbackResp
	if err := json.Unmarshal(callbackResp.Body(), &manualCallbackResp); err != nil {
		return fmt.Errorf("解析回调响应失败: %w", err)
	}
	f.manualCallbackResp = &manualCallbackResp

	fmt.Println("文件上传完毕, 开始获取 app 最新数据")
	appResp, err := f.FetchAppInfo()
	if err != nil {
		return err
	}

	var apiAppInfo ApiAppInfo
	if err := json.Unmarshal(appResp.Body(), &apiAppInfo); err != nil {
		return fmt.Errorf("解析 app info 失败: %w", err)
	}
	f.ApiAppInfo = &apiAppInfo

	if f.ForcePinHistory && manualCallbackResp.ReleaseId != "" {
		if err := f.forcePinRelease(apiAppInfo.Id, manualCallbackResp.ReleaseId); err != nil {
			fmt.Println("固定 release 到下载页失败:", err)
		} else {
			fmt.Println("已固定 release 到下载页 (--force-pin-history)")
		}
	}

	return nil
}

// uploadAppIcon 上传 app 图标到 CDN, 然后回调 fir.im 登记.
// 没有图标时静默跳过, 返回 (nil, nil).
func (f *FirApi) uploadAppIcon(uploadingInfo AppPrepareUploadData) (*resty.Response, error) {
	if f.appFileInfo.Icon == nil && f.CustomIconPath == "" {
		fmt.Println("没有图标, 跳过保存图标")
		return nil, nil
	}

	iconFile := "blob"
	if f.CustomIconPath != "" {
		iconFile = f.CustomIconPath
	} else {
		if err := f.uploadAppService.SaveImage(iconFile); err != nil {
			return nil, fmt.Errorf("保存图标失败: %w", err)
		}
		defer os.Remove(iconFile)
	}

	uploadFile, err := os.Open(iconFile)
	if err != nil {
		return nil, fmt.Errorf("打开图标文件失败: %w", err)
	}
	defer uploadFile.Close()

	client := resty.New()
	resp, err := client.R().
		SetBody(uploadFile).
		SetHeaders(uploadingInfo.Cert.Icon.CustomHeaders).
		Put(uploadingInfo.Cert.Icon.UploadUrl)
	if err != nil {
		return resp, fmt.Errorf("上传图标到 CDN 失败: %w", err)
	}

	iconStat, _ := os.Stat(iconFile)
	iconCallback := IconCallback{
		Key:      uploadingInfo.Cert.Icon.Key,
		Token:    uploadingInfo.Cert.Icon.Token,
		Origin:   "go-fir-cli",
		ParentId: uploadingInfo.Id,
		Fsize:    int(iconStat.Size()),
		Fname:    "blob",
	}
	body, _ := json.Marshal(iconCallback)
	return client.R().
		SetBody(body).
		SetHeader("User-Agent", constants.USER_AGENT).
		SetHeader("Content-Type", "application/json").
		Post(domain + "/auth/ali/callback")
}

// uploadAppFile PUT 二进制到 CDN, 带进度条; 支持 --user-download-file-name 覆盖 Content-Disposition.
func (f *FirApi) uploadAppFile(uploadingInfo AppPrepareUploadData, file string) (*http.Response, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("查看文件失败: %w", err)
	}
	fmt.Println("文件大小:", stat.Size())

	bar := pb.New64(stat.Size())
	bar.Start()
	defer bar.Finish()

	uploadFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer uploadFile.Close()

	req, err := http.NewRequest("PUT", uploadingInfo.Cert.Binary.UploadUrl, io.TeeReader(uploadFile, bar))
	if err != nil {
		return nil, err
	}
	for k, v := range uploadingInfo.Cert.Binary.CustomHeaders {
		req.Header.Set(k, v)
	}
	if f.UserDownloadFileName != "" {
		req.Header.Set("CONTENT-DISPOSITION",
			"attachment; filename="+url.QueryEscape(f.UserDownloadFileName))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("PUT 到 CDN 失败: %w", err)
	}
	defer resp.Body.Close()
	return resp, nil
}

// manualCallback 通知 fir.im "上传完了, 这是元信息", 拿到 release id.
func (f *FirApi) manualCallback(file string, appInfo *analysis.AppFileInfo, uploadingInfo AppPrepareUploadData) (*resty.Response, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("查看文件失败: %w", err)
	}

	name := appInfo.Name
	if f.SpecifyAppDisplayName != "" {
		name = f.SpecifyAppDisplayName
	}

	callbackData := CallbackData{
		Build:       appInfo.Build,
		Fsize:       int(stat.Size()),
		Fname:       path.Base(file),
		ReleaseTag:  "develop",
		Key:         uploadingInfo.Cert.Binary.Key,
		Name:        name,
		Origin:      "go-fir-cli",
		ParentId:    uploadingInfo.Id,
		ReleaseType: appInfo.ReleaseType,
		Token:       uploadingInfo.Cert.Binary.Token,
		Version:     appInfo.Version,
		Changelog:   f.AppChangelog,
		UserId:      uploadingInfo.AppUserId,
	}
	body, _ := json.Marshal(callbackData)

	resp, err := resty.New().R().
		SetBody(body).
		SetHeader("User-Agent", constants.USER_AGENT).
		SetHeader("Content-Type", "application/json").
		Post(domain + "/auth/ali/callback")
	if err != nil {
		return resp, err
	}
	if resp.StatusCode() >= 400 {
		return resp, errors.New(resp.Status() + ": " + string(resp.Body()))
	}
	return resp, nil
}
