package wind

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/before80/go/lg"
	"github.com/go-vgo/robotgo"
	"github.com/gonutz/w32/v3"
	"github.com/tailscale/win"
	"math"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type windowSearchParams struct {
	keyword string
	result  w32.HWND
}

var (
	searchParams *windowSearchParams
	cbOnce       sync.Once
	cbPtr        uintptr
)

func initEnumWindowsCallback() {
	cbPtr = syscall.NewCallback(func(hwnd, lparam uintptr) uintptr {
		h := w32.HWND(hwnd)
		if !w32.IsWindowVisible(h) {
			return 1
		}

		title, err := w32.GetWindowText(h)
		if err != nil {
			return 1
		}

		if strings.Contains(title, searchParams.keyword) {
			searchParams.result = h
			return 0 // stop
		}
		return 1 // continue
	})
}

// FindWindowByTitle 封装查找窗口句柄的函数
func FindWindowByTitle(keyword string) (win.HWND, error) {
	cbOnce.Do(initEnumWindowsCallback)

	searchParams = &windowSearchParams{keyword: keyword}

	err := w32.EnumWindows(cbPtr, 0)
	if err != nil {
		return 0, fmt.Errorf("EnumWindows failed: %v", err)
	}

	if searchParams.result == 0 {
		return 0, fmt.Errorf("no window found containing title: %q", keyword)
	}
	return win.HWND(searchParams.result), nil
}

// FindWindowByTitle2 封装查找窗口句柄的函数
// 反复创建了过多的 syscall.NewCallback（或 syscall.NewCallbackCDecl）对象，而它们没有被释放或重复创建导致耗尽系统资源。
// 最终导致执行时报错：fatal error: too many callback functions
func FindWindowByTitle2(keyword string) (win.HWND, error) {
	//fmt.Printf("keyword=%s\n", keyword)
	var result w32.HWND

	var cbFunc = func(hwnd uintptr, lparam uintptr) uintptr {
		h := w32.HWND(hwnd)
		if !w32.IsWindowVisible(h) {
			return 1
		}

		title, err := w32.GetWindowText(h)
		//fmt.Printf("title=%s,err=%v\n", title, err)
		if err != nil {
			return 1
		}

		//fmt.Printf("strings.Contains(%q, %q)=%v\n", title, keyword, strings.Contains(title, keyword))
		if strings.Contains(title, keyword) {
			result = h
			return 0 // stop
		}
		return 1 // continue
	}

	// 这里 callback 保存在变量中，防止作用域问题
	callback := syscall.NewCallback(cbFunc)

	_ = w32.EnumWindows(callback, 0)

	if result == 0 {
		return 0, fmt.Errorf("no window found containing title: %q", keyword)
	}
	return win.HWND(result), nil
}

// InChromePageDoCtrlAAndC 在浏览器页面中执行全选和复制操作
func InChromePageDoCtrlAAndC(tempHwnd win.HWND) (contentByes int, err error) {
	hwnd := w32.HWND(tempHwnd)
	// 清空剪贴板
	_ = clipboard.WriteAll("")

	// 激活主窗口
	w32.ShowWindow(hwnd, w32.SW_RESTORE)
	if err = w32.SetForegroundWindow(hwnd); err != nil {
		return 0, fmt.Errorf("SetForegroundWindow failed")
	}
	time.Sleep(800 * time.Millisecond) // 增加延迟

	// 定位内容区域
	contentHwnd := FindChromeBrowserContentWindow(hwnd)
	if contentHwnd == 0 {
		return 0, fmt.Errorf("内容窗口未找到")
	}

	// 设置焦点并发送虚拟鼠标事件
	_, _ = w32.SetFocus(contentHwnd)
	w32.SendMessage(contentHwnd, w32.WM_MOUSEMOVE, 0, 0)
	time.Sleep(300 * time.Millisecond)

	// 执行复制操作
	pressCtrlAndKey(VK_A)
	time.Sleep(100 * time.Millisecond)
	pressCtrlAndKey(VK_C)
	time.Sleep(300 * time.Millisecond)
	// 等待剪贴板数据
	if contentByes, err = waitForClipboard(); err != nil {
		return 0, err
	}
	return contentByes, nil
}

