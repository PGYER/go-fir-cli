package main

import (
	"fmt"

	"github.com/PGYER/go-fir-cli/api"
	"github.com/PGYER/go-fir-cli/utils"
	"gopkg.in/urfave/cli.v1"
)

// loginCommand: go-fir-cli login -t <api_token>
// 用 api token 换 email, 把 (email, token) 写到本地配置.
func loginCommand() cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "登录 fir.im",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "token, t", Usage: "fir.im 的 api token"},
		},
		Action: func(c *cli.Context) error {
			token := c.String("token")
			firApi := &api.FirApi{}
			if err := firApi.Login(token); err != nil {
				fmt.Println("登录失败, 请检查 token 是否正确")
				return nil
			}
			fmt.Println(firApi.Email + " 登录成功")
			return utils.SaveToLocal(firApi.Email, token)
		},
	}
}

// logoutCommand: 清理本地 token 配置.
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

// meCommand: 用本地 token 校验后打印当前登录邮箱.
func meCommand() cli.Command {
	return cli.Command{
		Name:  "me",
		Usage: "查看当前登录的 fir.im 用户",
		Action: func(c *cli.Context) error {
			token := c.GlobalString("token")
			if token == "" {
				token = utils.LoadLocalToken()
			}
			if token == "" {
				fmt.Println("请先登录: go-fir-cli login -t <api_token>")
				return nil
			}
			firApi := &api.FirApi{}
			if err := firApi.Login(token); err != nil {
				fmt.Println("登录失败, 请检查 token 是否正确:", err)
				return err
			}
			fmt.Println("当前登录用户:", firApi.Email)
			return nil
		},
	}
}
