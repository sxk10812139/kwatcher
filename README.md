## Example
```
watcher := kwatcher.NewWatcher()
watcher.SetInterval(time.Second)
watcher.AddFiles("path/to/file")
watcher.SetCallback(func(modifyFiles []kwatcher.Path) {
  for _, path := range modifyFiles {
    fmt.Println(path + "changed")
  }
})
watcher.Run()
```
## screenshot
![2018091715371192238739.png](http://pic.aipp.vip/2018091715371192238739.png)
