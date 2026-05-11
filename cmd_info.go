package main

import (
	"fmt"

	"github.com/PGYER/go-fir-cli/analysis"
	"gopkg.in/urfave/cli.v1"
)

// infoCommand: 不上传, 直接解析本地 ipa/apk 输出元信息.
// 用法: go-fir-cli info -f FILE 或 go-fir-cli info FILE
func infoCommand() cli.Command {
	return cli.Command{
		Name:      "info",
		ShortName: "i",
		Usage:     "查看 ipa/apk 应用信息, 用法: go-fir-cli info -f FILE 或 go-fir-cli info FILE",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "file, f", Usage: "apk 或者 ipa 的文件路径"},
		},
		Action: func(c *cli.Context) error {
			file := c.String("file")
			if file == "" && c.NArg() > 0 {
				file = c.Args().First()
			}
			if file == "" {
				fmt.Println("请使用 -f 设置文件路径, 或者 go-fir-cli info <文件路径>")
				return nil
			}

			svc, err := analysis.NewUploadAppService(file)
			if err != nil {
				fmt.Println("解析失败:", err)
				return err
			}
			info := svc.AppFileInfo
			fmt.Printf("type:         %s\n", info.Type)
			fmt.Printf("name:         %s\n", info.Name)
			fmt.Printf("bundle_id:    %s\n", info.BundleId)
			fmt.Printf("version:      %s\n", info.Version)
			fmt.Printf("build:        %s\n", info.Build)
			fmt.Printf("release_type: %s\n", info.ReleaseType)
			fmt.Printf("size:         %d\n", info.Size)
			if len(info.Udids) > 0 {
				fmt.Printf("udids:        %v\n", info.Udids)
			}
			fmt.Printf("file_path:    %s\n", file)
			return nil
		},
	}
}
