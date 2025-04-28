package pydTr

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
