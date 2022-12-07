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

	app.Commands = []cli.Command{
		initLogin(),
		readApkPackage(),
		readIpaPackage(),
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

			answer, _ := analysis.Apk(file)

			fmt.Println(answer)
			analysis.ApkIcon(file)

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
