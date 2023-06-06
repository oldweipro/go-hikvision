# go-hikvision
go语言写的类似SADPTool工具，搜索海康设备。

## 运行

```bash
go run main.go
```

## 原理

项目运行后会监听端口`37020`，本机会向局域网中广播`探测数据`

```go
uuidString := strings.ToUpper(uuid.NewString())
req := model.Probe{
    Uuid:  uuidString,
    Types: "inquiry",
}
sendBytes, err := xml.Marshal(req)
```

接收到探测数据的设备会向发起的udp组播的IP地址（也就是本机）发送设备信息。

```go
n, _, err := conn.ReadFromUDP(data)
if err != nil {
    panic(err)
}
var ipc model.Device
err = xml.Unmarshal(data[:n], &ipc)
```

## 起源

本项目代码是根据海康SADPTool设备网络搜索工具，通过wireshark抓包分析而来。

PS：主要还是作者的电脑是Mac，运行不起来《SADPTool设备网络搜索工具》

![udp抓包截图](images/udp_capture_screenshot.png)