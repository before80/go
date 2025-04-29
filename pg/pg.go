package pg

import (
	"bufio"
	"fmt"
	"github.com/before80/go/cfg"
	"os"
	"strings"
)

// InsertAnyPageData 插入页面数据
func InsertAnyPageData(fpDst string) (err error) {
	fpSrc := cfg.Default.UniqueMdFilepath
	var dstFile, srcFile *os.File
	dstFile, err = os.OpenFile(fpDst, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("打开文件 %s 时出错: %v\n", fpDst, err)
	}
	defer dstFile.Close()

	var dstFileSomeLines []string
	foundShouLu := false
	scanner := bufio.NewScanner(dstFile)
	for scanner.Scan() {
		line := scanner.Text()
		dstFileSomeLines = append(dstFileSomeLines, line)
		if strings.HasPrefix(line, "> 收录时间：") {
			foundShouLu = true
			break
		}
	}
	if !foundShouLu {
		return fmt.Errorf("未找到 %q 的起始行", "> 收录时间：")
	}

	// 读取uniqueMd文件中的内容
	srcFile, err = os.OpenFile(fpSrc, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("打开文件 %s 时出错: %v\n", fpSrc, err)
	}
	defer srcFile.Close()

	var srcFileTotalLines []string
	scanner = bufio.NewScanner(srcFile)
	for scanner.Scan() {
		srcFileTotalLines = append(srcFileTotalLines, scanner.Text())
	}

	var newTotalLines []string
	newTotalLines = append(newTotalLines, dstFileSomeLines...)
	newTotalLines = append(newTotalLines, []string{"", ""}...) // 插入两个空行
	newTotalLines = append(newTotalLines, srcFileTotalLines...)

	_ = dstFile.Truncate(0)   // 清空
	_, _ = dstFile.Seek(0, 0) // 从头开始写入
	writer := bufio.NewWriter(dstFile)
	for _, line := range newTotalLines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}

	return nil
}
