package mysqldEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
	"github.com/before80/go/pg/mysqldPg"
	"github.com/before80/go/pg/phpPg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"strconv"
)

func Do() {
	defer phpPg.CloseInitFiles()
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
	// 打开浏览器
	browser, err = bs.GetBrowser(strconv.Itoa(0))
	defer browser.MustClose()
	// 创建新页面
	page = browser.MustPage()
	browserHwnd = robotgo.GetHWND()

	var menuInfos []mysqldPg.MenuInfo
	menuInfos, err = mysqldPg.GetAllMenuInfo(page, "https://dev.mysql.com/doc/refman/8.0/en/")
	menuInfosLen := len(menuInfos)
	fmt.Println("menuInfos=", menuInfos)
	for i, menuInfo := range menuInfos {
		//if !slices.Contains([]string{""}, menuInfo.Filename) {
		//	continue
		//}
		surplus := menuInfosLen - i - 1
		err = mysqldPg.DealMenuMdFile(surplus, browserHwnd, "mysql", menuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v\n", err), 3)
			return
		}
	}
	mysqldPg.CloseInitFiles()
	lg.InfoToFileAndStdOut("已全部处理完成！")
}
