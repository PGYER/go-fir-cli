package main

import (
	"fmt"

	"github.com/PGYER/go-fir-cli/constants"
	"github.com/PGYER/go-fir-cli/notifiers"
	"github.com/PGYER/go-fir-cli/utils"
	"gopkg.in/urfave/cli.v1"
)

// testWebhookCommand: 用 dingtalk 机器人发一条"测试消息"验证 webhook 是否连通.
func testWebhookCommand() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "测试 webhook",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "token, t"},
			cli.StringFlag{Name: "secret, s"},
		},
		Action: func(c *cli.Context) error {
			token := c.String("token")
			if token == "" {
				token = utils.LoadLocalToken()
			}
			notifier := &notifiers.DingTalkNotifier{
				Key:         token,
				SecretToken: c.String("secret"),
			}
			if err := notifier.Notify("测试消息"); err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}
}

// versionCommand: 打印当前 go-fir-cli 版本号.
func versionCommand() cli.Command {
	return cli.Command{
		Name:      "version",
		ShortName: "v",
		Usage:     "查看 go-fir-cli 版本",
		Action: func(c *cli.Context) error {
			fmt.Println(constants.VERSION)
			return nil
		},
	}
}

// upgradeCommand: 提示用户去 Release 页面拉新版本.
func upgradeCommand() cli.Command {
	return cli.Command{
		Name:  "upgrade",
		Usage: "如何升级 go-fir-cli",
		Action: func(c *cli.Context) error {
			fmt.Println("请访问 https://github.com/PGYER/go-fir-cli/releases 下载对应版本, 并替换原有的 go-fir-cli 文件")
			return nil
		},
	}
}
