package main

import (
	"os"
	"testing"
)

func Test_unzip(t *testing.T) {
	type args struct {
		zipPath string
		destDir string
	}
	tests := []struct {
		name         string
		args         args
		wantExistDir bool
		wantDir      string
	}{
		{
			name: "解压成功 -> 创建或者覆盖目录和文件成功",
			args: args{
				zipPath: "./gpressdatadir/gpress-geekdoc-master.zip",
				destDir: "gpressdatadir/",
			},
			wantExistDir: true,
			wantDir:      "./gpressdatadir/gpress-geekdoc-master/",
		},
		{
			name: "解压失败-原目录存在 -> 不进行任何处理",
			args: args{
				zipPath: "./gpressdatadir/gpress-geekdoc-master2.zip",
				destDir: "gpressdatadir/",
			},
			wantExistDir: true,
			wantDir:      "./gpressdatadir/gpress-geekdoc-master/",
		},
		{
			name: "解压失败-原目录不存在 -> 删除原目录",
			args: args{
				zipPath: "./gpressdatadir/gpress-geekdoc-master_副本.zip", // 这个副本 把里面的一个.html后缀的改成.zip 然后压缩
				destDir: "tmp/",
			},
			wantExistDir: false,
			wantDir:      "./tmp/gpress-geekdoc-master/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = unzip(tt.args.zipPath, tt.args.destDir)
			if _, err := os.Stat(tt.wantDir); (err == nil) != tt.wantExistDir {
				t.Errorf("unzip() error = %v, wantExistDir %v", err, tt.wantExistDir)
			}
		})
	}
}
