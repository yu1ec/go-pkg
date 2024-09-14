package dirx

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateNestedDir 创建嵌套目录
func CreateNestedDir(dirpath string, modePerm os.FileMode) error {
	// 检查路径是否存在
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		// 创建目录
		err := os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

// CreateNestedDirFromFilepath 从文件路径创建嵌套目录
func CreateNestedDirFromFilepath(path string, modePerm os.FileMode) error {
	dirpath := filepath.Dir(path)
	return CreateNestedDir(dirpath, modePerm)
}
