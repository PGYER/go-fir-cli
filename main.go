package main

import (
	"fmt"
	"os"

	"betaqr.com/fir_cli/analysis"
	"betaqr.com/fir_cli/api"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	initCli()
}
func initCli() {
	app := cli.NewApp()
	app.Name = "fir_cli"

	app.Usage = "完成 fir.im 的命令行操作"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token, t",
			Usage: "fir.im 的 api token",
		},
	}

	app.Commands = []cli.Command{
		initLogin(),
		readApkPackage(),
		readIpaPackage(),
		uploadFile(),
	}
	app.Run(os.Args)
}

func initLogin() cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "登录 fir.im",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "token, t",
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
		Usage: "上传文件",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "apk 或者 ipa 的文件路径",
			},
		},
		Action: func(c *cli.Context) error {

			file := c.String("file")
			token := c.GlobalString("token")
			fmt.Println(token)

			api := api.FirApi{
				ApiToken: token,
			}

			api.Upload(file)
			fmt.Println("上传成功")
			fmt.Printf("下载页面: %s/%s  ,ReleaseID=%s\n", api.ApiAppInfo.DownloadDomain, api.ApiAppInfo.Short, api.ApiAppInfo.MasterReleaseId)

			return nil
		},
	}
}

func readApkPackage() cli.Command {
	return cli.Command{
		Name:  "apk",
		Usage: "读取 apk 包信息",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "apk 文件路径",
			},
		},
		Action: func(c *cli.Context) error {

			file := c.String("file")
			token := c.GlobalString("token")
			fmt.Println(token)

			analysis.Apk(file)

			// fmt.Println(answer)

			api := api.FirApi{
				ApiToken: token,
			}

			err := api.Upload(file)
			if err != nil {
				fmt.Println("s上传有错误: ", err.Error())
				os.Exit(1)
			}

			// analysis.ApkIcon(file)

			return nil
		},
	}
}

func readIpaPackage() cli.Command {
	return cli.Command{
		Name:  "ipa",
		Usage: "读取 ipa 包信息",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "ipa 文件路径",
			},
		},
		Action: func(c *cli.Context) error {

			file := c.String("file")

			answer, _ := analysis.Ipa(file)

			fmt.Println(answer)

			return nil
		},
	}
}
