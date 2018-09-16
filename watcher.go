package kwatcher

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Path string

type Watcher struct {
	Interval   time.Duration
	WatchFiles []File
	FileStatus map[Path]time.Time
	Callback   func([]Path)
	mux        *sync.Mutex
}

func (w *Watcher) AddFiles(f string) {
	file := File{Path: Path(f)}
	w.WatchFiles = append(w.WatchFiles, file)
}

func (w *Watcher) SetInterval(t time.Duration) {
	w.Interval = t
}

func (w *Watcher) SetCallback(f func([]Path)) {
	w.Callback = f
}

func (w *Watcher) init() {
	for _, f := range w.WatchFiles {
		fileInfo, err := os.Stat(string(f.Path))
		if err != nil {
			fmt.Println(err)
			continue
		}
		w.FileStatus[f.Path] = fileInfo.ModTime()
	}
}

func (w *Watcher) checkFiles() []Path {
	var files []Path
	for path, time := range w.FileStatus {
		fileInfo, err := os.Stat(string(path))
		if err != nil {
			fmt.Println(err)
			continue
		}
		if fileInfo.ModTime() != time {
			files = append(files, path)
			w.FileStatus[path] = fileInfo.ModTime()
		}
	}
	return files
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
					w.Callback(modifyFiles)
				} else {
					fmt.Println("watcher running")
				}
			}()
		}
	}

}

type File struct {
	Path Path
}

func NewWatcher() *Watcher {
	return &Watcher{Interval: time.Second, FileStatus: make(map[Path]time.Time)}
}
