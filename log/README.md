# log

A structured logger for Go, based on [zap](https://github.com/uber-go/zap).

## Quick Start

```go
opts := &log.Newoptions()
log.Configure(opts)
defer func() {
    _ = log.Close()
}()

log.Debug("debug message")
log.Info("info message")
log.Warn("warn message")
log.Error("error message")

log.Debugf("%s message", "debug")
log.Infof("%s message", "info")
log.Warnf("%s message", "warn")
log.Errorf("%s message", "error")

log.Debugt("debug message", log.String("key1", "value1"))
log.Infot("info message", log.Int32("key2", 10))
log.Warnt("warn message", log.Bool("key3", false))
log.Errort("error message", log.Any("key4", "any"))
```
