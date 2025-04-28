package bs

import (
	"fmt"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"github.com/before80/go/ext"
	"github.com/before80/go/lg"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type MyBrowser struct {
	Browser *rod.Browser
	Ok      bool
	Index   int
	Router  *rod.HijackRouter
}

var MyBrowserSlice []MyBrowser

func init() {
	var err error
	// 获取临时目录
	tempDir := os.TempDir()

	// 构造完整路径
	forCacheFolderPath := filepath.Join(tempDir, contants.ForChromeTempCacheFolderName, contants.AppName)

	// 删除 forCacheFolderPath 中的缓存
	err = os.RemoveAll(forCacheFolderPath)
	if err != nil {
		panic(fmt.Sprintf("删除浏览器缓存文件夹%s出现错误：%v\n", forCacheFolderPath, err))
	}
}

func GetTempFolderPath(folderName, appName, subFolderName string) (string, error) {
	system := runtime.GOOS
	// 获取临时目录
	tempDir := os.TempDir()
	// 删除
	_ = os.RemoveAll(filepath.Join(tempDir, folderName, appName))
	// 构造完整路径
	targetPath := filepath.Join(tempDir, folderName, appName, subFolderName)

	// 打印当前系统及路径
	lg.InfoToFile(fmt.Sprintf("操作系统类型: %s\n", system))
	lg.InfoToFile(fmt.Sprintf("用于缓存的临时目录路径： %s\n", targetPath))

	// 创建文件夹
	err := os.MkdirAll(targetPath, 0777)
	return targetPath, err
}

func GetBrowser(subCacheFolderName string) (browser *rod.Browser, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("在打开浏览器时发生panic：%v\n", r)
			if browser != nil {
				_ = browser.Close()
			}
		}
	}()

	folderPath, err := GetTempFolderPath(contants.ForChromeTempCacheFolderName, contants.AppName, subCacheFolderName)
	if err != nil {
		err = fmt.Errorf("获取用于临时文件夹路径发生错误：%v\n", err)
		return nil, err
	}
	//var absPathToJsHijackExtension string
	//absPathToJsHijackExtension, err = ext.CreateJsHijackExtension(folderPath, subCacheFolderName, true)
	//if err != nil {
	//	return nil, err
	//}

	var absPathToProxyAuthExtension string
	if cfg.Default.UseProxy == 1 {
		absPathToProxyAuthExtension, err = ext.CreateProxyAuthExtension(folderPath, subCacheFolderName, true)
		if err != nil {
			return nil, err
		}
	}

	var l *launcher.Launcher
	var extensionStr string
	//if cfg.Default.UseProxy == 1 {
	//	extensionStr = fmt.Sprintf("%s,%s", absPathToJsHijackExtension, absPathToProxyAuthExtension)
	//} else {
	//	extensionStr = absPathToJsHijackExtension
	//}

	if cfg.Default.UseProxy == 1 {
		extensionStr = fmt.Sprintf("%s", absPathToProxyAuthExtension)
	}

	// Set("disable-component-update", "true")    // 禁止插件自动更新
	if path, exists := launcher.LookPath(); exists {
		lg.InfoToFile(fmt.Sprintf("当前使用的浏览器所在路径是：%s\n", path))

		l = launcher.New().Bin(path).
			Set("window-size", fmt.Sprintf("%d,%d", cfg.Default.BrowserWidth, cfg.Default.BrowserHeight)).
			Set("user-data-dir", folderPath).
			Set("load-extension", extensionStr).
			Set("disable-extensions-http-throttling", "false").
			Set("allow-insecure-localhost", "1").
			Set("profile.default_content_setting_values.insecure_content", "1").
			//Set("auto-open-devtools-for-tabs", "true").
			Set("disable-features", "ExtensionsNetworkBlocking") // 禁用扩展网络限制)

		//Set("no-sandbox"). // 禁用沙盒
	} else {
		lg.InfoToFile(fmt.Sprintf("当前使用的是临时下载的浏览器\n"))
		l = launcher.New().
			Set("window-size", fmt.Sprintf("%d,%d", cfg.Default.BrowserWidth, cfg.Default.BrowserHeight)).
			Set("user-data-dir", folderPath).
			Set("load-extension", extensionStr).
			Set("disable-extensions-http-throttling", "false").
			Set("allow-insecure-localhost", "1").
			Set("profile.default_content_setting_values.insecure_content", "1").
			//Set("auto-open-devtools-for-tabs", "true").
			Set("disable-features", "ExtensionsNetworkBlocking") // 禁用扩展网络限制)
		//Set("no-sandbox"). // 禁用沙盒
	}

	//defaults.Show = false

	u := l.MustLaunch()
	browser = rod.New().ControlURL(u).SlowMotion(200 * time.Millisecond).MustConnect()
	// TODO 待验证会不会启动太多 goroutines

	// 打开一个空白页，防止关闭浏览器
	_ = browser.MustPage("about:blank")
	return browser, err
}
