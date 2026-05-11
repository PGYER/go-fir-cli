package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PGYER/go-fir-cli/api"
	"github.com/PGYER/go-fir-cli/notifiers"
	"github.com/PGYER/go-fir-cli/utils"
	"github.com/skip2/go-qrcode"
	"gopkg.in/urfave/cli.v1"
)

// uploadCommand: 发布 ipa/apk 到 fir.im, 主流程见 api.FirApi.Upload.
func uploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "上传文件, 例如 go-fir-cli -t FIR_TOKEN upload -f FILE_PATH -c CHANGELOG",
		Flags: uploadFlags(),
		Action: func(c *cli.Context) error {
			file := c.String("file")
			token := c.GlobalString("token")
			if token == "" {
				token = utils.LoadLocalToken()
			}
			if token == "" {
				fmt.Println("请先设置 token")
				return nil
			}
			if file == "" {
				fmt.Println("请使用 -f 设置文件路径")
				return nil
			}

			changelog := readChangelog(c.String("changelog"))

			firApi := api.FirApi{
				ApiToken:              token,
				CustomIconPath:        c.String("icon_path"),
				AppChangelog:          changelog,
				QrCodePngNeed:         c.Bool("qrcode"),
				QrCodeAsciiNeed:       c.Bool("qrcodeascii"),
				SaveUploadedInfo:      c.Bool("save-uploaded-info"),
				SaveUploadedPath:      c.String("save-uploaded-file"),
				SpecifyAppDisplayName: c.String("specify-app-display-name"),
				SkipUpdateIcon:        c.Bool("skip-update-icon"),
				UserDownloadFileName:  c.String("user-download-file-name"),
				ForcePinHistory:       c.Bool("force-pin-history"),
			}

			if err := firApi.Upload(file); err != nil {
				fmt.Println(err)
				return err
			}

			url := buildDownloadUrl(firApi.ApiAppInfo, c.Bool("specific_release"))
			firApi.ApiAppInfo.DownloadUrl = url
			fmt.Printf("下载页面: %s\nReleaseID: %s\n", url, firApi.ApiAppInfo.MasterReleaseId)

			if firApi.QrCodePngNeed {
				fmt.Println("二维码文件: qrcode.png")
				qrcode.WriteFile(url, qrcode.Medium, 256, "qr.png")
			}

			if firApi.SaveUploadedInfo {
				path := firApi.SaveUploadedPath
				if path == "" {
					path = "go-fir-cli-answer.json"
				}
				utils.SaveAnswer(path, firApi.ApiAppInfo)
			}

			for _, n := range buildNotifiersFromContext(c) {
				if err := n.Notify(n.BuildAppPubishedMessage(firApi.ApiAppInfo, url)); err != nil {
					fmt.Println("通知发送失败:", err)
				}
			}
			return nil
		},
	}
}

// readChangelog 把 --changelog 的取值规范化: 如果是个存在的文件, 读其内容;
// 否则原样作为 changelog 字符串.
func readChangelog(input string) string {
	if input == "" {
		return ""
	}
	if _, err := os.Stat(input); err != nil {
		return input
	}
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return input
	}
	return string(data)
}

// buildNotifiersFromContext 从 cli flag 里挑出用户配置了 token 的通知渠道,
// 返回统一抽象的 Notifier 列表; upload 流程结束后逐个 Notify.
func buildNotifiersFromContext(c *cli.Context) []notifiers.Notifier {
	var ns []notifiers.Notifier

	if c.String("dingtalkToken") != "" {
		var atPhones []string
		if raw := c.String("dingtalkAtPhones"); raw != "" {
			for _, p := range strings.Split(raw, ",") {
				if p = strings.TrimSpace(p); p != "" {
					atPhones = append(atPhones, p)
				}
			}
		}
		ns = append(ns, &notifiers.DingTalkNotifier{
			Key:         c.String("dingtalkToken"),
			SecretToken: c.String("dingtalkSecret"),
			CustomMsg:   c.String("dingtalkCustomMsg"),
			AtPhones:    atPhones,
			IsAtAll:     c.Bool("dingtalkAtAll"),
		})
	}

	if c.String("larkToken") != "" {
		ns = append(ns, &notifiers.LarkNotifier{
			Key:         c.String("larkToken"),
			SecretToken: c.String("larkSecret"),
			CustomMsg:   c.String("larkCustomMsg"),
			CustomTitle: c.String("larkCustomTitle"),
		})
	}

	if c.String("wecomToken") != "" {
		ns = append(ns, &notifiers.WeComNotifier{
			Key:       c.String("wecomToken"),
			CustomMsg: c.String("wecomCustomMsg"),
			PicUrl:    c.String("wecomPicUrl"),
		})
	}

	return ns
}