var copyPasteLock sync.Mutex

func DoCopyAndPaste(threadIndex int, absUniqueMdFilePath, typoraWindowTitle, chromePageWindowTitle, url string) (contentBytes int, err error) {
	copyPasteLock.Lock()
	defer copyPasteLock.Unlock()
	var typoraHwnd win.HWND
	browserHwnd := robotgo.FindWindow(chromePageWindowTitle)

	//_ = OpenTypora(absUniqueMdFilePath)
	timeoutChan := time.After(10 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 每隔 interval 时间检查一次条件
			hwnd1 := robotgo.FindWindow(typoraWindowTitle)
			lg.InfoToFile(fmt.Sprintf("%d - typoraHwnd=%v\n", threadIndex, hwnd1))
			if hwnd1 != 0 {
				typoraHwnd = hwnd1
				goto LabelForContinue
			}
			//hwnd1, err1 = wind.FindWindowByTitle(uniqueMdFilename + " - Typora")
		case <-timeoutChan:
			// 超时后退出循环
			goto LabelForContinue
		}
	}
LabelForContinue:
	lg.InfoToFile(fmt.Sprintf("线程%d中获取到的typoraHwnd=%v\n", threadIndex, typoraHwnd))
	var err1 error
	contentBytes, err1 = InChromePageDoCtrlAAndC(browserHwnd)
	lg.InfoToFile(fmt.Sprintf("在页面%s获取到的字节数为：%d\n", url, contentBytes))
	if err1 != nil {
		lg.ErrorToFile(fmt.Sprintf("在浏览器中进行复制遇到错误：%v\n", err1))
	}
	_ = DoCtrlVAndS(typoraHwnd, contentBytes)
	//_ = win.SendMessage(typoraHwnd, win.WM_CLOSE, 0, 0)
	_ = win.SendMessage(typoraHwnd, win.WM_SYSCOMMAND, win.SC_MINIMIZE, 0)
	//time.Sleep(time.Duration(cfg.Default.WaitTyporaCloseSeconds) * time.Second)
	return contentBytes, nil
}

func DoCtrlVAndS(tempHwnd win.HWND, contentBytes int) error {
	hwnd := w32.HWND(tempHwnd)

	// 激活主窗口
	w32.ShowWindow(hwnd, w32.SW_RESTORE)
	if err := w32.SetForegroundWindow(hwnd); err != nil {
		return fmt.Errorf("SetForegroundWindow failed")
	}
	time.Sleep(250 * time.Millisecond) // 增加延迟
	//time.Sleep(5 * time.Second)
	//// 定位内容区域
	//contentHwnd := FindChromeBrowserContentWindow(hwnd)
	//fmt.Printf("contentHwnd=%v\n", contentHwnd)
	//if contentHwnd == 0 {
	//	return fmt.Errorf("内容窗口未找到")
	//}
	//
	//// 设置焦点并发送虚拟鼠标事件
	//_, _ = w32.SetFocus(contentHwnd)
	//time.Sleep(100 * time.Millisecond) // 给焦点设置一点时间
	//currentFocus := w32.GetFocus()
	//if currentFocus != contentHwnd {
	//	time.Sleep(200 * time.Millisecond) // 给焦点设置一点时间
	//	if w32.GetFocus() != contentHwnd {
	//		fmt.Println("警告: 未能将焦点设置到Typora内容窗口")
	//	}
	//}

	// 模拟鼠标点击可以帮助某些应用正确接受键盘输入
	// 获取窗口客户区的一个点
	rect, _ := w32.GetClientRect(hwnd)
	clientX, clientY := (rect.Right-rect.Left)/2, (rect.Bottom-rect.Top)/2
	screenPoint, _ := w32.ClientToScreen(hwnd, w32.POINT{X: clientX, Y: clientY})

	_ = w32.SetCursorPos(screenPoint.X, screenPoint.Y)
	robotgo.Click()
	lg.InfoToFile(fmt.Sprintf("触发点击左键\n"))
	time.Sleep(50 * time.Millisecond) // 点击后等待

	//w32.SendMessage(contentHwnd, w32.WM_MOUSEMOVE, 0, 0)
	//time.Sleep(300 * time.Millisecond)

	// 执行粘贴保存操作
	pressCtrlAndKey(VK_V)
	lg.InfoToFileAndStdOut(fmt.Sprintf("contentBytes=%d\n", contentBytes))
	if contentBytes > 10000 {
		time.Sleep(time.Duration(int(math.Ceil(float64(contentBytes)/5000))*200) * time.Millisecond)
	} else {
		time.Sleep(time.Duration(int(math.Ceil(float64(contentBytes)/5000))*100) * time.Millisecond)
	}

	pressCtrlAndKey(VK_S)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("had v and s\n")
	return nil
}

