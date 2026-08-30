# ql-cli

一个用 Go 编写的命令行工具集，基于 [Cobra](https://github.com/spf13/cobra) 框架构建。

## 功能特性

- **rename** — 批量修改文件名，支持查找替换、正则匹配、添加前后缀、大小写转换、清理特殊字符、扩展名过滤等

## 安装

### 从源码构建

```bash
git clone <repository-url>
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

#### 所有参数

| 参数 | 缩写 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--dir` | `-d` | string | `.` | 要处理的目录路径 |
| `--recursive` | `-r` | bool | `false` | 是否递归处理子目录 |
| `--dry-run` | | bool | `false` | 仅预览修改，不实际执行 |
| `--from` | | string | | 要查找的文本 |
| `--to` | | string | | 替换为的文本 |
| `--regex` | | bool | `true` | 是否使用正则表达式匹配 |
| `--prefix` | | string | | 添加文件名前缀 |
| `--suffix` | | string | | 添加文件名后缀（在扩展名前） |
| `--to-lower` | | bool | `false` | 文件名转小写 |
| `--to-upper` | | bool | `false` | 文件名转大写 |
| `--clean` | | bool | `false` | 清理特殊字符（只保留英文字母、数字、下划线、连字符、中文汉字及括号） |
| `--remove-ext` | | bool | `false` | 移除扩展名 |
| `--filter` | | string | | 仅处理匹配扩展名的文件（如 `.jpg`） |

#### 使用示例

```bash
# 普通文本替换
ql rename -d ./files --from "IMG_" --to "Photo_"

# 正则表达式替换，$1 引用捕获组（单引号防止 shell 转义 $ 符号）
ql rename -d ./files --regex --from 'IMG_(\d+)' --to 'Photo_$1'

# 添加前缀
ql rename -d ./files --prefix "2024_"

# 添加后缀（在扩展名前）
ql rename -d ./files --suffix "_backup"

# 指定JPG文件名转小写
ql rename -d ./files --to-lower --filter ".JPG" 

# 清理特殊字符 + 移除扩展名
ql rename -d ./files --clean --remove-ext

# 递归处理子目录
ql rename -d ./files -r --from "old" --to "new"

# 预览模式，不实际执行
ql rename -d ./files --from "A" --to "B" --dry-run
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