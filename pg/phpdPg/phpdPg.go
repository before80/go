package phpdPg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/before80/go/bs"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/js/phpdJs"
	"github.com/before80/go/lg"
	"github.com/before80/go/next/phpdNext"
	"github.com/before80/go/pg"
	"github.com/before80/go/res"
	"github.com/before80/go/tr"
	"github.com/before80/go/tr/phpdTr"
	"github.com/before80/go/wind"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var didF *os.File
var didUrl []string

func init() {
	var err error
	err = os.MkdirAll(filepath.Join(contants.DidFolderName, "php"), 0777)
	if err != nil {
		panic(fmt.Sprintf("无法创建%s目录：%v\n", filepath.Join(contants.DidFolderName, "php"), err))
	}

	didF, err = os.OpenFile(filepath.Join(contants.DidFolderName, "php", "php.net.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("无法创建或打开文件%q: %v\n", "php.net.txt", err))
		panic(err)
	}

	scanner := bufio.NewScanner(didF)
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的文本
		//fmt.Println("line=", line)
		if strings.TrimSpace(line) != "" {
			didUrl = append(didUrl, line)
		}
	}
	//fmt.Println("didUrl=", didUrl)
	//_, _ = didF.Seek(2, 0)
	res.NewPHP(bufio.NewWriter(didF))
}

func CloseInitFiles() {
	if didF != nil {
		_ = didF.Close()
	}
}

