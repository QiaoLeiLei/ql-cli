package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func MoveFile(src, dst string) error {
	err := checkDirExist(src, dst)
	if err != nil {
		return err
	}

	// 先尝试快速rename（同分区优先）
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// rename失败，跨分区，走复制+删除
	err = CopyFile(src, dst)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	// 复制完成，删除源文件
	return os.Remove(src)
}

func CopyFile(src, dst string) error {
	err := checkDirExist(src, dst)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	// 刷新写入磁盘
	if err = out.Sync(); err != nil {
		return err
	}

	return nil
}

func checkDirExist(src, dst string) error {
	// 检查源文件是否存在
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("源文件不存在: %s", src)
	}
	// 检查目标目录是否存在
	if _, err := os.Stat(filepath.Dir(dst)); os.IsNotExist(err) {
		// 创建目标目录
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return fmt.Errorf("创建目标目录失败: %w", err)
		}
	}
	return nil
}
