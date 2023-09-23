package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/shadiestgoat/log"
)

func urlBase(protocol string) string {
	return protocol + SECURE_S + "://" + BASE_LOCATION
}

type Message struct {
	Author string `json:"author"`
	Content string `json:"content"`
	Prev *Message `json:"-"`
}

func (m Message) Color() string {
	return m.Author[len(m.Author) - 6:]
}

func SendMessage(content string) {
	m, _ := json.Marshal(Message{
		Author:  AUTHOR_SEND_INFO,
		Content: content,
	})

	http.Post(urlBase("http") + "/messages", "application/json", bytes.NewBuffer(m))
}

func GetMessages() []*Message {
	resp, err := http.Get(urlBase("http") + "/messages")
	log.FatalIfErr(err, "fetching msgs resp")

	messages := []*Message{}
	err = json.NewDecoder(resp.Body).Decode(&messages)
	log.FatalIfErr(err, "Decoding messages")
	return messages
}