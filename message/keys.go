package message 

import (
	"os"
	"io"
	"bytes"
)

type Keys struct{
	Public string
	Private string
	Passphrase []byte
}

func NewKeys (privatekeyPath string, publickeyPath string, passphrase []byte) (*Keys, error){
	k := new (Keys)
	err := k.LoadKey("public", publickeyPath)
	if err != nil {
		return nil, err
	}
	err = k.LoadKey("private", privatekeyPath)
	if err != nil {
		return nil, err
	}

	k.Passphrase = passphrase
	return k, nil
}
func (k *Keys) LoadKey(kind string, filepath string) error {
	buf := bytes.NewBuffer(nil)
	key, err := os.Open(filepath)
	if err != nil {
		return err
	}
	io.Copy(buf, key)
	key.Close()

	if kind == "public" {
		k.Public = string(buf.Bytes())
	}
	if kind == "private" {
		k.Private = string(buf.Bytes())
	}
	return nil
}
