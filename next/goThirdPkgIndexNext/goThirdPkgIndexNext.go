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
)

type ThirdPkgBaseInfo struct {
	PkgName string
	Url     string
	Weight  int
	CanUse  bool
	HadUse  bool
}

var ThirdPkgBaseInfos []ThirdPkgBaseInfo
var initWaitHandleInfoCount = 0

func init() {
	f, err := os.OpenFile("config/go_third_pkg.txt", os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			info := strings.Split(line, "|")
			if len(info) != 3 {
				panic(fmt.Sprintf("config/go_third_pkg.txt中%q写法有错，未发现两个|", line))
			}
			weight, _ := strconv.Atoi(info[2])
			ThirdPkgBaseInfos = append(ThirdPkgBaseInfos, ThirdPkgBaseInfo{PkgName: strings.TrimSpace(info[0]), Url: strings.TrimSpace(info[1]), Weight: weight})
		}
	}
	initWaitHandleInfoCount = len(ThirdPkgBaseInfos)
}

var forInfoStackLock sync.Mutex
var forInfoStack = arraystack.New()
var IsFirstTimeGetInfo = true

// ReverseInfoSlice 倒序排列
func ReverseInfoSlice(infoSlice []ThirdPkgBaseInfo) {
	for i, j := 0, len(infoSlice)-1; i < j; i, j = i+1, j-1 {
		infoSlice[i], infoSlice[j] = infoSlice[j], infoSlice[i]
	}
}

// ShuffleInfoSlice 乱序排列
func ShuffleInfoSlice(infoSlice []ThirdPkgBaseInfo) {
	for i := len(infoSlice) - 1; i > 0; i-- {
		// 生成一个从 0 到 i 的随机整数 j
		j := rand.IntN(i + 1)
		// 交换 infoSlice[i] 和 infoSlice[j]
		infoSlice[i], infoSlice[j] = infoSlice[j], infoSlice[i]
	}
}

// PushWaitDealInfoToStack 放入栈中
func PushWaitDealInfoToStack(infoSlice []ThirdPkgBaseInfo) {
	for _, v := range infoSlice {
		forInfoStack.Push(v)
	}
}

func GetNextInfoFromStack() (index int, info ThirdPkgBaseInfo, isEnd bool) {
	forInfoStackLock.Lock()
	defer forInfoStackLock.Unlock()
	elCount := forInfoStack.Size()
	lg.InfoToFile(fmt.Sprintf("elCount=%d\n", elCount))
	if IsFirstTimeGetInfo {
		IsFirstTimeGetInfo = false
	}

	v, ok := forInfoStack.Pop()
	lg.InfoToFile(fmt.Sprintf("v=%v,ok=%v\n", v, ok))
	if !ok || v == nil {
		return initWaitHandleInfoCount - elCount, ThirdPkgBaseInfo{}, true
	}

	return initWaitHandleInfoCount - elCount, v.(ThirdPkgBaseInfo), false
}

var isFirstTimeGetInfo = true
var currentInfoIndex = -1
var forInfoSliceLock sync.Mutex

func GetNextInfoFromSlice() (info ThirdPkgBaseInfo, err error) {
	forInfoSliceLock.Lock()
	defer forInfoSliceLock.Unlock()
	l := len(ThirdPkgBaseInfos)
	var index1 int
	totalTempNum := 0
LabelForContinue:
	totalTempNum++
	// 当循环两轮仍没有可用用户的情况下
	if totalTempNum >= 2*l {
		return ThirdPkgBaseInfo{}, fmt.Errorf("获取不到可用信息")
	}

	if isFirstTimeGetInfo {
		isFirstTimeGetInfo = false
		currentInfoIndex = -1
	}

	index1 = currentInfoIndex + 1

	if index1 >= l {
		currentInfoIndex = 0

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
		currentInfoIndex = index1
		// 若该已经不能登录
		if !ThirdPkgBaseInfos[currentInfoIndex].CanUse {
			goto LabelForContinue
		}

		if ThirdPkgBaseInfos[currentInfoIndex].HadUse {
			goto LabelForContinue
		}

		// 设置为已经使用
		ThirdPkgBaseInfos[currentInfoIndex].HadUse = true
		return ThirdPkgBaseInfos[currentInfoIndex], nil
	}
}
