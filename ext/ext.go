package ext

import (
	"archive/zip"
	"fmt"
	"github.com/before80/go/cfg"
	"github.com/before80/go/contants"
	"io"
	"os"
	"path/filepath"
)

// CreateProxyAuthExtension 创建代理扩展
func CreateProxyAuthExtension(absExtensionParentDir, subCacheFolderName string, forceGenerate bool) (extensionDir string, err error) {
	tempDir := contants.AppName + "_proxy_auth_extension"
	tempExtensionDir := filepath.Join(absExtensionParentDir, tempDir)

	if forceGenerate {
		_ = os.RemoveAll(tempExtensionDir)
	} else {
		exist, _ := pathExists(tempExtensionDir)
		if exist {
			return tempExtensionDir, nil
		}
	}

	// 创建临时目录
	if err = os.MkdirAll(tempExtensionDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 1. 生成 manifest.json
	manifest := []byte(fmt.Sprintf(`{
	"name": "%s Proxy Auth",
	"version": "1.0.0",
	"manifest_version": 3,
	"permissions": [
		"proxy",
		"webRequest",
		"webRequestAuthProvider"
	],
	"host_permissions": ["<all_urls>"],
	"background": {
		"service_worker": "background.js"
	}
}`, subCacheFolderName))

	if err = os.WriteFile(filepath.Join(tempExtensionDir, "manifest.json"), manifest, 0644); err != nil {
		return "", fmt.Errorf("写入manifest.json失败: %w", err)
	}

	// 2. 生成 background.js
	bgJS := []byte(fmt.Sprintf(`
		const authCredentials = {
            username: "%s",
            password: "%s"
        };

        const config = {
            mode: "fixed_servers",
            rules: {
                singleProxy: {
                    scheme: "%s",
                    host: "%s",
                    port: %d
                }
            }
        };

        chrome.proxy.settings.set({value: config, scope: 'regular'});

        chrome.webRequest.onAuthRequired.addListener((details, callback) => {
            callback({
              authCredentials: authCredentials
            });
          },
          { urls: ["<all_urls>"] },
          ['asyncBlocking']
        );`,
		cfg.Default.ProxyUsername, cfg.Default.ProxyPassword,
		cfg.Default.ProxyScheme, cfg.Default.ProxyHost, cfg.Default.ProxyPort))

	if err = os.WriteFile(filepath.Join(tempExtensionDir, "background.js"), bgJS, 0644); err != nil {
		return "", fmt.Errorf("写入background.js失败: %w", err)
	}

	//// 4. 打包为ZIP
	//if err = createZip(tempDir, zipPath); err != nil {
	//	return "",fmt.Errorf("创建ZIP失败: %w", err)
	//}
	//
	//// 5. 解压ZIP到目标目录
	//if err = unzip(zipPath, unpackedPath); err != nil {
	//	return "",fmt.Errorf("解压ZIP失败: %w", err)
	//}

	return tempExtensionDir, nil
}

// CreateJsHijackExtension 创建JavaScript劫持替换插件
func CreateJsHijackExtension(absExtensionParentDir, subCacheFolderName string, forceGenerate bool) (extensionDir string, err error) {
	tempDir := contants.AppName + "_js_hijack_extension"
	tempExtensionDir := filepath.Join(absExtensionParentDir, tempDir)

	if forceGenerate {
		_ = os.RemoveAll(tempExtensionDir)
	} else {
		exist, _ := pathExists(tempExtensionDir)
		if exist {
			return tempExtensionDir, nil
		}
	}

	// 创建临时目录
	if err = os.MkdirAll(tempExtensionDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 1. 生成 manifest.json
	manifest := []byte(fmt.Sprintf(`{
  "name": "%s JS Hijacker",
  "manifest_version": 3,
  "version": "1.0.0",
  "permissions": [
	"declarativeNetRequest"
  ],
  "host_permissions": [	
	"https://www.google-analytics.com/*",	
	"https://www.gstatic.com/*",	
	"https://www.google.com/*",
	"https://shop.dovolks.jp/*",
	"https://fonts.gstatic.com/*",
	"https://ajax.googleapis.com/*",
	"https://pi.pardot.com/*"
  ],
  "background": {
	"service_worker": "background.js"
  },
  "web_accessible_resources": [
	{
	  "resources": ["token.js"],
	  "matches": ["<all_urls>"]
	}
  ]
}`, subCacheFolderName))

	if err = os.WriteFile(filepath.Join(tempExtensionDir, "manifest.json"), manifest, 0644); err != nil {
		return "", fmt.Errorf("写入manifest.json失败: %w", err)
	}

	// 2. 生成 background.js
	bgJS := []byte(`chrome.declarativeNetRequest.updateDynamicRules({
  removeRuleIds: [1,2,3,4,5,6,7,8],  // 移除可能已存在的同 ID 规则
  addRules: [	
	{
	  "id": 2,
	  "priority": 1,
	  "action": {
		"type": "block"		
	  },
	  "condition": {
		"urlFilter": "https://www.google-analytics.com/analytics.js",
		"resourceTypes": ["script"]
	  }
	},
	{
	  "id": 3,
	  "priority": 1,
	  "action": {
		"type": "block"		
	  },
	  "condition": {
		"urlFilter": "https://www.google.com/recaptcha/api.js",
		"resourceTypes": ["script"]
	  }
	},
	{
	  "id": 4,
	  "priority": 1,
	  "action": {
		"type": "block"		
	  },
	  "condition": {
		"urlFilter": "https://fonts.gstatic.com",
		"resourceTypes": ["font"]
	  }
	},
	{
	  "id": 5,
	  "priority": 1,
	  "action": {
		"type": "block"		
	  },
	  "condition": {
		"urlFilter": "https://shop.dovolks.jp",
		"resourceTypes": ["image"]
	  }
	},
	{
	  "id": 6,
	  "priority": 1,
	  "action": {
		"type": "block"		
	  },
	  "condition": {
		"urlFilter": "https://ajax.googleapis.com",
		"resourceTypes": ["script"]
	  }
	},
	{
	  "id": 7,
	  "priority": 1,
	  "action": {
		"type": "block"		
	  },
	  "condition": {
		"urlFilter": "https://pi.pardot.com",
		"resourceTypes": ["script"]
	  }
	},
	{
	  "id": 8,
	  "priority": 1,
	  "action": {
		"type": "block"		
	  },
	  "condition": {
		"urlFilter": "https://www.gstatic.com",
		"resourceTypes": ["script"]
	  }
	}
  ]
}, () => {
  if (chrome.runtime.lastError) {
	console.log("Error updating rules:", chrome.runtime.lastError);
  } else {
	console.log("Redirect rule added successfully");
  }
});`)

	if err = os.WriteFile(filepath.Join(tempExtensionDir, "background.js"), bgJS, 0644); err != nil {
		return "", fmt.Errorf("写入background.js失败: %w", err)
	}

	//// 3. 复制依赖文件
	//for _, file := range []string{filepath.Join(contants.ReplaceHijackJsDir, "token.js")} {
	//	if err = copyFile(file, filepath.Join(tempExtensionDir, filepath.Base(file))); err != nil {
	//		return "", fmt.Errorf("复制文件失败: %w", err)
	//	}
	//}

	//// 4. 打包为ZIP
	//if err = createZip(tempDir, zipPath); err != nil {
	//	return "",fmt.Errorf("创建ZIP失败: %w", err)
	//}
	//
	//// 5. 解压ZIP到目标目录
	//if err = unzip(zipPath, unpackedPath); err != nil {
	//	return "",fmt.Errorf("解压ZIP失败: %w", err)
	//}

	return tempExtensionDir, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // 路径存在
	}
	if os.IsNotExist(err) {
		return false, nil // 路径不存在
	}
	return false, err // 发生了其他错误
}

// 复制文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

// 创建ZIP文件
func createZip(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 创建文件头
		relPath, _ := filepath.Rel(source, path)
		header, _ := zip.FileInfoHeader(info)
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

// 解压ZIP文件
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
