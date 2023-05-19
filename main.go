package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/PGYER/go-fir-cli/api"
	"github.com/PGYER/go-fir-cli/constants"
	"github.com/PGYER/go-fir-cli/notifiers"
	"github.com/PGYER/go-fir-cli/utils"
	"github.com/skip2/go-qrcode"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	initCli()
}
func initCli() {
	app := cli.NewApp()
	app.Name = "go-fir-cli"

	app.Usage = "完成 fir.im 的命令行操作"
	app.Version = constants.VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token, t",
			Usage: "fir.im 的 api token",
		},
	}

	app.Commands = []cli.Command{
		initLogin(),
		logoutCommand(),
		testWebhook(),
		uploadFile(),
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "查看 go-fir-cli 版本",
			Action: func(c *cli.Context) error {
				fmt.Println(constants.VERSION)
				return nil
			},
		},
		{
			Name:  "upgrade",
			Usage: "如何升级 go-fir-cli",
			Action: func(c *cli.Context) error {
				fmt.Println("请访问 https://github.com/PGYER/go-fir-cli/releases 下载对应版本, 并替换原有的 go-fir-cli 文件")
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func initLogin() cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "登录 fir.im",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "token, t",

				Usage: "fir.im 的 api token",
			},
		},

		Action: func(c *cli.Context) error {
			api_token := c.String("token")
			fir_api := &api.FirApi{}
			if err := fir_api.Login(api_token); err != nil {
				fmt.Println("登录失败, 请检查 token 是否正确")
			} else {
				fmt.Println(fir_api.Email + " 登录成功")
				return utils.SaveToLocal(fir_api.Email, api_token)
			}
			return nil
		},
	}
}

func logoutCommand() cli.Command {
	return cli.Command{
		Name:  "logout",
		Usage: "退出登录",
		Action: func(c *cli.Context) error {
			utils.DelConfig()
			fmt.Println("已经退出登录")
			return nil
		},
	}
}

func testWebhook() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "测试 webhook",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "token, t",
			},
			cli.StringFlag{
				Name: "secret, s",
			},
		},
		Action: func(c *cli.Context) error {
			token := c.String("token")
			secret := c.String("secret")

			if token == "" {
				token = utils.LoadLocalToken()
			}

			notifier := &notifiers.DingTalkNotifier{
				Key:         token,
				SecretToken: secret,
			}
			err := notifier.Notify("测试消息")
			if err != nil {
				fmt.Println(err)
			}

			return nil
		},
	}
}

