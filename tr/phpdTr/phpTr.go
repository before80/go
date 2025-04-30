package phpdTr

import (
	"os"
	"regexp"
)

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
		{regexp.MustCompile(`\\]]`), "]]"},
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
