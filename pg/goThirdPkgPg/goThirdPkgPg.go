package goThirdPkgPg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/goThirdPkgJs"
	"github.com/before80/go/pg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MenuInfo struct {
	MenuName string `json:"menu_name"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
	Index    int    `json:"index"`
}

func InitThirdPkgMdFile(baseDirname string, menuInfo MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, baseDirname, menuInfo.Filename)
	return preInitMdFile(true, dir, menuInfo)
}

func InsertPkgDetailPageData(browserHwnd win.HWND, baseDirname string, menuInfo MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入detailPage=%s数据时遇到错误：%v", menuInfo.Url, r)
		}
	}()

	page.MustNavigate(menuInfo.Url)
	page.MustWaitLoad()
	mdFp := filepath.Join(contants.OutputFolderName, baseDirname, menuInfo.Filename, "_index.md")
	err = InsertVersionInfo(mdFp, menuInfo, page)
	if err != nil {
		return err
	}

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, goThirdPkgJs.ReplaceJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行goThirdPkgJs.ReplaceJs遇到错误：%v", menuInfo.Url, err)
	}

	err = pg.DealUniqueMd(browserHwnd, menuInfo.Url, "detailPage")
	if err != nil {
		return err
	}

	err = pg.InsertAnyPageData(mdFp, "> 仓库网址：")
	return
}

type VersionInfo struct {
	Version    string `json:"version"`
	CommitTime string `json:"commit_time"`
	Repo       string `json:"repo"`
}

func InsertVersionInfo(fpDst string, menuInfo MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入版本信息等=%s时遇到错误：%v", menuInfo.Url, r)
		}
	}()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(goThirdPkgJs.GetVersionInfoJs)
	if err != nil {
		return fmt.Errorf("在网页%s中执行goThirdPkgJs.GetVersionInfoJs遇到错误：%v", menuInfo.Url, err)
	}

	var versionInfo VersionInfo
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", menuInfo.Url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &versionInfo)
	if err != nil {
		return fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", menuInfo.Url, err)
	}

	var dstFile *os.File
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

	var newTotalLines []string
	newTotalLines = append(newTotalLines, dstFileSomeLines...)
	newTotalLines = append(newTotalLines, `>`)                                              // 插入两个空行
	newTotalLines = append(newTotalLines, fmt.Sprintf(`> 版本：%s`, versionInfo.Version))      // 插入两个空行
	newTotalLines = append(newTotalLines, `>`)                                              // 插入两个空行
	newTotalLines = append(newTotalLines, fmt.Sprintf(`> 发布时间：%s`, versionInfo.CommitTime)) // 插入两个空行
	newTotalLines = append(newTotalLines, `>`)
	newTotalLines = append(newTotalLines, fmt.Sprintf(`> 仓库网址：[%s](%s)`, versionInfo.Repo, versionInfo.Repo)) // 插入两个空行
	newTotalLines = append(newTotalLines, []string{"", ""}...)                                                // 插入两个空行

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

func preInitMdFile(useUnderlineIndexMd bool, dir string, menuInfo MenuInfo) (err error) {
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("无法创建%s目录：%v\n", dir, err)
	}
	var filename string
	if useUnderlineIndexMd {
		filename = "_index.md"
	} else {
		filename = menuInfo.Filename + ".md"
	}

	mdFp := filepath.Join(dir, filename)
	var mdF *os.File
	_, err1 := os.Stat(mdFp)

	// 当文件不存在的情况下，新建文件并初始化该文件
	if err1 != nil && errors.Is(err1, fs.ErrNotExist) {
		//fmt.Println("err=", err1)
		mdF, err = os.OpenFile(mdFp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("创建文件 %s 时出错: %w", mdFp, err)
		}
		defer mdF.Close()
		date := time.Now().Format(time.RFC3339)

		_, err = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++

> 原文：[%s](%s)
>
> 收录时间：%s
`, menuInfo.MenuName, date, menuInfo.Index*10, "", menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))

		if err != nil {
			return fmt.Errorf("初始化%s文件时出错: %v", mdFp, err)
		}
	}
	return
}
