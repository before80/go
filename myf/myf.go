package myf

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TruncFileContent 清空指定文件中的内容
func TruncFileContent(filePath string) (err error) {
	var f *os.File
	// 以写入模式打开文件，若文件不存在则创建，若存在则清空内容
	f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("清空%s文件内容出现错误：%v", filePath, err)
	}
	// 确保文件在函数结束时关闭
	defer func() {
		_ = f.Close()
	}()
	return nil
}

// AddVersionInfoToMenu 添加版本信息到菜单中
func AddVersionInfoToMenu(menu string, code string) (versionInfo string) {
	// 从 code 中寻找版本信息
	// (C11 前) 替换成 -> 11 F
	// (C11 起)  替换成 -> 11+
	// (C11 弃用) 替换成 -> 11 D
	// (C11 移除) 替换成 -> 11 R
	// (C23 前) 替换成 -> 23 F
	// (C23 起)	替换成 -> 23+
	// (C23 弃用) 替换成 -> 23 D
	// (C23 移除) 替换成 -> 23 R

	return
}

// preLine 将所有以 >=2 个空白符开头的行，拼接到上一行。
// 保留其他行原样，最后用 "\n" 重新组装。
func preLine(src string) string {
	//fmt.Printf("src=%v\n\n", src)
	var out []string
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
			// 认为这是上一行的续行，去掉前导空白，拼接
			if len(out) > 0 {
				out[len(out)-1] += " " + strings.TrimSpace(line)
				continue
			}
		}

		// 否则，新起一行
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// FindIdentifierVersion 查找标识符的版本号
func FindIdentifierVersion(identifier, src string) (string, error) {
	src = preLine(src)
	// 1) 找到含有 identifier 且带注释的那一行
	// (?m)：开启多行模式，让 ^ … $ 分别匹配每一行的行首/行尾
	//
	// \b ... \b：确保完整的单词匹配
	//
	// .*//.*：该行必须包含 // 注释
	lineRe := regexp.MustCompile(`(?m)^.*\b` + regexp.QuoteMeta(identifier) + `\b.*//.*$`)
	line := lineRe.FindString(src)
	if line == "" {
		return "", errors.New("identifier not found or no comment")
	}
	//fmt.Printf("line=%s\n", line)
	//fmt.Println()

	// 2) 拆出注释部分（"// ..."）
	idx := strings.Index(line, "//")
	if idx < 0 {
		return "", errors.New("no comment found on matched line")
	}
	comment := line[idx+2:] // 去掉 "//"
	//fmt.Printf("comment=%s\n", comment)
	//fmt.Println()

	// 3) 提取所有 "(…)"
	parenRe1 := regexp.MustCompile(`\([^)]*\)`)
	all := parenRe1.FindAllString(comment, -1)
	if len(all) == 0 {
		return "", errors.New("no parenthesized version found")
	}
	if len(all) > 1 {
		parenRe2 := regexp.MustCompile(`\(\d+\)`)
		if len(parenRe2.FindAllString(all[0], -1)) > 0 {
			return strings.TrimSpace(strings.Join(all[1:], " ")), nil
		}
	}
	return all[len(all)-1], nil
}

// ReplaceMarkdownFileContent 替换文件中的特定内容
func ReplaceMarkdownFileContent(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	replacements := []struct {
		pattern     *regexp.Regexp
		replacement string
	}{
		{regexp.MustCompile("```\\s*?//"), "```python\n//"},
		{regexp.MustCompile("```\\s*?>>>"), "```python\n>>>"},
		{regexp.MustCompile("```\\s*?\ndef"), "```python\ndef"},
		{regexp.MustCompile("```\\s*?\nimport"), "```python\nimport"},
		{regexp.MustCompile("```\\s*?\ncase"), "```python\ncase"},
		{regexp.MustCompile("```\\s*?\nif"), "```python\nif"},
		{regexp.MustCompile("```\\s*?\nclass"), "```python\nclass"},
		{regexp.MustCompile("```\\s*?\nPoint"), "```python\nPoint"},
		{regexp.MustCompile("```\\s*?\nmatch"), "```python\nmatch"},
		{regexp.MustCompile("```\\s*?\nfrom"), "```python\nfrom"},
		{regexp.MustCompile("```\\s*?\nparrot"), "```python\nparrot"},
		{regexp.MustCompile("```\\s*?\ncheeseshop"), "```python\ncheeseshop"},
		{regexp.MustCompile("```\\s*?\n\\$"), "```sh\n$"},
		{regexp.MustCompile("&zeroWidthSpace;"), "​\t"},
	}

	modified := false
	newContent := string(content)
	for _, r := range replacements {
		if r.pattern.MatchString(newContent) {
			newContent = r.pattern.ReplaceAllString(newContent, r.replacement)
			if !modified {
				modified = true
			}
		}
	}

	if modified {
		err = os.WriteFile(filePath, []byte(newContent), 0666)
		//fmt.Println("1")
		if err != nil {
			return false, err
		}
		//fmt.Println("2")
	}

	return modified, nil
}
