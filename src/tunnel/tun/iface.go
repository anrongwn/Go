package tun

import (
	"fmt"
	"net"

	//"os/exec"
	"../../share/logtion"
	"../../share/cmdtion"

	//"bufio"
	//"bytes"
	"github.com/songgao/water"
	//"io"
	"os"
	//"strings"
	//"time"
	//"sync"
)


func NewTun(MTU int,DevName string) (DevFace *water.Interface, err error) {

	DevFace, err = water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams:water.PlatformSpecificParams{
			Name:DevName,
			Persist:false,
			Permissions:nil,
		},
	})
	if err != nil {
		logtion.JLogger.Fatalf("create interface error %s", err)
	}

	logtion.JLogger.Infof("interface %s created", DevFace.Name())
	args := fmt.Sprintf("ip link set dev %s up mtu %d qlen 100", DevName, MTU)
	err = cmdtion.CommandRun(args)
	if err != nil {
		logtion.JLogger.Fatalf("ip %s == %s", args, err)
		return nil, err
	}

	//Logger.Info("face addr",DevAddr,len(DevAddr))
	//ip, subnet, err := net.ParseCIDR(DevAddr)
	//
	//if err != nil {
	//	Logger.Fatalf("ParseCIDR %s",err)
	//}
	//
	//err = SetTunIP(iface, ip, subnet)
	//if err != nil {
	//	Logger.Fatalf("set %s ip  error %s", iface.Name(), err,ip,subnet)
	//}

	return DevFace, nil
}

func SetTunNet(DevFace *water.Interface, ip net.IP, subnet *net.IPNet) (err error) {
	//if ip == nil {
	//	Logger.Fatalf("set tun ip ip = nil")
	//}
	ip = ip.To4()
	logtion.JLogger.Debugf("%v", ip)
	//if ip[3]%2 == 0 {
	//	return invalidAddr
	//}

	peer := net.IP(make([]byte, 4))
	copy([]byte(peer), []byte(ip))
	peer[3]++
	//	tun_peer = peer

	args := fmt.Sprintf("ip addr add dev %s local %s peer %s", DevFace.Name(), ip, peer)
	err = cmdtion.CommandRun(args)
	if err != nil {
		os.Exit(1)
		return err
	}

	args = fmt.Sprintf("ip route add %s via %s dev %s", subnet, peer, DevFace.Name())
	//err = share.common.
	err = cmdtion.CommandRun(args)
	if err != nil {
		logtion.JLogger.Fatalf("ip addr add dev via peer error %s", err)
		return err
	}

	return err
}
//4e:8a:21:87:39:8a
//ifconfig eth0 hw ether 00:0C:29:36:97:20
//func SetTunMac()  {
//
//	err := cmdtion.CommandRun(args)
//}