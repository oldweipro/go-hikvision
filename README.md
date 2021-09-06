# Go-SADPTool

## 首先启动server，监听端口37020

接收设备发送到本机的udp请求

## 启动client，发送udp组播

设备会向发起的udp组播的IP地址（也就是本机）发送设备信息，由最先启动的服务端server去接收设备发出udp请求。