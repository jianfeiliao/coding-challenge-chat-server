package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type ChatServer struct {
	Users       map[string]*User
	ChatHistory []*ChatMessage
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		Users: make(map[string]*User, 0),
	}
}

func (cs *ChatServer) Start() {
	log.Println("Starting chat server")

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error on listen: %s\n", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Error on accept: %s\n", err)
		}
		go cs.handleConnection(conn)
	}
}

func (cs *ChatServer) handleConnection(conn net.Conn) {
	// the currUser is only set after the user is connected with a unique nickname
	var currUser *User
	conn.Write([]byte("< Welcome to my chat server! What's your name?\n"))

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if _, ok := err.(*net.OpError); ok {
			cs.handleUserDisconnect(currUser)
			break
		}
		if err != nil {
			log.Printf("Error on read: %s\n", err)
		}

		// strip out the line ending at the end
		input := string(buffer[:n-1])

		if currUser != nil {
			cs.handleChatMessages(currUser, input)
		} else {
			currUser = cs.handleUserConnect(conn, input)
		}
	}
}

func (cs *ChatServer) handleUserConnect(conn net.Conn, name string) *User {
	_, found := cs.Users[name]

	// user not found, so the nickname is unique
	if !found {
		user := NewUser(name, conn)
		cs.Users[user.NickName] = user

		cs.listOnlineUsers(user)
		cs.sendChatHistory(user)
		cs.broadcastUserConnect(user)

		log.Printf("User %s has connected\n", name)
		return user
	}

	conn.Write([]byte("< That nickname is already taken, please pick a different one:\n"))
	return nil
}

func (cs *ChatServer) handleUserDisconnect(user *User) {
	if user == nil {
		log.Println("User disconnected before establishing a valid connection")
		return
	}

	cs.broadcastUserDisconnect(user)
	log.Printf("User %s has disconnected\n", user.NickName)
}

func (cs *ChatServer) listOnlineUsers(user *User) {
	allNames := make([]string, 0)
	for name := range cs.Users {
		if name != user.NickName {
			allNames = append(allNames, name)
		}
	}

	if len(allNames) > 0 {
		usersList := fmt.Sprintf("< You are connected with users: [%s]\n", strings.Join(allNames, ", "))
		user.Conn.Write([]byte(usersList))
	} else {
		user.Conn.Write([]byte("< You are the only one online right now\n"))
	}
}

func (cs *ChatServer) sendChatHistory(user *User) {
	history := cs.ChatHistory
	if len(cs.ChatHistory) > 10 {
		history = cs.ChatHistory[len(cs.ChatHistory)-10:]
	}

	for _, chatMsg := range history {
		user.Conn.Write([]byte(chatMsg.ToString()))
	}
}

func (cs *ChatServer) broadcastUserConnect(user *User) {
	msg := fmt.Sprintf("%s has joined the chat", user.NickName)
	chatMsg := NewChatMessage("", msg)
	cs.sendBroadcast(user, chatMsg)
}

func (cs *ChatServer) broadcastUserDisconnect(user *User) {
	msg := fmt.Sprintf("%s has left the chat", user.NickName)
	chatMsg := NewChatMessage("", msg)
	cs.sendBroadcast(user, chatMsg)
}

func (cs *ChatServer) sendBroadcast(user *User, msg *ChatMessage) {
	cs.ChatHistory = append(cs.ChatHistory, msg)
	for _, u := range cs.Users {
		if u.NickName != user.NickName {
			// also send a bell code if user is being "@mentioned"
			mention := fmt.Sprintf("@%s", u.NickName)
			if strings.Contains(msg.Text, mention) {
				u.Conn.Write([]byte("\a"))
			}

			u.Conn.Write([]byte(msg.ToString()))
		}
	}
}

func (cs *ChatServer) handleChatMessages(user *User, msgText string) {
	chatMsg := NewChatMessage(user.NickName, msgText)

	// replay the message to other users, basically as a broadcast
	cs.sendBroadcast(user, chatMsg)
}