// buildDownloadUrl 拼出 fir.im 下载页 URL; includeRelease=true 时带 ?release_id=
// 锁定到本次发布的版本, 否则是 app 的常驻短链.
func buildDownloadUrl(apiAppInfo *api.ApiAppInfo, includeRelease bool) string {
	if includeRelease {
		return fmt.Sprintf("http://%s/%s?release_id=%s", apiAppInfo.DownloadDomain, apiAppInfo.Short, apiAppInfo.MasterReleaseId)
	}
	return fmt.Sprintf("http://%s/%s", apiAppInfo.DownloadDomain, apiAppInfo.Short)
}

// uploadFlags 是 upload 子命令的所有 flag, 单独抽出来防止 uploadCommand 太长.
func uploadFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "file, f", Usage: "apk 或者 ipa 的文件路径"},
		cli.StringFlag{Name: "icon_path, ip", Usage: "如果需要自定义 icon, 则这里传 icon 的路径"},
		cli.StringFlag{Name: "changelog, c", Usage: "app 的更新日志, 可以是文件路径, 也可以是字符串"},
		cli.BoolFlag{Name: "specific_release, s", Usage: "生成的下载地址是否精确指定到 release, 默认为 false"},
		cli.BoolFlag{Name: "qrcode, Q", Usage: "输出二维码文件 qrcode.png, 用于下载, 默认为 false"},

		cli.BoolFlag{Name: "save-uploaded-info, sui", Usage: "上传成功后, 保存上传信息到本地 json 文件 (默认 go-fir-cli-answer.json)"},
		cli.StringFlag{Name: "save-uploaded-file, suf", Usage: "指定上传成功后, 保存的上传文件路径, 默认为当前目录的 go-fir-cli-answer.json"},

		// 对齐 Ruby fir-cli publish 同名参数
		cli.StringFlag{Name: "specify-app-display-name, sadn", Usage: "指定 app 在 fir.im 上显示的名称, 覆盖 ipa/apk 内的 display name"},
		cli.BoolFlag{Name: "skip-update-icon, sui-icon", Usage: "跳过 app 图标上传, 保留服务器上已有的图标"},
		cli.StringFlag{Name: "user-download-file-name, udfn", Usage: "自定义下载时的文件名 (设置到 CDN Content-Disposition)"},
		cli.BoolFlag{Name: "force-pin-history, fph", Usage: "上传完成后将本次 release 固定到下载页面 (超出最大固定数会挤掉最老的)"},

		// DingTalk
		cli.StringFlag{Name: "dingtalkToken, dt", Usage: "dingtalk 的机器人的 token, 用于发送通知"},
		cli.StringFlag{Name: "dingtalkSecret, ds", Usage: "dingtalk 的机器人的 secret, 用于发送通知时的校验"},
		cli.StringFlag{Name: "dingtalkCustomMsg, dcm", Usage: "dingtalk 的机器人的自定义消息, 用于发送通知增加关键字"},
		cli.StringFlag{Name: "dingtalkAtPhones, dap", Usage: "dingtalk 通知 @ 的手机号, 多个用逗号分隔"},
		cli.BoolFlag{Name: "dingtalkAtAll, daa", Usage: "dingtalk 通知 @ 所有人"},

		// Lark
		cli.StringFlag{Name: "larkToken, lt", Usage: "飞书的机器人的 token (url hook 后面那段)"},
		cli.StringFlag{Name: "larkSecret, ls", Usage: "飞书的机器人的 secret"},
		cli.StringFlag{Name: "larkCustomMsg, lcm", Usage: "飞书的机器人的自定义消息, 用于发送通知增加关键字"},
		cli.StringFlag{Name: "larkCustomTitle, lct", Usage: "飞书的机器人通知的自定义标题, 默认 \"<name> uploaded\""},

		// WeCom
		cli.StringFlag{Name: "wecomToken, wt", Usage: "企业微信的机器人的 token"},
		cli.StringFlag{Name: "wecomCustomMsg, wcm", Usage: "企业微信的机器人的自定义消息, 用于发送通知增加关键字"},
		cli.StringFlag{Name: "wecomPicUrl, wpu", Usage: "企业微信通知卡片的封面图 URL, 默认是下载页二维码"},
	}
}
