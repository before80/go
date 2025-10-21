使用:  go install github.com/go-rod/rod/lib/utils/get-browser@latest 下载到 get-browser

查看比较好用的Chromium版本：https://github.com/GoogleChromeLabs/chrome-for-testing#json-api-endpoints

指定Chromium 版本和路径：
get-browser -v 137.0.7138.0 -p ./browser/
get-browser -v 138.0.7204.49 -p ./browser/

或者：.\get-browser.exe -v 137.0.7138.0 -p .\browser\
或者：.\get-browser.exe -v 138.0.7204.49 -p .\browser\


_ = ctx.LoadResponse(http.DefaultClient, true)