type MenuInfo struct {
	MenuName string `json:"menu_name"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
	Index    int    `json:"index"`
}

// GetAllFirstMenuInfo 获取所有第一级菜单信息
func GetAllFirstMenuInfo(page *rod.Page, url string) (firstMenuInfos []phpdNext.MenuInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("获取allStdPkgInfo时遇到错误：%v", r)
		}
	}()

	page.MustNavigate(url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(phpdJs.GetAllFirstMenuJs)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行phpdJs.GetAllFirstMenuJs遇到错误：%v", url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &firstMenuInfos)
	if err != nil {
		return nil, fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", url, err)
	}
	return
}

func DealWithMenuPageData(threadIndex int, wg *sync.WaitGroup) {
	var err error
	hadWgDone := false
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
	preDir := cfg.Default.PHPdPreFolderName
	browser := bs.MyBrowserSlice[threadIndex].Browser
	page := browser.MustPage()

	defer func() {
		_ = page.Close()
	}()
	var pageTitle, chromePageWindowTitle string
	var fpDst string
	var result *proto.RuntimeRemoteObject
	var subMenuInfos []phpdNext.MenuInfo
	uniqueMdFilename := "do" + strconv.Itoa(threadIndex) + ".md"
	relUniqueMdFilePath := filepath.Join("markdown", uniqueMdFilename)
	typoraWindowTitle := uniqueMdFilename + " - Typora"
	_, _ = pg.CreateFileIfNotExists(relUniqueMdFilePath)
	absUniqueMdFilePath, _ := filepath.Abs(relUniqueMdFilePath)
LabelForContinue:
	subMenuInfos = nil
	_ = tr.TruncFileContent(relUniqueMdFilePath)
	date := time.Now().Format(time.RFC3339)
	curMenu, isEnd := phpdNext.GetNextMenuInfoFromQueue()
	lg.InfoToFile(fmt.Sprintf("线程%d正要处理的menu=%v\n", threadIndex, curMenu))
	if isEnd {
		if !hadWgDone {
			hadWgDone = true
			lg.InfoToFile(fmt.Sprintf("线程%d中已设置hadWgDone = true，且调用了wg.Done()\n", threadIndex))
			wg.Done()
		}
		return
	}

	page.MustNavigate(curMenu.Url)
	page.MustWaitLoad()

	result, err = page.Eval(phpdJs.InDetailPageGetMenuJs)
	if err != nil {
		panic(fmt.Sprintf("线程%d在网页%s中执行phpdJs.InDetailPageGetMenuJs遇到错误：%v", threadIndex, curMenu.Url, err))
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		panic(fmt.Sprintf("线程%d在处理网页%s时执行json.Marshal遇到错误: %v", threadIndex, curMenu.Url, err))
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &subMenuInfos)
	if err != nil {
		panic(fmt.Sprintf("线程%d在处理网页%s时执行json.Unmarshal遇到错误: %v", threadIndex, curMenu.Url, err))
	}

	lg.InfoToFile(fmt.Sprintf("线程%d网页%s的子菜单=%v\n", threadIndex, curMenu.Url, subMenuInfos))

	if curMenu.IsTopMenu == 1 {
		fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, "_index.md")
	} else {
		fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, curMenu.Filename+".md")
	}

	if len(subMenuInfos) > 0 {
		for i, _ := range subMenuInfos {
			if curMenu.IsTopMenu == 1 {
				subMenuInfos[i].Dir = curMenu.Filename
			} else {
				subMenuInfos[i].Dir = curMenu.Dir + "/" + curMenu.Filename
			}
		}
		if curMenu.IsTopMenu != 1 {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, curMenu.Filename, "_index.md")
		} else {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Filename, "_index.md")
		}
	} else {
		if curMenu.IsTopMenu != 1 {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Dir, curMenu.Filename+".md")
		} else {
			fpDst = filepath.Join(contants.OutputFolderName, preDir, curMenu.Filename+".md")
		}
	}

	hadExist, _ := pg.CreateFileIfNotExists(fpDst)
	_ = hadExist
	mdF, _ := os.OpenFile(fpDst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if curMenu.IsTopMenu == 1 {
		_, _ = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
linkTitle = "%s"
date = %s
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false
[menu.main]
	weight = %d

+++

> 原文：[%s](%s)
>
> 收录时间：%s
`, curMenu.MenuName, curMenu.MenuName, date, "", curMenu.Weight, curMenu.Url, curMenu.Url, fmt.Sprintf("`%s`", date)))
	} else {
		_, _ = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++

> 原文：[%s](%s)
>
> 收录时间：%s
`, curMenu.MenuName, date, curMenu.Weight, "", curMenu.Url, curMenu.Url, fmt.Sprintf("`%s`", date)))
	}

	_ = mdF.Close()

	_, err = page.Eval(fmt.Sprintf(`() => { %s }`, phpdJs.ReplaceJs))
	if err != nil {
		panic(fmt.Errorf("线程%d在网页%s中执行ReplaceJs遇到错误：%v", threadIndex, curMenu.Url, err))
	}

	result, err = page.Eval(fmt.Sprintf(`() => { return document.querySelector("#layout-content").textContent.trim() }`, phpdJs.ReplaceJs))
	if err != nil {
		panic(fmt.Sprintf("线程%d在网页%s中执行js获取复制区域的内容遇到错误：%v", threadIndex, curMenu.Url, err))
	}

	if strings.TrimSpace(result.Value.String()) != "" {
		// 获取当前网页的title，在后面会用来查找该网页所在窗口的操作句柄
		result, _ = page.Eval(`() => { return document.title }`)
		pageTitle = result.Value.String()
		chromePageWindowTitle = pageTitle + " - Google Chrome"

		// 再次清空
		_ = tr.TruncFileContent(relUniqueMdFilePath)

		contentBytes, _ := wind.DoCopyAndPaste(threadIndex, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, curMenu.Url)

		if contentBytes == 0 {
			lg.InfoToFile(fmt.Sprintf("线程%d发现复制网页%s的字节数为0，将加入到下一次进行重试", threadIndex, curMenu.Url))
			curMenu.Retry = curMenu.Retry + 1
			phpdNext.PushWaitDealMenuInfoToQueue([]phpdNext.MenuInfo{curMenu})
		} else {
			_, err = phpdTr.ReplaceMarkdownFileContent(absUniqueMdFilePath)
			if err != nil {
				panic(fmt.Errorf("线程%d在替换网页%s的内容对应的md文件时出现错误：%v", threadIndex, curMenu.Url, err))
			}
			lg.InfoToFile(fmt.Sprintf("线程%d正要处理Insert", threadIndex))
			err = pg.InsertAnyPageData(fpDst, relUniqueMdFilePath, "> 收录时间：")
			if err != nil {
				panic(fmt.Errorf("线程%d在将网页%s中的内容插入到目标md文件时遇到错误：%v", threadIndex, curMenu.Url, err))
			}
		}
	}

	if curMenu.Retry == 0 && len(subMenuInfos) > 0 {
		// 将 subMenuInfos 推入
		phpdNext.PushWaitDealMenuInfoToQueue(subMenuInfos)
	}
	goto LabelForContinue
}

func InitFirstMenuMdFile(browserHwnd win.HWND, firstMenuInfo MenuInfo, page *rod.Page) (err error) {
	var dir, curDir string
	var subMenuInfos []MenuInfo
	//isBar := false
	useUnderlineIndexMd := false
	prefixDirname := "php_"

	// 判断是否还有二级菜单
	page.MustNavigate(firstMenuInfo.Url)
	page.MustWaitLoad()
	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(phpdJs.InDetailPageGetMenuJs)
	if err != nil {
		return fmt.Errorf("在网页%s中执行phpdJs.InDetailPageGetMenuJs遇到错误：%v", firstMenuInfo.Url, err)
	}

	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", firstMenuInfo.Url, err)
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &subMenuInfos)
	if err != nil {
		return fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", firstMenuInfo.Url, err)
	}

	if len(subMenuInfos) > 0 {
		if slices.Contains(didUrl, firstMenuInfo.Url) {
			lg.InfoToFileAndStdOut(fmt.Sprintf("%s 之前已处理过 - %s\n", firstMenuInfo.MenuName, firstMenuInfo.Url))
			time.Sleep(2 * time.Second)
		} else {
			useUnderlineIndexMd = true
			dir = filepath.Join(contants.OutputFolderName, prefixDirname+firstMenuInfo.Filename)
			curDir = dir
			err = preInitMdFile(true, useUnderlineIndexMd, dir, firstMenuInfo)
			if err != nil {
				return
			}

			_, err = page.Eval(fmt.Sprintf(`() => { %s }`, phpdJs.ReplaceJs))
			if err != nil {
				return fmt.Errorf("在网页%s中执行phpdJs.ReplaceJs遇到错误：%v", firstMenuInfo.Url, err)
			}

			var result1 *proto.RuntimeRemoteObject
			result1, err = page.Eval(phpdJs.GetLayoutContentJs)
			if result1.Value.String() == "" {
				lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", firstMenuInfo.Url))
			} else {
				// 插入内容
				mdFilePath := filepath.Join(dir, "_index.md")
				err = InsertDetailPageData(browserHwnd, mdFilePath, firstMenuInfo)
				if err != nil {
					return
				}
				// 记录已经处理的url
				res.PHP.WriteStringAndFlush(fmt.Sprintf("%s\n", firstMenuInfo.Url))
			}
		}
	} else {
		if slices.Contains(didUrl, firstMenuInfo.Url) {
			lg.InfoToFileAndStdOut(fmt.Sprintf("%s 之前已处理过 - %s\n", firstMenuInfo.MenuName, firstMenuInfo.Url))
			time.Sleep(2 * time.Second)
		} else {
			useUnderlineIndexMd = true
			dir = filepath.Join(contants.OutputFolderName, prefixDirname+firstMenuInfo.Filename)
			curDir = dir
			err = preInitMdFile(true, useUnderlineIndexMd, dir, firstMenuInfo)
			if err != nil {
				return
			}

			_, err = page.Eval(fmt.Sprintf(`() => { %s }`, phpdJs.ReplaceJs))
			if err != nil {
				return fmt.Errorf("在网页%s中执行phpdJs.ReplaceJs遇到错误：%v", firstMenuInfo.Url, err)
			}

			var result1 *proto.RuntimeRemoteObject
			result1, err = page.Eval(phpdJs.GetLayoutContentJs)
			vStr := strings.TrimSpace(result1.Value.String())
			//lg.InfoToFileAndStdOut(fmt.Sprintf("%v -> %s", vStr, firstMenuInfo.Url))
			if vStr == "" {
				lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", firstMenuInfo.Url))
			} else {
				// 插入内容
				mdFilePath := filepath.Join(dir, "_index.md")
				err = InsertDetailPageData(browserHwnd, mdFilePath, firstMenuInfo)
				if err != nil {
					return
				}
				// 记录已经处理的url
				res.PHP.WriteStringAndFlush(fmt.Sprintf("%s\n", firstMenuInfo.Url))
			}
		}
	}

	if len(subMenuInfos) > 0 {
		culMenuLevel := 2
		err = DealSubMenuInfo(browserHwnd, subMenuInfos, curDir, culMenuLevel, page)
		if err != nil {
			return
		}
	}
	return
}

func DealSubMenuInfo(browserHwnd win.HWND, subMenuInfos []MenuInfo, curDir string, menuLevel int, page *rod.Page) (err error) {
	if len(subMenuInfos) > 0 {
		var subSubMenuInfos []MenuInfo
		useUnderlineIndexMd := false
		var dir, subCurDir string
		subMenuInfosLen := len(subMenuInfos)
		for i, subMenuInfo := range subMenuInfos {
			if slices.Contains([]string{"refs_basic_php", "refs_compression", "refs_remote_auth", "refs_utilspec_audio", "refs_utilspec_cmdline", "refs_crypto"}, subMenuInfo.Filename) {
				continue
			}
			subSubMenuInfos = []MenuInfo{}
			lg.InfoToFileAndStdOut(fmt.Sprintf("正在处理第%d层(当前层还有%d个菜单待处理) %s - %s\n", menuLevel, subMenuInfosLen-i-1, subMenuInfo.MenuName, subMenuInfo.Url))

			// 判断是否还有二级菜单
			page.MustNavigate(subMenuInfo.Url)
			page.MustWaitLoad()
			var result *proto.RuntimeRemoteObject
			result, err = page.Eval(phpdJs.InDetailPageGetMenuJs)
			if err != nil {
				return fmt.Errorf("在网页%s中执行phpdJs.InDetailPageGetMenuJs遇到错误：%v", subMenuInfo.Url, err)
			}

			// 将结果序列化为 JSON 字节
			var jsonBytes []byte
			jsonBytes, err = json.Marshal(result.Value)
			if err != nil {
				return fmt.Errorf("在网页%s中执行json.Marshal遇到错误: %v", subMenuInfo.Url, err)
			}

			// 将 JSON 数据反序列化到结构体中
			err = json.Unmarshal(jsonBytes, &subSubMenuInfos)
			if err != nil {
				return fmt.Errorf("在网页%s中执行json.Unmarshal遇到错误: %v", subMenuInfo.Url, err)
			}

			if len(subSubMenuInfos) > 0 {
				if slices.Contains(didUrl, subMenuInfo.Url) {
					lg.InfoToFileAndStdOut(fmt.Sprintf("%s 之前已处理过 - %s\n", subMenuInfo.MenuName, subMenuInfo.Url))
					time.Sleep(2 * time.Second)
				} else {
					useUnderlineIndexMd = true
					dir = filepath.Join(curDir, subMenuInfo.Filename)
					subCurDir = dir
					err = preInitMdFile(false, useUnderlineIndexMd, dir, subMenuInfo)
					if err != nil {
						return
					}
					_, err = page.Eval(fmt.Sprintf(`() => { %s }`, phpdJs.ReplaceJs))
					if err != nil {
						return fmt.Errorf("在网页%s中执行phpdJs.ReplaceJs遇到错误：%v", subMenuInfo.Url, err)
					}

					var result1 *proto.RuntimeRemoteObject
					result1, err = page.Eval(phpdJs.GetLayoutContentJs)
					vStr := strings.TrimSpace(result1.Value.String())
					//lg.InfoToFileAndStdOut(fmt.Sprintf("%v -> %s", vStr, subMenuInfo.Url))
					if vStr == "" {
						lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", subMenuInfo.Url))
					} else {
						// 插入内容
						mdFilePath := filepath.Join(dir, "_index.md")
						err = InsertDetailPageData(browserHwnd, mdFilePath, subMenuInfo)
						if err != nil {
							return
						}
						// 记录已经处理的url
						res.PHP.WriteStringAndFlush(fmt.Sprintf("%s\n", subMenuInfo.Url))
					}
				}
			} else {
				if slices.Contains(didUrl, subMenuInfo.Url) {
					lg.InfoToFileAndStdOut(fmt.Sprintf("%s 之前已处理过 - %s\n", subMenuInfo.MenuName, subMenuInfo.Url))
					time.Sleep(2 * time.Second)
				} else {
					useUnderlineIndexMd = false
					dir = curDir
					subCurDir = dir
					err = preInitMdFile(false, useUnderlineIndexMd, dir, subMenuInfo)
					if err != nil {
						return
					}

					_, err = page.Eval(fmt.Sprintf(`() => { %s }`, phpdJs.ReplaceJs))
					if err != nil {
						return fmt.Errorf("在网页%s中执行phpdJs.ReplaceJs遇到错误：%v", subMenuInfo.Url, err)
					}

					var result1 *proto.RuntimeRemoteObject
					result1, err = page.Eval(phpdJs.GetLayoutContentJs)
					vStr := strings.TrimSpace(result1.Value.String())
					//lg.InfoToFileAndStdOut(fmt.Sprintf("%v -> %s", vStr, subMenuInfo.Url))
					if vStr == "" {
						lg.InfoToFileAndStdOut(fmt.Sprintf("%s执行phpdJs.GetLayoutContentJs后，页面内容为空", subMenuInfo.Url))
					} else {
						// 插入内容
						mdFilePath := filepath.Join(dir, subMenuInfo.Filename+".md")
						err = InsertDetailPageData(browserHwnd, mdFilePath, subMenuInfo)
						if err != nil {
							return
						}
						// 记录已经处理的url
						res.PHP.WriteStringAndFlush(fmt.Sprintf("%s\n", subMenuInfo.Url))
					}
				}
			}

			if err != nil {
				return
			}

			if len(subSubMenuInfos) > 0 {
				curMenuLevel := menuLevel + 1
				err = DealSubMenuInfo(browserHwnd, subSubMenuInfos, subCurDir, curMenuLevel, page)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func InsertDetailPageData(browserHwnd win.HWND, mdFilePath string, menuInfo MenuInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("插入detailPage=%s数据时遇到错误：%v", menuInfo.Url, r)
		}
	}()
	step := "detailPage"
	err = pg.DealUniqueMd(browserHwnd, menuInfo.Url, step)
	if err != nil {
		return err
	}
	_, err = phpdTr.ReplaceMarkdownFileContent(cfg.Default.UniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理%s=%s时，替换出现错误：%v", step, menuInfo.Url, err)
	}

	err = pg.InsertAnyPageData(mdFilePath, cfg.Default.UniqueMdFilepath, "> 收录时间：")
	return
}

func preInitMdFile(isBar, useUnderlineIndexMd bool, dir string, menuInfo MenuInfo) (err error) {
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("无法创建%s目录：%v\n", dir, err)
	}
	var filename string
	if useUnderlineIndexMd {
		filename = "_index.md"
	} else {
		filename = menuInfo.Filename + ".md"
	}

	mdFp := filepath.Join(dir, filename)
	var mdF *os.File
	_, err1 := os.Stat(mdFp)

	// 当文件不存在的情况下，新建文件并初始化该文件
	if err1 != nil && errors.Is(err1, fs.ErrNotExist) {
		//fmt.Println("err=", err1)
		mdF, err = os.OpenFile(mdFp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("创建文件 %s 时出错: %w", mdFp, err)
		}
		defer mdF.Close()
		date := time.Now().Format(time.RFC3339)
		if isBar {
			_, err = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
linkTitle = "%s"
date = %s
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false
[menu.main]
	weight = %d
+++

> 原文：[%s](%s)
>
> 收录时间：%s
`, menuInfo.MenuName, menuInfo.MenuName, date, "", menuInfo.Index*10, menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))
		} else {
			_, err = mdF.WriteString(fmt.Sprintf(`+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++

> 原文：[%s](%s)
>
> 收录时间：%s
`, menuInfo.MenuName, date, menuInfo.Index*10, "", menuInfo.Url, menuInfo.Url, fmt.Sprintf("`%s`", date)))
		}

		if err != nil {
			return fmt.Errorf("初始化%s文件时出错: %v", mdFp, err)
		}
	}
	return
}
