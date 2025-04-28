package lg

import (
	"fmt"
	"github.com/before80/go/contants"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var logger *slog.Logger
var logFile *os.File
var removeLink bool

func init() {
	var err error
	removeLink = false
	err = os.MkdirAll(contants.LogFolderName, 0777)
	if err != nil {
		panic(fmt.Sprintf("无法创建%s目录：%v\n", contants.LogFolderName, err))
	}

	logFileName := fmt.Sprintf("%s_%s%s", contants.LogFileNamePrefix, time.Now().Format("2006_0102_1504"), ".txt")
	logFile, err = os.OpenFile(filepath.Join(contants.LogFolderName, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// 若不能打开日志文件，则不能使用 LogErrorToFile() 函数
		panic(fmt.Sprintf("在log模块中发生错误：%v", err))
	}

	// 配置日志记录器，将日志写入文件
	logger = slog.New(slog.NewJSONHandler(logFile, nil))
}

func SetRemoveLink(b bool) {
	removeLink = b
}

var urlRegex = regexp.MustCompile(`(?i)\b(https?://|www\.)[a-z0-9.-]+\.[a-z]{2,}(:\d+)?(/[-\w.%&?=/]*)?`)

func removeURLs(s string) string {
	cleaned := urlRegex.ReplaceAllStringFunc(s, func(match string) string {
		if u, err := url.Parse(match); err == nil && (u.Scheme == "http" || u.Scheme == "https" || strings.Contains(match, "www.")) {
			return "<URL>" //替换 URL
		}
		return match
	})
	//cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ") // 规范化空白
	//return strings.TrimSpace(cleaned)
	return cleaned
}

func InfoToFile(str string) {
	//logLock.Lock()
	//defer logLock.Unlock()
	// 将日志写入文件
	// 记录到日志文件中的字符串，需要去除尾部的换行符
	if removeLink {
		str = removeURLs(str)
		logger.Info(strings.TrimSuffix(str, "\n"))
	} else {
		logger.Info(strings.TrimSuffix(str, "\n"))
	}
}

func InfoToFileAndStdOut(str string) {
	if removeLink {
		str = removeURLs(str)
		logger.Info(strings.TrimSuffix(str, "\n"))
		fmt.Printf(str)
	} else {
		logger.Info(strings.TrimSuffix(str, "\n"))
		fmt.Printf(str)
	}
}

func ErrorToFile(str string) {
	//logLock.Lock()
	//defer logLock.Unlock()
	// 将日志写入文件
	// 记录到日志文件中的字符串，需要去除尾部的换行符
	if removeLink {
		str = removeURLs(str)
		logger.Error(strings.TrimSuffix(str, "\n"))
	} else {
		logger.Error(strings.TrimSuffix(str, "\n"))
	}
}

func ErrorToFileAndStdOutWithSleepSecond(str string, seconds int) {
	if removeLink {
		str = removeURLs(str)
		logger.Error(strings.TrimSuffix(str, "\n"))
		fmt.Printf(str)
	} else {
		logger.Error(strings.TrimSuffix(str, "\n"))
		fmt.Printf(str)
	}
	time.Sleep(time.Duration(seconds) * time.Second)
}
