package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-share/cmdtion"
	"go-share/jsontion"
	"go-share/logtion"
	"go-share/pathtion"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"../general"
	"../tun"
	"github.com/FlowerWrong/netstack/tcpip/header"
	"github.com/songgao/water"
)

type VTunServerConfig struct {
	ListenIp      string
	ListenPort    int
	DevName       string
	DevAddr       string
	PotAddr       string
	ClientTunAddr string
	MTU           int
}

type TunnelServer struct {
	// config
	Config *VTunServerConfig
	// interface
	DevFace *water.Interface
	// client udp socket list
	ClientSocketListMap sync.Map
	// client udp socket list peer
	ClientSocketListPeerMap sync.Map
	// server conn
	ServerConn *net.UDPConn
	//mutex sync.Mutex
	// attack host list
	//AttackHostList route.AttackHost
	// for
	StoreKey bytes.Buffer
	LoadKey  bytes.Buffer
}

func SetInterFaceRoute(DevName string) error {

	//ip route add default dev ppp0 table 10
	args := fmt.Sprintf("ip route add default dev %s table 10", DevName)
	err := cmdtion.CommandRun(args)
	if err != nil {
		logtion.JLogger.Fatalf("ip %s %s", args, err)
		return err
	}
	// ip rule
	output, err := cmdtion.CommandOutput("ip rule")
	if err != nil {
		logtion.JLogger.Fatalf("ip rule %s", err)
		return err
	}

	//fmt.Println(output)
	if strings.Contains(output, "from all fwmark 0x1 lookup") {
		logtion.JLogger.Info("have find fwmark 0x1 flags ")
		return nil
	}

	//ip rule add fwmark 0x1 table 10
	err = cmdtion.CommandRun("ip rule add fwmark 0x1 table 10")
	if err != nil {
		logtion.JLogger.Fatalf("rule add fwmark 0x1 table 10 %s", err)
		return err
	}
	return cmdtion.CommandRun("ip route flush cache")
}

