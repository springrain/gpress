package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// unzip 用于解压ZIP文件到指定目录
func unzip(zipPath, destDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建解压目录，如果不存在
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, 0644); err != nil {
			return err
		}
	}

	// 遍历ZIP文件中的所有文件
	for _, f := range r.File {
		// 构造文件路径
		filePath := filepath.Join(destDir, f.Name)

		// 如果文件是目录，则创建目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, 0644)
			continue
		}

		// 打开ZIP文件中的文件
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 创建解压后的文件,不可执行
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// 将文件内容写入到解压后的文件中
		if _, err := io.Copy(outFile, rc); err != nil {
			return err
		}
	}

	return nil
}
