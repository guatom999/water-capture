package utils

import (
	"fmt"
	"time"
)

func GenerateFileName() string {
	return fmt.Sprintf("pathum_snap_%s.png", time.Now().Format("2006-01-02_150405"))
}
