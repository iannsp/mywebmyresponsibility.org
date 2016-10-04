package message

import (
	"github.com/mywebmyresponsibility.org/message"
	"testing"
)

func TestEncode(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("ainda nao sei se eh por aqui que esta a resposta")
	expected := "--NEWMESSAGE--\n" +
		"KIND nota\n" +
		"CONTENT ainda nao sei se eh por aqui que esta a resposta\n"
	if message.Encode(m) != expected {
		t.Error("mensagem formatada incorretamente", message.Encode(m))
	}
}

func TestEncodeMultiLineContent(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("linha 1\nv2 linha2\nv3 linha3")
	expected := "--NEWMESSAGE--\n" +
		"KIND nota\n" +
		"CONTENT linha 1\nv2 linha2\nv3 linha3\n"
	if message.Encode(m) != expected {
		t.Error("mensagem formatada incorretamente\n", expected,"\n---\n", message.Encode(m), "\n---")
	}
}

func TestEncodeMultiLineContentFirstLineEmpty(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("\nlinha 1\nv2 linha2\nv3 linha3")
	expected := "--NEWMESSAGE--\n" +
		"KIND nota\n" +
		"CONTENT \nlinha 1\nv2 linha2\nv3 linha3\n"
	if message.Encode(m) != expected {
		t.Error("mensagem formatada incorretamente\n", expected,"\n---\n", message.Encode(m), "\n---")
	}
}



func TestDecode(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("ainda nao sei se eh por aqui que esta a resposta")
	messageString := message.Encode(m)
	decodedMessage := message.Decode(messageString)
	if m.Kind != decodedMessage.Kind || m.Content != decodedMessage.Content {
		t.Error("decode fail", decodedMessage)
	}
}

func TestDecodeMultiple(t *testing.T) {
	var messagesString string
	m := message.NewMessage("nota")
	m.SetContent("ainda nao sei se eh por aqui que esta a resposta")
	messagesString = message.Encode(m)+message.Encode(m)+ message.Encode(m)
	t.Log(messagesString)
	messages:= message.DecodePost(messagesString)
	for _,onemessage:= range messages {
		if m.Kind != onemessage.Kind || onemessage.Content != m.Content {
			t.Error("decode fail")
		}
	}

}
