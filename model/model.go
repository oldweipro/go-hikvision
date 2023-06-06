package model

import (
	"encoding/xml"
	"sync"
)

type Probe struct {
	Uuid  string `xml:"Uuid"`
	Types string `xml:"Types"`
}

type Device struct {
	XMLName                 xml.Name `xml:"ProbeMatch"`
	Uuid                    string   `xml:"Uuid"`
	Types                   string   `xml:"Types"`
	DeviceType              int      `xml:"DeviceType"`
	DeviceDescription       string   `xml:"DeviceDescription"`
	DeviceSN                string   `xml:"DeviceSN"`
	CommandPort             int      `xml:"CommandPort"`
	HttpPort                int8     `xml:"HttpPort"`
	EHomeVer                string   `xml:"EHomeVer"`
	IPv4Address             string   `xml:"IPv4Address"`
	IPv4SubnetMask          string   `xml:"IPv4SubnetMask"`
	IPv4Gateway             string   `xml:"IPv4Gateway"`
	IPv6Address             string   `xml:"IPv6Address"`
	IPv6Gateway             string   `xml:"IPv6Gateway"`
	IPv6MaskLen             string   `xml:"IPv6MaskLen"`
	DHCP                    bool     `xml:"DHCP"`
	AnalogChannelNum        int8     `xml:"AnalogChannelNum"`
	DigitalChannelNum       int8     `xml:"DigitalChannelNum"`
	SoftwareVersion         string   `xml:"SoftwareVersion"`
	DSPVersion              string   `xml:"DSPVersion"`
	Encrypt                 bool     `xml:"Encrypt"`
	Salt                    string   `xml:"Salt"`
	BootTime                string   `xml:"BootTime"`
	DiskNumber              int8     `xml:"DiskNumber"`
	OEMInfo                 string   `xml:"OEMInfo"`
	Activated               bool     `xml:"Activated"`
	PasswordResetAbility    bool     `xml:"PasswordResetAbility"`
	ResetAbility            bool     `xml:"ResetAbility"`
	SyncIPCPassword         bool     `xml:"SyncIPCPassword"`
	PasswordResetModeSecond bool     `xml:"PasswordResetModeSecond"`
	DeviceLock              bool     `xml:"DeviceLock"`
	DHCPAbility             bool     `xml:"DHCPAbility"`
	SupportGUID             bool     `xml:"SupportGUID"`
	SupportSecurityQuestion bool     `xml:"SupportSecurityQuestion"`
	SupportIPv6             string   `xml:"supportIPv6"`
	SupportModifyIPv6       string   `xml:"supportModifyIPv6"`
	SupportHCPlatform       bool     `xml:"SupportHCPlatform"`
	//xml转换结构体异常: `strconv.ParseBool: parsing "flase": invalid syntax`
	//设备报文"flase"无法转换bool导致报错,无奈只能改为string,这真是海康的一个低级错误
	HCPlatformEnable         string `xml:"HCPlatformEnable"`
	IsModifyVerificationCode bool   `xml:"IsModifyVerificationCode"`
	SupportMailBox           bool   `xml:"SupportMailBox"`
	SupportEzvizUnbind       string `xml:"supportEzvizUnbind"`
}

type DeviceList struct {
	// 使用 sync.Mutex 进行加锁，避免多个 goroutine 并发写入 deviceList
	sync.Mutex
	Devices []Device
}
