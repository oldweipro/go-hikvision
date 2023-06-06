package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/oldweipro/go-hikvision-scan/model"
	"net"
	"os"
	"strings"
)

const OUTPUT = "device.json"

func main() {
	//准备广播地址
	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:37020")
	if err != nil {
		panic(err)
	}

	//准备监听地址
	listenAddr, err := net.ResolveUDPAddr("udp4", ":37020")
	if err != nil {
		panic(err)
	}

	//创建连接
	conn, err := net.ListenUDP("udp4", listenAddr)
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
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			panic(err)
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
