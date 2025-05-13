package godPg

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
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
		return nil, fmt.Errorf("在网页%s中执行GetBarMenusJs遇到错误：%v", url, err)
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

	defer func() {
		_ = page.Close()
	}()
	var pageTitle, chromePageWindowTitle string
	var fpDst string
	var result *proto.RuntimeRemoteObject
	uniqueMdFilename := "do" + strconv.Itoa(threadIndex) + ".md"
	relUniqueMdFilePath := filepath.Join("markdown", uniqueMdFilename)
	typoraWindowTitle := uniqueMdFilename + " - Typora"
	_, _ = pg.CreateFileIfNotExists(relUniqueMdFilePath)
	absUniqueMdFilePath, _ := filepath.Abs(relUniqueMdFilePath)
LabelForContinue:
	_ = tr.TruncFileContent(relUniqueMdFilePath)
	date := time.Now().Format(time.RFC3339)
	curMenu, isEnd := godNext.GetNextMenuInfoFromStack()
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
`, curMenu.MenuName, curMenu.MenuName, date, curMenu.Desc, curMenu.Weight, curMenu.Url, curMenu.Url, fmt.Sprintf("`%s`", date)))
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
`, curMenu.MenuName, date, curMenu.Weight, curMenu.Desc, curMenu.Url, curMenu.Url, fmt.Sprintf("`%s`", date)))
	}

	_ = mdF.Close()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, godJs.ReplaceJs))
	if err != nil {
		panic(fmt.Errorf("线程%d在网页%s中执行godJs.ReplaceJs遇到错误：%v", threadIndex, curMenu.Url, err))
	}
	// 获取当前网页的title，在后面会用来查找该网页所在窗口的操作句柄
	result, _ = page.Eval(`() => { return document.title }`)
	pageTitle = result.Value.String()
	chromePageWindowTitle = pageTitle + " - Google Chrome"

	// 再次清空
	_ = tr.TruncFileContent(relUniqueMdFilePath)

	_ = DoCopyAndPaste(threadIndex, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, curMenu.Url)

	lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
	err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 收录时间：")
	if err != nil {
		panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, curMenu.Url, err))
	}

	goto LabelForContinue
}

var copyPasteLock sync.Mutex

func DoCopyAndPaste(threadIndex int, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, url string) (err error) {
	copyPasteLock.Lock()
	defer copyPasteLock.Unlock()
	var typoraHwnd win.HWND
	browserHwnd := robotgo.FindWindow(chromePageWindowTitle)

	//_ = wind.OpenTypora(absUniqueMdFilePath)
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
	//_ = win.SendMessage(typoraHwnd, win.WM_CLOSE, 0, 0)
	_ = win.SendMessage(typoraHwnd, win.WM_SYSCOMMAND, win.SC_MINIMIZE, 0)
	//time.Sleep(time.Duration(cfg.Default.WaitTyporaCloseSeconds) * time.Second)
	return nil
}

type MenuInfo struct {
	MenuName             string   `json:"menu_name"`
	Filename             string   `json:"filename"`
	Url                  string   `json:"url"`
	Desc                 string   `json:"desc"`
	IsTop                int      `json:"is_top"`
	Index                int      `json:"index"`
	PFilename            string   `json:"p_filename"`
	ChildrenMenuFilename []string `json:"children"`
}

// GetAllStdPkgInfo0 获取所有标准库的pkg信息
func GetAllStdPkgInfo0(page *rod.Page, url string) (stdPkgMenuInfos []MenuInfo, err error) {
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
		return nil, fmt.Errorf("在网页%s中执行GetBarMenusJs遇到错误：%v", url, err)
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

func InitStdPkgMdFile(pkgMenu MenuInfo) (err error) {
	var dir string
	//isBar := false
	useUnderlineIndexMd := false
	baseDirname := "go_std_pkg"
	if pkgMenu.IsTop == 1 {
		if len(pkgMenu.ChildrenMenuFilename) > 0 {
			useUnderlineIndexMd = true
			dir = filepath.Join(contants.OutputFolderName, baseDirname, pkgMenu.Filename)
		} else {
			dir = filepath.Join(contants.OutputFolderName, baseDirname)
		}
	} else {
		dir = filepath.Join(contants.OutputFolderName, baseDirname, pkgMenu.PFilename)
	}
	return preInitMdFile(false, useUnderlineIndexMd, dir, pkgMenu)
}

func InsertPkgDetailPageData(browserHwnd win.HWND, pkgMenu MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入detailPage=%s数据时遇到错误：%v", pkgMenu.Url, r)
		}
	}()

	page.MustNavigate(pkgMenu.Url + "?GOOS=windows")
	page.MustWaitLoad()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, godJs.ReplaceJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行godJs.ReplaceJs遇到错误：%v", pkgMenu.Url, err)
	}

	err = pg.DealUniqueMd(browserHwnd, pkgMenu.Url, "detailPage")
	if err != nil {
		return err
	}
	var mdFp string
	baseDirname := "go_std_pkg"
	if pkgMenu.IsTop == 1 {
		if len(pkgMenu.ChildrenMenuFilename) > 0 {
			mdFp = filepath.Join(contants.OutputFolderName, baseDirname, pkgMenu.Filename, "_index.md")
		} else {
			mdFp = filepath.Join(contants.OutputFolderName, baseDirname, pkgMenu.Filename+".md")
		}
	} else {
		mdFp = filepath.Join(contants.OutputFolderName, baseDirname, pkgMenu.PFilename, pkgMenu.Filename+".md")
	}
	err = pg.InsertAnyPageData(mdFp, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
}

func preInitMdFile(isBar, useUnderlineIndexMd bool, dir string, menuInfo MenuInfo) (err error) {
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
			return fmt.Errorf("初始化%s文件时出错: %v", mdFp, err)
		}
	}
	return
}
