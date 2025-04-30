package phpdEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
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

	//page.MustNavigate("https://www.php.net/manual/zh/language.control-structures.php")
	//page.MustWaitLoad()
	//page.Eval(fmt.Sprintf(`() => { %s }`, phpdJs.ReplaceJs))
	//
	//var result1 *proto.RuntimeRemoteObject
	//result1, err = page.Eval(phpdJs.GetLayoutContentJs)
	//fmt.Printf("%q\n", result1.Value)
	//fmt.Printf("%q\n", result1.Value.String())
	//time.Sleep(200 * time.Second)
	//return

	var firstMenuInfos []phpPg.MenuInfo
	firstMenuInfos, err = phpPg.GetAllFirstMenuInfo(page, "https://www.php.net/manual/zh/index.php")
	firstMenuInfosLen := len(firstMenuInfos)
	fmt.Println("firstMenuInfos=", firstMenuInfos)
	for i, firstMenuInfo := range firstMenuInfos {
		//if slices.Contains([]string{"copyright", "getting-started", "install", "preface"}, firstMenuInfo.Filename) {
		//	continue
		//}
		lg.InfoToFileAndStdOut(fmt.Sprintf("正在处理第%d层(当前层还有%d个菜单待处理) %s - %s\n", 1, firstMenuInfosLen-i-1, firstMenuInfo.MenuName, firstMenuInfo.Url))
		err = phpPg.InitFirstMenuMdFile(browserHwnd, firstMenuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v\n", err), 3)
			return
		}
	}
	lg.InfoToFileAndStdOut("已全部处理完成！")
}
