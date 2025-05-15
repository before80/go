package wind

import (
	"github.com/gonutz/w32/v3"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	VK_BACK    = 0x08 // Backspace 键
	VK_TAB     = 0x09 // Tab 键
	VK_RETURN  = 0x0D // Enter 键
	VK_SHIFT   = 0x10 // Shift 键
	VK_CONTROL = 0x11 // Ctrl 键
	VK_MENU    = 0x12 // Alt 键
	VK_PAUSE   = 0x13 // Pause 键
	VK_CAPITAL = 0x14 // Caps Lock 键
	VK_ESCAPE  = 0x1B // Esc 键
	VK_SPACE   = 0x20 // 空格键
	VK_PRIOR   = 0x21 // Page Up 键
	VK_NEXT    = 0x22 // Page Down 键
	VK_END     = 0x23 // End 键
	VK_HOME    = 0x24 // Home 键
	VK_LEFT    = 0x25 // 左箭头键
	VK_UP      = 0x26 // 上箭头键
	VK_RIGHT   = 0x27 // 右箭头键
	VK_DOWN    = 0x28 // 下箭头键
	VK_INSERT  = 0x2D // Insert 键
	VK_DELETE  = 0x2E // Delete 键

	VK_0 = 0x30 // 数字 0 键
	VK_1 = 0x31 // 数字 1 键
	VK_2 = 0x32 // 数字 2 键
	VK_3 = 0x33 // 数字 3 键
	VK_4 = 0x34 // 数字 4 键
	VK_5 = 0x35 // 数字 5 键
	VK_6 = 0x36 // 数字 6 键
	VK_7 = 0x37 // 数字 7 键
	VK_8 = 0x38 // 数字 8 键
	VK_9 = 0x39 // 数字 9 键

	VK_A = 0x41 // 字母 A 键
	VK_B = 0x42 // 字母 B 键
	VK_C = 0x43 // 字母 C 键
	VK_D = 0x44 // 字母 D 键
	VK_E = 0x45 // 字母 E 键
	VK_F = 0x46 // 字母 F 键
	VK_G = 0x47 // 字母 G 键
	VK_H = 0x48 // 字母 H 键
	VK_I = 0x49 // 字母 I 键
	VK_J = 0x4A // 字母 J 键
	VK_K = 0x4B // 字母 K 键
	VK_L = 0x4C // 字母 L 键
	VK_M = 0x4D // 字母 M 键
	VK_N = 0x4E // 字母 N 键
	VK_O = 0x4F // 字母 O 键
	VK_P = 0x50 // 字母 P 键
	VK_Q = 0x51 // 字母 Q 键
	VK_R = 0x52 // 字母 R 键
	VK_S = 0x53 // 字母 S 键
	VK_T = 0x54 // 字母 T 键
	VK_U = 0x55 // 字母 U 键
	VK_V = 0x56 // 字母 V 键
	VK_W = 0x57 // 字母 W 键
	VK_X = 0x58 // 字母 X 键
	VK_Y = 0x59 // 字母 Y 键
	VK_Z = 0x5A // 字母 Z 键

	VK_NUMPAD0  = 0x60 // 小键盘 0 键
	VK_NUMPAD1  = 0x61 // 小键盘 1 键
	VK_NUMPAD2  = 0x62 // 小键盘 2 键
	VK_NUMPAD3  = 0x63 // 小键盘 3 键
	VK_NUMPAD4  = 0x64 // 小键盘 4 键
	VK_NUMPAD5  = 0x65 // 小键盘 5 键
	VK_NUMPAD6  = 0x66 // 小键盘 6 键
	VK_NUMPAD7  = 0x67 // 小键盘 7 键
	VK_NUMPAD8  = 0x68 // 小键盘 8 键
	VK_NUMPAD9  = 0x69 // 小键盘 9 键
	VK_MULTIPLY = 0x6A // 小键盘 * 键
	VK_ADD      = 0x6B // 小键盘 + 键
	VK_SUBTRACT = 0x6D // 小键盘 - 键
	VK_DECIMAL  = 0x6E // 小键盘 . 键
	VK_DIVIDE   = 0x6F // 小键盘 / 键

	VK_F1  = 0x70 // 功能键 F1
	VK_F2  = 0x71 // 功能键 F2
	VK_F3  = 0x72 // 功能键 F3
	VK_F4  = 0x73 // 功能键 F4
	VK_F5  = 0x74 // 功能键 F5
	VK_F6  = 0x75 // 功能键 F6
	VK_F7  = 0x76 // 功能键 F7
	VK_F8  = 0x77 // 功能键 F8
	VK_F9  = 0x78 // 功能键 F9
	VK_F10 = 0x79 // 功能键 F10
	VK_F11 = 0x7A // 功能键 F11
	VK_F12 = 0x7B // 功能键 F12

	VK_LWIN = 0x5B // 左 Windows 键
	VK_RWIN = 0x5C // 右 Windows 键
	VK_APPS = 0x5D // 应用程序键（菜单键）

	KEYEVENTF_KEYUP = 0x0002 // 表示按键弹起

	VK_OEM_1      = 0xBA // 分号 (;) 或 冒号 (:) 键
	VK_OEM_PLUS   = 0xBB // 等号 (=) 或 加号 (+) 键
	VK_OEM_COMMA  = 0xBC // 逗号 (,) 键
	VK_OEM_MINUS  = 0xBD // 减号 (-) 键
	VK_OEM_PERIOD = 0xBE // 句号 (.) 键
	VK_OEM_2      = 0xBF // 斜杠 (/) 或 问号 (?) 键
	VK_OEM_3      = 0xC0 // 波浪号 (~) 或 反引号 (`) 键
	VK_OEM_4      = 0xDB // 左中括号 ([) 键
	VK_OEM_5      = 0xDC // 反斜杠 (\) 键
	VK_OEM_6      = 0xDD // 右中括号 (]) 键
	VK_OEM_7      = 0xDE // 单引号 (') 或 双引号 (") 键

	VK_NUMLOCK  = 0x90 // Num Lock 键
	VK_SCROLL   = 0x91 // Scroll Lock 键
	VK_LSHIFT   = 0xA0 // 左 Shift 键
	VK_RSHIFT   = 0xA1 // 右 Shift 键
	VK_LCONTROL = 0xA2 // 左 Ctrl 键
	VK_RCONTROL = 0xA3 // 右 Ctrl 键
	VK_LMENU    = 0xA4 // 左 Alt 键
	VK_RMENU    = 0xA5 // 右 Alt 键

	VK_BROWSER_BACK      = 0xA6 // 浏览器后退键
	VK_BROWSER_FORWARD   = 0xA7 // 浏览器前进键
	VK_BROWSER_REFRESH   = 0xA8 // 浏览器刷新键
	VK_BROWSER_STOP      = 0xA9 // 浏览器停止键
	VK_BROWSER_SEARCH    = 0xAA // 浏览器搜索键
	VK_BROWSER_FAVORITES = 0xAB // 浏览器收藏键
	VK_BROWSER_HOME      = 0xAC // 浏览器主页键

	VK_VOLUME_MUTE         = 0xAD // 静音键
	VK_VOLUME_DOWN         = 0xAE // 音量减键
	VK_VOLUME_UP           = 0xAF // 音量加键
	VK_MEDIA_NEXT_TRACK    = 0xB0 // 媒体下一曲键
	VK_MEDIA_PREV_TRACK    = 0xB1 // 媒体上一曲键
	VK_MEDIA_STOP          = 0xB2 // 媒体停止键
	VK_MEDIA_PLAY_PAUSE    = 0xB3 // 媒体播放/暂停键
	VK_LAUNCH_MAIL         = 0xB4 // 启动邮件客户端键
	VK_LAUNCH_MEDIA_SELECT = 0xB5 // 启动媒体选择器键
	VK_LAUNCH_APP1         = 0xB6 // 启动应用程序 1 键
	VK_LAUNCH_APP2         = 0xB7 // 启动应用程序 2 键
)

