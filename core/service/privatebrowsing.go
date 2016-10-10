package service

import (
	"errors"
	"github.com/mywebmyresponsibility.org/core/message"
)

var operations = []string{}

type PrivateBorwsing struct {
}

func (s PrivateBorwsing) CheckMessage(msg *message.Message) bool {
	return msg.Kind == "PRIVATEBROWSING"
}
func (s PrivateBorwsing) Process(msg *message.Message) error {
	return errors.New("Service PrivateBorwsing can't process")
}

func NewPrivateBorwsing() PrivateBorwsing {
	return PrivateBorwsing{}
}

/*

request a token to access private site with preferences
kind PRIVATEBROWSING
CONTENT generatetoken

revoke a token
kind PRIVATEBROWSING
CONTENT revoke token token


----- return
tokenstring
*/
