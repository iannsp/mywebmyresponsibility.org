package main

import(
	"fmt"
	"github.com/mywebmyresponsibility.org/message"
	"flag"
)


func main() {
	messageStr := flag.String("message", "","--NEWMESSAGE--")
	flag.Parse()
	
	m := message.EncryptPost(*messageStr)
	fmt.Printf("%s",m)
}
