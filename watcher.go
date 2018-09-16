package kwatcher

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Watcher struct {
	Interval     time.Duration
	WatchFiles   []File
	WatchDirs    []Dir
	FileStatus   map[string]time.Time
	DirStatus    map[string](map[string]time.Time)
	FileCallback func([]string)
	DirCallback  func(map[string]([]string))
	mux          *sync.Mutex
}

func (w *Watcher) AddDir(d string) {
	dir := Dir{Path: d}
	w.WatchDirs = append(w.WatchDirs, dir)
}

func (w *Watcher) AddDirs(dirs []string) {
	for _, dir := range dirs {
		w.AddDir(dir)
	}
}

func (w *Watcher) AddFile(f string) {
	file := File{Path: f}
	w.WatchFiles = append(w.WatchFiles, file)
}

func (w *Watcher) AddFiles(files []string) {
	for _, file := range files {
		w.AddFile(file)
	}
}

func (w *Watcher) SetInterval(t time.Duration) {
	w.Interval = t
}

func (w *Watcher) SetFileCallback(f func([]string)) {
	w.FileCallback = f
}

func (w *Watcher) SetDirCallback(f func(map[string]([]string))) {
	w.DirCallback = f
}

func (w *Watcher) init() {
	w.FileStatus = ParseFiles(w.WatchFiles)
	w.DirStatus = ParseDirs(w.WatchDirs)
}

func (w *Watcher) checkFiles() []string {
	var files []string
	for path, time := range w.FileStatus {
		fileInfo, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if fileInfo.ModTime() != time {
			w.mux.Lock()
			files = append(files, path)
			w.FileStatus[path] = fileInfo.ModTime()
			w.mux.Unlock()
		}
	}
	return files
}

func (w *Watcher) checkDirs() map[string]([]string) {
	res := make(map[string]([]string))

	for dirPath, dirFileList := range w.DirStatus {
		dirChangedFiles := []string{}
		for filePath, fileModTime := range dirFileList {
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if fileInfo.ModTime() != fileModTime {
				dirChangedFiles = append(dirChangedFiles, filePath)
				w.mux.Lock()
				w.DirStatus[dirPath][filePath] = fileInfo.ModTime()
				w.mux.Unlock()
			}
		}

		if len(dirChangedFiles) > 0 {
			res[dirPath] = dirChangedFiles
		}

	}
	return res
}

func (w *Watcher) Run() {
	w.init()
	ticker := time.Tick(w.Interval)
	for {
		select {
		case <-ticker:
			go func() {
				modifyFiles := w.checkFiles()
				if len(modifyFiles) > 0 {
					w.FileCallback(modifyFiles)
				} else {
					fmt.Println("file watcher running")
				}
			}()
			go func() {
				modifyDirFiles := w.checkDirs()
				if len(modifyDirFiles) > 0 {
					w.DirCallback(modifyDirFiles)
				} else {
					fmt.Println("dir watcher running")
				}
			}()
		}
	}

}

func NewWatcher() *Watcher {
	return &Watcher{Interval: time.Second, FileStatus: make(map[string]time.Time), mux: new(sync.Mutex)} //mutex需要初始化
}