// 定义 Windows API 调用
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procKeyBdEvent       = user32.NewProc("keybd_event")
	procEnumChildWindows = user32.NewProc("EnumChildWindows")
	procGetClassName     = user32.NewProc("GetClassNameW")
	procSetFocus         = user32.NewProc("SetFocus")
	procSendMessage      = user32.NewProc("SendMessageW")
)

var (
	onceCallback sync.Once
	enumCallback uintptr
)

// 存储结果的结构体，通过 lparam 传递
type findWindowResult struct {
	targetClass string
	hwnd        w32.HWND
}

// 初始化回调函数（只执行一次）
func ensureCallbackInitialized() {
	onceCallback.Do(func() {
		enumCallback = syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
			result := (*findWindowResult)(unsafe.Pointer(lparam))
			className := getClassName(hwnd)
			if className == result.targetClass {
				result.hwnd = w32.HWND(hwnd)
				return 0 // 停止枚举
			}
			return 1 // 继续枚举
		})
	})
}

// 获取窗口类名
func getClassName(hwnd syscall.Handle) string {
	buffer := make([]uint16, 256)
	procGetClassName.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	return syscall.UTF16ToString(buffer)
}

// FindChromeBrowserContentWindow 查找 Chrome 浏览器内容窗口
func FindChromeBrowserContentWindow(parent w32.HWND) w32.HWND {
	ensureCallbackInitialized()

	result := &findWindowResult{
		targetClass: "Chrome_RenderWidgetHostHWND",
		hwnd:        0,
	}

	procEnumChildWindows.Call(
		uintptr(parent),
		enumCallback,
		uintptr(unsafe.Pointer(result)),
	)

	return result.hwnd
}

