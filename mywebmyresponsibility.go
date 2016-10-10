package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"github.com/mywebmyresponsibility.org/core/message"
	"bufio"
	"os"

)
 var myKeys *message.Keys
func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
    	fmt.Fprintf(w, "<h1>Ola Enfermeira</h1><div></div>",)
}
// to-do adicionar ums requisicao de token via mensagem que seja uma senha de acesso do site privado


func MessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	var byteSizeLimit int64 = 1048576
	var messageStr string
	if(r.ContentLength < byteSizeLimit) {
		messageStr, err := ioutil.ReadAll(r.Body);
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err, mStr := message.DecryptPost(string(messageStr), myKeys)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err, m := message.DecodePost(mStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		for _,msg := range m {
			fmt.Fprintf(w,"%s",msg)
		}
	}
	fmt.Fprintf(w,"%s",messageStr)
}

func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	var byteSizeLimit int64 = 1048576
	if(r.ContentLength < byteSizeLimit) {
		messageStr, err := ioutil.ReadAll(r.Body);
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		err, m := message.EncryptPost(string(messageStr), myKeys)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Fprintf(w,"%s",m)
	}
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
 var err error
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")
	pwd, _ := reader.ReadString('\n')
 myKeys, err = message.NewKeys("./privatekey.asc", "./publickey.asc", []byte(pwd[0:len(pwd)-1]))

 if err!= nil {
	log.Fatal(err)
 }

 fmt.Println("My web my responsibility")
 http.HandleFunc("/view", viewHandler)
 http.HandleFunc("/message", MessageHandler)
 http.HandleFunc("/crypt", EncryptHandler)
 err = http.ListenAndServeTLS(":443", "./cert.pem", "./key.pem", nil)
 if err != nil {
	log.Fatal("ListenAndServe: ", err)
 }
}
