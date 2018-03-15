package main

import (
	"fmt"
	"time"
)

type ChatMessage struct {
	Timestamp time.Time
	FromUser  string
	Text      string
}

func NewChatMessage(user, text string) *ChatMessage {
	return &ChatMessage{
		Timestamp: time.Now(),
		FromUser:  user,
		Text:      text,
	}
}

func (cm *ChatMessage) ToString() string {
	// format the timestamp to HH:mm:ss
	ts := cm.Timestamp.Format("15:04:05")

	// it's a message from a connected user
	if cm.FromUser != "" {
		return fmt.Sprintf("< [%s] <%s> %s\n", ts, cm.FromUser, cm.Text)
	}

	// it's a boardcast from the server
	return fmt.Sprintf("< [%s] *%s*\n", ts, cm.Text)
}
