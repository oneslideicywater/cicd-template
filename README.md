# cicd-template


## 用法

```bash
NAME:
   cicd-template - ci/cd configuration generator

USAGE:
   cicd-template [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --namespace value, -n value    project namespace
   --maintainers value, -m value  project maintainers
   --path value, -p value         project path,default is current directory
   --help, -h                     show help (default: false)

```

## 下载

`cicd-template` 是一个二进制可执行文件。

- windows下载`cicd-template.exe`
- linux下载`cicd-template`

下载链接:
http://172.16.66.37:9001/buckets/cicd-template/browse

## 例子

- windows

```bash
$> cicd-template.exe -n cicd-dev -m lifengqian -p D:\geoway\workbunch\sample-vue-app\demo-cicd-vuejs
project type: npm
generate success,happy ci/cd :) !
```

- linux

```bash
$ cicd-template -n cicd-dev -m lifengqian -p demo-cicd 
project type: maven 
generate success,happy ci/cd :) !

```

## 开发 cicd-template

### 如何编译二进制

```bash
go-bindata templates/...
# windows build on windows
# Build
# linux binary
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o cicd-template main.go
# windows binary
RUN CGO_ENABLED=0 GOOS=windows GO111MODULE=on GOARCH=amd64 go build -a -o cicd-template.exe main.go
```