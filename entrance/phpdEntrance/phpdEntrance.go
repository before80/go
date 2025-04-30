package phpdEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
	"github.com/before80/go/pg/phpPg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/tailscale/win"
	"strconv"
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
	var firstMenuInfos []phpPg.MenuInfo
	firstMenuInfos, err = phpPg.GetAllFirstMenuInfo(page, "https://www.php.net/manual/zh/index.php")

	for _, firstMenuInfo := range firstMenuInfos {
		err = phpPg.InitFirstMenuMdFile(browserHwnd, firstMenuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v\n", err), 3)
			return
		}
	}
	lg.InfoToFileAndStdOut("已全部处理完成！")
}
