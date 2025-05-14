package phpdNext

import (
	"fmt"
	"github.com/before80/go/lg"
	"github.com/emirpasic/gods/queues/arrayqueue"
	"math/rand/v2"
	"sync"
	"time"
)

type MenuInfo struct {
	MenuName  string `json:"menu_name"`
	IsTopMenu int    `json:"is_top_menu"`
	Filename  string `json:"filename"`
	Url       string `json:"url"`
	Dir       string `json:"dir"`
	Weight    int    `json:"weight"`
	Retry     int    `json:"retry"`
}

var forMenuInfoQueueLock sync.Mutex
var forMenuInfoQueue = arrayqueue.New()
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

// PushWaitDealMenuInfoToQueue 放入栈中
func PushWaitDealMenuInfoToQueue(infoSlice []MenuInfo) {
	forMenuInfoQueueLock.Lock()
	defer forMenuInfoQueueLock.Unlock()
	for _, v := range infoSlice {
		forMenuInfoQueue.Enqueue(v)
	}
}

func ReversePushWaitDealMenuInfoToQueue(infoSlice []MenuInfo) {
	forMenuInfoQueueLock.Lock()
	defer forMenuInfoQueueLock.Unlock()
	ReverseMenuInfoSlice(infoSlice)
	for _, v := range infoSlice {
		forMenuInfoQueue.Enqueue(v)
	}
}

func GetNextMenuInfoFromQueue() (info MenuInfo, isEnd bool) {
	forMenuInfoQueueLock.Lock()
	defer forMenuInfoQueueLock.Unlock()
	elCount := forMenuInfoQueue.Size()
	lg.InfoToFileAndStdOut(fmt.Sprintf("待处理还有%d个\n", elCount))
	if IsFirstTimeGetMenuInfo {
		IsFirstTimeGetMenuInfo = false
	} else {
		time.Sleep(2 * time.Second)
	}

	v, ok := forMenuInfoQueue.Dequeue()
	lg.InfoToFile(fmt.Sprintf("v=%v,ok=%v\n", v, ok))
	if !ok || v == nil {
		return MenuInfo{}, true
	}

	return v.(MenuInfo), false
}
