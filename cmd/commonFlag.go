package cmd

import "github.com/spf13/cobra"

// 共享的命令行选项变量
var (
	dir       string
	recursive bool
	dryRun    bool
	filter    []string
	outDir    string
)

// addCommonDirFlags 注册通用的目录相关 flags 到指定命令
func addCommonDirFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "目标/要处理的目录路径")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "是否递归处理子目录")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览模式(不会实际修改)")
	cmd.Flags().StringSliceVarP(&filter, "filter", "f", nil, "筛选指定文件类型,支持多种-f=png -f=txt，也支持逗号--f=png,txt")
}

func addOutFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outDir, "out", "o", ".", "输出目录")
}