//
//type enumChildCallback func(hwnd syscall.Handle, lparam uintptr) uintptr
//
//func enumChildWindows(parent syscall.Handle, callback enumChildCallback, lparam uintptr) bool {
//	cb := syscall.NewCallback(callback)
//	ret, _, _ := procEnumChildWindows.Call(
//		uintptr(parent),
//		cb,
//		lparam,
//	)
//	return ret != 0
//}
//
//func getClassName(hwnd syscall.Handle) string {
//	buffer := make([]uint16, 1024)
//	_, _, _ = procGetClassName.Call(
//		uintptr(hwnd),
//		uintptr(unsafe.Pointer(&buffer[0])),
//		uintptr(len(buffer)),
//	)
//	return syscall.UTF16ToString(buffer)
//}
//
//// FindChromeBrowserContentWindow 查找Chrome浏览器指定网页的内容窗口
//func FindChromeBrowserContentWindow(parent w32.HWND) w32.HWND {
//	var contentHwnd w32.HWND
//
//	callback := func(hwnd syscall.Handle, lparam uintptr) uintptr {
//		className := getClassName(hwnd)
//		// Chrome：类名 Chrome_RenderWidgetHostHWND
//		// Firefox：类名 MozillaWindowClass 的子窗口
//		if className == "Chrome_RenderWidgetHostHWND" {
//			contentHwnd = w32.HWND(hwnd)
//			return 0
//		}
//		return 1
//	}
//
//	enumChildWindows(syscall.Handle(parent), callback, 0)
//	return contentHwnd
//}

func keybdEvent(bVk, bScan, dwFlags, dwExtraInfo byte) {
	_, _, _ = procKeyBdEvent.Call(
		uintptr(bVk),
		uintptr(bScan),
		uintptr(dwFlags),
		uintptr(dwExtraInfo),
	)
}

func pressCtrlAndKey(vk byte) {
	// Ctrl down
	keybdEvent(VK_CONTROL, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	// Key down
	keybdEvent(vk, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	// Key up
	keybdEvent(vk, 0, KEYEVENTF_KEYUP, 0)
	time.Sleep(50 * time.Millisecond)

	// Ctrl up
	keybdEvent(VK_CONTROL, 0, KEYEVENTF_KEYUP, 0)
}

func pressCtrlShiftAndKey(vk byte) {
	// Ctrl down
	keybdEvent(VK_CONTROL, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	// Shift down
	keybdEvent(VK_SHIFT, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	// Key down
	keybdEvent(vk, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	// Key up
	keybdEvent(vk, 0, KEYEVENTF_KEYUP, 0)
	time.Sleep(50 * time.Millisecond)

	// Shift up
	keybdEvent(VK_SHIFT, 0, KEYEVENTF_KEYUP, 0)
	time.Sleep(50 * time.Millisecond)

	// Ctrl up
	keybdEvent(VK_CONTROL, 0, KEYEVENTF_KEYUP, 0)
	time.Sleep(50 * time.Millisecond)
}
