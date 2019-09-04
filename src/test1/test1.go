package main

import (
	"fmt"
	"liba"
	"libb"
)

type name interface {
	call()
}

/*
id :
name:
*/
type Books struct {
	id    int
	name  string
	price int
	sn    string
	add   string
	pub   string
}

type Pics struct {
	name string
	path string
}

func (book Books) call() {
	fmt.Println("book.name")
}
func (pic Pics) call() {
	fmt.Println("pic.path")
}

func main() {
	var id int
	id = 1000
	var str = ", wangjr"
	str += ", hello go"

	fmt.Print("hello go, ", id, str)

	str = printHello("wangjr, ", "hello golang...")
	go printHello("ddd", "dddd")
	//fmt.Println(str)

	str = liba.Print("wangjr, ", "hello golang...")
	fmt.Println(str)
	str = libb.Print("wangjr, ", "hello golang...")
	fmt.Println(str)

	fmt.Println(Books{1, "golang", 100, "989292", "gz", "wangjr"})

	var nm name
	nm = new(Books)
	nm.call()

	nm = new(Pics)
	nm.call()

}
