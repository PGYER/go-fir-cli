# go-fir-cli

一个基于 go 的 go-fir-cli 程序


## 说明


这是一个 fir-cli [https://github.com/PGYER/fir-cli](https://github.com/PGYER/fir-cli) 的 go 版本, 用于上传文件到 [betaqr.com](https://www.betaqr.com) (原fir.im)

这个版本的主要目的是 当你不想或者不能在机器上安装 ruby 环境时, 无法使用 RUBY 版本的 fir-cli 时, 可以使用这个 go 版本的 fir-cli, 特别是在某些jenkins 上.

go-fir-cli 只实现了 fir-cli 的部分功能, 并无计划实现全部功能, 仅供参考.

由于作者不善 golang, 所以大部分代码皆来自于 Copilot 和 chatGPT 生成, 作者仅在此基础上做了一些修改与调试, 以便于使用. 如果您在使用中发现任何问题, 欢迎提 issue 或者 pr.

## 更新说明
- 0.0.9 upload 增加了参数 --save-uploaded-info, 将结果存入当前目录下的 go-fir-cli-answer.json
- 0.0.8 upload 增加了自定义图标参数 --icon_path
- 0.0.7 支持了持久化登录
- 0.0.6 支持了windows平台
- 0.0.5 支持了上传进度条显示
- 0.0.4 支持了 arm64 的二进制文件, 上传时会有命令行的提示
- 0.0.3 upload 支持了 企业微信、飞书和钉钉的通知机器人
- 0.0.2 upload 支持了生成二维码文件, 支持精确到版本的下载地址
- 0.0.1 支持了上传文件

## 安装



进入[https://github.com/PGYER/go-fir-cli/releases](https://github.com/PGYER/go-fir-cli/releases)下载自己对应的系统的文件, 并给予可执行权限. 若想在全局使用,请将其放到 path 里即可.

- macOS 使用 darwin-amd64
- Linux 使用 amd64
- Windows 使用 windows-386 或 windows-amd64




## 使用

假设您正确下载的您操作系统的 go-fir-cli 到您app 文件的当前目录, 并将其命名为了 go-fir-cli (您也可以放进环境变量里)

### 查看帮助

```bash
./go-fir-cli -h # 查看能使用的命令

# 查看某个命令的帮助, 如
./go-fir-cli login -h # 查看 login 命令的帮助
./go-fir-cli upload -h # 查看 upload 命令的帮助

```

### 检测API 是否可用

```bash
./go-fir-cli login -t 您的API_TOKEN

# 如
# ./go-fir-cli login -t 1234567890abcdefg

```
如果正常, 则返回用户邮件, 如果不正常, 则返回错误信息


### 上传文件

```bash

./go-fir-cli -t 您的API_TOKEN upload -f apk或者ipa文件路径

# 如
# ./go-fir-cli -t 1234567890abcdefg upload -f /Users/xxx/Desktop/xxx.apk
# 具体参数见  ./go-fir-cli upload -h

```




### 自行编译

下载好代码 安装好依赖即可运行 go build
