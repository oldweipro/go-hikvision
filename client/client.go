package main

import (
	"encoding/xml"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"strings"
)

type Probe struct {
	Uuid  string `xml:"Uuid"`
	Types string `xml:"Types"`
}

func main() {
	lAddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 37020,
	}
	rAddr := &net.UDPAddr{
		IP:   net.IPv4(239, 255, 255, 250),
		Port: 37020,
	}
	socket, err := net.DialUDP("udp", lAddr, rAddr)
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}
	defer func(socket *net.UDPConn) {
		err := socket.Close()
		if err != nil {
			fmt.Println("err:", err)
		}
	}(socket)
	uid := strings.ToUpper(uuid.NewV4().String())
	req := Probe{
		Uuid:  uid,
		Types: "inquiry",
	}
	marshal, err := xml.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	sendData := marshal
	_, err = socket.Write(sendData) // 发送数据
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}
	fmt.Println("扫描设备完成！")
}
