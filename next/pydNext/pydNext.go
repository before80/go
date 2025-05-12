package pydNext

import (
	"fmt"
	"github.com/before80/go/lg"
	"github.com/emirpasic/gods/stacks/arraystack"
	"math/rand/v2"
	"sync"
	"time"
)

type MenuInfo struct {
	MenuName    string `json:"menu_name"`
	IsTopMenu   int    `json:"is_top_menu"`
	TopMenuName string `json:"top_menu_name"`
	Filename    string `json:"filename"`
	Url         string `json:"url"`
	Dir         string `json:"dir"`
	Weight      int    `json:"weight"`
	Desc        string `json:"desc"`
}

var forMenuInfoStackLock sync.Mutex
var forMenuInfoStack = arraystack.New()
var IsFirstTimeGetMenuInfo = true

// ReverseMenuInfoSlice 倒序排列
func ReverseMenuInfoSlice(infoSlice []MenuInfo) {
	for i, j := 0, len(infoSlice)-1; i < j; i, j = i+1, j-1 {
		infoSlice[i], infoSlice[j] = infoSlice[j], infoSlice[i]
	}
}

// ShuffleMenuInfoSlice 乱序排列
func ShuffleMenuInfoSlice(infoSlice []MenuInfo) {
	for i := len(infoSlice) - 1; i > 0; i-- {
		// 生成一个从 0 到 i 的随机整数 j
		j := rand.IntN(i + 1)
		// 交换 infoSlice[i] 和 infoSlice[j]
		infoSlice[i], infoSlice[j] = infoSlice[j], infoSlice[i]
	}
}

// PushWaitDealMenuInfoToStack 放入栈中
func PushWaitDealMenuInfoToStack(infoSlice []MenuInfo) {
	forMenuInfoStackLock.Lock()
	defer forMenuInfoStackLock.Unlock()
	for _, v := range infoSlice {
		forMenuInfoStack.Push(v)
	}
}

func ReversePushWaitDealMenuInfoToStack(infoSlice []MenuInfo) {
	forMenuInfoStackLock.Lock()
	defer forMenuInfoStackLock.Unlock()
	ReverseMenuInfoSlice(infoSlice)
	for _, v := range infoSlice {
		forMenuInfoStack.Push(v)
	}
}

func GetNextMenuInfoFromStack() (info MenuInfo, isEnd bool) {
	forMenuInfoStackLock.Lock()
	defer forMenuInfoStackLock.Unlock()
	//elCount := forMenuInfoStack.Size()
	//lg.InfoToFile(fmt.Sprintf("elCount=%d\n", elCount))
	if IsFirstTimeGetMenuInfo {
		IsFirstTimeGetMenuInfo = false
	} else {
		time.Sleep(2 * time.Second)
	}

	v, ok := forMenuInfoStack.Pop()
	lg.InfoToFile(fmt.Sprintf("v=%v,ok=%v\n", v, ok))
	if !ok || v == nil {
		return MenuInfo{}, true
	}

	return v.(MenuInfo), false
}
