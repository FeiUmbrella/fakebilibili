package location

import (
	"log"
	"os"
)

// IsDir 判断传入的目录路径是否存在
func IsDir(disPath string) bool {
	st, err := os.Stat(disPath)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return st.IsDir()
}
