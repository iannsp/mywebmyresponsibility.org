package message

import(
	"fmt"
	"golang.org/x/crypto/openpgp"
	"encoding/base64"
	"strings"
	"io/ioutil"
	"bytes"
	"log"
	"time"
	"strconv"

)
type Message struct{
	// kind define the route of the message into the system
	Kind string
	Content string
	DateTime time.Time
}


func NewMessage(kind string) *Message {
	return &Message{kind, "", time.Now()}
}

// SetContent add content to the message.
func (m *Message)SetContent(content string) {
	m.Content = strings.TrimRight(content, "\n")
}

// Encode return the message string and need be prepared by the data service. 
func Encode(m *Message) string {
	return fmt.Sprintf("--NEWMESSAGE--\nDATE %s\nKIND %s\nCONTENT %s\n",m.DateTime.Format("20060102150405"), m.Kind, m.Content)

}

func Decode(message string) *Message {
	log.Print(message+"\n---\n")
        parts := strings.Split(message,"\n")
	if strings.Split(parts[1]," ")[0] != "DATE" ||
	   strings.Split(parts[2]," ")[0] != "KIND" ||
	   strings.Split(parts[3]," ")[0] != "CONTENT" {
		log.Fatal("vix")
	}
	intdate, err := strconv.ParseInt(strings.Split(parts[1]," ")[1], 10, 64)
	if err != nil {
		log.Println("Fail parse Date:"+strings.Split(parts[1]," ")[1])
	}
	unixDate := time.Unix(intdate, 0)
	kind := strings.Split(parts[2]," ")[1]
	content := strings.Replace(message,"--NEWMESSAGE--\n"+parts[1]+"\n"+parts[2]+"\nCONTENT ","",1)

	m := NewMessage(kind)
	m.DateTime = unixDate 
	m.SetContent(content)
	return m
}

func DecodePost(postMessage string) []*Message{
	var DecodedMessages  []*Message
	messages := strings.Split(postMessage,"--NEWMESSAGE--\n")
	for _,onemessage:= range messages {
		if onemessage == "" { continue }
		DecodedMessages = append(DecodedMessages, 
			Decode("--NEWMESSAGE--\n"+onemessage))
	}
	return DecodedMessages
}

func DecryptPost(publicPostMessage string, keys *Keys) (error, string) {
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(keys.Private))
	if err != nil {
		log.Fatal(err)
		return err, "" 
	}
	entity := entitylist[0]

	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		err := entity.PrivateKey.Decrypt(keys.Passphrase)
		if err != nil {
			log.Println(err)
			return err, ""
		}
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			err := subkey.PrivateKey.Decrypt(keys.Passphrase)
			if err != nil {
				fmt.Println("failed to decrypt subkey")
				log.Println(err)
				return err, ""
			}
		}
	}

	dec, err := base64.StdEncoding.DecodeString(publicPostMessage)
	if err != nil {
		log.Println(err, nil)
		return err, ""
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entitylist, nil, nil)
	if err != nil {
		log.Println(err, nil)
		return err, "" 
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	return  nil, string(bytes)
}

func EncryptPost(postMessage string, keys *Keys) string {
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(keys.Public))
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entitylist, nil, nil, nil)
	if err != nil {
	}

	_, err = w.Write([]byte(postMessage))
	if err != nil {
	}
	err = w.Close()
	if err != nil {
	}

	bytes, err := ioutil.ReadAll(buf)
	return base64.StdEncoding.EncodeToString(bytes)
}