func main() {
	LogFile := pathtion.GetCurrentDirectory() + "/tun.log"
	if pathtion.FileExists(LogFile) {
		os.Rename(LogFile, fmt.Sprintf("%s_%s", LogFile, time.Now().Format("2006-01-02 15:04:05")))
	}

	f, err := os.Create(LogFile)

	if err != nil {
		fmt.Println(f)
		return
	}

	defer f.Close()
	logtion.InitJLogger(logtion.DEBUG, "glog", f)
	runtime.GOMAXPROCS(runtime.NumCPU())

	TunnelServer := new(TunnelServer)
	//TunnelServer.AddKey.Bytes() =

	TunnelServer.Config = new(VTunServerConfig)
	err = jsontion.ParseJsonFile(pathtion.GetCurrentDirectory()+"/vtun.cfg", TunnelServer.Config)
	//fmt.Println("1111111")

	if err != nil {
		logtion.JLogger.Fatal("json file error", err)
	}

	//	route.GetLocalNetWorkCardList()
	//	go route.GetLocalRouteList()
	//	go route.GetLocalArpList()

	TunnelServer.DevFace, err = tun.NewTun(TunnelServer.Config.MTU, TunnelServer.Config.DevName)
	//	RouteList = make(map[string]int64)
	//GetRouteList(ServerFace.Name())
	if err != nil {
		logtion.JLogger.Fatal(err)
	}

	SetInterFaceRoute(TunnelServer.Config.DevName)

	go TunnelServer.ListenAndServe()
	go TunnelServer.DeleteIdleResource()

	packet := make([]byte, 2048)
	offset := 14
	IpHeader := (header.IPv4)(packet[offset:])

	var ClientSocket *net.UDPAddr

	TcpHeader := (header.TCP)(packet[20+offset:])
	UdpHeader := (header.UDP)(packet[20+offset:])
	// var header
	for {
		//copy(packet,(0))
		//packet[0:] = 0
		n, err := TunnelServer.DevFace.Read(packet[offset:])
		if err != nil {
			logtion.JLogger.Errorf("recv tun data %s error", TunnelServer.Config.DevName, err)
		}
		if n < general.READMINLEN {
			continue
		}

		PayloadLen := IpHeader.PayloadLength()

		if PayloadLen < general.IPPAYLOADMINLEN {
			//Logger.Error("tun data IpHeader.PayloadLength is too small",PayloadLen,IpHeader.SourceAddress(),IpHeader.DestinationAddress())
			continue
		}

		if PayloadLen > general.READMAXLEN {
			//Logger.Error("tun data IpHeader.PayloadLength is too big",PayloadLen,IpHeader.SourceAddress(),IpHeader.DestinationAddress())
			continue
		}

		//ClientSocket = nil

		//Logger.Infof(" tun recv srcip: %s dstip %s Protocol %d\n",
		//	IpHeader.SourceAddress(),
		//	IpHeader.DestinationAddress(),
		//	IpHeader.Protocol())

		switch IpHeader.Protocol() {
		case general.TCP:
			ClientSocket = TunnelServer.LoadClientSocket(packet[12+offset:16+offset], TcpHeader.SourcePort(), packet[16+offset:20+offset])
			break
		case general.UDP:
			ClientSocket = TunnelServer.LoadClientSocket(packet[12+offset:16+offset], UdpHeader.SourcePort(), packet[16+offset:20+offset])
			break
			//case general.ICMP: //icmp
			//	ClientSocket = TunnelServer.SearchClientSocket(packet[12+offset:16+offset], 0, packet[16+offset:20+offset])
			//	break
		}

		if ClientSocket != nil {
			_, err = TunnelServer.ServerConn.WriteToUDP(packet[offset:n+offset], ClientSocket)
			if err != nil {
				logtion.JLogger.Errorf("send tun data error %s", err)
			}

		} else {
			//TunnelServer.DevFace.Write(packet[offset : n+offset])
			//	TunnelServer.PotDataHandle(packet[0:n+offset], packet[12+offset:16+offset], UdpHeader.SourcePort(), packet[16+offset:20+offset], UdpHeader.DestinationPort())
		}

	}

}

func (tl *TunnelServer) DeleteIdleResource() {
	for {
		now := time.Now().Unix()
		tl.ClientSocketListPeerMap.Range(func(key, value interface{}) bool {
			if now-value.(int64) > 30 {
				tl.ClientSocketListPeerMap.Delete(key)
				tl.ClientSocketListMap.Delete(key)

			}
			return true
		})

		time.Sleep(time.Duration(5) * time.Second)
	}
}

func (tl *TunnelServer) StoreClientSocket(SrcIP net.IP, SrcPort uint16, DstIP net.IP, ClientSocket *net.UDPAddr) {

	tl.StoreKey.Reset()
	tl.StoreKey.Write(SrcIP)
	tl.StoreKey.WriteString(strconv.FormatInt(int64(SrcPort), 10))
	tl.StoreKey.Write(DstIP)
	tmp := tl.StoreKey.String()
	tl.ClientSocketListMap.Store(tmp, ClientSocket)
	tl.ClientSocketListPeerMap.Store(tmp, time.Now().Unix())
}

func (tl *TunnelServer) LoadClientSocket(SrcIP net.IP, SrcPort uint16, DstIP net.IP) *net.UDPAddr {
	tl.LoadKey.Reset()
	tl.LoadKey.Write(SrcIP)
	tl.LoadKey.WriteString(strconv.FormatInt(int64(SrcPort), 10))
	tl.LoadKey.Write(DstIP)

	tmp := tl.LoadKey.String()
	if v, ok := tl.ClientSocketListMap.Load(tmp); ok {
		tl.ClientSocketListPeerMap.Store(tmp, time.Now().Unix())
		return v.(*net.UDPAddr)
	}
	//logtion.JLogger.Error("not find value for key ", SrcIP, SrcPort, DstIP)
	return nil
}

