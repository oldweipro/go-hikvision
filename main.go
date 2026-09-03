package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/oldweipro/go-hikvision/model"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"strings"
	"time"
)

const OUTPUT = "device.json"

// newListenConn 创建 UDP 监听连接。
// ifaceName 不为空时监听该网卡的 IPv4 地址，并强制组播探测从该网卡发出；
// 为空时监听所有网卡（默认行为）。
func newListenConn(ifaceName string) (*net.UDPConn, error) {
	if ifaceName == "" {
		listenAddr, err := net.ResolveUDPAddr("udp4", ":37020")
		if err != nil {
			return nil, err
		}
		return net.ListenUDP("udp4", listenAddr)
	}

	ifi, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}

	var ifaceIP net.IP
	addrs, err := ifi.Addrs()
	if err != nil {
		return nil, err
	}
	for _, a := range addrs {
		ip, _, e := net.ParseCIDR(a.String())
		if e == nil && ip.To4() != nil {
			ifaceIP = ip
			break
		}
	}
	if ifaceIP == nil {
		return nil, fmt.Errorf("网卡 %s 没有 IPv4 地址", ifaceName)
	}

	//只绑定该网卡，只接收发往该网卡 IP 的回复
	listenAddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(ifaceIP.String(), "37020"))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		return nil, err
	}

	//强制组播探测从指定网卡发出，而不是系统默认路由的网卡
	pc := ipv4.NewPacketConn(conn)
	if err := pc.SetMulticastInterface(ifi); err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	//命令行参数指定搜索使用的网卡，例如 -i eth3；不指定则使用所有网卡
	ifaceName := flag.String("i", "", "指定搜索使用的网卡名称，例如 eth3")
	flag.Parse()

	//准备广播地址
	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:37020")
	if err != nil {
		panic(err)
	}

	//创建连接
	conn, err := newListenConn(*ifaceName)
	if err != nil {
		panic(err)
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	//向广播地址发送探测数据
	uuidString := strings.ToUpper(uuid.NewString())
	req := model.Probe{
		Uuid:  uuidString,
		Types: "inquiry",
	}
	sendBytes, err := xml.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = conn.WriteToUDP(sendBytes, addr)
	if err != nil {
		panic(err)
	}

	//接收回复数据
	deviceList := &model.DeviceList{}
	// 判断文件是否存在
	_, err = os.Stat(OUTPUT)
	if os.IsNotExist(err) { // 不存在则创建文件
		err := os.WriteFile(OUTPUT, nil, 0644)
		if err != nil {
			return
		}
	}
	for {
		data := make([]byte, 2048)
		// 设置读取超时时间,否则会持续阻塞，一般情况下2-3秒就接收完了，如果时间太短，可按需分配
		err := conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		if err != nil {
			return
		}
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			return
		}
		var ipc model.Device
		err = xml.Unmarshal(data[:n], &ipc)
		if err != nil {
			fmt.Println("xml转换结构体异常：", err.Error())
		}
		//打印回复数据
		deviceList.Lock()
		deviceList.Devices = append(deviceList.Devices, ipc)
		deviceList.Unlock()
		// 编码为JSON格式
		jsonBytes, err := json.MarshalIndent(deviceList.Devices, "", "    ")
		if err != nil {
			fmt.Println("编码 JSON 时出错:", err)
			return
		}
		// 写入文件
		err = os.WriteFile(OUTPUT, jsonBytes, 0644)
		if err != nil {
			fmt.Println("写入文件错误:", err)
			return
		}
	}
}
