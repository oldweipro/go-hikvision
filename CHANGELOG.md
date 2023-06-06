## 2023-06-06

设置读取超时时间,否则会持续阻塞，一般情况下2-3秒就接收完了，如果时间太短，可按需分配

```go
err := conn.SetReadDeadline(time.Now().Add(time.Second * 2))
```