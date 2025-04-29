package phpdEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
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

	_ = page

}
