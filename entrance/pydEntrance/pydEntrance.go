package pydEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/pydNext"
	"github.com/before80/go/pg/pydPg"
	"github.com/before80/go/plg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/go-vgo/robotgo"
	"github.com/spf13/cobra"
	"github.com/tailscale/win"
	"strconv"
	"sync"
	"time"
)

func Do0() {
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
	var barMenuInfos []pydNext.MenuInfo
	barMenuInfos, err = pydPg.GetBarMenus(page, "https://docs.python.org/zh-cn/3.13/index.html")

	var secondMenuInfos []pydNext.MenuInfo
	var thirdMenuInfos []pydNext.MenuInfo
	var fourthMenuInfos []pydNext.MenuInfo
	for i, barMenuInfo := range barMenuInfos {
		//if !slices.Contains([]string{"glossary"}, barMenuInfo.Filename) {
		//	continue
		//}

		plg.InfoToFileAndStdOut("bar", "正要处理", barMenuInfo)
		err = pydPg.InitBarIndexMdFile(i, barMenuInfo)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		plg.InfoToFileAndStdOut("bar", "初始化完成", barMenuInfo)

		secondMenuInfos, err = pydPg.InsertBarMenuPageData(browserHwnd, barMenuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		plg.InfoToFileAndStdOut("bar", "插入数据完成", barMenuInfo)

		if len(secondMenuInfos) <= 0 {
			continue
		}

		plg.InfoToFileAndStdOut("second", "处理二级菜单中", barMenuInfo)

		for j, secondMenuInfo := range secondMenuInfos {
			//if barMenuInfo.Filename == "library" &&
			//	slices.Contains([]string{"constants", "allos", "binary", "crypto", "datatypes", "fileformats", "filesys", "functional", "numeric", "persistence", "text", "constants", "exceptions", "functions", "intro"}, secondMenuInfo.Filename) {
			//	continue
			//}

			thirdMenuInfos, err = pydPg.GetThirdLevelMenu(secondMenuInfo, page)
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
				return
			}
			plg.InfoToFileAndStdOut("second", "获取第三级菜单完成", barMenuInfo, secondMenuInfo)

			// 存在第三级菜单的情况
			if len(thirdMenuInfos) > 0 {
				err = pydPg.InitSecondIndexMdFile(j, barMenuInfo, secondMenuInfo)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				plg.InfoToFileAndStdOut("second", "初始化完成1", barMenuInfo, secondMenuInfo)

				err = pydPg.InsertSecondMenuPageData(browserHwnd, barMenuInfo, secondMenuInfo, page)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				plg.InfoToFileAndStdOut("second", "插入数据完成1", barMenuInfo, secondMenuInfo)

				for k, thirdMenuInfo := range thirdMenuInfos {
					fourthMenuInfos, err = pydPg.GetFourthLevelMenu(thirdMenuInfo, page)
					if err != nil {
						lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
						return
					}
					plg.InfoToFileAndStdOut("third", "获取第四级菜单完成", barMenuInfo, secondMenuInfo, thirdMenuInfo)

					//存在第四级菜单的情况
					if len(fourthMenuInfos) > 0 {
						err = pydPg.InitThirdIndexMdFile(k, barMenuInfo, secondMenuInfo, thirdMenuInfo)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "初始化完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo)

						err = pydPg.InsertThirdMenuPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, page)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "插入数据完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo)

						for l, fourthMenuInfo := range fourthMenuInfos {
							err = pydPg.InitFourthDetailPageMdFile(l, barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo)
							if err != nil {
								lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
								return
							}
							plg.InfoToFileAndStdOut("fourth", "初始化完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo)

							err = pydPg.InsertFourthDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo, page)
							if err != nil {
								lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
								return
							}
							plg.InfoToFileAndStdOut("fourth", "插入数据完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo)
						}
					} else {
						//不存在第四级菜单的情况
						err = pydPg.InitThirdDetailPageMdFile(k, barMenuInfo, secondMenuInfo, thirdMenuInfo)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "初始化完成2", barMenuInfo, secondMenuInfo, thirdMenuInfo)

						err = pydPg.InsertThirdDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, page)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "插入数据完成2", barMenuInfo, secondMenuInfo, thirdMenuInfo)
					}
				}
			} else {
				// 不存在第三级菜单的情况
				err = pydPg.InitSecondDetailPageMdFile(j, barMenuInfo, secondMenuInfo)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				plg.InfoToFileAndStdOut("second", "初始化完成2", barMenuInfo, secondMenuInfo)

				err = pydPg.InsertSecondDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, page)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				plg.InfoToFileAndStdOut("second", "插入数据完成2", barMenuInfo, secondMenuInfo)
			}
		}
	}
	lg.InfoToFileAndStdOut("已全部完成")
	_ = browser.Close()
}

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
	barMenuInfos, err = pydPg.GetBarMenus(page, "https://docs.python.org/zh-cn/3.13/index.html")
	pydNext.PushWaitDealMenuInfoToStack(barMenuInfos)
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
