package godEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
	"github.com/before80/go/pg/godPg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"strconv"
	"time"
)

func Do() {
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
	//page.MustNavigate("https://pkg.go.dev/archive/tar@go1.24.2")
	//page.MustWaitLoad()
	//
	//time.Sleep(2000 * time.Second)

	browserHwnd = robotgo.GetHWND()
	var stdPkgMenuInfos []godPg.MenuInfo
	stdPkgMenuInfos, err = godPg.GetAllStdPkgInfo(page, "https://pkg.go.dev/std")
	//fmt.Println(stdPkgMenuInfos)

	for _, stdPkgMenuInfo := range stdPkgMenuInfos {
		err = godPg.InitStdPkgMdFile(stdPkgMenuInfo)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("初始化完成%s-%s\n", stdPkgMenuInfo.Filename, stdPkgMenuInfo.Url))

		if stdPkgMenuInfo.Url == "" {
			continue
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("准备插入数据%s-%s\n", stdPkgMenuInfo.Filename, stdPkgMenuInfo.Url))

		err = godPg.InsertPkgDetailPageData(browserHwnd, stdPkgMenuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("插入数据完成%s-%s\n", stdPkgMenuInfo.Filename, stdPkgMenuInfo.Url))

	}
	lg.InfoToFileAndStdOut("已全部完成\n")
	time.Sleep(1000 * time.Second)
}
