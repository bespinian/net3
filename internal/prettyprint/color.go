package prettyprint

const colorReset = string("\033[0m")

func asBold(str string) string {
	return string("\033[1m") + str + colorReset
}

func asGreen(str string) string {
	return string("\033[32m") + str + colorReset
}

func asRed(str string) string {
	return string("\033[31m") + str + colorReset
}
