package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"betaqr.com/go_fir_cli/api"
	"betaqr.com/go_fir_cli/constants"
	"github.com/skip2/go-qrcode"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	initCli()
}
func initCli() {
	app := cli.NewApp()
	app.Name = "go_fir_cli"

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

		uploadFile(),
		cli.Command{
			Name:      "version",
			ShortName: "v",
			Usage:     "查看 go_fir_cli 版本",
			Action: func(c *cli.Context) error {
				fmt.Println(constants.VERSION)
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
			fir_api.Login(api_token)
			fmt.Println(fir_api.Email)

			return nil
		},
	}
}

func uploadFile() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "上传文件, 例如 go_fir_cli -t FIR_TOKEN upload -f FILE_PATH -c CHANGELOG",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "apk 或者 ipa 的文件路径",
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
		},
		Action: func(c *cli.Context) error {

			file := c.String("file")
			token := c.GlobalString("token")

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
				ApiToken:        token,
				AppChangelog:    changelog,
				QrCodePngNeed:   c.Bool("qrcode"),
				QrCodeAsciiNeed: c.Bool("qrcodeascii"),
			}

			api.Upload(file)
			fmt.Println("上传成功")
			url := buildDownloadUrl(api.ApiAppInfo, c.Bool("specific_release"))
			fmt.Printf("下载页面: %s\nReleaseID: %s\n", url, api.ApiAppInfo.MasterReleaseId)

			if api.QrCodePngNeed {
				fmt.Println("二维码文件: qrcode.png")
				qrcode.WriteFile(url, qrcode.Medium, 256, "qr.png")
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
