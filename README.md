# kwatcher
File watcher

## Example
```
watcher := kwatcher.NewWatcher()
watcher.AddFiles("/Users/sunxiangke/project/go/src/local/kwatcher/test")
watcher.SetCallback(func(modifyFiles []kwatcher.Path) {
  for _, path := range modifyFiles {
    fmt.Println(path + "changed")
  }
})
watcher.Run()
```