func (tl *TunnelServer) ClientReply(pkt []byte, ClientSocket *net.UDPAddr) {

	switch pkt[1] {
	case 0x01:
		tl.ServerConn.WriteToUDP(pkt, ClientSocket)
		break
	default:
		logtion.JLogger.Error("not support version ", pkt[1])
		return
	}

	//copy(packet[20:],FaceAddr)

}

func (tl *TunnelServer) ListenAndServe() {
	// 创建监听
	var err error
	tl.ServerConn, err = net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.ParseIP(tl.Config.ListenIp),
		Port: tl.Config.ListenPort,
	})
	if err != nil {
		logtion.JLogger.Fatalf("listen error %s", err)
	}
	packet := make([]byte, 2048)
	var ClientSocket *net.UDPAddr
	var n int
	offset := 14
	IpHeader := (header.IPv4)(packet[offset:])
	TcpHeader := (header.TCP)(packet[20+offset:])
	UdpHeader := (header.UDP)(packet[20+offset:])
	BeatByte := make([]byte, 8)

	for {
		n, ClientSocket, err = tl.ServerConn.ReadFromUDP(packet[offset:])
		if err != nil {
			logtion.JLogger.Errorf("rcv udp data from udp error %s !", err)
			continue
		}

		if n == 8 {
			binary.BigEndian.PutUint64(BeatByte, rand.Uint64())
			tl.ServerConn.WriteToUDP(BeatByte, ClientSocket)
			continue
		}

		if n < general.READMINLEN {
			continue
		}
		flag := packet[offset]
		if flag == 0xFF && n == 64 {
			tl.ClientReply(packet[offset:n+offset], ClientSocket)
			continue
		}

		if flag != 0x45 {
			logtion.JLogger.Error("ip packet is error", packet[offset:16])
			continue
		}

		PayloadLen := IpHeader.PayloadLength()
		if PayloadLen < general.IPPAYLOADMINLEN {
			logtion.JLogger.Error("udp data IpHeader.PayloadLength is too small", PayloadLen, IpHeader.SourceAddress(), IpHeader.DestinationAddress())
			continue
		}

		if PayloadLen > general.READMAXLEN {
			logtion.JLogger.Error("udp data IpHeader.PayloadLength is too big", PayloadLen, IpHeader.SourceAddress(), IpHeader.DestinationAddress())
			continue
		}
		//Logger.Infof(" udp recv srcip: %s dstip %s Protocol %d\n",
		//	IpHeader.SourceAddress(),
		//	IpHeader.DestinationAddress(),
		//	IpHeader.Protocol(),ClientAddress,ClientAddress.IP.Equal(packet[12:16]))

		switch IpHeader.Protocol() {
		case general.TCP:
			//tl.ClientDataHandle(packet[0:n+offset], packet[12+offset:16+offset], TcpHeader.SourcePort(), packet[16+offset:20+offset], TcpHeader.DestinationPort(), ClientSocket)
			tl.StoreClientSocket(packet[16+offset:20+offset], TcpHeader.DestinationPort(), packet[12+offset:16+offset], ClientSocket)
			tl.DevFace.Write(packet[offset : n+offset])
			break
		case general.UDP:
			//tl.ClientDataHandle(packet[0:n+offset], packet[12+offset:16+offset], UdpHeader.SourcePort(), packet[16+offset:20+offset], UdpHeader.DestinationPort(), ClientSocket)
			tl.StoreClientSocket(packet[16+offset:20+offset], UdpHeader.DestinationPort(), packet[12+offset:16+offset], ClientSocket)
			tl.DevFace.Write(packet[offset : n+offset])
			break
			//case general.ICMP:
			//	tl.ClientDataHandle(packet[0:n+14], packet[12+offset:16+offset], 0, packet[16+offset:20+offset], 0, ClientSocket)
			//	break
		}

	}
	return
}
