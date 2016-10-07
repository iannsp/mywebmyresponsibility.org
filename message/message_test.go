package message

import (
	"github.com/mywebmyresponsibility.org/message"
	"time"
	"testing"
)
var  tests = []struct{
	Kind string
	Content string
	Dt	time.Time
	}{

	{"firstlineempty","\nsegunda linha\nainda nao sei se eh por aqui que esta a resposta", time.Now()},
	{"note","ainda nao sei se eh por aqui que esta a resposta", time.Now()},
	{"link","http://zeromq.org/", time.Now()},
	{"calendar","god speed you, Pieter Hintjens\nDATE 08:35 - 4 de out de 2016", time.Now()},
}
func formatExpectedMessage(kind string, content string, dt time.Time) string {
	return  "--NEWMESSAGE--\n" +
		"DATE "+ dt.Format("20060102150405")+"\n"+
		"KIND "+kind+"\n" +
		"CONTENT "+content+"\n"
}

func createMessage(kind string, content string, dt time.Time) *message.Message {
	m := message.NewMessage(kind)
	m.SetContent(content)
	m.DateTime = dt
	return m
}
func TestEncode(t *testing.T) {
	test := tests[1]
	m := createMessage(test.Kind, test.Content, test.Dt)
	expected := formatExpectedMessage(test.Kind, test.Content, test.Dt)
	if message.Encode(m) != expected {
		t.Error("mensagem formatada incorretamente", message.Encode(m))
	}
}

func TestEncodeMultiLineContent(t *testing.T) {
	expected :=""
	encoded  :=""
	for _,msg := range tests {
		m := createMessage(msg.Kind, msg.Content, msg.Dt)
		expected = expected+ formatExpectedMessage(msg.Kind, msg.Content, msg.Dt)
		encoded = encoded+ message.Encode(m)
	}
	if encoded != expected {
		t.Error("mensagem formatada incorretamente\n", expected,"\n---\n", encoded, "\n---")
	}
}

func TestEncodeMultiLineContentFirstLineEmpty(t *testing.T) {
	test := tests[0]
	m := createMessage(test.Kind, test.Content, test.Dt)
	expected := formatExpectedMessage(test.Kind, test.Content, test.Dt)
	if message.Encode(m) != expected {
		t.Error("mensagem formatada incorretamente\n", expected,"\n---\n", message.Encode(m), "\n---")
	}
}



func TestDecode(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("ainda nao sei se eh por aqui que esta a resposta")
	messageString := message.Encode(m)
	err, decodedMessage := message.Decode(messageString)
	if err != nil || m.Kind != decodedMessage.Kind || m.Content !=  decodedMessage.Content {
		t.Error("decode fail", decodedMessage.Kind,"\n", decodedMessage.Content, err)
	}
}

func TestDecodePostMessage(t *testing.T) {
	msgs := make([]*message.Message,0)
	encoded  :=""
	for _,msg := range tests {
		m := createMessage(msg.Kind, msg.Content, msg.Dt)
		err, mPost := message.DecodePost(message.Encode(m))
		t.Log(mPost[0].DateTime)
		t.Log(m.DateTime)
		if err != nil {
			t.Error(err)
		}
		if !mPost[0].Equals(m) { 
			t.Error("mensagem formatada incorretamente\n"/*, m , mPost[0]*/)
		}
		encoded = encoded+ message.Encode(m)
		msgs = append(msgs,m)
	}

	err, messages:= message.DecodePost(encoded)
	if err != nil {
		t.Error(err)
	}

	for i,onemessage:= range messages {
		if !onemessage.Equals(msgs[i]) {
				t.Error("mensagem multipla formatada incorretamente\n", onemessage , msgs[i])
		}
	}

}
