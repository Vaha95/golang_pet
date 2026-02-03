package main

var Urls map[string]string

func init() {
	Urls = make(map[string]string)
}

func main() {
	getEndpoints()
}
