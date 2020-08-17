## 安装golang插件

快捷键 ctrl+shift+p 或者 F1 打开执行命令：

```
ext install Go
```

或者直接在扩展面板中搜索“go”安装

插件安装完成后，退出vscode，重新启动vscode

## 安装 go tools

快捷键 ctrl+shift+p 或者 F1 打开执行命令：

```
Go Install/Update Tools
```

如果安装出错，很有可能需要梯子，参考“让go get命令使用socks5代理”

## 自定义配置

根据自己实际情况修改，特别是 goroot 和 gopath:

```
//golang
"go.gopath": "/Users/richen/Workspace/go",
"go.goroot": "/usr/local/Cellar/go/1.9.2/libexec",
"go.buildOnSave": "workspace", //在保存代码时自动编译代码
"go.lintOnSave": "workspace", //在保存代码时自动检查代码可以优化的地方，并给出建议
"go.vetOnSave": "workspace", //在保存代码时自动检查潜在的错误
"go.formatOnSave": true, //在保存代码时自动格式化代码
"go.formatTool": "goimports", //使用 goimports 工具进行代码格式化，或者使用 goreturns和gofmt
"go.formatFlags": [],
"go.lintTool": "golint", //使用 golint 库来执行 lint操作，你也可以选择使用gometalinter
"go.lintFlags": [],
"go.vetFlags": []
```
