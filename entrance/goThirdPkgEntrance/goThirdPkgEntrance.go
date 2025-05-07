package goThirdPkgEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/contants"
	"github.com/before80/go/lg"
	"github.com/before80/go/pg"
	"github.com/before80/go/pg/goThirdPkgPg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/go-vgo/robotgo"
	"github.com/spf13/cobra"
	"github.com/tailscale/win"
	"path/filepath"
	"strconv"
)

func Do(cmd *cobra.Command) {
	var err error
	defer func() {
		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("%v", err))
		}
	}()
	defaults.ResetWith("show=true")
	_ = err
	var browser *rod.Browser
	var page *rod.Page
	var browserHwnd win.HWND
	_ = browserHwnd

	// 打开浏览器
	browser, err = bs.GetBrowser(strconv.Itoa(0))
	defer browser.MustClose()
	// 创建新页面
	page = browser.MustPage()
	browserHwnd = robotgo.GetHWND()

	url, err := cmd.Flags().GetString("url")
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("获取url标志时出错:%v\n", err))
		return
	}
	_ = url
	pkgName, err := cmd.Flags().GetString("pkg-name")
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("获取pkg-name包名标志时出错:%v\n", err))
		return
	}
	_ = pkgName

	weight, err := cmd.Flags().GetInt("weight")
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("获取weight权重标志时出错:%v\n", err))
		return
	}

	_ = weight
	baseDirname := "go_third_pkg"
	_ = pg.GenBarIndexMdFile(filepath.Join(contants.OutputFolderName, baseDirname), "Go 第三方包", 100)

	menuInfo := goThirdPkgPg.MenuInfo{MenuName: pkgName, Filename: pkgName, Url: url, Index: weight}
	err = goThirdPkgPg.InitThirdPkgMdFile(baseDirname, menuInfo)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
		return
	}

	err = goThirdPkgPg.InsertPkgDetailPageData(browserHwnd, baseDirname, menuInfo, page)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	} else {
		lg.InfoToFileAndStdOut("已完成处理！")
	}
}
