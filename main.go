package main

import (
	"fmt"
	"os"

	"betaqr.com/go_fir_cli/api"
	"betaqr.com/go_fir_cli/constants"
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
			Usage:     "查看版本",
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
				Usage: "app 的更新日志",
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

			api := api.FirApi{
				ApiToken:     token,
				AppChangelog: changelog,
			}

			api.Upload(file)
			fmt.Println("上传成功")
			fmt.Printf("下载页面: http://%s/%s\nReleaseID: %s\n", api.ApiAppInfo.DownloadDomain, api.ApiAppInfo.Short, api.ApiAppInfo.MasterReleaseId)

			return nil
		},
	}
}
