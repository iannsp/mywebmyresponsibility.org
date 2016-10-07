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

func DecodePost(postMessage string) (error, []*Message) {
	var DecodedMessages  []*Message
	messages := strings.Split(postMessage,"--NEWMESSAGE--\n")
	for _,onemessage:= range messages {
		if onemessage == "" { continue }
		err, decoded :=Decode("--NEWMESSAGE--\n"+onemessage)
		if err != nil {
			return err, DecodedMessages
		}
		DecodedMessages = append(DecodedMessages, decoded)
	}
	return nil, DecodedMessages
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

func EncryptPost(postMessage string, keys *Keys) (error, string) {
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(keys.Public))
	if err != nil {
		return err, ""
	}

	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entitylist, nil, nil, nil)
	if err != nil {
		return err, ""
	}
	_, err = w.Write([]byte(postMessage))
	if err != nil {
		return err, ""
	}
	err = w.Close()
	if err != nil {
		return err, ""
	}

	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return err, ""
	}


	return nil, base64.StdEncoding.EncodeToString(bytes)
}
