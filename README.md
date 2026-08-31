# ql-cli

一个用 Go 编写的命令行工具集，基于 [Cobra](https://github.com/spf13/cobra) 框架构建。

## 功能特性

- **rename** — 批量修改文件名，支持查找替换、正则匹配、添加前后缀、大小写转换、清理特殊字符、扩展名过滤等

## 安装

### 从源码构建

```bash
git clone https://github.com/QiaoLeiLei/ql-cli.git
cd ql-cli
go build -o ql .
```

### 直接运行

```bash
go run main.go rename --help
```

## 命令列表

### `rename` — 批量重命名

遍历指定目录，批量修改文件名。支持查找替换、正则匹配、添加前后缀、大小写转换等操作。



#### 基本用法

```bash
ql rename -d ./files --from "IMG_" --to "Photo_"
```
### `move` — 批量移动文件

遍历指定目录，批量移动文件指定类型文件到目标目录。如果目标目录不存在，会自动创建。

#### 基本用法
```bash
ql move -d ./files --filter ".jpg" --dst "./backup"
```

### `completion` — 生成 Shell 补全脚本

为不同 Shell 生成命令补全脚本，安装后可通过 Tab 键自动补全子命令和参数。

#### 基本用法

```bash
ql completion [bash|zsh|fish|powershell]
```

#### 安装补全脚本

```bash
# Bash（Linux/macOS 默认）
ql completion bash > /etc/bash_completion.d/ql-cli

# Zsh
ql completion zsh > "${fpath[1]}/_ql"

# Fish
ql completion fish > ~/.config/fish/completions/ql-cli.fish

# PowerShell
ql completion powershell > ql.ps1
```

## 技术栈

- [Go 1.26+](https://go.dev/)
- [Cobra](https://github.com/spf13/cobra) — CLI 框架

## License

MIT
