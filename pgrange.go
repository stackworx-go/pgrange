package pgrange

func unquote(str string) string {
	return str[1 : len(str)-1]
}
