package goThirdPkgIndexEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/goThirdPkgIndexNext"
	"github.com/before80/go/pg/goThirdPkgIndexPg"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/spf13/cobra"
	"strconv"
	"sync"
)

func Do(cmd *cobra.Command) {
	defaults.ResetWith("show=true")

	goThirdPkgIndexNext.ReverseBaseInfoSlice(goThirdPkgIndexNext.ThirdPkgBaseInfos)
	goThirdPkgIndexNext.PushWaitDealBaseInfoToStack(goThirdPkgIndexNext.ThirdPkgBaseInfos)

	//fmt.Println("thirdPkgBaseInfos")
	threadNum, err := cmd.Flags().GetInt("thread-num")
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("获取线程数标志时出错：%v\n", err))
		return
	}

	bs.MyBrowserSlice = make([]bs.MyBrowser, threadNum)
	// 实例化多个 *rod.Browser 实例
	for j := 0; j < threadNum; j++ {
		browser, err1 := bs.GetBrowser(strconv.Itoa(j))
		if err1 != nil {
			if len(bs.MyBrowserSlice) > 0 {
				for _, mb := range bs.MyBrowserSlice {
					if mb.Browser != nil {
						_ = mb.Browser.Close()
					}
				}
			}
		}
		bs.MyBrowserSlice[j] = bs.MyBrowser{Browser: browser, Ok: true, Index: j}
	}

	var wg sync.WaitGroup
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go goThirdPkgIndexPg.DealWithPkgBaseInfo(i, &wg)
	}
	wg.Wait()

	goThirdPkgIndexPg.SortAndGenAllPkgInfos()
	err = goThirdPkgIndexPg.TruncWriteLinesToFile()
	if err != nil {
		panic(err)
	}
	lg.InfoToFile(fmt.Sprintf("AllPkgInfos=%v\n", goThirdPkgIndexNext.AllPkgInfos))
	goThirdPkgIndexNext.InitWaitHandlePkgInfoCount()
	goThirdPkgIndexNext.ReversePkgInfoSlice(goThirdPkgIndexNext.AllPkgInfos)
	goThirdPkgIndexNext.PushWaitDealPkgInfoToStack(goThirdPkgIndexNext.AllPkgInfos)

	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go goThirdPkgIndexPg.DealWithPkgPageData(i, &wg)
	}
	wg.Wait()

	// 关闭打开的浏览器
	for _, myBrowser := range bs.MyBrowserSlice {
		if myBrowser.Browser != nil {
			_ = myBrowser.Browser.Close()
		}
	}

	lg.InfoToFileAndStdOut("已完成处理！")
}
