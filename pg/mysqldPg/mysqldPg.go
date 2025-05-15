package mysqldPg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/mysqldJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/mysqldNext"
	"github.com/before80/go/pg"
	"github.com/before80/go/res"
	"github.com/before80/go/tr"
	"github.com/before80/go/wind"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var didF *os.File
var didUrl []string

func init() {
	var err error
	err = os.MkdirAll(filepath.Join(contants.DidFolderName, "mysql"), 0777)
	if err != nil {
		panic(fmt.Sprintf("无法创建%s目录：%v\n", filepath.Join(contants.DidFolderName, "mysql"), err))
	}

	didF, err = os.OpenFile(filepath.Join(contants.DidFolderName, "mysql", "mysql.com.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("无法创建或打开文件%q: %v\n", "mysql.com.txt", err))
		panic(err)
	}

	scanner := bufio.NewScanner(didF)
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的文本
		//fmt.Println("line=", line)
		if strings.TrimSpace(line) != "" {
			didUrl = append(didUrl, line)
		}
	}
	//fmt.Println("didUrl=", didUrl)
	//_, _ = didF.Seek(2, 0)
	res.NewMySQL(bufio.NewWriter(didF))
}

func CloseInitFiles() {
	if didF != nil {
		_ = didF.Close()
	}
}

func GetAllMenuInfo(page *rod.Page, url string) (menuInfos []mysqldNext.MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("获取allMenuInfo时遇到错误：%v", r)
		}
	}()

	page.MustNavigate(url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(mysqldJs.ExpandMenusJs)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行mysqldJs.ExpandMenusJs遇到错误：%v", url, err)
	}
	time.Sleep(3 * time.Second)

	result, err = page.Eval(mysqldJs.GetAllMenuInfoJs)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行mysqldJs.GetAllMenuInfoJs遇到错误：%v", url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &menuInfos)
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
	preDir := cfg.Default.MySQLdPreFolderName
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()

	defer func() {
		_ = page.Close()
	}()
	var pageTitle, chromePageWindowTitle string
	var fpDst string
	var result *proto.RuntimeRemoteObject
	var hadInsetPageData bool
	uniqueMdFilename := "do" + strconv.Itoa(threadIndex) + ".md"
	relUniqueMdFilePath := filepath.Join("markdown", uniqueMdFilename)
	typoraWindowTitle := uniqueMdFilename + " - Typora"
	_, _ = pg.CreateFileIfNotExists(relUniqueMdFilePath)
	absUniqueMdFilePath, _ := filepath.Abs(relUniqueMdFilePath)
LabelForContinue:
	hadInsetPageData = false
	_ = tr.TruncFileContent(relUniqueMdFilePath)
	date := time.Now().Format(time.RFC3339)
	curMenu, isEnd := mysqldNext.GetNextMenuInfoFromQueue()
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
		if curMenu.HaveSub == 1 {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Filename, "_index.md")
		} else {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Filename+".md")
		}
	} else {
		if curMenu.HaveSub == 1 {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, curMenu.Filename, "_index.md")
		} else {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, curMenu.Filename+".md")
		}
	}

	// 判断md文件中是否已经插入内容
	hadInsetPageData, err = pg.JudgeHadInsertPageData(fpDst, "> 收录时间：")
	if hadInsetPageData {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%s之前已经插入过数据\n", curMenu.MenuName))
		goto LabelForContinue
	}

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

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, mysqldJs.ReplaceJs))
	if err != nil {
		panic(fmt.Errorf("线程%d在网页%s中执行mysqldJs.ReplaceJs遇到错误：%v", threadIndex, curMenu.Url, err))
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
		mysqldNext.PushWaitDealMenuInfoToQueue([]mysqldNext.MenuInfo{curMenu})
	} else {
		lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
		err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 收录时间：")
		if err != nil {
			panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, curMenu.Url, err))
		}
	}

	goto LabelForContinue
}
