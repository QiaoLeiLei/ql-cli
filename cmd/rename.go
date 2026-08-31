package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	from      string
	to        string
	useRegex  bool
	prefix    string
	suffix    string
	toLower   bool
	toUpper   bool
	clean     bool
	removeExt bool
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "批量修改文件名",
	Long: `遍历指定目录,批量修改文件名。支持查找替换、正则匹配、添加前后缀、大小写转换等操作。

示例:
  ql-rename-cli rename -d ./files --from "IMG_" --to "Photo_"
  ql-rename-cli rename -d ./files --regex --from "IMG_(\d+)" --to "Photo_$1"
  ql-rename-cli rename -d ./files --prefix "2024_" --dry-run
  ql-rename-cli rename -d ./files --to-lower --filter ".JPG"`,
	RunE: runRename,
}

func init() {
	addCommonDirFlags(renameCmd)
	renameCmd.Flags().StringVarP(&dir, "dir", "d", ".", "要处理的目录路径")
	renameCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "是否递归处理子目录")
	renameCmd.Flags().BoolVar(&dryRun, "dry-run", false, "仅预览修改,不实际执行")
	renameCmd.Flags().StringVar(&from, "from", "", "要查找的文本")
	renameCmd.Flags().StringVar(&to, "to", "", "替换为的文本")
	renameCmd.Flags().BoolVar(&useRegex, "regex", true, "是否使用正则表达式匹配")
	renameCmd.Flags().StringVar(&prefix, "prefix", "", "添加文件名前缀")
	renameCmd.Flags().StringVar(&suffix, "suffix", "", "添加文件名后缀(在扩展名前)")
	renameCmd.Flags().BoolVar(&toLower, "to-lower", false, "文件名转小写")
	renameCmd.Flags().BoolVar(&toUpper, "to-upper", false, "文件名转大写")
	renameCmd.Flags().BoolVar(&clean, "clean", false, "清理特殊字符(只保留英文字母、数字、下划线、连字符和中文汉字以及括号)")
	renameCmd.Flags().BoolVar(&removeExt, "remove-ext", false, "移除扩展名")

	rootCmd.AddCommand(renameCmd)

	_ = renameCmd.RegisterFlagCompletionFunc("dir", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	})

	_ = renameCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".svg", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md", ".csv", ".json", ".xml", ".html", ".css", ".js", ".ts", ".go", ".py", ".java", ".c", ".cpp", ".h", ".rs", ".rb", ".php", ".zip", ".tar", ".gz", ".rar", ".7z", ".mp3", ".mp4", ".mov", ".avi", ".mkv"}
		return extensions, cobra.ShellCompDirectiveNoFileComp
	})
}

// 白名单正则：只保留英文字母、数字、下划线、连字符和中文汉字以及括号
var pattern = regexp.MustCompile(`[^0-9a-zA-Z_\x{4e00}-\x{9fa5}()（）\-]`)

func computeNewName(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	newBase := base

	// 使用白名单清理：移除所有非英文字母、数字、下划线、连字符和中文汉字以及括号的字符
	if clean {
		newBase = pattern.ReplaceAllString(newBase, "")
	}

	if removeExt {
		ext = ""
	}

	if from != "" {
		if useRegex {
			re, err := regexp.Compile(from)
			if err == nil {
				newBase = re.ReplaceAllString(newBase, to)
			}
		} else {
			newBase = strings.ReplaceAll(newBase, from, to)
		}
	}

	if toLower {
		newBase = strings.ToLower(newBase)
	}
	if toUpper {
		newBase = strings.ToUpper(newBase)
	}

	if prefix != "" {
		newBase = prefix + newBase
	}
	if suffix != "" {
		newBase = newBase + suffix
	}

	return newBase + ext
}

func runRename(cmd *cobra.Command, args []string) error {
	// 检查目录
	absDir, err := getAbsDir()
	if err != nil {
		return err
	}

	// 检查正则表达式
	err = checkRegex()
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Println("===== 预览模式(不会实际修改) =====")
	}

	err = traverseDir(absDir)
	if err != nil {
		return err
	}

	return nil
}

func getAbsDir() (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("获取目录绝对路径失败: %w", err)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		return "", fmt.Errorf("目录不存在: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("路径不是目录: %s", absDir)
	}
	return absDir, nil
}

func checkRegex() error {
	if useRegex && from != "" {
		_, err := regexp.Compile(from)
		if err != nil {
			return fmt.Errorf("正则表达式无效: %w", err)
		}
	}
	return nil
}

// 遍历目录，处理所有文件
func traverseDir(absDir string) error {
	count := 0
	err := filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		// 如果访问当前路径时发生错误，直接返回错误终止遍历
		if err != nil {
			return err
		}
		// 如果是目录：
		//   - 非根目录且未启用递归模式，则跳过该目录，不再向下遍历
		//   - 否则继续遍历
		if info.IsDir() {
			if path != absDir && !recursive {
				return filepath.SkipDir
			}
			return nil
		}

		name := info.Name()

		// 如果指定了文件扩展名过滤器，则仅处理匹配的文件
		if filter != nil {
			// 检查文件扩展名是否在过滤器列表中
			for _, ext := range filter {
				// 确保过滤器以 "." 开头，例如 "txt" → ".txt"
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				if strings.EqualFold(filepath.Ext(name), ext) {
					return nil
				}
			}
			return filepath.SkipDir // 如果不在列表中，跳过该目录，不再向下遍历
		}

		newName := computeNewName(name)

		if newName == name {
			return nil
		}

		newPath := filepath.Join(filepath.Dir(path), newName)

		if dryRun {
			fmt.Printf("  %s -> %s\n", path, newPath)
		} else {
			if err := os.Rename(path, newPath); err != nil {
				// 将重命名失败的路径及错误信息输出到标准错误
				_, _ = fmt.Fprintf(os.Stderr, "重命名失败 %s: %v\n", path, err)
				return nil
			}
			fmt.Printf("已重命名: %s -> %s\n", path, newPath)
		}
		count++
		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历目录时出错: %w", err)
	}

	if count == 0 {
		fmt.Println("没有匹配的文件需要修改")
	} else if dryRun {
		fmt.Printf("共 %d 个文件将被修改\n", count)
	} else {
		fmt.Printf("共修改了 %d 个文件\n", count)
	}
	return nil
}
