package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"betaqr.com/go_fir_cli/analysis"
	"betaqr.com/go_fir_cli/constants"
	"github.com/go-resty/resty/v2"
)

const domain = "https://api.appmeta.cn"

type UserInfo struct {
	Email string `json:"email"`
}

type FirApi struct {
	ApiToken             string
	AppChangelog         string
	Email                string
	ApiAppInfo           *ApiAppInfo
	uploadAppService     *analysis.UploadAppService
	appFileInfo          *analysis.AppFileInfo
	appPrepareUploadData *AppPrepareUploadData
	manualCallbackResp   *ManualCallbackResp
}

type ApiAppInfo struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Short           string `json:"short"`
	BundleId        string `json:"bundle_id"`
	DownloadDomain  string `json:"download_domain"`
	MasterReleaseId string `json:"master_release_id"`
}

type IconCallback struct {
	Key      string `json:"key"`
	Token    string `json:"token"`
	Origin   string `json:"origin"`
	ParentId string `json:"parent_id"`
	Fsize    int    `json:"fsize"`
	Fname    string `json:"fname"`
}

type CallbackData struct {
	Build            string `json:"build"`
	Fname            string `json:"fname"`
	Key              string `json:"key"`
	Name             string `json:"name"`
	Origin           string `json:"origin"`
	ParentId         string `json:"parent_id"`
	ReleaseTag       string `json:"release_tag"`
	Fsize            int    `json:"fsize"`
	ReleaseType      string `json:"release_type"`
	DistributionName string `json:"distribution_name"`
	Token            string `json:"token"`
	Version          string `json:"version"`
	Changelog        string `json:"changelog"`
	UserId           string `json:"user_id"`
}

type FirUploadPrepare struct {
	ApiToken       string `json:"api_token"`
	Type           string `json:"type"`
	BundleId       string `json:"bundle_id"`
	Fname          string `json:"fname"`
	SkipIconUpload bool   `json:"skip_icon_upload"`
	ManualCallback bool   `json:"manual_callback"`
	Protocol       string `json:"protocol"`
	ForceUpload    string `json:"force_upload"`
}

type UploadFile struct {
	Key                string            `json:"key"`
	Token              string            `json:"token"`
	UploadUrl          string            `json:"upload_url"`
	CustomHeaders      map[string]string `json:"custom_headers"`
	CustomCallbackData map[string]string `json:"custom_callback_data"`
}

type UploadCert struct {
	Icon   UploadFile `json:"icon"`
	Binary UploadFile `json:"binary"`
}

type ManualCallbackResp struct {
	ReleaseId string `json:"release_id"`
}

type AppPrepareUploadData struct {
	UserSystemDefaultDownloadDomain string     `json:"user_system_default_download_domain"`
	Id                              string     `json:"id"`
	Type                            string     `json:"type"`
	Short                           string     `json:"short"`
	DownloadDomain                  string     `json:"download_domain"`
	DownloadDomainHttpsReady        bool       `json:"download_domain_https_ready"`
	AppUserId                       string     `json:"app_user_id"`
	Storage                         string     `json:"storage"`
	FormMethod                      string     `json:"form_method"`
	Cert                            UploadCert `json:"cert"`
}

func (f *FirApi) FetchAppInfo() (*resty.Response, error) {
	url := domain + "/apps/" + f.appPrepareUploadData.Id
	client := resty.New()
	resp, err := client.R().SetHeader("User-Agent", constants.USER_AGENT).SetQueryParam("api_token", f.ApiToken).SetHeader("Content-Type", "application/json").Get(url)

	if err != nil {
		fmt.Println("获取app 最新内容失败")
		return nil, err
	}
	fmt.Println(string(resp.Body()))

	return resp, nil
}

func (f *FirApi) Login(token string) error {
	url := domain + "/user"
	client := resty.New()
	// body := `{"api_token":` + token + `}`

	resp, err := client.R().SetQueryParam("api_token", token).SetHeader("Content-Type", "application/json").Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("登录失败, 请检查api token是否正确")

	}
	var userInfo UserInfo

	json.Unmarshal(resp.Body(), &userInfo)
	f.Email = userInfo.Email
	f.ApiToken = token
	return nil
}

