package godPg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/before80/go/bs"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/godJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/godNext"
	"github.com/before80/go/pg"
	"github.com/before80/go/tr"
	"github.com/before80/go/wind"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// GetAllStdPkgInfo 获取所有标准库的pkg信息
func GetAllStdPkgInfo(page *rod.Page, url string) (stdPkgMenuInfos []godNext.MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("获取allStdPkgInfo时遇到错误：%v", r)
		}
	}()

	page.MustNavigate(url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(godJs.FromTableGetAllStdPkgInfoJs)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行godJs.FromTableGetAllStdPkgInfoJs遇到错误：%v", url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &stdPkgMenuInfos)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", url, err)
	}
	return
}

func DealWithMenuPageData(threadIndex int, wg *sync.WaitGroup) {
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
	preDir := "go_std"
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()
	page.MustSetViewport(cfg.Default.BrowserWidth, cfg.Default.BrowserHeight, 1, false)
	//page.SetViewport()

	defer func() {
		_ = page.Close()
	}()
	var pageTitle, chromePageWindowTitle string
	var fpDst string
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
	curMenu, isEnd := godNext.GetNextMenuInfoFromQueue()
	lg.InfoToFile(fmt.Sprintf("线程%d正要处理的menu=%v\n", threadIndex, curMenu))
	if isEnd {
		if !hadWgDone {
			hadWgDone = true
			lg.InfoToFile(fmt.Sprintf("线程%d中已设置hadWgDone = true，且调用了wg.Done()\n", threadIndex))
			wg.Done()
		}
		return
	}

	if curMenu.IsTopMenu == 1 {
		if len(curMenu.Children) > 0 {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Filename, "_index.md")
		} else {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Filename+".md")
		}
	} else {
		if len(curMenu.Children) > 0 {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.PFilename, "_index.md")
		} else {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.PFilename, curMenu.Filename+".md")
		}
	}

	// 判断md文件中是否已经插入内容
	hadInsetPageData, err = pg.JudgeHadInsertPageData(fpDst, "> 收录时间：")
	if hadInsetPageData {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%s之前已经插入过数据\n", curMenu.MenuName))
		goto LabelForContinue
	}

	if curMenu.Url == "" {
		articleUrl = ""
	} else {
		articleUrl = fmt.Sprintf(`[%s](%s)`, curMenu.Url, curMenu.Url)
	}

	desc = strings.ReplaceAll(curMenu.Desc, "\"", "'")
	desc = strings.ReplaceAll(desc, "\n", " ")
	page.MustNavigate(curMenu.Url)
	page.MustWaitLoad()

	hadExist, _ := pg.CreateFileIfNotExists(fpDst)
	_ = hadExist
	mdF, _ := os.OpenFile(fpDst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if curMenu.IsTopMenu == 1 {

		_, _ = mdF.WriteString(fmt.Sprintf(`+++
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

> 原文：%s
>
> 收录时间：%s
`, curMenu.MenuName, curMenu.MenuName, date, curMenu.Desc, curMenu.Weight, articleUrl, fmt.Sprintf("`%s`", date)))
	} else {
		_, _ = mdF.WriteString(fmt.Sprintf(`+++
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
`, curMenu.MenuName, date, curMenu.Weight, desc, articleUrl, fmt.Sprintf("`%s`", date)))
	}

	_ = mdF.Close()
	//time.Sleep(10 * time.Second)

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, `
	const des = document.querySelectorAll(".Documentation-exampleDetails");
    const noTopHLevel = 2;
    const maxIterations = 100; // 最大迭代次数
    if (des.length > 0) {
        console.log("des.length=",des.length)
        des.forEach(de => {
            const deh = de.querySelector(":scope > .Documentation-exampleDetailsHeader");
            deh.click();
            const deb = de.querySelector(":scope > .Documentation-exampleDetailsBody");
			console.log("deh=",deh);
			console.log("deb=",deb);        
		})
	}`))
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond("出现错误1", 1)
		panic(fmt.Errorf("线程%d在网页%s中执行godJs.ReplaceJs遇到错误：%v", threadIndex, curMenu.Url, err))
	}

	// time.Sleep(1000 * time.Second)
	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, godJs.ReplaceJs))
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond("出现错误2", 10000)
		panic(fmt.Errorf("线程%d在网页%s中执行godJs.ReplaceJs遇到错误：%v", threadIndex, curMenu.Url, err))
	}

	//time.Sleep(1000 * time.Second)
	// 获取当前网页的title，在后面会用来查找该网页所在窗口的操作句柄
	result, _ = page.Eval(`() => { return document.title }`)
	pageTitle = result.Value.String()
	chromePageWindowTitle = pageTitle + " - Google Chrome"

	// 再次清空
	_ = tr.TruncFileContent(relUniqueMdFilePath)

	contentBytes, _ := wind.DoCopyAndPaste(threadIndex, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, curMenu.Url)
	if contentBytes == 0 {
		lg.InfoToFile(fmt.Sprintf("线程%d发现复制网页%s的字节数为0，将加入到下一次进行重试", threadIndex, curMenu.Url))
		godNext.PushWaitDealMenuInfoToQueue([]godNext.MenuInfo{curMenu})
	} else {
		lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
		err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 收录时间：")
		if err != nil {
			panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, curMenu.Url, err))
		}
	}

	goto LabelForContinue
}
