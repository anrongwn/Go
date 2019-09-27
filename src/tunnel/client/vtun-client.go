package main

import (
	"../general"
	"../tun"
	"encoding/binary"
	"fmt"
	"github.com/FlowerWrong/netstack/tcpip/header"
	"github.com/songgao/water"
	"go-share/jsontion"
	"go-share/logtion"
	"go-share/pathtion"
	"math/rand"
	"net"
	"os"
	"runtime"
	"time"
)

// Client Config
type VTunClientConfig struct {
	ConnectIp   string
	ConnectPort int
	DevName     string
	NodeAddr    string
	MTU         int
}

type TunnelClient struct {
	// config
	Config *VTunClientConfig
	// interface
	DevFace *water.Interface
	// server conn
	ClientConn *net.UDPConn
	BeatCount int
}

func main() {
	LogFile := pathtion.GetCurrentDirectory() + "/tun.log"
	//if common.PathExists(LogFile) {
	//	os.Rename(LogFile,fmt.Sprintf("%s_%s",LogFile,time.Now().Format("2006-01-02 15:04:05")))
	//}
	f, err := os.Create(LogFile)

	if err != nil {
		fmt.Println(f)
		return
	}
	logtion.InitJLogger(logtion.DEBUG, "jx", f)
	runtime.GOMAXPROCS(4)

	TunnelClient := new(TunnelClient)
	TunnelClient.Config = new(VTunClientConfig)
	err = jsontion.ParseJsonFile(pathtion.GetCurrentDirectory()+"/vtun.cfg", TunnelClient.Config)
	if err != nil {
		logtion.JLogger.Fatal("json file error", err)
	}

	TunnelClient.GetTunIPFromSever()

	go TunnelClient.ReadFromServe()

	TunnelClient.DevFace, err = tun.NewTun(TunnelClient.Config.MTU, TunnelClient.Config.DevName)
	//	RouteList = make(map[string]int64)
	//GetRouteList(ServerFace.Name())
	if err != nil {
		logtion.JLogger.Fatal(err)
	}

	packet := make([]byte, 2048)
	//packet[0] = 1

	//	IpHeader := (*header.IPv4)(unsafe.Pointer(&packet))
	var n int
	// var header
	for {
		n, err = TunnelClient.DevFace.Read(packet)
		if err != nil {
			logtion.JLogger.Errorf("recv tun  data error", TunnelClient.DevFace.Name(), err)
			continue
		}
		if n < general.READMINLEN {
			logtion.JLogger.Error("recv tun data is too small", n)
			continue
		}

		//Logger.Infof("tun recv srcip: %s dstip %s Protocol %d\n",
		//	IpHeader.SourceAddress(),
		//	IpHeader.DestinationAddress(),
		//	IpHeader.Protocol())

		_, err = TunnelClient.ClientConn.Write(packet[:n])

		if err != nil {
			logtion.JLogger.Errorf("send tun data error %s", err)
			//	return err
		}
		TunnelClient.BeatCount ++

	}
}

func (tc *TunnelClient) GetTunIPFromSever() string {
	//var err error
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(tc.Config.ConnectIp),
		Port: tc.Config.ConnectPort,
	})

	if err != nil {
		logtion.JLogger.Fatalf("connect error %s", err)
		return ""
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Duration(30) * time.Second))

	packet := make([]byte, 64)
	packet[0] = 0xFF //flags
	packet[1] = 0x01 // ver
	//packet[2] = uint8(len(NodeAddr)) // ver
	//copy(packet[3:len(NodeAddr)+3],NodeAddr)
	//var n int
	_, err = conn.Write(packet)
	if err != nil {
		logtion.JLogger.Fatalf("get tun ip From server write %s ", err)
	}
	_, err = conn.Read(packet)
	if err != nil {
		logtion.JLogger.Fatalf("get tun ip From server read %s", err)
	}
	return ""
}

func (tc *TunnelClient) BindLocalUDPPort() bool {
	var err error
	Ret := false
	count := 0
	for {
		if count > 5{
			break
		}

		//Port := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(50000)
		//logtion.JLogger.Info("bind port",Port)
		//Port +=  5000 + count
		//logtion.JLogger.Info("bind port",Port)

		tc.ClientConn, err = net.DialUDP("udp4",
			&net.UDPAddr{IP: nil, Port: rand.New(rand.NewSource(time.Now().UnixNano())).Intn(50000) + 5000 + count },
			&net.UDPAddr{
			IP:   net.ParseIP(tc.Config.ConnectIp),
			Port: tc.Config.ConnectPort,
		})

		count++

		if err == nil {
			Ret = true
			break
		}

		logtion.JLogger.Error("DialUDP error ", err,count)
	}
	return  Ret
}

func (tc *TunnelClient) ReadFromServe() {

	if !tc.BindLocalUDPPort() {	logtion.JLogger.Fatal("exit") }
	var err error

	go func() {
		//保持UDP会话
		test := make([]byte, 8)
		//test[0] = bytetion.IntToByte(rand.Intn(50000))

		for {
			binary.BigEndian.PutUint64(test, rand.Uint64())
			_, err = tc.ClientConn.Write(test)
			tc.BeatCount ++
			if err != nil {
				logtion.JLogger.Errorf("beat send udp data error %s", err)
			}

			time.Sleep(time.Duration(10) * time.Second)

			// if two minute beat count  have not be set 0
			if tc.BeatCount > 12 {
				logtion.JLogger.Fatal("have not recv udp packet in two minute")
			}
		}

	}()

	packet := make([]byte, 2048)
	//IpHeader := (*header.IPv4)(unsafe.Pointer(&packet))
	IpHeader := (header.IPv4)(packet)

	var n int

	for {
		n, err = tc.ClientConn.Read(packet)
		if err != nil {
			//Logger.Errorf("recv udp data error %s", err)
			continue
		}

		tc.BeatCount = 0 //set beat count  0

		if n < general.READMINLEN {
			//Logger.Error("recv udp data is too small",n)
			continue
		}

		if packet[0] != 0x45 {
			logtion.JLogger.Error("ip packet is error")
			continue
		}

		//Logger.Infof("udp recv srcip: %s dstip %s Protocol %d\n",
		//	IpHeader.SourceAddress(),
		//	IpHeader.DestinationAddress(),
		//	IpHeader.Protocol())
		_, err = tc.DevFace.Write(packet[:n])
		if err != nil {
			logtion.JLogger.Error(IpHeader.SourceAddress(), IpHeader.DestinationAddress(), "iface send data error", err)
			break
			//	return err
		}

	}

	return
}
