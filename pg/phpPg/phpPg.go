package phpPg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/phpdJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/pg"
	"github.com/before80/go/tr"
	"github.com/before80/go/tr/phpdTr"
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

type MenuInfo struct {
	MenuName string `json:"menu_name"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
	Index    int    `json:"index"`
}

// GetAllFirstMenuInfo 获取所有第一级菜单信息
func GetAllFirstMenuInfo(page *rod.Page, url string) (firstMenuInfos []MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("获取allStdPkgInfo时遇到错误：%v", r)
		}
	}()

	page.MustNavigate(url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(phpdJs.GetAllFirstMenuJs)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行phpdJs.GetAllFirstMenuJs遇到错误：%v", url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &firstMenuInfos)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", url, err)
	}
	return
}

func InitFirstMenuMdFile(browserHwnd win.HWND, firstMenuInfo MenuInfo, page *rod.Page) (err error) {
	var dir, curDir string
	var subMenuInfos []MenuInfo
	//isBar := false
	useUnderlineIndexMd := false
	prefixDirname := "php_"

	// 判断是否还有二级菜单
	page.MustNavigate(firstMenuInfo.Url)
	page.MustWaitLoad()
	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(phpdJs.InDetailPageGetMenuJs)
	if err != nil {
		return fmt.Errorf("在网页%s中执行phpdJs.InDetailPageGetMenuJs遇到错误：%v", firstMenuInfo.Url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", firstMenuInfo.Url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &subMenuInfos)
	if err != nil {
		return fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", firstMenuInfo.Url, err)
	}

	if len(subMenuInfos) > 0 {
		useUnderlineIndexMd = true
		dir = filepath.Join(contants.OutputFolderName, prefixDirname+firstMenuInfo.Filename)
		curDir = dir
		err = preInitMdFile(firstMenuInfo.Index, true, useUnderlineIndexMd, dir, firstMenuInfo)
		if err != nil {
			return
		}
		var result1 *proto.RuntimeRemoteObject
		result1, err = page.Eval(phpdJs.GetLayoutContentJs)
		if result1.Value.String() == "" {
			lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", firstMenuInfo.Url))
		} else {
			// 插入内容
			mdFilePath := filepath.Join(dir, "_index.md")
			err = InsertDetailPageData(browserHwnd, mdFilePath, firstMenuInfo, page)
			if err != nil {
				return
			}
		}
	} else {
		useUnderlineIndexMd = true
		dir = filepath.Join(contants.OutputFolderName, prefixDirname+firstMenuInfo.Filename)
		curDir = dir
		err = preInitMdFile(firstMenuInfo.Index, true, useUnderlineIndexMd, dir, firstMenuInfo)
		if err != nil {
			return
		}

		var result1 *proto.RuntimeRemoteObject
		result1, err = page.Eval(phpdJs.GetLayoutContentJs)
		vStr := strings.TrimSpace(result1.Value.String())
		//lg.InfoToFileAndStdOut(fmt.Sprintf("%v -> %s", vStr, firstMenuInfo.Url))
		if vStr == "" {
			lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", firstMenuInfo.Url))
		} else {
			// 插入内容
			mdFilePath := filepath.Join(dir, "_index.md")
			err = InsertDetailPageData(browserHwnd, mdFilePath, firstMenuInfo, page)
			if err != nil {
				return
			}
		}
	}

	if len(subMenuInfos) > 0 {
		culMenuLevel := 2
		err = DealSubMenuInfo(browserHwnd, subMenuInfos, curDir, culMenuLevel, page)
		if err != nil {
			return
		}
	}
	return
}

func DealSubMenuInfo(browserHwnd win.HWND, subMenuInfos []MenuInfo, curDir string, menuLevel int, page *rod.Page) (err error) {
	if len(subMenuInfos) > 0 {
		var subSubMenuInfos []MenuInfo
		useUnderlineIndexMd := false
		var dir, subCurDir string
		subMenuInfosLen := len(subMenuInfos)
		for i, subMenuInfo := range subMenuInfos {
			subSubMenuInfos = []MenuInfo{}
			lg.InfoToFileAndStdOut(fmt.Sprintf("正在处理第%d层(当前层还有%d个菜单待处理) %s - %s\n", menuLevel, subMenuInfosLen-i-1, subMenuInfo.MenuName, subMenuInfo.Url))

			// 判断是否还有二级菜单
			page.MustNavigate(subMenuInfo.Url)
			page.MustWaitLoad()
			var result *proto.RuntimeRemoteObject
			result, err = page.Eval(phpdJs.InDetailPageGetMenuJs)
			if err != nil {
				return fmt.Errorf("在网页%s中执行phpdJs.InDetailPageGetMenuJs遇到错误：%v", subMenuInfo.Url, err)
			}

			// 将结果序列化为 JSON 字节
			var jsonBytes []byte
			jsonBytes, err = json.Marshal(result.Value)
			if err != nil {
				return fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", subMenuInfo.Url, err)
			}

			// 将 JSON 数据反序列化到结构体中
			err = json.Unmarshal(jsonBytes, &subSubMenuInfos)
			if err != nil {
				return fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", subMenuInfo.Url, err)
			}

			if len(subSubMenuInfos) > 0 {
				useUnderlineIndexMd = true
				dir = filepath.Join(curDir, subMenuInfo.Filename)
				subCurDir = dir
				err = preInitMdFile(subMenuInfo.Index, false, useUnderlineIndexMd, dir, subMenuInfo)
				if err != nil {
					return
				}
				var result1 *proto.RuntimeRemoteObject
				result1, err = page.Eval(phpdJs.GetLayoutContentJs)
				vStr := strings.TrimSpace(result1.Value.String())
				//lg.InfoToFileAndStdOut(fmt.Sprintf("%v -> %s", vStr, subMenuInfo.Url))
				if vStr == "" {
					lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", subMenuInfo.Url))
				} else {
					// 插入内容
					mdFilePath := filepath.Join(dir, "_index.md")
					err = InsertDetailPageData(browserHwnd, mdFilePath, subMenuInfo, page)
					if err != nil {
						return
					}
				}
			} else {
				useUnderlineIndexMd = false
				dir = curDir
				subCurDir = dir
				err = preInitMdFile(subMenuInfo.Index, false, useUnderlineIndexMd, dir, subMenuInfo)
				if err != nil {
					return
				}

				var result1 *proto.RuntimeRemoteObject
				result1, err = page.Eval(phpdJs.GetLayoutContentJs)
				vStr := strings.TrimSpace(result1.Value.String())
				//lg.InfoToFileAndStdOut(fmt.Sprintf("%v -> %s", vStr, subMenuInfo.Url))
				if vStr == "" {
					lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", subMenuInfo.Url))
				} else {
					// 插入内容
					mdFilePath := filepath.Join(dir, subMenuInfo.Filename+".md")
					err = InsertDetailPageData(browserHwnd, mdFilePath, subMenuInfo, page)
					if err != nil {
						return
					}
				}
			}

			if err != nil {
				return
			}

			if len(subSubMenuInfos) > 0 {
				curMenuLevel := menuLevel + 1
				err = DealSubMenuInfo(browserHwnd, subSubMenuInfos, subCurDir, curMenuLevel, page)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func InsertDetailPageData(browserHwnd win.HWND, mdFilePath string, menuInfo MenuInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入detailPage=%s数据时遇到错误：%v", menuInfo.Url, r)
		}
	}()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, phpdJs.ReplaceJs))
	if err != nil {
		return fmt.Errorf("在网页%s中执行phpdJs.ReplaceJs遇到错误：%v", menuInfo.Url, err)
	}

	err = dealUniqueMd(browserHwnd, menuInfo.Url, "detailPage")
	if err != nil {
		return err
	}

	err = pg.InsertAnyPageData(mdFilePath)
	return
}

func dealUniqueMd(browserHwnd win.HWND, curUrl, step string) (err error) {
	uniqueMdFilepath := cfg.Default.UniqueMdFilepath
	// 获取文件名
	spSlice := strings.Split(uniqueMdFilepath, "\\")
	mdFilename := spSlice[len(spSlice)-1]

	// 清空唯一共用的markdown文件的文件内容
	err = tr.TruncFileContent(uniqueMdFilepath)
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
	time.Sleep(200 * time.Microsecond)
	wind.SelectAllAndDelete(typoraHwnd)
	wind.CtrlV(typoraHwnd)
	time.Sleep(time.Duration(cfg.Default.WaitTyporaCopiedToSaveSeconds) * time.Second)
	wind.CtrlS(typoraHwnd)
	time.Sleep(time.Duration(cfg.Default.WaitTyporaSaveSeconds) * time.Second)
	robotgo.CloseWindow()
	time.Sleep(time.Duration(cfg.Default.WaitTyporaCloseSeconds) * time.Second)

	_, err = phpdTr.ReplaceMarkdownFileContent(uniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", step, curUrl, err)
	}
	return nil
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
