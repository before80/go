package goThirdPkgIndexNext

import (
	"bufio"
	"fmt"
	"github.com/before80/go/lg"
	"github.com/emirpasic/gods/stacks/arraystack"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ThirdPkgBaseInfo struct {
	PkgName string
	Url     string
	Weight  int
	CanUse  bool
	HadUse  bool
}

var ThirdPkgBaseInfos []ThirdPkgBaseInfo
var initWaitHandleBaseInfoCount = 0

func init() {
	f, err := os.OpenFile("config/go_third_pkg.txt", os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			info := strings.Split(line, "|")
			if len(info) != 3 {
				panic(fmt.Sprintf("config/go_third_pkg.txt中%q写法有错，未发现两个|", line))
			}
			weight, _ := strconv.Atoi(info[2])
			ThirdPkgBaseInfos = append(ThirdPkgBaseInfos, ThirdPkgBaseInfo{PkgName: strings.TrimSpace(info[0]), Url: strings.TrimSpace(info[1]), Weight: weight})
		}
	}
	initWaitHandleBaseInfoCount = len(ThirdPkgBaseInfos)
}

var forBaseInfoStackLock sync.Mutex
var forBaseInfoStack = arraystack.New()
var IsFirstTimeGetBaseInfo = true

// ReverseBaseInfoSlice 倒序排列
func ReverseBaseInfoSlice(infoSlice []ThirdPkgBaseInfo) {
	for i, j := 0, len(infoSlice)-1; i < j; i, j = i+1, j-1 {
		infoSlice[i], infoSlice[j] = infoSlice[j], infoSlice[i]
	}
}

// ShuffleBaseInfoSlice 乱序排列
func ShuffleBaseInfoSlice(infoSlice []ThirdPkgBaseInfo) {
	for i := len(infoSlice) - 1; i > 0; i-- {
		// 生成一个从 0 到 i 的随机整数 j
		j := rand.IntN(i + 1)
		// 交换 infoSlice[i] 和 infoSlice[j]
		infoSlice[i], infoSlice[j] = infoSlice[j], infoSlice[i]
	}
}

// PushWaitDealBaseInfoToStack 放入栈中
func PushWaitDealBaseInfoToStack(infoSlice []ThirdPkgBaseInfo) {
	for _, v := range infoSlice {
		forBaseInfoStack.Push(v)
	}
}

func GetNextBaseInfoFromStack() (index int, info ThirdPkgBaseInfo, isEnd bool) {
	forBaseInfoStackLock.Lock()
	defer forBaseInfoStackLock.Unlock()
	elCount := forBaseInfoStack.Size()
	lg.InfoToFile(fmt.Sprintf("elCount=%d\n", elCount))
	if IsFirstTimeGetBaseInfo {
		IsFirstTimeGetBaseInfo = false
	} else {
		time.Sleep(2 * time.Second)
	}

	v, ok := forBaseInfoStack.Pop()
	lg.InfoToFile(fmt.Sprintf("v=%v,ok=%v\n", v, ok))
	if !ok || v == nil {
		return initWaitHandleBaseInfoCount - elCount, ThirdPkgBaseInfo{}, true
	}

	return initWaitHandleBaseInfoCount - elCount, v.(ThirdPkgBaseInfo), false
}

var isFirstTimeGetBaseInfo = true
var currentBaseInfoIndex = -1
var forBaseInfoSliceLock sync.Mutex

func GetNextBaseInfoFromSlice() (info ThirdPkgBaseInfo, err error) {
	forBaseInfoSliceLock.Lock()
	defer forBaseInfoSliceLock.Unlock()
	l := len(ThirdPkgBaseInfos)
	var index1 int
	totalTempNum := 0
LabelForContinue:
	totalTempNum++
	// 当循环两轮仍没有可用用户的情况下
	if totalTempNum >= 2*l {
		return ThirdPkgBaseInfo{}, fmt.Errorf("获取不到可用信息")
	}

	if isFirstTimeGetBaseInfo {
		isFirstTimeGetBaseInfo = false
		currentBaseInfoIndex = -1
	}

	index1 = currentBaseInfoIndex + 1

	if index1 >= l {
		currentBaseInfoIndex = 0

		if !ThirdPkgBaseInfos[0].CanUse {
			goto LabelForContinue
		}

		if ThirdPkgBaseInfos[0].HadUse {
			goto LabelForContinue
		}

		// 设置为已经使用
		ThirdPkgBaseInfos[0].HadUse = true
		return ThirdPkgBaseInfos[0], nil
	} else {
		currentBaseInfoIndex = index1
		// 若该已经不能登录
		if !ThirdPkgBaseInfos[currentBaseInfoIndex].CanUse {
			goto LabelForContinue
		}

		if ThirdPkgBaseInfos[currentBaseInfoIndex].HadUse {
			goto LabelForContinue
		}

		// 设置为已经使用
		ThirdPkgBaseInfos[currentBaseInfoIndex].HadUse = true
		return ThirdPkgBaseInfos[currentBaseInfoIndex], nil
	}
}

type PkgInfo struct {
	PkgName            string `json:"pkg_name"`
	Filename           string `json:"filename"`
	Url                string `json:"url"`
	Dir                string `json:"dir"`
	Weight             int    `json:"weight"`
	NeedPreCreateIndex int    `json:"need_pre_create_index"`
	Desc               string `json:"desc"`
}

var AllPkgInfos []PkgInfo
var initWaitHandlePkgInfoCount int
var forPkgInfoStackLock sync.Mutex
var forPkgInfoStack = arraystack.New()
var IsFirstTimeGetPkgInfo = true

func InitWaitHandlePkgInfoCount() {
	initWaitHandlePkgInfoCount = len(AllPkgInfos)
}

func ReversePkgInfoSlice(infoSlice []PkgInfo) {
	for i, j := 0, len(infoSlice)-1; i < j; i, j = i+1, j-1 {
		infoSlice[i], infoSlice[j] = infoSlice[j], infoSlice[i]
	}
}

func PushWaitDealPkgInfoToStack(infoSlice []PkgInfo) {
	for _, v := range infoSlice {
		forPkgInfoStack.Push(v)
	}
}

func GetNextPkgInfoFromStack() (index int, info PkgInfo, isEnd bool) {
	forPkgInfoStackLock.Lock()
	defer forPkgInfoStackLock.Unlock()
	elCount := forPkgInfoStack.Size()
	lg.InfoToFile(fmt.Sprintf("elCount=%d\n", elCount))
	if IsFirstTimeGetPkgInfo {
		IsFirstTimeGetPkgInfo = false
	}

	v, ok := forPkgInfoStack.Pop()
	lg.InfoToFile(fmt.Sprintf("v=%v,ok=%v\n", v, ok))
	if !ok || v == nil {
		return initWaitHandlePkgInfoCount - elCount, PkgInfo{}, true
	}

	return initWaitHandlePkgInfoCount - elCount, v.(PkgInfo), false
}
