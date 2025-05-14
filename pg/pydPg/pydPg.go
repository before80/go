package pydPg

import (
	"encoding/json"
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/pydJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/pydNext"
	"github.com/before80/go/pg"
	"github.com/before80/go/tr"
	"github.com/before80/go/tr/pydTr"
	"github.com/before80/go/wind"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func GetBarMenus(page *rod.Page, url string) (barMenuInfos []pydNext.MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("获取barmenu时遇到错误：%v", r)
		}
	}()
	// https://docs.python.org/zh-cn/3.13/index.html
	page.MustNavigate(url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(fmt.Sprintf(pydJs.GetBarMenusJs, url))
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行GetBarMenusJs遇到错误：%v", url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &barMenuInfos)
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
	preDir := cfg.Default.PydPreFolderName
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()

	defer func() {
		_ = page.Close()
	}()
	var pageTitle, chromePageWindowTitle string
	var fpDst string
	var result *proto.RuntimeRemoteObject
	var subMenuInfos []pydNext.MenuInfo
	uniqueMdFilename := "do" + strconv.Itoa(threadIndex) + ".md"
	relUniqueMdFilePath := filepath.Join("markdown", uniqueMdFilename)
	typoraWindowTitle := uniqueMdFilename + " - Typora"
	_, _ = pg.CreateFileIfNotExists(relUniqueMdFilePath)
	absUniqueMdFilePath, _ := filepath.Abs(relUniqueMdFilePath)
LabelForContinue:
	subMenuInfos = nil
	_ = tr.TruncFileContent(relUniqueMdFilePath)
	date := time.Now().Format(time.RFC3339)
	curMenu, isEnd := pydNext.GetNextMenuInfoFromQueue()
	lg.InfoToFile(fmt.Sprintf("线程%d正要处理的menu=%v\n", threadIndex, curMenu))
	if isEnd {
		if !hadWgDone {
			hadWgDone = true
			lg.InfoToFile(fmt.Sprintf("线程%d中已设置hadWgDone = true，且调用了wg.Done()\n", threadIndex))
			wg.Done()
		}
		return
	}

	page.MustNavigate(curMenu.Url)
	page.MustWaitLoad()

	if curMenu.IsTopMenu == 1 && curMenu.Filename == "howto" {
		result, err = page.Eval(pydJs.GetMenusJs2)
	} else {
		result, err = page.Eval(pydJs.GetMenusJs)
	}

	if err != nil {
		panic(fmt.Sprintf("线程%d在网页%s中执行pydJs.GetMenusJs遇到错误：%v", threadIndex, curMenu.Url, err))
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		panic(fmt.Sprintf("线程%d在处理网页%s时执行json.Marshal遇到错误: %v", threadIndex, curMenu.Url, err))
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &subMenuInfos)
	if err != nil {
		panic(fmt.Sprintf("线程%d在处理网页%s时执行json.Unmarshal遇到错误: %v", threadIndex, curMenu.Url, err))
	}
	lg.InfoToFile(fmt.Sprintf("线程%d网页%s的子菜单=%v\n", threadIndex, curMenu.Url, subMenuInfos))

	if curMenu.IsTopMenu == 1 {
		fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, "_index.md")
	} else {
		fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, curMenu.Filename+".md")
	}

	if len(subMenuInfos) > 0 {
		for i, _ := range subMenuInfos {
			subMenuInfos[i].TopMenuName = curMenu.TopMenuName
			if curMenu.IsTopMenu == 1 {
				subMenuInfos[i].Dir = curMenu.Dir
			} else {
				subMenuInfos[i].Dir = curMenu.Dir + "/" + curMenu.Filename
			}
		}
		if curMenu.IsTopMenu != 1 {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, curMenu.Filename, "_index.md")
		}
	}

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

> 原文：[%s](%s)
>
> 收录时间：%s
`, curMenu.MenuName, curMenu.MenuName, date, "", curMenu.Weight, curMenu.Url, curMenu.Url, fmt.Sprintf("`%s`", date)))
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

> 原文：[%s](%s)
>
> 收录时间：%s
`, curMenu.MenuName, date, curMenu.Weight, "", curMenu.Url, curMenu.Url, fmt.Sprintf("`%s`", date)))
	}

	_ = mdF.Close()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		panic(fmt.Errorf("线程%d在网页%s中执行GetDetailPageDataJs遇到错误：%v", threadIndex, curMenu.Url, err))
	}
	// 获取当前网页的title，在后面会用来查找该网页所在窗口的操作句柄
	result, _ = page.Eval(`() => { return document.title }`)
	pageTitle = result.Value.String()
	chromePageWindowTitle = pageTitle + " - Google Chrome"

	// 再次清空
	_ = tr.TruncFileContent(relUniqueMdFilePath)

	contentBytes, _ := wind.DoCopyAndPaste(threadIndex, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, curMenu.Url)

	if contentBytes == 0 {
		lg.InfoToFile(fmt.Sprintf("线程%d发现复制网页%s的字节数为0，将加入到下一次进行重试", threadIndex, curMenu.Url))
		pydNext.PushWaitDealMenuInfoToQueue([]pydNext.MenuInfo{curMenu})
	} else {
		_, err = pydTr.ReplaceMarkdownFileContent(absUniqueMdFilePath)
		if err != nil {
			panic(fmt.Errorf("线程%d在替换网页%s的内容对应的md文件时出现错误：%v", threadIndex, curMenu.Url, err))
		}
		lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
		err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 收录时间：")
		if err != nil {
			panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, curMenu.Url, err))
		}
	}

	if len(subMenuInfos) > 0 {
		// 将 subMenuInfos 推入
		pydNext.PushWaitDealMenuInfoToQueue(subMenuInfos)
	}
	goto LabelForContinue
}