func waitForClipboard() (int, error) {
	start := time.Now()
	for time.Since(start) < 6*time.Second {
		if v, err := clipboard.ReadAll(); err == nil {
			//fmt.Printf("等待%v后获取到剪贴板的值:%v\n", time.Since(start), v)
			return len(v), nil
		}
		time.Sleep(80 * time.Millisecond)
	}
	return 0, fmt.Errorf("剪贴板超时")
}

func FindWindowHwndByWindowTitle(windowTitle string) (hwnd win.HWND, err error) {
	hwnd = robotgo.FindWindow(windowTitle)
	if hwnd == 0 {
		return 0, fmt.Errorf(`未找到 '%s' 窗口`, windowTitle)
	}
	return hwnd, nil
}

func OpenTypora(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", "-a", "typora", filePath)
	case "windows": // Windows
		cmd = exec.Command("cmd", "/c", "start", "typora", filePath)
	default: // Linux 或其他
		cmd = exec.Command("typora", filePath)
	}

	return cmd.Run()
}

func OpenDevToolToConsole(hwnd win.HWND) {
	setActiveAndForeg(hwnd)
	_ = robotgo.KeyTap("j", "ctrl", "shift")
}

func SelectAll(hwnd win.HWND) {
	setActiveAndForeg(hwnd)
	_ = robotgo.KeyTap("a", "ctrl")
	robotgo.MilliSleep(200)
}

func CtrlC(hwnd win.HWND) {
	setActiveAndForeg(hwnd)
	//var err error
	_ = robotgo.KeyTap("c", "ctrl")
	//if err != nil {
	//	fmt.Printf("ctrl + c出现错误：%v\n", err)
	//}
	robotgo.MilliSleep(200)
}

func setActiveAndForeg(hwnd win.HWND) {
	robotgo.SetActiveWindow(hwnd)
	robotgo.MilliSleep(100)
	robotgo.SetForeg(hwnd)
	robotgo.MilliSleep(800)
}

func CtrlV(hwnd win.HWND) {
	setActiveAndForeg(hwnd)
	//var err error
	_ = robotgo.KeyTap("v", "ctrl")
	//if err != nil {
	//	fmt.Printf("ctrl + v出现错误：%v\n", err)
	//}
	robotgo.MilliSleep(200)
}

func CtrlS(hwnd win.HWND) {
	setActiveAndForeg(hwnd)
	var err error
	err = robotgo.KeyTap("s", "ctrl")

	if err != nil {
		fmt.Printf("ctrl + s出现错误：%v\n", err)
	}
	robotgo.MilliSleep(200)
	err = robotgo.KeyTap("s", "ctrl")

	if err != nil {
		fmt.Printf("ctrl + s出现错误：%v\n", err)
	}
	robotgo.MilliSleep(200)
}

func SelectAllAndCtrlC(hwnd win.HWND) {
	setActiveAndForeg(hwnd)
	_ = robotgo.KeyTap("a", "ctrl")
	robotgo.MilliSleep(200)
	_ = robotgo.KeyTap("c", "ctrl")
	robotgo.MilliSleep(200)
}

func SelectAllAndDelete(hwnd win.HWND) {
	setActiveAndForeg(hwnd)
	_ = robotgo.KeyTap("a", "ctrl")
	robotgo.MilliSleep(200)
	_ = robotgo.KeyTap("delete")
	robotgo.MilliSleep(200)
	_ = robotgo.KeyTap("a", "ctrl")
	robotgo.MilliSleep(200)
	_ = robotgo.KeyTap("delete")
	robotgo.MilliSleep(200)
}
