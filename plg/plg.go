package plg

import (
	"fmt"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/pydNext"
)

func getStr(menuInfos []pydNext.MenuInfo) string {
	str := ""
	curUrl := ""
	for i, menuInfo := range menuInfos {
		if i == 0 {
			str += fmt.Sprintf(" barfile=%s barmenu=%s", menuInfo.Filename, menuInfo.MenuName)
		}

		if i == 1 {
			str += fmt.Sprintf(" secondfile=%s secondmenu=%s", menuInfo.Filename, menuInfo.MenuName)
		}
		if i == 2 {
			str += fmt.Sprintf(" thirdfile=%s thirdmenu=%s", menuInfo.Filename, menuInfo.MenuName)
		}
		if i == 3 {
			str += fmt.Sprintf(" fourthfile=%s fourthmenu=%s", menuInfo.Filename, menuInfo.MenuName)
		}
		curUrl = menuInfo.Url
	}
	str += fmt.Sprintf(" curUrl=%s", curUrl)
	return str
}

func InfoToFileAndStdOut(process, step string, menuInfos ...pydNext.MenuInfo) {
	str := getStr(menuInfos)
	lg.InfoToFileAndStdOut(fmt.Sprintf("%s%s->%s\n", process, step, str))
}
