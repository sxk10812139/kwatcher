package kwatcher

import (
	"fmt"
	"os"
	"time"
)

type File struct {
	Path string
}

func ParseFiles(files []File) map[string]time.Time {
	res := make(map[string]time.Time)
	for _, f := range files {
		fileInfo, err := os.Stat(f.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		res[f.Path] = fileInfo.ModTime()
	}
	return res
}
