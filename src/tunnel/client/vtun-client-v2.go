package main

import (
	"fmt"
	"github.com/FlowerWrong/netstack/tcpip/header"
	"net"
	"syscall"
)
import "../../share/logtion"

// Client Config

func main() {

	//SendHandler, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW) //syscall.ETH_P_ALL
	//SendHandler, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_IP) //syscall.ETH_P_ALL
	SendHandler, err := syscall.Socket(syscall.AF_INET,syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Println("Socket() error: ", err)
		return
	}
	err = syscall.SetsockoptInt(SendHandler, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)

	if err != nil {
		fmt.Println("SetsockoptInt error: ", err)
		return
	}

	//err = syscall.SetsockoptString(SendHandler,syscall.SOL_SOCKET,syscall.SO_BINDTODEVICE, "eth1")

	addr := syscall.SockaddrInet4{Port: 2000}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())

	err = syscall.Bind(SendHandler,&addr)

	if err != nil {
		fmt.Println("Bind() error: ", err)
		return
	}

	//err = syscall.Listen(SendHandler, 10)


	if err != nil {
		fmt.Println("Listen() error: ", err)
		return
	}


	packet := make([]byte, 2048)
	offset := 0
	IpHeader := (header.IPv4)(packet[offset:])

	TcpHeader := (header.TCP)(packet[20+offset:])
	// var header
	for {
		//copy(packet,(0))
		//packet[0:] = 0
	//	syscall.Accept(SendHandler)
		n, err := syscall.Read(SendHandler, packet[offset:])
		if err != nil {
			fmt.Printf("recv tun data %s error %s", "eth0", err)
		}
		if n < 20 {
			logtion.JLogger.Error("n <  20 ")
			continue
		}

		PayloadLen := IpHeader.PayloadLength()

		if PayloadLen < 8 {
			//Logger.Error("tun data IpHeader.PayloadLength is too small",PayloadLen,IpHeader.SourceAddress(),IpHeader.DestinationAddress())
			continue
		}

		if PayloadLen > 1680 {
			//Logger.Error("tun data IpHeader.PayloadLength is too big",PayloadLen,IpHeader.SourceAddress(),IpHeader.DestinationAddress())
			continue
		}

		switch IpHeader.Protocol() {
		case 6:
			if TcpHeader.DestinationPort() != 22 {

				fmt.Printf(" recv srcip: %s srcport %d dstip %s dstport %d\n",
					IpHeader.SourceAddress(),
					TcpHeader.SourcePort(),
					IpHeader.DestinationAddress(),
					TcpHeader.DestinationPort())
			}

			break
		}
		//ClientSocket = nil

	}
}


