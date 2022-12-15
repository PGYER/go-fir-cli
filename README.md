# go-fir-cli

一个基于 go 的 go-fir-cli 程序

## 说明

这是一个 fir-cli [https://github.com/PGYER/fir-cli](https://github.com/PGYER/fir-cli) 的 go 版本, 用于上传文件到 betaqr.com (原fir.im)

## 安装

下载自己对应的系统的文件, 赋值可执行. 若想在全局使用,放到 path 里即可.




## 使用

假设您已经将 go-fir-cli 放到了当前目录

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

```


### 自行编译

下载好代码 安装好依赖即可运行 go build

