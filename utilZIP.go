package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// unzip 用于解压ZIP文件到指定目录
func unzip(zipPath, destDir string) (err error) {
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

	var (
		once sync.Once
		dir  string
	)
	defer func() {
		if err != nil && dir != "" {
			_ = os.RemoveAll(dir)
		}
	}()
	// 遍历ZIP文件中的所有文件
	for _, f := range r.File {
		// 构造文件路径
		filePath := filepath.Join(destDir, f.Name)

		// 如果文件是目录，则创建目录
		if f.FileInfo().IsDir() {
			once.Do(func() {
				if _, err := os.Stat(filePath); err != nil {
					// 说明是本次解压创建的文件夹
					dir = filePath
				}
			})
			os.MkdirAll(filePath, 0644)
			continue
		}

		// 打开ZIP文件中的文件
		if err := func(f *zip.File) error {
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
			return nil
		}(f); err != nil {
			return err
		}
	}

	return nil
}
