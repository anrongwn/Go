package main

import "fmt"

func printHello(param1 string, param2 string) string {
	var ret = (param1 + param2)

	fmt.Print(ret)

	return ret

}
