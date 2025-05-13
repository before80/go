package pydPg

import (
	"encoding/json"
	"errors"
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
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	preDir := "pyd"
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
	curMenu, isEnd := pydNext.GetNextMenuInfoFromStack()
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
	lg.InfoToFile(fmt.Sprintf("线程%d发现网页%s存在子菜单=%v\n", threadIndex, curMenu.Url, subMenuInfos))

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

	_ = DoCopyAndPaste(threadIndex, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, curMenu.Url)

	_, err = pydTr.ReplaceMarkdownFileContent(absUniqueMdFilePath)
	if err != nil {
		panic(fmt.Errorf("线程%d在替换网页%s的内容对应的md文件时出现错误：%v", threadIndex, curMenu.Url, err))
	}
	lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
	err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 收录时间：")
	if err != nil {
		panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, curMenu.Url, err))
	}

	if len(subMenuInfos) > 0 {
		// 将 subMenuInfos 推入
		pydNext.PushWaitDealMenuInfoToStack(subMenuInfos)
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

func InitBarIndexMdFile(index int, barMenuInfo pydNext.MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename)
	return preInitMdFile(index, true, true, dir, barMenuInfo)
}

func InitSecondIndexMdFile(index int, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename)
	return preInitMdFile(index, false, true, dir, secondMenuInfo)
}

func InitThirdIndexMdFile(index int, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, thirdMenuInfo pydNext.MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename)
	return preInitMdFile(index, false, true, dir, thirdMenuInfo)
}

func InitSecondDetailPageMdFile(index int, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename)
	return preInitMdFile(index, false, false, dir, secondMenuInfo)
}

func InitThirdDetailPageMdFile(index int, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, thirdMenuInfo pydNext.MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename)
	return preInitMdFile(index, false, false, dir, thirdMenuInfo)
}

func InitFourthDetailPageMdFile(index int, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, thirdMenuInfo pydNext.MenuInfo, fourthMenuInfo pydNext.MenuInfo) (err error) {
	dir := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename)
	return preInitMdFile(index, false, false, dir, fourthMenuInfo)
}

func InsertBarMenuPageData(browserHwnd win.HWND, barMenuInfo pydNext.MenuInfo, page *rod.Page) (secondMenus []pydNext.MenuInfo, err error) {
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

	err = pg.DealUniqueMd(browserHwnd, barMenuInfo.Url, "barmenu")
	if err != nil {
		return nil, err
	}

	_, err = pydTr.ReplaceMarkdownFileContent(cfg.Default.UniqueMdFilepath)
	if err != nil {
		return nil, fmt.Errorf("在处理%s=%s时，替换出现错误：%v", "detailPage", barMenuInfo.Url, err)
	}

	indexMdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, "_index.md")
	err = pg.InsertAnyPageData(indexMdFp, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
}

func InsertSecondMenuPageData(browserHwnd win.HWND, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, page *rod.Page) (err error) {
	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", secondMenuInfo.Url, err)
	}

	err = pg.DealUniqueMd(browserHwnd, secondMenuInfo.Url, "second")
	if err != nil {
		return err
	}
	_, err = pydTr.ReplaceMarkdownFileContent(cfg.Default.UniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", "detailPage", secondMenuInfo.Url, err)
	}

	indexMdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, "_index.md")
	err = pg.InsertAnyPageData(indexMdFp, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
}

func InsertThirdMenuPageData(browserHwnd win.HWND, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, thirdMenuInfo pydNext.MenuInfo, page *rod.Page) (err error) {
	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, pydJs.GetDetailPageDataJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行GetDetailPageDataJs遇到错误：%v", thirdMenuInfo.Url, err)
	}

	err = pg.DealUniqueMd(browserHwnd, thirdMenuInfo.Url, "third")
	if err != nil {
		return err
	}
	_, err = pydTr.ReplaceMarkdownFileContent(cfg.Default.UniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", "detailPage", thirdMenuInfo.Url, err)
	}

	indexMdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename, "_index.md")
	err = pg.InsertAnyPageData(indexMdFp, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
}

func InsertSecondDetailPageData(browserHwnd win.HWND, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, page *rod.Page) (err error) {
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

	err = pg.DealUniqueMd(browserHwnd, secondMenuInfo.Url, "secondDetailPage")
	if err != nil {
		return err
	}
	_, err = pydTr.ReplaceMarkdownFileContent(cfg.Default.UniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", "detailPage", secondMenuInfo.Url, err)
	}

	mdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename+".md")
	err = pg.InsertAnyPageData(mdFp, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
}

func InsertThirdDetailPageData(browserHwnd win.HWND, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, thirdMenuInfo pydNext.MenuInfo, page *rod.Page) (err error) {
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

	err = pg.DealUniqueMd(browserHwnd, thirdMenuInfo.Url, "thirdDetailPage")
	if err != nil {
		return err
	}

	_, err = pydTr.ReplaceMarkdownFileContent(cfg.Default.UniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", "detailPage", thirdMenuInfo.Url, err)
	}

	mdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename+".md")
	err = pg.InsertAnyPageData(mdFp, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
}

func InsertFourthDetailPageData(browserHwnd win.HWND, barMenuInfo pydNext.MenuInfo, secondMenuInfo pydNext.MenuInfo, thirdMenuInfo pydNext.MenuInfo, fourthMenuInfo pydNext.MenuInfo, page *rod.Page) (err error) {
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

	err = pg.DealUniqueMd(browserHwnd, fourthMenuInfo.Url, "fourthDetailPage")
	if err != nil {
		return err
	}
	_, err = pydTr.ReplaceMarkdownFileContent(cfg.Default.UniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", "detailPage", fourthMenuInfo.Url, err)
	}

	mdFp := filepath.Join(contants.OutputFolderName, barMenuInfo.Filename, secondMenuInfo.Filename, thirdMenuInfo.Filename, fourthMenuInfo.Filename+".md")
	err = pg.InsertAnyPageData(mdFp, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
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

func GetThirdLevelMenu(secondMenuInfo pydNext.MenuInfo, page *rod.Page) (thirdMenuInfos []pydNext.MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("在第二级菜单%s中获取第三级菜单时遇到错误：%v", secondMenuInfo.Url, r)
		}
	}()

	page.MustNavigate(secondMenuInfo.Url)
	page.MustWaitLoad()
	return evalJsGetSubMenuInfos(page, "GetThirdMenusJs", pydJs.GetThirdMenusJs, secondMenuInfo.Url)
}

func GetFourthLevelMenu(thirdMenuInfo pydNext.MenuInfo, page *rod.Page) (fourthMenuInfos []pydNext.MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("在第二级菜单%s中获取第三级菜单时遇到错误：%v", thirdMenuInfo.Url, r)
		}
	}()

	page.MustNavigate(thirdMenuInfo.Url)
	page.MustWaitLoad()
	return evalJsGetSubMenuInfos(page, "GetFourthMenusJs", pydJs.GetFourthMenusJs, thirdMenuInfo.Url)
}

func evalJsGetSubMenuInfos(page *rod.Page, jsName, js, pageUrl string) (subMenuInfos []pydNext.MenuInfo, err error) {
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

func preInitMdFile(index int, isBar, useUnderlineIndexMd bool, dir string, menuInfo pydNext.MenuInfo) (err error) {
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
`, menuInfo.MenuName, menuInfo.MenuName, date, "", index*10, menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))
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
`, menuInfo.MenuName, date, index*10, "", menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))
		}

		if err != nil {
			return fmt.Errorf("初始化%s文件时出错: %v", mdFp, err)
		}
	}
	return
}
