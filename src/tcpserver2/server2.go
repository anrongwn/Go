package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	gHostName   string
	signChannel = make(chan os.Signal, 1)
	//exitChannel = make(chan int)
	wg     sync.WaitGroup
	logout *log.Logger
)

func init() {
	flag.StringVar(&gHostName, "hostname", "127.0.0.1:9090", "ip:port")
}

func main() {
	go installSignalHandler()
	flag.Parse()

	/*
		fruitarray := [...]string{"apple", "orange", "grape", "mango", "water melon", "pine apple", "chikoo"}
		fruitslice := fruitarray[1:2]

		for i, v := range fruitarray {
			log.Println(i, v)
		}
		for i, v := range fruitslice {
			log.Println(i, v)
		}
	*/

	log.Println("==== start ", gHostName, " service ====")
	logfile, err := os.OpenFile("server2.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Println("create log file fail :", err.Error())
		return
	}
	logout = log.New(logfile, "[info] ", log.LstdFlags|log.Lshortfile)
	defer func() {
		logfile.Close()
	}()
	logout.Println("==== start ", gHostName, " service ====")

	//
	tcpAddr, err := net.ResolveTCPAddr("tcp4", gHostName)
	if err != nil {
		tmp := "ip addr error,"
		tmp += err.Error()
		logout.Println(tmp)
		log.Fatalln(tmp)
	}

	log.Println("=== Listening ", gHostName, "....")
	logout.Println("=== Listening ", gHostName, "....")
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		tmp := "net.listen error,"
		tmp += err.Error()
		logout.Println(tmp)
		log.Fatalln(tmp)
	}

	defer func() {
		listener.Close()
	}()

	wg.Add(1)
	go accept(listener)

	///
	for {
		select {
		case s := <-signChannel:
			log.Println("Get system signal:", s)
			logout.Println("Get system signal:", s)

			listener.Close()

			goto EXIT
		}
	}

EXIT:
	log.Println("Waiting gorouting exit ....")
	wg.Wait()

}

func installSignalHandler() {
	signal.Notify(signChannel, os.Interrupt, os.Kill)
}

func accept(listener *net.TCPListener) {
	defer wg.Done()
	ctx, cancleServer := context.WithCancel(context.Background())

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			tmp := "accept error : "
			tmp += err.Error()
			log.Println(tmp)
			logout.Println(tmp)

			cancleServer()
			return
		} else {
			wg.Add(1)
			go connHandler(ctx, connection)
		}

	}
}

func connHandler(ctx context.Context, conn *net.TCPConn) {
	defer func() {
		conn.Close()
		wg.Done()
	}()

	tmp := "new connection sesion :"
	tmp += conn.RemoteAddr().String()

	log.Println(tmp)
	logout.Println(tmp)

	for {
		select {
		case <-(ctx).Done():
			log.Println("connHandler ctx.Done exit")
			return
		default:
			//设置读取timeout
			conn.SetReadDeadline(time.Now().Add(10 * time.Microsecond))

			//
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				ec := err.Error()
				if err == io.EOF {
					tmp = conn.RemoteAddr().String()
					tmp += " error:"
					tmp += ec
					log.Println(tmp)
					logout.Println(tmp)
					break
					//neterr, ok := err.(net.Error) 用来判断err interface 变量是否为net.Error 类型，如果是
					//neterr 为net.Error 的值，ok=true
				} else if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					//log.Println("read timeout")
					continue
				}

			} else if len > 0 {
				data := string(buf)
				log.Println("recive ", conn.RemoteAddr().String(), " data: ", data)
				logout.Println(tmp)
			}
		}

	}

}
