// Package main 是 go-fir-cli 的入口.
//
// 各子命令实现按主题拆到独立文件:
//
//	cmd_auth.go    - login / logout / me
//	cmd_info.go    - info
//	cmd_upload.go  - upload + helpers
//	cmd_misc.go    - test / version / upgrade
package main

import (
	"os"

	"github.com/PGYER/go-fir-cli/constants"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-fir-cli"
	app.Usage = "完成 fir.im 的命令行操作"
	app.Version = constants.VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "token, t", Usage: "fir.im 的 api token"},
	}
	app.Commands = []cli.Command{
		loginCommand(),
		logoutCommand(),
		meCommand(),
		infoCommand(),
		testWebhookCommand(),
		uploadCommand(),
		versionCommand(),
		upgradeCommand(),
	}
	_ = app.Run(os.Args)
}
