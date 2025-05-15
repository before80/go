package pydEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/cfg"
	"github.com/before80/go/entrance"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/pydNext"
	"github.com/before80/go/pg/pydPg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/spf13/cobra"
	"strconv"
	"sync"
	"time"
)

func Do(cmd *cobra.Command) {
	startTime := time.Now()
	lg.InfoToFileAndStdOut(fmt.Sprintf("开始时间：%v\n", startTime))

	var err error
	defer func() {
		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("%v", err))
		}
	}()
	defaults.ResetWith("show=true")
	var browser *rod.Browser
	var page *rod.Page
	// 打开浏览器
	browser, err = bs.GetBrowser(strconv.Itoa(0))
	defer browser.MustClose()
	// 创建新页面
	page = browser.MustPage()
	var barMenuInfos []pydNext.MenuInfo
	barMenuInfos, err = pydPg.GetBarMenus(page, cfg.Default.PydEntranceUrl)
	pydNext.PushWaitDealMenuInfoToQueue(barMenuInfos)
	_ = browser.Close()
	//fmt.Println("thirdPkgBaseInfos")
	threadNum, err := cmd.Flags().GetInt("thread-num")
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("获取线程数标志时出错：%v\n", err))
		return
	}

	bs.MyBrowserSlice = make([]bs.MyBrowser, threadNum)
	// 实例化多个 *rod.Browser 实例
	for j := 0; j < threadNum; j++ {
		browser1, err1 := bs.GetBrowser(strconv.Itoa(j))
		if err1 != nil {
			if len(bs.MyBrowserSlice) > 0 {
				for _, mb := range bs.MyBrowserSlice {
					if mb.Browser != nil {
						_ = mb.Browser.Close()
					}
				}
			}
		}
		entrance.OpenUniqueMdFile(j)
		bs.MyBrowserSlice[j] = bs.MyBrowser{Browser: browser1, Ok: true, Index: j}
	}

	var wg sync.WaitGroup
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go pydPg.DealWithMenuPageData(i, &wg)
	}
	wg.Wait()
	lg.InfoToFileAndStdOut(fmt.Sprintf("结束时间：%v\n", time.Now()))
	lg.InfoToFileAndStdOut(fmt.Sprintf("用时：%.2f分钟\n", time.Since(startTime).Minutes()))
	lg.InfoToFileAndStdOut("已完成处理！")

}
