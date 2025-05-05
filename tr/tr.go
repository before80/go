package tr

import (
	"fmt"
	"os"
	"path/filepath"
)

func TruncFileContent(filePath string) (err error) {
	var f *os.File
	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("创建目录时出错: %w", err)
	}
	// 以写入模式打开文件，若文件不存在则创建，若存在则清空内容
	f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("清空%s文件内容出现错误：%v", filePath, err)
	}
	// 确保文件在函数结束时关闭
	defer func() {
		_ = f.Close()
	}()
	return nil
}
