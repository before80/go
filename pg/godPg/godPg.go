package godPg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/godJs"
	"github.com/before80/go/pg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

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

// GetAllStdPkgInfo 获取所有标准库的pkg信息
func GetAllStdPkgInfo(page *rod.Page, url string) (stdPkgMenuInfos []MenuInfo, err error) {
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
	err = pg.InsertAnyPageData(mdFp, "> 收录时间：")
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
