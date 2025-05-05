package pg

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/before80/go/cfg"
	"github.com/before80/go/tr"
	"github.com/before80/go/wind"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"strings"
	"time"
)

func JudgeFileExist(mdFilePath string) bool {
	_, err := os.Stat(mdFilePath)

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

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

func DealUniqueMd(browserHwnd win.HWND, curUrl, step string) (err error) {
	uniqueMdFilepath := cfg.Default.UniqueMdFilepath
	// 获取文件名
	spSlice := strings.Split(uniqueMdFilepath, "\\")
	mdFilename := spSlice[len(spSlice)-1]

	// 清空唯一共用的markdown文件的文件内容
	err = tr.TruncFileContent(uniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，清空%q文件内容出现错误：%v", step, curUrl, uniqueMdFilepath, err)
	}

	// 打开 唯一共用的markdown文件
	err = wind.OpenTypora(uniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，打开窗口名为%q时出现错误：%v", step, curUrl, uniqueMdFilepath, err)
	}

	// 适当延时保证能打开 typora
	time.Sleep(time.Duration(cfg.Default.WaitTyporaOpenSeconds) * time.Second)

	var typoraHwnd win.HWND
	typoraWindowName := fmt.Sprintf("%s - Typora", mdFilename)
	typoraHwnd, err = wind.FindWindowHwndByWindowTitle(typoraWindowName)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，找不到%q窗口：%v", step, curUrl, typoraWindowName, err)
	}

	wind.SelectAllAndCtrlC(browserHwnd)
	time.Sleep(200 * time.Microsecond)
	wind.SelectAllAndDelete(typoraHwnd)
	wind.CtrlV(typoraHwnd)
	time.Sleep(time.Duration(cfg.Default.WaitTyporaCopiedToSaveSeconds) * time.Second)
	wind.CtrlS(typoraHwnd)
	time.Sleep(time.Duration(cfg.Default.WaitTyporaSaveSeconds) * time.Second)
	robotgo.CloseWindow()
	time.Sleep(time.Duration(cfg.Default.WaitTyporaCloseSeconds) * time.Second)

	return nil
}
