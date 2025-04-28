package pydPg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/pydJs"
	"github.com/before80/go/myf"
	"github.com/before80/go/wind"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetBarMenus(page *rod.Page, url string) (barMenuInfos []MenuInfo, err error) {
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

type MenuInfo struct {
	MenuName string `json:"menu_name"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
}

func InitBarIndexMdFile(index int, barMenuInfo MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename)
	return preInitMdFile(index, true, true, dir, barMenuInfo)
}

func InitSecondIndexMdFile(index int, barMenuInfo MenuInfo, secondMenuInfo MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename)
	return preInitMdFile(index, false, true, dir, secondMenuInfo)
}

func InitThirdIndexMdFile(index int, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, thirdMenuInfo MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename)
	return preInitMdFile(index, false, true, dir, thirdMenuInfo)
}

func InitSecondDetailPageMdFile(index int, barMenuInfo MenuInfo, secondMenuInfo MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename)
	return preInitMdFile(index, false, false, dir, secondMenuInfo)
}

func InitThirdDetailPageMdFile(index int, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, thirdMenuInfo MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename)
	return preInitMdFile(index, false, false, dir, thirdMenuInfo)
}

func InitFourthDetailPageMdFile(index int, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, thirdMenuInfo MenuInfo, fourthMenuInfo MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename)
	return preInitMdFile(index, false, false, dir, fourthMenuInfo)
}

func InsertBarMenuPageData(browserHwnd win.HWND, barMenuInfo MenuInfo, page *rod.Page) (secondMenus []MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("初始化barmenu=%s时遇到错误：%v", barMenuInfo.Url, r)
		}
	}()

	page.MustNavigate(barMenuInfo.Url)
	page.MustWaitLoad()

	// 判断是否还有第二级菜单
	var result *proto.RuntimeRemoteObject

	result, err = page.Eval(fmt.Sprintf(pydJs.GetSecondMenusJs, barMenuInfo.Url))
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行GetSecondMenusJs遇到错误：%v", barMenuInfo.Url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", barMenuInfo.Url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &secondMenus)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", barMenuInfo.Url, err)
	}

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", barMenuInfo.Url, err)
	}

	err = dealUniqueMd(browserHwnd, barMenuInfo.Url, "barmenu")
	if err != nil {
		return nil, err
	}

	indexMdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, "_index.md")
	err = insertAnyPageData(indexMdFp)
	return
}

func InsertSecondMenuPageData(browserHwnd win.HWND, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, page *rod.Page) (err error) {
	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", secondMenuInfo.Url, err)
	}

	err = dealUniqueMd(browserHwnd, secondMenuInfo.Url, "second")
	if err != nil {
		return err
	}
	indexMdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, "_index.md")
	err = insertAnyPageData(indexMdFp)
	return
}

func InsertThirdMenuPageData(browserHwnd win.HWND, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, thirdMenuInfo MenuInfo, page *rod.Page) (err error) {
	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", thirdMenuInfo.Url, err)
	}

	err = dealUniqueMd(browserHwnd, thirdMenuInfo.Url, "third")
	if err != nil {
		return err
	}

	indexMdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename, "_index.md")
	err = insertAnyPageData(indexMdFp)
	return
}

func InsertSecondDetailPageData(browserHwnd win.HWND, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入detailPage=%s数据时遇到错误：%v", secondMenuInfo.Url, r)
		}
	}()
	page.MustNavigate(secondMenuInfo.Url)
	page.MustWaitLoad()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", secondMenuInfo.Url, err)
	}

	err = dealUniqueMd(browserHwnd, secondMenuInfo.Url, "secondDetailPage")
	if err != nil {
		return err
	}
	mdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename+".md")
	err = insertAnyPageData(mdFp)
	return
}

func InsertThirdDetailPageData(browserHwnd win.HWND, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, thirdMenuInfo MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入thirdDetailPage=%s数据时遇到错误：%v", thirdMenuInfo.Url, r)
		}
	}()
	page.MustNavigate(thirdMenuInfo.Url)
	page.MustWaitLoad()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", thirdMenuInfo.Url, err)
	}

	err = dealUniqueMd(browserHwnd, thirdMenuInfo.Url, "thirdDetailPage")
	if err != nil {
		return err
	}

	mdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename+".md")
	err = insertAnyPageData(mdFp)
	return
}

func InsertFourthDetailPageData(browserHwnd win.HWND, barMenuInfo MenuInfo, secondMenuInfo MenuInfo, thirdMenuInfo MenuInfo, fourthMenuInfo MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入fourthDetailPage=%s数据时遇到错误：%v", fourthMenuInfo.Url, r)
		}
	}()
	page.MustNavigate(fourthMenuInfo.Url)
	page.MustWaitLoad()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", fourthMenuInfo.Url, err)
	}

	err = dealUniqueMd(browserHwnd, fourthMenuInfo.Url, "fourthDetailPage")
	if err != nil {
		return err
	}

	mdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename, fourthMenuInfo.Filename+".md")
	err = insertAnyPageData(mdFp)
	return
}

// dealUniqueMd 处理唯一的md文件中的内容
func dealUniqueMd(browserHwnd win.HWND, curUrl, step string) (err error) {
	uniqueMdFilepath := cfg.Default.UniqueMdFilepath
	// 获取文件名
	spSlice := strings.Split(uniqueMdFilepath, "\\")
	mdFilename := spSlice[len(spSlice)-1]

	// 清空唯一共用的markdown文件的文件内容
	err = myf.TruncFileContent(uniqueMdFilepath)
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
	wind.SelectAllAndDelete(typoraHwnd)
	wind.CtrlV(typoraHwnd)
	time.Sleep(time.Duration(cfg.Default.WaitTyporaCopiedToSaveSeconds) * time.Second)
	wind.CtrlS(typoraHwnd)
	time.Sleep(time.Duration(cfg.Default.WaitTyporaSaveSeconds) * time.Second)
	robotgo.CloseWindow()
	time.Sleep(time.Duration(cfg.Default.WaitTyporaCloseSeconds) * time.Second)
	_, err = myf.ReplaceMarkdownFileContent(uniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", step, curUrl, err)
	}
	return nil
}

// insertAnyPageData 插入页面数据
func insertAnyPageData(fpDst string) (err error) {
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
		return fmt.Errorf("未找到类型为 %s 的起始行", "> 收录时间：")
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

// findShouLuStart 找到 “收录时间：”所在行
func findShouLuStart(lines []string, shouLu string) (start int, err error) {
	for i, line := range lines {
		if strings.HasPrefix(line, shouLu) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("未找到%q所在行", shouLu)
}

func GetThirdLevelMenu(secondMenuInfo MenuInfo, page *rod.Page) (thirdMenuInfos []MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("在第二级菜单%s中获取第三级菜单时遇到错误：%v", secondMenuInfo.Url, r)
		}
	}()

	page.MustNavigate(secondMenuInfo.Url)
	page.MustWaitLoad()
	return evalJsGetSubMenuInfos(page, "GetThirdMenusJs", pydJs.GetThirdMenusJs, secondMenuInfo.Url)
}

func GetFourthLevelMenu(thirdMenuInfo MenuInfo, page *rod.Page) (fourthMenuInfos []MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("在第二级菜单%s中获取第三级菜单时遇到错误：%v", thirdMenuInfo.Url, r)
		}
	}()

	page.MustNavigate(thirdMenuInfo.Url)
	page.MustWaitLoad()
	return evalJsGetSubMenuInfos(page, "GetFourthMenusJs", pydJs.GetFourthMenusJs, thirdMenuInfo.Url)
}

func evalJsGetSubMenuInfos(page *rod.Page, jsName, js, pageUrl string) (subMenuInfos []MenuInfo, err error) {
	// 判断是否还有第三级菜单
	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(fmt.Sprintf(js, pageUrl))

	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行%s遇到错误：%v", pageUrl, jsName, err)
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", pageUrl, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &subMenuInfos)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", pageUrl, err)
	}
	return
}

func preInitMdFile(index int, isBar, useUnderlineIndexMd bool, dir string, menuInfo MenuInfo) (err error) {
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
		if isBar {
			_, err = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
linkTitle = "%s"
date = %s
type="docs"
description = "%s"
isCJKLanguage = true
draft = false
[menu.main]
	weight = %d
+++

> 原文：[%s](%s)
>
> 收录时间：%s
`, menuInfo.MenuName, menuInfo.MenuName, date, "", index*10, menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))
		} else {
			_, err = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
date = %s
weight = %d
type="docs"
description = "%s"
isCJKLanguage = true
draft = false

+++

> 原文：[%s](%s)
>
> 收录时间：%s
`, menuInfo.MenuName, date, index*10, "", menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))
		}

		if err != nil {
			return fmt.Errorf("初始化%s文件时出错: %v", mdFp, err)
		}
	}
	return
}
