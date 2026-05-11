package api

import (
	"github.com/PGYER/go-fir-cli/analysis"
)

// domain 是 fir.im 的 API 地址 (历史上叫 appmeta.cn).
const domain = "https://api.appmeta.cn"

// UserInfo 是 /user 接口返回的用户信息.
type UserInfo struct {
	Email string `json:"email"`
}

// FirApi 是 fir.im API 客户端。一次 Upload 流程会复用同一个实例,
// 内部 cache 解析出来的 ipa/apk 元信息和 prepare 接口返回的 token.
type FirApi struct {
	ApiToken     string
	AppChangelog string
	Email        string

	// 图标 / 文件名相关
	CustomIconPath       string
	UserDownloadFileName string // 对应 Ruby --user_download_file_name, 写入 CDN Content-Disposition
	SkipUpdateIcon       bool   // 对应 Ruby --skip_update_icon

	// 发布行为
	SpecifyAppDisplayName string // 对应 Ruby --specify_app_display_name, 覆盖 callback name
	ForcePinHistory       bool   // 对应 Ruby --force_pin_history, 上传完后钉住 release

	// 二维码与本地存档
	QrCodePngNeed    bool
	QrCodeAsciiNeed  bool
	SaveUploadedInfo bool
	SaveUploadedPath string

	// 调用过程产物 (外部只读 ApiAppInfo)
	ApiAppInfo *ApiAppInfo

	uploadAppService     *analysis.UploadAppService
	appFileInfo          *analysis.AppFileInfo
	appPrepareUploadData *AppPrepareUploadData
	manualCallbackResp   *ManualCallbackResp
}

// ApiAppInfo 是 GET /apps/:id 返回的 app 视图.
type ApiAppInfo struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Short           string `json:"short"`
	BundleId        string `json:"bundle_id"`
	DownloadDomain  string `json:"download_domain"`
	MasterReleaseId string `json:"master_release_id"`
	DownloadUrl     string
}

// IconCallback 是图标上传后回调 fir.im 时的 body.
type IconCallback struct {
	Key      string `json:"key"`
	Token    string `json:"token"`
	Origin   string `json:"origin"`
	ParentId string `json:"parent_id"`
	Fsize    int    `json:"fsize"`
	Fname    string `json:"fname"`
}

// CallbackData 是二进制上传后回调 fir.im 时的 body.
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

// FirUploadPrepare 是 POST /apps 的 body, 向服务器申请上传凭证.
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

// UploadFile 是单个文件的 CDN 上传凭证 (icon / binary).
type UploadFile struct {
	Key                string            `json:"key"`
	Token              string            `json:"token"`
	UploadUrl          string            `json:"upload_url"`
	CustomHeaders      map[string]string `json:"custom_headers"`
	CustomCallbackData map[string]string `json:"custom_callback_data"`
}

// UploadCert 是上传凭证集合.
type UploadCert struct {
	Icon   UploadFile `json:"icon"`
	Binary UploadFile `json:"binary"`
}

// ManualCallbackResp 是回调成功后的返回, 含 release id.
type ManualCallbackResp struct {
	ReleaseId string `json:"release_id"`
}

// AppPrepareUploadData 是 POST /apps 的响应体.
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
