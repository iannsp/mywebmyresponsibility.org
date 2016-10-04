package message

import(
	"fmt"
	"golang.org/x/crypto/openpgp"
	"encoding/base64"
	"strings"
	"io/ioutil"
	"bytes"
	"log"

)
type Message struct{
	// kind define the route of the message into the system
	Kind string
	Content string
}


func NewMessage(kind string) *Message {
	return &Message{kind, ""}
}

// SetContent add content to the message.
func (m *Message)SetContent(content string) {
	m.Content = content
}

// Encode return the message string and need be prepared by the data service. 
func Encode(m *Message) string {
	return fmt.Sprintf("--NEWMESSAGE--\nKIND %s\nCONTENT %s\n",m.Kind, m.Content)

}

func Decode(message string) *Message {
	parts := strings.Split(message,"\n")
	m := NewMessage(strings.Split(parts[1]," ")[1])
	m.SetContent(strings.Replace(parts[2], "CONTENT ", "",1))
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

func DecryptPost(publicPostMessage string) (error, string) {
	// Read armored private key into type EntityList
	// An EntityList contains one or more Entities.
	// This assumes there is only one Entity involved
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(privateKey))
	if err != nil {
		log.Fatal(err)
		return err, "" 
	}
	entity := entitylist[0]
	fmt.Println("Private key from armored string:", entity.Identities)

	// Decrypt private key using passphrase
	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		fmt.Println("Decrypting private key using passphrase")
		err := entity.PrivateKey.Decrypt(passphrase)
		if err != nil {
			log.Println(err)
			return err, ""
		}
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			err := subkey.PrivateKey.Decrypt(passphrase)
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

	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entitylist, nil /* no prompt */, nil)
	if err != nil {
		log.Println(err, nil)
		return err, "" 
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	return  nil, string(bytes)
}

func EncryptPost(postMessage string) string {
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(publicKey))
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

	// Output as base64 encoded string
	bytes, err := ioutil.ReadAll(buf)
	return base64.StdEncoding.EncodeToString(bytes)
}

const passphrase = []byte("")
const publicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
-----END PGP PUBLIC KEY BLOCK-----`

const privateKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
-----END PGP PRIVATE KEY BLOCK-----`
