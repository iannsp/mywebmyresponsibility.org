package message

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	// kind define the route of the message into the system
	Kind     string
	Content  string
	DateTime time.Time
}

// NewMessage create an instance, return a ptr to it
func NewMessage(kind string) *Message {
	return &Message{kind, "", time.Now()}
}

// SetContent add content to the message.
func (m *Message) SetContent(content string) {
	m.Content = strings.TrimRight(content, "\n")
}

// Equals compare it instance to another
func (m *Message) Equals(msgToCompare *Message) bool {
	return msgToCompare.Kind == m.Kind &&
		msgToCompare.Content == m.Content /* ||
		mPost[0].DateTime != m.DateTime*/

}

// Encode return the message string and need be prepared by the data service.
func Encode(m *Message) string {
	return fmt.Sprintf("--NEWMESSAGE--\nDATE %s\nKIND %s\nCONTENT %s\n", m.DateTime.Format("20060102150405"), m.Kind, m.Content)

}

//Decode get a messagem string of a single *Message and create it.
// string message format:
// has 4 or more lines
//--NEWMESSAGE--
//DATE unixtimestamp
//KIND string
//CONTENT info
//[OPTONAL] info...
// BLANK LINE
func Decode(msg string) (error, *Message) {
	parts := strings.Split(msg, "\n")
	if strings.Split(parts[1], " ")[0] != "DATE" ||
		strings.Split(parts[2], " ")[0] != "KIND" ||
		strings.Split(parts[3], " ")[0] != "CONTENT" {
		log.Fatal("vix")
	}
	intdate, err := strconv.ParseInt(strings.Split(parts[1], " ")[1], 10, 64)
	if err != nil {
		return err, &Message{}
		log.Println("Fail parse Date:" + strings.Split(parts[1], " ")[1])
	}
	unixDate := time.Unix(intdate, 0)
	kind := strings.Split(parts[2], " ")[1]
	content := strings.Replace(msg, "--NEWMESSAGE--\n"+parts[1]+"\n"+parts[2]+"\nCONTENT ", "", 1)

	m := NewMessage(kind)
	m.DateTime = unixDate
	m.SetContent(content)
	return nil, m
}
