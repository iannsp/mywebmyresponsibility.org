package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mywebmyresponsibility.org/core/message"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")
	pwd, _ := reader.ReadString('\n')
	myKeys, err := message.NewKeys("./privatekey.asc", "./publickey.asc", []byte(pwd))
	if err != nil {
		fmt.Println("Fail on read the Keys")
		log.Fatal(err)
	}
	messageStr := flag.String("message", "", "--NEWMESSAGE--")
	flag.Parse()
	err, m := message.EncryptPost(*messageStr, myKeys)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", m)
}
