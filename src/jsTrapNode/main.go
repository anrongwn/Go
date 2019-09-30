package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

var (
	gHostName   string
	signChannel = make(chan os.Signal, 1) //系统信号通道
	wg          sync.WaitGroup            //goroutine 计数器
)

func init() {
	flag.StringVar(&gHostName, "hostname", "127.0.0.1:9090", "ip:port")
}

func installSignalHandler() {
	signal.Notify(signChannel, os.Interrupt, os.Kill)
}

func main() {
	//分析输入参数
	flag.Parse()

	//安装系统信号handler
	go installSignalHandler()

	//test tun
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Interface name : ", ifce.Name())
	}
	var frame ethernet.Frame

	for {
		select {
		case s := <-signChannel:
			log.Println("Get system signal:", s)
			//logout.Println("Get system signal:", s)

			//listener.Close()

			goto EXIT
		default:
			frame.Resize(1500)
			n, err := ifce.Read([]byte(frame))
			if err != nil {
				log.Fatal(err)
			}
			frame = frame[:n]

			log.Println("frame len: ", n)
			log.Printf("Dst: %s\n", frame.Destination())
			log.Printf("Src: %s\n", frame.Source())
			log.Printf("Ethertype: % x\n", frame.Ethertype())
			log.Printf("Payload: % x\n", frame.Payload())
		}
	}

	//
	return

	//1、Int
	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(128))
		if err == nil {
			fmt.Println("rand.Int：", n, n.BitLen())
			//key[i] = n.Bytes()
			//fmt.Println("", string(n.Bytes()))
		}
	}
	fmt.Println("key=", string(key))

	//2、Prime
	//for i := 0; i < 32; i++ {
	p, err := rand.Prime(rand.Reader, 5)
	if err == nil {
		fmt.Println("rand.Prime：", p)
	}

	//}

	//data := []byte(`{"name=","wangjr"}`)
	data, err := GetAesKey(32)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
		strdata := string(data[:])
		fmt.Println(strdata)
	}

	code, err := RsaEncrypt(data)
	if err != nil {
		fmt.Println(err)
	} else {
		code64 := Base64URLCode(code)
		fmt.Println(code64)
		code64 = Base64Code(code)
		fmt.Println(code64)
	}

	for {
		select {
		case s := <-signChannel:
			log.Println("Get system signal:", s)
			//logout.Println("Get system signal:", s)

			//listener.Close()

			goto EXIT
		}
	}

EXIT:
	log.Println("Waiting gorouting exit ....")
	wg.Wait()
}
