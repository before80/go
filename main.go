package main

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
	"github.com/before80/go/pg"
	"github.com/before80/go/plg"
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"strconv"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("%v", err))
		}
	}()

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
	var barMenuInfos []pg.MenuInfo
	barMenuInfos, err = pg.GetBarMenus(page, "https://docs.python.org/zh-cn/3.13/index.html")

	var secondMenuInfos []pg.MenuInfo
	var thirdMenuInfos []pg.MenuInfo
	var fourthMenuInfos []pg.MenuInfo
	for i, barMenuInfo := range barMenuInfos {
		//if !slices.Contains([]string{"glossary"}, barMenuInfo.Filename) {
		//	continue
		//}

		plg.InfoToFileAndStdOut("bar", "正要处理", barMenuInfo)
		err = pg.InitBarIndexMdFile(i, barMenuInfo)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		plg.InfoToFileAndStdOut("bar", "初始化完成", barMenuInfo)

		secondMenuInfos, err = pg.InsertBarMenuPageData(browserHwnd, barMenuInfo, page)
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

			thirdMenuInfos, err = pg.GetThirdLevelMenu(secondMenuInfo, page)
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
				return
			}
			plg.InfoToFileAndStdOut("second", "获取第三级菜单完成", barMenuInfo, secondMenuInfo)

			// 存在第三级菜单的情况
			if len(thirdMenuInfos) > 0 {
				err = pg.InitSecondIndexMdFile(j, barMenuInfo, secondMenuInfo)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				plg.InfoToFileAndStdOut("second", "初始化完成1", barMenuInfo, secondMenuInfo)

				err = pg.InsertSecondMenuPageData(browserHwnd, barMenuInfo, secondMenuInfo, page)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				plg.InfoToFileAndStdOut("second", "插入数据完成1", barMenuInfo, secondMenuInfo)

				for k, thirdMenuInfo := range thirdMenuInfos {
					fourthMenuInfos, err = pg.GetFourthLevelMenu(thirdMenuInfo, page)
					if err != nil {
						lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
						return
					}
					plg.InfoToFileAndStdOut("third", "获取第四级菜单完成", barMenuInfo, secondMenuInfo, thirdMenuInfo)

					//存在第四级菜单的情况
					if len(fourthMenuInfos) > 0 {
						err = pg.InitThirdIndexMdFile(k, barMenuInfo, secondMenuInfo, thirdMenuInfo)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "初始化完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo)

						err = pg.InsertThirdMenuPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, page)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "插入数据完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo)

						for l, fourthMenuInfo := range fourthMenuInfos {
							err = pg.InitFourthDetailPageMdFile(l, barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo)
							if err != nil {
								lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
								return
							}
							plg.InfoToFileAndStdOut("fourth", "初始化完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo)

							err = pg.InsertFourthDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo, page)
							if err != nil {
								lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
								return
							}
							plg.InfoToFileAndStdOut("fourth", "插入数据完成1", barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo)
						}
					} else {
						//不存在第四级菜单的情况
						err = pg.InitThirdDetailPageMdFile(k, barMenuInfo, secondMenuInfo, thirdMenuInfo)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "初始化完成2", barMenuInfo, secondMenuInfo, thirdMenuInfo)

						err = pg.InsertThirdDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, page)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						plg.InfoToFileAndStdOut("third", "插入数据完成2", barMenuInfo, secondMenuInfo, thirdMenuInfo)
					}
				}
			} else {
				// 不存在第三级菜单的情况
				err = pg.InitSecondDetailPageMdFile(j, barMenuInfo, secondMenuInfo)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				plg.InfoToFileAndStdOut("second", "初始化完成2", barMenuInfo, secondMenuInfo)

				err = pg.InsertSecondDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, page)
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
