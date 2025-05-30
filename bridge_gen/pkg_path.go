package bridgegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParsePkgPath 返回给定文件路径的包路径，查找最近的go.mod文件并提取模块名称
// 返回:
//   - module: 来自go.mod的Go模块名称
//   - pkgPath: 完整的包导入路径
//   - err: 过程中遇到的任何错误
func ParsePkgPath(path string) (module, pkgPath string, err error) {
	// 转换为绝对路径
	path, err = filepath.Abs(path)
	if err != nil {
		err = fmt.Errorf("failed to get absolute path: %w", err)
		return
	}

	// 检查路径是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("failed to stat path: %w", err)
		return
	}

	// 确定要搜索的目录
	var dirToSearch string
	if fileInfo.IsDir() {
		dirToSearch = path
	} else {
		dirToSearch = filepath.Dir(path)
	}

	// 在当前目录和父目录中搜索go.mod文件
	currentDir := dirToSearch
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, statErr := os.Stat(goModPath); statErr == nil {
			// 找到go.mod，读取其内容
			data, readErr := os.ReadFile(goModPath)
			if readErr != nil {
				err = fmt.Errorf("failed to read go.mod file: %w", readErr)
				return
			}

			// 提取模块名称
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if strings.HasPrefix(trimmedLine, "module ") {
					module = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "module "))

					// 计算从模块根目录到目标目录的相对路径
					relPath, relErr := filepath.Rel(currentDir, dirToSearch)
					if relErr == nil && relPath != "." {
						pkgPath = filepath.ToSlash(filepath.Join(module, relPath))
						return
					}

					// 如果目标就是模块根目录本身
					pkgPath = module
					return
				}
			}

			err = fmt.Errorf("module declaration not found in go.mod")
			return
		}

		// 移动到父目录
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// 到达文件系统根目录仍未找到go.mod
			err = fmt.Errorf("no go.mod file found in directory hierarchy")
			break
		}
		currentDir = parentDir
	}

	return
}
