package main

import(
	"fmt"
	"github.com/mywebmyresponsibility.org/message"
	"flag"
	"bufio"
	"os"
	"log"
)


func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")
	pwd, _ := reader.ReadString('\n')
	myKeys,err := message.NewKeys("./privatekey.asc", "./publickey.asc", []byte(pwd))
	if err!= nil {
		fmt.Println("Fail on read the Keys")
		log.Fatal(err)
	}
	messageStr := flag.String("message", "","--NEWMESSAGE--")
	flag.Parse()
	m := message.EncryptPost(*messageStr, myKeys)
	fmt.Printf("%s",m)
}
