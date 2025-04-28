package wind

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"
)

// GetWindowText 使用 GetWindowTextW 获取窗口标题
func GetWindowText(hwnd win.HWND) (string, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	procGetWindowText := user32.NewProc("GetWindowTextW")

	// 分配缓冲区
	const maxChars = 256
	buffer := make([]uint16, maxChars)

	// 调用 GetWindowTextW
	ret, _, err := procGetWindowText.Call(
		uintptr(hwnd),                       // 窗口句柄
		uintptr(unsafe.Pointer(&buffer[0])), // 缓冲区
		uintptr(maxChars),                   // 缓冲区大小
	)
	if ret == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return "", fmt.Errorf("GetWindowText 失败: %v", err)
		}
		return "", nil // 窗口没有标题
	}

	// 将宽字符转换为 Go 字符串
	return syscall.UTF16ToString(buffer), nil
}

func FindWindowHwndByWindowTitle(windowTitle string) (hwnd win.HWND, err error) {
	hwnd = robotgo.FindWindow(windowTitle)
	if hwnd == 0 {
		return 0, fmt.Errorf(`未找到 '%s' 窗口`, windowTitle)
	}
	return hwnd, nil
}

func FindWindowHwndByPId(pid int32) {

}

func SetChromeWindowsName(hwnd win.HWND, windowTitle string) {
	robotgo.SetActiveWindow(hwnd)
	robotgo.SetForeg(hwnd)
	_ = robotgo.KeyTap("e", "alt")
	robotgo.MilliSleep(500)
	_ = robotgo.KeyTap("l")
	robotgo.MilliSleep(500)
	_ = robotgo.KeyTap("w")
	robotgo.MilliSleep(500)
	robotgo.TypeStr(windowTitle)
	_ = robotgo.KeyTap("enter")
	robotgo.MilliSleep(500)
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
	robotgo.MilliSleep(100)
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
