package main

import (
	"bufio"
	"bytes"
	"datapackage"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

//var g_wg sync.WaitGroup //全局

func startServer(port int, logout *log.Logger) int {
	//var r int = 0
	/*
		addr := "127.0.0.1:"
		addr += fmt.Sprintf("%d", port)

		listener, err := net.Listen("tcp", addr)
	*/
	listener, err := net.Listen("tcp", gHostName)
	if err != nil {
		fmt.Println("Error lister : ", err.Error())
		logout.Println("Error lister : ", err.Error())
		return -1
	}
	defer listener.Close()

	//
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		select {
		case s := <-c:
			log.Println("Got signal:", s)

			listener.Close()
			//os.Exit(0)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			logout.Println("Error accepting: ", err.Error())

			//
			return -2
		}

		//g_wg.Add(1)
		go doWork(conn, logout)

	}

	return 0
}

func doWork(conn net.Conn, logout *log.Logger) {
	log.Println("a new connection from:", conn.RemoteAddr())
	logout.Println("a new connection from:", conn.RemoteAddr())
	defer conn.Close()
	//defer g_wg.Done()

	for {
		/*
			buf := make([]byte, 1024)
			length, err := conn.Read(buf)
			if err != nil {
				fmt.Println("error read:", err.Error())
				logout.Println("error read:", err.Error())
				conn.Close()
				return
			}

			fmt.Println("Receive data from client:", string(buf[:length]))
			logout.Println("Receive data from client:", string(buf[:length]))
		*/

		scanner := bufio.NewScanner(conn) // reader为实现了io.Reader接口的对象，如net.Conn
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if !atEOF && data[0] == 'V' { // 由于我们定义的数据包头最开始为两个字节的版本号，所以只有以V开头的数据包才处理
				if len(data) > 4 { // 如果收到的数据>4个字节(2字节版本号+2字节数据包长度)
					length := int16(0)
					binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &length) // 读取数据包第3-4字节(int16)=>数据部分长度
					if int(length)+4 <= len(data) {                                    // 如果读取到的数据正文长度+2字节版本号+2字节数据长度不超过读到的数据(实际上就是成功完整的解析出了一个包)
						return int(length) + 4, data[:int(length)+4], nil
					}
				}
			} else if atEOF {
				//log.Println("Receive data from client EOF!")
				//panic("client EOF")

				return 0, nil, errors.New("client EOF")
			}
			return 0, nil, nil
		})
		// 打印接收到的数据包
		for scanner.Scan() {
			scannedPack := new(datapackage.Package)
			scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
			//log.Println(scannedPack)

			log.Println("Receive data from client:", scannedPack.String())
			logout.Println("Receive data from client:", scannedPack.String())

			/*
				e := scanner.Err()
				if e != nil {
					log.Println("scan err:", e.Error())
					break
				}
			*/

			/*
				if err := recover(); err != nil {
					log.Println(err) // 这里的err其实就是panic传入的内容，55
					break
				}
			*/
		}
	}

	log.Println("doWork end.")
	logout.Println("doWork end.")
}

var gHostName string

func init() {
	flag.StringVar(&gHostName, "hostname", "127.0.0.1:9090", "ip:port")
}

func main() {

	flag.Parse()

	port := 9090
	fmt.Println("====start port: ", gHostName, "=====")
	logfile, err := os.OpenFile("server.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend)

	if err != nil {
		fmt.Println("create log file fail :", err.Error())
		return
	}
	defer logfile.Close()

	logout := log.New(logfile, "[info] ", log.LstdFlags|log.Lshortfile)
	logout.Println("====start port: ", port, "=====")

	startServer(port, logout)

	//g_wg.Wait()
}
