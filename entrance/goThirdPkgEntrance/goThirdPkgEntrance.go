package goThirdPkgEntrance

import (
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/lg"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/spf13/cobra"
	"strconv"
	"sync"
)

func Do(cmd *cobra.Command) {
	defaults.ResetWith("show=true")

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

	}
	wg.Wait()
}
