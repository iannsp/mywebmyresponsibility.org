package message

import (
	"github.com/mywebmyresponsibility.org/core/message"
	"testing"
	"time"
)

func getMessage(kind string, date time.Time, content string) string {
	return "--NEWMESSAGE--\n" +
		"DATE " + date.Format("20060102150405") + "\n" +
		"KIND nota\n" +
		"CONTENT " + content + "\n"

}
func TestEncode(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("ainda nao sei se eh por aqui que esta a resposta")
	if message.Encode(m) != getMessage(m.Kind, m.DateTime, m.Content) {
		t.Error("mensagem formatada incorretamente", message.Encode(m))
	}
}

func TestEncodeMultiLineContent(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("linha 1\nv2 linha2\nv3 linha3")
	expected := getMessage(m.Kind, m.DateTime, m.Content)
	if message.Encode(m) != expected {
		t.Error("mensagem formatada incorretamente\n", expected, "\n---\n", message.Encode(m), "\n---")
	}
}

func TestEncodeMultiLineContentFirstLineEmpty(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("\nlinha 1\nv2 linha2\nv3 linha3")
	expected := getMessage(m.Kind, m.DateTime, m.Content)
	if message.Encode(m) != expected {
		t.Error("mensagem formatada incorretamente\n", expected, "\n---\n", message.Encode(m), "\n---")
	}
}

func TestDecode(t *testing.T) {
	m := message.NewMessage("nota")
	m.SetContent("ainda nao sei se eh por aqui que esta a resposta")
	messageString := message.Encode(m)
	err, decodedMessage := message.Decode(messageString)
	if err != nil || m.Kind != decodedMessage.Kind || m.Content != decodedMessage.Content {
		t.Error("decode fail", decodedMessage)
	}
}

func TestDecodeMultiple(t *testing.T) {
	var messagesString string
	m := message.NewMessage("nota")
	m.SetContent("ainda nao sei se eh por aqui que esta a resposta")
	messagesString = message.Encode(m) + message.Encode(m) + message.Encode(m)
	t.Log(messagesString)
	err, messages := message.DecodePost(messagesString)
	if err != nil {
		t.Error("fail on decode")
	}
	for _, onemessage := range messages {
		if m.Kind != onemessage.Kind || onemessage.Content != m.Content {
			t.Error("decode fail")
		}
	}

}
