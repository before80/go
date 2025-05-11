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
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
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
	_, info, isEnd := goThirdPkgIndexNext.GetNextBaseInfoFromStack()
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
			line := fmt.Sprintf("%s||%s||%s||%s||%d||%d||%s\n", info.Url, info.PkgName, info.Dir, info.Filename, info.NeedPreCreateIndex, info.Weight, info.Desc)
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
	preDir := "go_third_pkg"
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()

	defer func() {
		_ = page.Close()
	}()
	var pageTitle, chromePageWindowTitle string
	var originArticle string
	var fpDst string
	var versionInfo VersionInfo
	var result *proto.RuntimeRemoteObject
	uniqueMdFilename := "do" + strconv.Itoa(threadIndex) + ".md"
	relUniqueMdFilePath := filepath.Join("markdown", uniqueMdFilename)
	typoraWindowTitle := uniqueMdFilename + " - Typora"
	_, _ = pg.CreateFileIfNotExists(relUniqueMdFilePath)
	absUniqueMdFilePath, _ := filepath.Abs(relUniqueMdFilePath)
LabelForContinue:
	date := time.Now().Format(time.RFC3339)
	_, pkg, isEnd := goThirdPkgIndexNext.GetNextPkgInfoFromStack()
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
	_ = tr.TruncFileContent(relUniqueMdFilePath)

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
			if pkg.Url == "" {
				originArticle = ""
			} else {
				originArticle = fmt.Sprintf("[%s](%s)", pkg.Url, pkg.Url)
			}
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
`, preIndexMdTitle, date, pkg.Weight, "", originArticle, fmt.Sprintf("`%s`", date)))
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

> 原文：[%s](%s)
>
> 收录时间：%s
>
> 版本：%s
>
> 发布时间：%s
>
> 仓库网址：[%s](%s)
`, pkg.PkgName, date, pkg.Weight, "", pkg.Url, pkg.Url, fmt.Sprintf("`%s`", date), versionInfo.Version, versionInfo.CommitTime, versionInfo.Repo, versionInfo.Repo))
	_ = fpDstF.Close()

	// 获取当前网页的title，在后面会用来查找该网页所在窗口的操作句柄
	result, _ = page.Eval(`() => { return document.title }`)
	pageTitle = result.Value.String()
	chromePageWindowTitle = pageTitle + " - Google Chrome"

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, goThirdPkgJs.ReplaceJs))
	if err != nil {
		panic(fmt.Errorf("线程%d在网页%s中执行goThirdPkgJs.ReplaceJs遇到错误：%v", threadIndex, pkg.Url, err))
	}

	_ = DoCopyAndPaste(threadIndex, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, pkg.Url)
	lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
	err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 仓库网址：")
	if err != nil {
		panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, pkg.Url, err))
	}

	goto LabelForContinue
}

var copyPasteLock sync.Mutex

func DoCopyAndPaste(threadIndex int, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, url string) (err error) {
	copyPasteLock.Lock()
	defer copyPasteLock.Unlock()
	var typoraHwnd win.HWND
	browserHwnd := robotgo.FindWindow(chromePageWindowTitle)

	_ = wind.OpenTypora(absUniqueMdFilePath)
	timeoutChan := time.After(10 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 每隔 interval 时间检查一次条件
			hwnd1 := robotgo.FindWindow(typoraWindowTitle)
			lg.InfoToFile(fmt.Sprintf("%d - typoraHwnd=%v\n", threadIndex, hwnd1))
			if hwnd1 != 0 {
				typoraHwnd = hwnd1
				goto LabelForContinue
			}
			//hwnd1, err1 = wind.FindWindowByTitle(uniqueMdFilename + " - Typora")
		case <-timeoutChan:
			// 超时后退出循环
			goto LabelForContinue
		}
	}
LabelForContinue:
	lg.InfoToFile(fmt.Sprintf("线程%d中获取到的typoraHwnd=%v\n", threadIndex, typoraHwnd))

	contentBytes, err1 := wind.InChromePageDoCtrlAAndC(browserHwnd)
	lg.InfoToFile(fmt.Sprintf("在页面%s获取到的字节数为：%d\n", url, contentBytes))
	if err1 != nil {
		lg.ErrorToFile(fmt.Sprintf("在浏览器中进行复制遇到错误：%v\n", err1))
	}
	_ = wind.DoCtrlVAndS(typoraHwnd, contentBytes)
	_ = win.SendMessage(typoraHwnd, win.WM_CLOSE, 0, 0)
	time.Sleep(time.Duration(cfg.Default.WaitTyporaCloseSeconds) * time.Second)
	return nil
}
