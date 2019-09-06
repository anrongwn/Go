package main

import (
	"bufio"
	"bytes"
	"datapackage"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

func waitSignal(ch chan bool) {
	log.Println("begin wait signal....")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c

	log.Println("Got signal:", s)
	//log.Fatalln("Got signal:", s)
	//ch <- true
	os.Exit(1)
}

func main() {
	log.Println("start tcp client.")

	logfile, err := os.OpenFile("client.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend)

	//
	stopch := make(chan bool)
	go waitSignal(stopch)

	if err != nil {
		fmt.Println("create log file fail :", err.Error())
		return
	}
	defer logfile.Close()

	logout := log.New(logfile, "[info] ", log.LstdFlags|log.Lshortfile)
	logout.Println("start tcp client.")

	inputReader := bufio.NewReader(os.Stdin)
	/*
		fmt.Println("Please input your server ip:")
		ip, _ := inputReader.ReadString('\n')
		ip = strings.Trim(ip, "\n")

		port, _ := inputReader.ReadString('\n')
		port = strings.Trim(port, "\n")

		addr := ip
		addr += fmt.Sprintf("%s", port)
		fmt.Println(addr)
	*/
	addr := "127.0.0.1:9090"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial error : ", err.Error())
		return
	}

	//inputReader = bufio.NewReader(os.Stdin)
	fmt.Println("Please input your name:")
	srcCoder := mahonia.NewDecoder("utf-8")
	strFS := srcCoder.ConvertString("\r\n") //注意windows 下用\r\n 换行，linux下用\n
	clientName, _ := inputReader.ReadString('\n')
	inputClientName := strings.Trim(clientName, strFS)

	/*
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
	*/

	//send info to server until Quit
	for {
		fmt.Println("What do you send to the server? Type Q to quit.")

		content, _ := inputReader.ReadString('\n')
		inputContent := strings.Trim(content, strFS)

		if inputContent == srcCoder.ConvertString("Q") {
			fmt.Println("input Q to quit!")
			return
		}

		pack := &datapackage.Package{
			Version:        [2]byte{'V', '1'},
			Timestamp:      time.Now().Unix(),
			HostnameLength: int16(len(inputClientName)),
			Hostname:       []byte(inputClientName),
			TagLength:      4,
			Tag:            []byte("demo"),
			Msg:            []byte(("nowtime:" + time.Now().Format("2006-01-02 15:04:05") + "---" + inputContent)),
		}
		pack.Length = 8 + 2 + pack.HostnameLength + 2 + pack.TagLength + int16(len(pack.Msg))

		buf := new(bytes.Buffer)
		pack.Pack(buf)
		//pack.Pack(buf)

		_, err := conn.Write([]byte(buf.String()))
		if err != nil {
			fmt.Println("Error Write:", err.Error())
			return
		}

	}
}
