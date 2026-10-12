package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/QiaoLeiLei/ql-cli/util"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "移动文件到指定目录",
	Long:  "移动文件到指定目录，支持筛选指定文件类型、递归，保留文件结构。",
	RunE:  runMove,
}

func init() {
	addCommonDirFlags(moveCmd)
	addOutFlag(moveCmd)
	addKeepDirFlag(moveCmd)
	rootCmd.AddCommand(moveCmd)
}

func runMove(cmd *cobra.Command, args []string) error {
	// 检查目录
	absDir, err := getAbsDir()
	if err != nil {
		return err
	}
	// 检查输出目录
	ok := filepath.IsAbs(outDir)
	if !ok {
		twd, _ := os.Getwd()
		outDir = filepath.Join(twd, outDir)
	}

	if absDir == outDir {
		return fmt.Errorf("目标目录不能与源目录相同")
	}

	if dryRun {
		fmt.Println("===== 预览模式(不会实际移动) =====")
	}

	err = execMove(absDir)
	if err != nil {
		return err
	}

	return nil
}

func execMove(absDir string) error {
	err := filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != absDir && !recursive {
				return filepath.SkipDir
			}
			return nil
		}
		if filter != nil {
			for _, f := range filter {
				if strings.Contains(f, ",") {
					for _, subExt := range strings.Split(f, ",") {
						if !strings.HasPrefix(subExt, ".") {
							subExt = "." + subExt
						}
						if !strings.EqualFold(filepath.Ext(info.Name()), subExt) {
							return nil
						}
					}
				}
				if !strings.HasPrefix(f, ".") {
					f = "." + f
				}
				if !strings.EqualFold(filepath.Ext(info.Name()), f) {
					return nil
				}
			}
		}
		newPath := computeNewPath(absDir, path, info)
		if dryRun {
			fmt.Printf("  %s -> %s\n", path, newPath)
		} else {
			if err := util.MoveFile(path, newPath); err != nil {
				return fmt.Errorf("移动文件失败: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历目录失败: %w", err)
	}
	return nil
}

func computeNewPath(absDir, path string, info os.FileInfo) string {
	if keepDir {
		return filepath.Join(outDir, strings.TrimPrefix(path, absDir))
	}
	return filepath.Join(outDir, info.Name())
}
