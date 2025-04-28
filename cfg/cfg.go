package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var canUseUniqueKeySlice = []string{"A"}

// DefaultConfig 定义整体 JSON 文件的结构
type DefaultConfig struct {
	WaitOpenNextPageSeconds       int    `mapstructure:"wait_open_next_page_seconds"`
	WaitTyporaSaveSeconds         int    `mapstructure:"wait_typora_save_seconds"`
	WaitTyporaCloseSeconds        int    `mapstructure:"wait_typora_close_seconds"`
	WaitTyporaOpenSeconds         int    `mapstructure:"wait_typora_open_seconds"`
	WaitTyporaCopiedToSaveSeconds int    `mapstructure:"wait_typora_copied_to_save_seconds"`
	UseProxy                      int    `mapstructure:"use_proxy"`
	ProxyScheme                   string `mapstructure:"proxy_scheme"`
	ProxyHost                     string `mapstructure:"proxy_host"`
	ProxyPort                     int    `mapstructure:"proxy_port"`
	ProxyUsername                 string `mapstructure:"proxy_username"`
	ProxyPassword                 string `mapstructure:"proxy_password"`
	UniqueMdFilepath              string `mapstructure:"unique_md_filepath"`
	BrowserWidth                  int    `mapstructure:"browser_width"`
	BrowserHeight                 int    `mapstructure:"browser_height"`
}

var Default DefaultConfig

func init() {
	var err error
	Default, err = getDefaultConfigInfo()
	if err != nil {
		panic(fmt.Sprintf("获取默认配置信息出现错误：%v", err))
	}
}

// GetDefaultConfigInfo 获取默认配置信息
func getDefaultConfigInfo() (defaultConfig DefaultConfig, err error) {
	viper.SetConfigName("Default")  // 配置文件名称（不包含扩展名）
	viper.SetConfigType("toml")     // 配置文件类型
	viper.AddConfigPath("./config") // 配置文件所在目录

	// 读取配置文件
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("配置文件未找到: %v\n", err)
		} else {
			fmt.Printf("读取配置文件时出错: %v\n", err)
		}
		return
	}
	if err = viper.Unmarshal(&defaultConfig); err != nil {
		fmt.Printf("解析配置信息到对象时出错: %v\n", err)
		return
	}

	return defaultConfig, err
}

//func JudgeUniqueKeyIsOk() bool {
//	return slices.Contains(canUseUniqueKeySlice, Default.UniqueKey)
//}

var noAuthStartTime time.Time
var noAuthPlanEndTime time.Time

func SetAuthTime(canUseMinutes int) {
	noAuthStartTime = time.Now()
	noAuthPlanEndTime = time.Now().Add(time.Duration(canUseMinutes) * time.Minute)
}

func GetNoAuthStartTime() time.Time {
	return noAuthStartTime
}

func JudgeIsAuthTimeout() bool {
	return time.Now().After(noAuthPlanEndTime)
}