func uploadFile() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "上传文件, 例如 go-fir-cli -t FIR_TOKEN upload -f FILE_PATH -c CHANGELOG",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "apk 或者 ipa 的文件路径",
			},

			cli.StringFlag{
				Name:  "icon_path, ip",
				Usage: "如果需要自定义icon, 则这里传 icon 的路径",
			},

			cli.StringFlag{
				Name:  "changelog, c",
				Usage: "app 的更新日志, 可以是文件路径, 也可以是字符串",
			},
			cli.BoolFlag{
				Name:  "specific_release, s",
				Usage: "生成的下载地址是否精确指定到 release, 默认为 false",
			},
			cli.BoolFlag{
				Name:  "qrcode, Q",
				Usage: "输出二维码文件 qrcode.png, 用于下载, 默认为 false",
			},

			cli.BoolFlag{
				Name:  "save-uploaded-info, sui",
				Usage: "上传成功后, 保存上传信息到本地json文件（默认为go-fir-cli-answer.json）, 用于用户的其他集成操作, 默认为 false",
			},

			cli.StringFlag{
				Name: "save-uploaded-file, suf",
				Usage: "指定上传成功后, 保存的上传文件路径，默认为当前目录的go-fir-cli-answer.json",
			},

			cli.StringFlag{
				Name:  "dingtalkToken, dt",
				Usage: "dingtalk 的机器人的 token, 用于发送通知",
			},
			cli.StringFlag{
				Name:  "dingtalkSecret, ds",
				Usage: "dingtalk 的机器人的 secret, 用于发送通知时的校验",
			},
			cli.StringFlag{
				Name:  "dingtalkCustomMsg, dcm",
				Usage: "dingtalk 的机器人的自定义消息, 用于发送通知增加关键字",
			},

			cli.StringFlag{
				Name:  "larkToken, lt",
				Usage: "飞书的机器人的 token (url hook 后面那段), 用于发送通知",
			},
			cli.StringFlag{
				Name:  "larkSecret, ls",
				Usage: "飞书的机器人的 secret, 用于发送通知时的校验",
			},
			cli.StringFlag{
				Name:  "larkCustomMsg, lcm",
				Usage: "飞书的机器人的自定义消息, 用于发送通知增加关键字",
			},

			cli.StringFlag{
				Name:  "wecomToken, wt",
				Usage: "企业微信的机器人的 token, 用于发送通知",
			},

			cli.StringFlag{
				Name:  "wecomCustomMsg, wcm",
				Usage: "企业微信的机器人的自定义消息, 用于发送通知增加关键字",
			},
		},
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

			changelog := c.String("changelog")

			// 检测 changelog 文件path是否存在
			if changelog != "" {
				_, err := os.Stat(changelog)
				if err != nil {
					// 文件不存在, 说明changlog 就是 changlog 字符串
				} else {
					//
					str, e := ioutil.ReadFile(changelog)
					if e != nil {

					} else {
						changelog = string(str)
					}
				}
			}

			api := api.FirApi{
				ApiToken:         token,
				CustomIconPath:   c.String("icon_path"),
				AppChangelog:     changelog,
				QrCodePngNeed:    c.Bool("qrcode"),
				QrCodeAsciiNeed:  c.Bool("qrcodeascii"),
				SaveUploadedInfo: c.Bool("save-uploaded-info"),
				SaveUploadedPath: c.String("save-uploaded-file"),
			}

			api.Upload(file)
			fmt.Println("上传成功")
			url := buildDownloadUrl(api.ApiAppInfo, c.Bool("specific_release"))
			api.ApiAppInfo.DownloadUrl = url
			fmt.Printf("下载页面: %s\nReleaseID: %s\n", url, api.ApiAppInfo.MasterReleaseId)

			if api.QrCodePngNeed {
				fmt.Println("二维码文件: qrcode.png")
				qrcode.WriteFile(url, qrcode.Medium, 256, "qr.png")
			}

			if api.SaveUploadedInfo {
				if api.SaveUploadedPath != "" {
					utils.SaveAnswer(api.SaveUploadedPath, api.ApiAppInfo)
				} else {
					utils.SaveAnswer("go-fir-cli-answer.json", api.ApiAppInfo)
				}
			}

			if c.String("dingtalkToken") != "" {
				notifier := &notifiers.DingTalkNotifier{
					Key:         c.String("dingtalkToken"),
					SecretToken: c.String("dingtalkSecret"),
				}

				json := notifier.BuildAppPubishedMessage(api.ApiAppInfo, c.String("dingtalkCustomMsg"), url)
				notifier.Notify(json)
			}

			if c.String("larkToken") != "" {
				notifier := &notifiers.LarkNotifier{
					Key:         c.String("larkToken"),
					SecretToken: c.String("larkSecret"),
				}

				json := notifier.BuildAppPubishedMessage(api.ApiAppInfo, c.String("larkCustomMsg"), url)
				notifier.Notify(json)
			}

			if c.String("wecomToken") != "" {
				notifier := &notifiers.WeComNotifier{
					Key: c.String("wecomToken"),
				}
				json := notifier.BuildAppPubishedMessage(api.ApiAppInfo, c.String("wecomCustomMsg"), url)
				notifier.Notify(json)
			}
			return nil
		},
	}
}

func buildDownloadUrl(apiAppInfo *api.ApiAppInfo, includeRelease bool) string {
	if includeRelease {
		return fmt.Sprintf("http://%s/%s?release_id=%s", apiAppInfo.DownloadDomain, apiAppInfo.Short, apiAppInfo.MasterReleaseId)
	}
	return fmt.Sprintf("http://%s/%s", apiAppInfo.DownloadDomain, apiAppInfo.Short)
}