// 获得上传需要的服务器的相关信息
func (f *FirApi) UploadPrepare(file string) (AppPrepareUploadData, error) {
	var err error

	if f.uploadAppService == nil {
		// uploadAppService, err := analysis.NewUploadAppService(file)
		f.uploadAppService, err = analysis.NewUploadAppService(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		f.appFileInfo, err = f.uploadAppService.GetAppFileInfo()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	appInfoService, err := analysis.NewUploadAppService(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	appFileInfo, err := appInfoService.GetAppFileInfo()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	firUploadPrepare := &FirUploadPrepare{
		ApiToken:       f.ApiToken,
		Type:           appFileInfo.Type,
		BundleId:       appFileInfo.BundleId,
		Fname:          path.Base(file),
		ForceUpload:    "ali",
		SkipIconUpload: false,
		ManualCallback: true,
		Protocol:       "https",
	}
	jsonBytes, _ := json.Marshal(firUploadPrepare)

	resp, err := resty.New().R().SetBody(jsonBytes).SetHeader("User-Agent", constants.USER_AGENT).SetHeader("Content-Type", "application/json").Post(domain + "/apps")
	var apiUploadJson AppPrepareUploadData

	if err != nil {
		return apiUploadJson, err
	}
	if resp.StatusCode() >= 400 {
		return apiUploadJson, errors.New("请求失败 " + resp.Status() + "," + string(resp.Body()))
	}

	err = json.Unmarshal(resp.Body(), &apiUploadJson)
	if err == nil {
		f.appPrepareUploadData = &apiUploadJson
	}

	return apiUploadJson, err

}

func (f *FirApi) Upload(file string) error {
	err := f.Login(f.ApiToken)
	if err != nil {
		fmt.Println("登录失败", err.Error())
		os.Exit(1)
	}

	// 获得上传需要的数据
	uploadingInfo, err := f.UploadPrepare(file)
	if err != nil {
		fmt.Println("获取上传信息失败", err.Error())
		os.Exit(1)
	}

	// 开始上传

	f.uploadAppIcon(file, uploadingInfo)
	fmt.Println("图标上传完毕, 开始上传App文件...")

	// 1.1 上传 具体文件
	resp, e := f.uploadAppFile(uploadingInfo, file)
	if e != nil {
		fmt.Println("上传失败", e.Error())
		os.Exit(1)
	}

	fmt.Println(resp)

	// 进行回调
	resp, e = manualCallback(file, f.appFileInfo, uploadingInfo)
	if e != nil {
		fmt.Println("上传成功, 但是回调失败", e.Error())
		os.Exit(1)
	}
	var manualCallbackResp ManualCallbackResp
	json.Unmarshal(resp.Body(), &manualCallbackResp)
	f.manualCallbackResp = &manualCallbackResp

	fmt.Println("文件上传完毕, 开始获取app 最新数据")

	resp, _ = f.FetchAppInfo()

	var apiAppInfo ApiAppInfo
	err = json.Unmarshal(resp.Body(), &apiAppInfo)
	if err != nil {
		fmt.Println("解析失败", resp.Body())
		return err
	}

	f.ApiAppInfo = &apiAppInfo

	return nil
}

func (f *FirApi) uploadAppIcon(file string, uploadingInfo AppPrepareUploadData) (*resty.Response, error) {
	if f.appFileInfo.Icon == nil {
		fmt.Println("没有图标, 跳过保存图标")
		return nil, nil
	}
	appUploadConfigInfo := uploadingInfo.Cert.Icon
	upload_url := appUploadConfigInfo.UploadUrl

	client := resty.New()

	e := f.uploadAppService.SaveImage("blob")
	defer os.Remove("blob")

	if e != nil {
		fmt.Println("保存图片失败, 跳过保存图标", e.Error())
		return nil, e
	}
	iconFile := "blob"

	uploadFile, _ := os.Open(iconFile)
	defer uploadFile.Close()

	headers := uploadingInfo.Cert.Icon.CustomHeaders

	fmt.Println(headers)
	resp, e := client.R().SetBody(uploadFile).SetHeaders(uploadingInfo.Cert.Icon.CustomHeaders).Put(upload_url)

	if e != nil {
		fmt.Println("上传图片失败")
		return resp, e
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

	str, _ := json.Marshal(iconCallback)
	url := domain + "/auth/ali/callback"

	resp, e = client.R().SetBody(str).SetHeader("User-Agent", constants.USER_AGENT).SetHeader("Content-Type", "application/json").Post(url)

	return resp, e
}

func (f *FirApi) uploadAppFile(uploadingInfo AppPrepareUploadData, file string) (*resty.Response, error) {
	appUploadConfigInfo := uploadingInfo.Cert.Binary
	upload_url := appUploadConfigInfo.UploadUrl

	info, _ := os.Stat(file)

	fileSize := info.Size()
	fmt.Println("文件大小: ", fileSize)
	uploadFile, _ := os.Open(file)
	defer uploadFile.Close()

	// uploadFile.on

	// var percent chan int64 = make(chan int64, 1)

	client := resty.New()
	// 上传时候显示百分比

	resp, e := client.R().SetBody(uploadFile).SetHeaders(uploadingInfo.Cert.Binary.CustomHeaders).Put(upload_url)

	return resp, e
}

func manualCallback(file string, appInfo *analysis.AppFileInfo, uploadingInfo AppPrepareUploadData) (*resty.Response, error) {
	// manual callback
	client := resty.New()

	url := domain + "/auth/ali/callback"

	fi, _ := os.Stat(file)

	callbackData := CallbackData{
		Build:       appInfo.Build,
		Fsize:       int(fi.Size()),
		Fname:       path.Base(file),
		ReleaseTag:  "develop",
		Key:         uploadingInfo.Cert.Binary.Key,
		Name:        appInfo.Name,
		Origin:      "go-fir-cli",
		ParentId:    uploadingInfo.Id,
		ReleaseType: appInfo.ReleaseType,
		Token:       uploadingInfo.Cert.Binary.Token,
		Version:     appInfo.Version,
		Changelog:   appInfo.Changelog,
		UserId:      uploadingInfo.AppUserId,
	}

	jsonStr, _ := json.Marshal(callbackData)

	resp, e := client.R().SetBody(jsonStr).SetHeader("User-Agent", constants.USER_AGENT).SetHeader("Content-Type", "application/json").Post(url)
	fmt.Println(resp)

	return resp, e

}
