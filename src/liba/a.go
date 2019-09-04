package liba

func Print(param1 string, param2 string) string {
	var str = "liba"
	var str2 = param1
	str2 += ", "
	str2 += param2

	str += str2

	return str
}
