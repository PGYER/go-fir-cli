// Package notifiers 实现 fir.im 发布完成后的第三方通知 (钉钉/飞书/企业微信).
//
// 每个具体 notifier (DingTalkNotifier / LarkNotifier / WeComNotifier)
// 都实现 Notifier interface, 上层在 upload 流程里以统一方式触发,
// 避免 main.go 里写一连串 if c.String("dingtalkToken") != "" 这样的判断.
package notifiers

import "github.com/PGYER/go-fir-cli/api"

// Notifier 是单个通知渠道的统一抽象.
//
// 用法:
//
//	for _, n := range notifiers {
//	    n.Notify(n.BuildAppPubishedMessage(apiAppInfo, downloadUrl))
//	}
//
// 每个实现都有自己的 CustomMsg 字段, 在构造时一次性设入.
type Notifier interface {
	// BuildAppPubishedMessage 根据 app 信息构造该渠道接受的 message payload (JSON 字符串).
	BuildAppPubishedMessage(apiAppInfo *api.ApiAppInfo, downloadUrl string) string
	// Notify 把 payload 推送给具体渠道, 网络层失败返回 error.
	Notify(message string) error
}
