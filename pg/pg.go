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
	"path/filepath"
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

// GenBarIndexMdFile 生成栏目菜单中的_index.md文件
func GenBarIndexMdFile(dir, title string, weight int) (err error) {
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("无法创建%s目录：%v\n", dir, err)
	}
	mdFp := filepath.Join(dir, "_index.md")

	_, err = os.Stat(mdFp)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		var mdF *os.File
		mdF, err = os.OpenFile(mdFp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("创建文件 %s 时出错: %w", mdFp, err)
		}
		defer mdF.Close()
		date := time.Now().Format(time.RFC3339)
		_, err = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
linkTitle = "%s"
date = %s
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false
[menu.main]
	weight = %d
+++

`, title, title, date, "", weight))

		return err
	}
	return nil
}

// InsertAnyPageData 插入页面数据
func InsertAnyPageData(fpDst, fpSrc string, waitFindLineStr string) (err error) {
	var dstFile, srcFile *os.File
	dstFile, err = os.OpenFile(fpDst, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("打开文件 %s 时出错: %v\n", fpDst, err)
	}
	defer dstFile.Close()

	var dstFileSomeLines []string
	foundLine := false
	scanner := bufio.NewScanner(dstFile)
	for scanner.Scan() {
		line := scanner.Text()
		dstFileSomeLines = append(dstFileSomeLines, line)
		if strings.HasPrefix(line, waitFindLineStr) {
			foundLine = true
			break
		}
	}
	if !foundLine {
		return fmt.Errorf("未找到 %q 的起始行", waitFindLineStr)
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

// CreateFileIfNotExists 如果文件不存在，则创建该文件及其所有上级目录
func CreateFileIfNotExists(filePath string) (hadExist bool, err error) {
	// 检查文件是否已经存在
	if _, err = os.Stat(filePath); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		// 如果错误不是“不存在”，说明有其他问题
		return false, fmt.Errorf("检查文件是否存在时出错: %w", err)
	}

	// 文件不存在，继续创建目录和文件
	// 获取父目录路径
	dir := filepath.Dir(filePath)

	// 创建所有必要的父目录
	if dir != "." && dir != "" {
		if err = os.MkdirAll(dir, 0666); err != nil {
			return false, fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 创建新文件（如果已存在则不会覆盖）
	file, err := os.Create(filePath)
	if err != nil {
		return false, fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	//fmt.Printf("成功创建文件: %s\n", filePath)
	return false, nil
}
