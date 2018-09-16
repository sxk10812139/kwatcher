package kwatcher

import (
	"fmt"
	"io/ioutil"
	"time"
)

type Dir struct {
	Path string
}

func ParseDirs(dirs []Dir) map[string](map[string]time.Time) {
	res := make(map[string](map[string]time.Time))
	for _, d := range dirs {
		fileList, err := ioutil.ReadDir(d.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dirFileStatus := make(map[string]time.Time)
		for _, fileInfo := range fileList {
			dirFileStatus[d.Path+"/"+fileInfo.Name()] = fileInfo.ModTime()
		}
		res[d.Path] = dirFileStatus
	}
	return res
}
