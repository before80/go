package goThirdPkgIndexPg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/goThirdPkgIndexJs"
	"github.com/before80/go/js/mysqldJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/goThirdPkgIndexNext"
	"github.com/go-rod/rod/lib/proto"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

type PkgInfo struct {
	PkgName            string `json:"pkg_name"`
	Filename           string `json:"filename"`
	Url                string `json:"url"`
	Dir                string `json:"dir"`
	Weight             int    `json:"weight"`
	NeedPreCreateIndex int    `json:"need_pre_create_index"`
	Desc               string `json:"desc"`
}

var Weight2PkgInfosMap = make(map[int][]PkgInfo)

func DealWithPkg(threadIndex int, wg *sync.WaitGroup) {
	var err error
	hadWgDone := false
	var pkgInfos []PkgInfo
	defer func() {
		if r := recover(); r != nil {
			lg.ErrorToFile(fmt.Sprintf("线程%d出现异常：%v\n", threadIndex, r))
			lg.ErrorToFile(fmt.Sprintf("线程%d将退出\n", threadIndex))
			if !hadWgDone {
				lg.InfoToFile(fmt.Sprintf("在线程%d的defer中调用了wg.Done()\n", threadIndex))
				wg.Done()
			}
		}
	}()
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()
LabelForContinue:
	pkgInfos = nil
	_, info, isEnd := goThirdPkgIndexNext.GetNextInfoFromStack()
	if isEnd {
		if !hadWgDone {
			hadWgDone = true
			lg.InfoToFile(fmt.Sprintf("在线程%d中已设置hadWgDone = true，且调用了wg.Done()\n", threadIndex))
			wg.Done()
		}
		return
	}
	page.MustNavigate(info.Url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(mysqldJs.ExpandMenusJs)
	if err != nil {
		panic(fmt.Errorf("在网页%s中执行goThirdPkgIndexJs.FromTableGetAllPkgInfoJs遇到错误：%v", info.Url, err))
	}
	time.Sleep(3 * time.Second)

	result, err = page.Eval(fmt.Sprintf(goThirdPkgIndexJs.FromTableGetAllPkgInfoJs, info.PkgName, info.Weight))
	if err != nil {
		panic(fmt.Errorf("在网页%s中执行goThirdPkgIndexJs.FromTableGetAllPkgInfoJs遇到错误：%v", info.Url, err))
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		panic(fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", info.Url, err))
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &pkgInfos)
	if err != nil {
		panic(fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", info.Url, err))
	}

	Weight2PkgInfosMap[pkgInfos[0].Weight] = pkgInfos

	//err = appendLinesToFile(pkgInfos)
	//if err != nil {
	//	panic(fmt.Sprintf("在处理%s遇到错误\n", info.PkgName))
	//}
	goto LabelForContinue
}

var insertLock sync.Mutex

func AppendLinesToFile() (err error) {
	insertLock.Lock()
	defer insertLock.Unlock()
	var file *os.File
	file, err = os.OpenFile(filepath.Join(contants.OutputFolderName, "go_third_pkg_info.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	var sPkgInfos [][]PkgInfo
	var weights []int
	for k, _ := range Weight2PkgInfosMap {
		weights = append(weights, k)
	}
	slices.Sort(weights)
	for _, k := range weights {
		sPkgInfos = append(sPkgInfos, Weight2PkgInfosMap[k])
	}

	// 创建一个写入器
	writer := bufio.NewWriter(file)

	// 先添加一个换行符
	if _, err = file.WriteString("\n"); err != nil {
		return fmt.Errorf("写入换行符时出错: %w", err)
	}

	for _, pkgInfos := range sPkgInfos {
		// 遍历要追加的每一行内容
		for _, info := range pkgInfos {
			line := fmt.Sprintf("%s||%s||%s||%s||%d||%d||%s\n", info.Url, info.PkgName, info.Dir, info.Filename, info.NeedPreCreateIndex, info.Weight, info.Desc)
			// 写入当前行
			if _, err = writer.WriteString(line); err != nil {
				return fmt.Errorf("写入行时出错: %w", err)
			}
		}
		if _, err = writer.WriteString("\n"); err != nil {
			return fmt.Errorf("写入行时出错: %w", err)
		}
	}

	// 将缓冲区的内容刷新到文件
	if err = writer.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区时出错: %w", err)
	}
	return nil
}
