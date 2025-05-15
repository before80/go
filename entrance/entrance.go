package entrance

import (
	"github.com/before80/go/tr"
	"github.com/before80/go/wind"
	"path/filepath"
	"strconv"
	"time"
)

func OpenUniqueMdFile(j int) {
	uniqueMdFilename := "do" + strconv.Itoa(j) + ".md"
	relUniqueMdFilePath := filepath.Join("markdown", uniqueMdFilename)
	absUniqueMdFilePath, _ := filepath.Abs(relUniqueMdFilePath)
	_ = tr.TruncFileContent(relUniqueMdFilePath)
	_ = wind.OpenTypora(absUniqueMdFilePath)
	time.Sleep(2 * time.Second)
}
