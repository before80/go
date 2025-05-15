package goThirdPkgIndexPg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/goThirdPkgIndexJs"
	"github.com/before80/go/js/goThirdPkgJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/goThirdPkgIndexNext"
	"github.com/before80/go/pg"
	"github.com/before80/go/tr"
	"github.com/before80/go/wind"
	"github.com/go-rod/rod/lib/proto"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Weight2PkgInfosMap = make(map[int][]goThirdPkgIndexNext.PkgInfo)

func DealWithPkgBaseInfo(threadIndex int, wg *sync.WaitGroup) {
	var err error
	hadWgDone := false
	var pkgInfos []goThirdPkgIndexNext.PkgInfo
	defer func() {
		if r := recover(); r != nil {
			lg.ErrorToFile(fmt.Sprintf("线程%d出现异常：%v\n", threadIndex, r))
			lg.ErrorToFile(fmt.Sprintf("线程%d将退出\n", threadIndex))
			if !hadWgDone {
				lg.InfoToFile(fmt.Sprintf("在线程%d的defer中调用了wg.Done()\n", threadIndex))
				wg.Done()
			}
		}
	}()
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()
	defer func() {
		_ = page.Close()
	}()
LabelForContinue:
	pkgInfos = nil
	_, info, isEnd := goThirdPkgIndexNext.GetNextBaseInfoFromQueue()
	if isEnd {
		if !hadWgDone {
			hadWgDone = true
			lg.InfoToFile(fmt.Sprintf("在线程%d中已设置hadWgDone = true，且调用了wg.Done()\n", threadIndex))
			wg.Done()
		}
		return
	}
	page.MustNavigate(info.Url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject

	result, err = page.Eval(fmt.Sprintf(goThirdPkgIndexJs.FromTableGetAllPkgInfoJs, info.PkgName, info.Weight))
	if err != nil {
		panic(fmt.Errorf("在网页%s中执行goThirdPkgIndexJs.FromTableGetAllPkgInfoJs遇到错误：%v", info.Url, err))
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		panic(fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", info.Url, err))
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &pkgInfos)
	if err != nil {
		panic(fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", info.Url, err))
	}

	Weight2PkgInfosMap[pkgInfos[0].Weight] = pkgInfos

	//err = appendLinesToFile(pkgInfos)
	//if err != nil {
	//	panic(fmt.Sprintf("在处理%s遇到错误\n", info.PkgName))
	//}
	goto LabelForContinue
}

func SortAndGenAllPkgInfos() {
	var sPkgInfos [][]goThirdPkgIndexNext.PkgInfo
	var weights []int
	for k, _ := range Weight2PkgInfosMap {
		weights = append(weights, k)
	}
	slices.Sort(weights)
	for _, k := range weights {
		sPkgInfos = append(sPkgInfos, Weight2PkgInfosMap[k])
		goThirdPkgIndexNext.AllPkgInfos = append(goThirdPkgIndexNext.AllPkgInfos, Weight2PkgInfosMap[k]...)
	}
}

func TruncWriteLinesToFile() (err error) {
	var file *os.File
	file, err = os.OpenFile(filepath.Join(contants.OutputFolderName, "go_third_pkg_info.txt"), os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	var sPkgInfos [][]goThirdPkgIndexNext.PkgInfo
	var weights []int
	for k, _ := range Weight2PkgInfosMap {
		weights = append(weights, k)
	}
	slices.Sort(weights)
	for _, k := range weights {
		sPkgInfos = append(sPkgInfos, Weight2PkgInfosMap[k])
	}

	// 创建一个写入器
	writer := bufio.NewWriter(file)

	// 先添加一个换行符
	if _, err = file.WriteString("\n"); err != nil {
		return fmt.Errorf("写入换行符时出错: %w", err)
	}

	for _, pkgInfos := range sPkgInfos {
		// 遍历要追加的每一行内容
		for _, info := range pkgInfos {
			line := fmt.Sprintf("%s||%s||%s||%s||%d||%d||%d||%s\n", info.Url, info.PkgName, info.Dir, info.Filename, info.NeedPreCreateIndex, info.Weight, info.PreCreateIndexWeight, info.Desc)
			// 写入当前行
			if _, err = writer.WriteString(line); err != nil {
				return fmt.Errorf("写入行时出错: %w", err)
			}
		}
		if _, err = writer.WriteString("\n"); err != nil {
			return fmt.Errorf("写入行时出错: %w", err)
		}
	}

	// 将缓冲区的内容刷新到文件
	if err = writer.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区时出错: %w", err)
	}
	return nil
}

type VersionInfo struct {
	Version    string `json:"version"`
	CommitTime string `json:"commit_time"`
	Repo       string `json:"repo"`
}

func DealWithPkgPageData(threadIndex int, wg *sync.WaitGroup) {
	var err error
	hadWgDone := false
	defer func() {
		if r := recover(); r != nil {
			lg.ErrorToFile(fmt.Sprintf("线程%d出现异常：%v\n", threadIndex, r))
			lg.ErrorToFile(fmt.Sprintf("线程%d将退出\n", threadIndex))
			if !hadWgDone {
				lg.InfoToFile(fmt.Sprintf("在线程%d的defer中调用了wg.Done()\n", threadIndex))
				wg.Done()
			}
		}
	}()
	preDir := cfg.Default.GoThirdPkgPreFolderName
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()

	defer func() {
		_ = page.Close()
	}()
	var pageTitle, chromePageWindowTitle string
	var fpDst string
	var versionInfo VersionInfo
	var result *proto.RuntimeRemoteObject
	var hadInsetPageData bool
	articleUrl := ""
	desc := ""
	uniqueMdFilename := "do" + strconv.Itoa(threadIndex) + ".md"
	relUniqueMdFilePath := filepath.Join("markdown", uniqueMdFilename)
	typoraWindowTitle := uniqueMdFilename + " - Typora"
	_, _ = pg.CreateFileIfNotExists(relUniqueMdFilePath)
	absUniqueMdFilePath, _ := filepath.Abs(relUniqueMdFilePath)
LabelForContinue:
	_ = tr.TruncFileContent(relUniqueMdFilePath)
	date := time.Now().Format(time.RFC3339)
	_, pkg, isEnd := goThirdPkgIndexNext.GetNextPkgInfoFromQueue()
	lg.InfoToFile(fmt.Sprintf("线程%d正要处理的pkg=%v\n", threadIndex, pkg))
	if isEnd {
		if !hadWgDone {
			hadWgDone = true
			lg.InfoToFile(fmt.Sprintf("在线程%d中已设置hadWgDone = true，且调用了wg.Done()\n", threadIndex))
			wg.Done()
		}
		return
	}

	fpDst = filepath.Join(contants.OutputFolderName, preDir, pkg.Dir, pkg.Filename+".md")
	// 判断md文件中是否已经插入内容
	hadInsetPageData, err = pg.JudgeHadInsertPageData(fpDst, "> 仓库网址：")
	if hadInsetPageData {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%s之前已经插入过数据\n", pkg.PkgName))
		goto LabelForContinue
	}

	if pkg.Url == "" {
		articleUrl = ""
	} else {
		articleUrl = fmt.Sprintf(`[%s](%s)`, pkg.Url, pkg.Url)
	}

	desc = strings.ReplaceAll(pkg.Desc, "\"", "'")
	desc = strings.ReplaceAll(desc, "\n", " ")
	page.MustNavigate(pkg.Url)
	page.MustWaitLoad()

	result, err = page.Eval(goThirdPkgJs.GetVersionInfoJs)
	if err != nil {
		panic(fmt.Sprintf("线程%d在网页%s中执行goThirdPkgJs.GetVersionInfoJs遇到错误：%v", threadIndex, pkg.Url, err))
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		panic(fmt.Sprintf("线程%d在处理网页%s时执行json.Marshal遇到错误: %v", threadIndex, pkg.Url, err))
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &versionInfo)
	if err != nil {
		panic(fmt.Sprintf("线程%d在处理网页%s时执行json.Unmarshal遇到错误: %v", threadIndex, pkg.Url, err))
	}

	// 创建 md文件和相关目录
	if pkg.NeedPreCreateIndex == 1 {
		preIndexMd := filepath.Join(contants.OutputFolderName, preDir, pkg.Dir, "_index.md")
		preIndexMdTitles := strings.Split(pkg.Dir, "/")
		preIndexMdTitle := preIndexMdTitles[len(preIndexMdTitles)-1]
		hadExist, _ := pg.CreateFileIfNotExists(preIndexMd)

		if !hadExist {
			preIndexMdF, _ := os.OpenFile(preIndexMd, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			_, _ = preIndexMdF.WriteString(fmt.Sprintf(`+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++

> 原文：%s
>
> 收录时间：%s
`, preIndexMdTitle, date, pkg.PreCreateIndexWeight, "", "", fmt.Sprintf("`%s`", date)))
			_ = preIndexMdF.Close()
		}
	}

	_, _ = pg.CreateFileIfNotExists(fpDst)
	fpDstF, _ := os.OpenFile(fpDst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	_, _ = fpDstF.WriteString(fmt.Sprintf(`+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++

> 原文：%s
>
> 收录时间：%s
>
> 版本：%s
>
> 发布时间：%s
>
> 仓库网址：[%s](%s)
`, pkg.PkgName, date, pkg.Weight, desc, articleUrl, fmt.Sprintf("`%s`", date), versionInfo.Version, versionInfo.CommitTime, versionInfo.Repo, versionInfo.Repo))
	_ = fpDstF.Close()

	// 获取当前网页的title，在后面会用来查找该网页所在窗口的操作句柄
	result, _ = page.Eval(`() => { return document.title }`)
	pageTitle = result.Value.String()
	chromePageWindowTitle = pageTitle + " - Google Chrome"

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, goThirdPkgJs.ReplaceJs))
	if err != nil {
		panic(fmt.Errorf("线程%d在网页%s中执行goThirdPkgJs.ReplaceJs遇到错误：%v", threadIndex, pkg.Url, err))
	}

	// 再次清空
	_ = tr.TruncFileContent(relUniqueMdFilePath)

	contentBytes, _ := wind.DoCopyAndPaste(threadIndex, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, pkg.Url)
	if contentBytes == 0 {
		lg.InfoToFile(fmt.Sprintf("线程%d发现复制网页%s的字节数为0，将加入到下一次进行重试", threadIndex, pkg.Url))
		goThirdPkgIndexNext.PushWaitDealPkgInfoToQueue([]goThirdPkgIndexNext.PkgInfo{pkg})
	} else {
		lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
		err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 仓库网址：")
		if err != nil {
			panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, pkg.Url, err))
		}
	}

	goto LabelForContinue
}
