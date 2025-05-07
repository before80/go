package mysqldPg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/mysqldJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/pg"
	"github.com/before80/go/res"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tailscale/win"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

type MenuInfo struct {
	MenuName string `json:"menu_name"`
	Filename string `json:"filename"`
	FilePath string `json:"file_path"`
	Url      string `json:"url"`
	Index    int    `json:"index"`
	IsTop    int    `json:"is_top"`
}

func GetAllMenuInfo(page *rod.Page, url string) (menuInfos []MenuInfo, err error) {
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

func DealMenuMdFile(surplus int, browserHwnd win.HWND, dirPrefix string, menuInfo MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("获取页面内容时遇到错误：%v", r)
		}
	}()
	if slices.Contains(didUrl, menuInfo.Url) {
		lg.InfoToFileAndStdOut(fmt.Sprintf("之前已处理 %s - %s\n", menuInfo.MenuName, menuInfo.Url))
		return
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("还有%d 正在处理 %s - %s\n", surplus, menuInfo.MenuName, menuInfo.Url))

	err = preInitMdFile(dirPrefix, menuInfo)
	if err != nil {
		return
	}
	page.MustNavigate(menuInfo.Url)
	page.MustWaitLoad()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, mysqldJs.ReplaceJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行mysqldJs.ReplaceJs遇到错误：%v", menuInfo.Url, err)
	}

	err = InsertDetailPageData(browserHwnd, dirPrefix, menuInfo)
	if err != nil {
		return
	}

	// 记录已经处理的url
	res.MySQL.WriteStringAndFlush(fmt.Sprintf("%s\n", menuInfo.Url))
	lg.InfoToFileAndStdOut(fmt.Sprintf("处理完毕 %s - %s\n", menuInfo.MenuName, menuInfo.Url))
	return nil
}

func InsertDetailPageData(browserHwnd win.HWND, dirPrefix string, menuInfo MenuInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入detailPage=%s数据时遇到错误：%v", menuInfo.Url, r)
		}
	}()

	err = pg.DealUniqueMd(browserHwnd, menuInfo.Url, "detailPage")
	if err != nil {
		return err
	}
	mdFilePath := filepath.Join(contants.OutputFolderName, dirPrefix, menuInfo.FilePath)
	err = pg.InsertAnyPageData(mdFilePath, "> 收录时间：")
	return
}

func preInitMdFile(dirPrefix string, menuInfo MenuInfo) (err error) {
	var mdF *os.File
	mdFilePath := filepath.Join(contants.OutputFolderName, dirPrefix, menuInfo.FilePath)
	dir := filepath.Dir(mdFilePath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("创建目录时出错: %w", err)
	}
	mdF, err = os.OpenFile(mdFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("创建文件 %s 时出错: %w", mdFilePath, err)
	}
	defer mdF.Close()
	date := time.Now().Format(time.RFC3339)
	if menuInfo.IsTop == 1 {
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

> 原文：[%s](%s)
>
> 收录时间：%s
`, menuInfo.MenuName, menuInfo.MenuName, date, "", menuInfo.Index*10, menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))
	} else {
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
	}

	if err != nil {
		return fmt.Errorf("初始化%s文件时出错: %v", mdFilePath, err)
	}
	return nil
}
