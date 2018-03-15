package main

import (
	"fmt"
	"net"
	"time"
)

type ChatServer struct {
	Users        map[string]*User
	ChatMessages []*ChatMessage
}

type User struct {
	NickName string
	Conn     net.Conn
}

type ChatMessage struct {
	Timestamp time.Time
	FromUser  string
	Message   string
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		Users: make(map[string]*User, 0),
	}
}

func (cs *ChatServer) StartListening() {
	fmt.Println("starting chat server")

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("error on listen: %s\n", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error on accept: %s\n", err)
		}
		go cs.handleConnection(conn)
	}
}

func (cs *ChatServer) checkUniqueUserName(name string) bool {
	_, found := cs.Users[name]
	if found {
		return false
	}

	return true
}

func (cs *ChatServer) addUser(user *User) {
	// TODO check input
	cs.Users[user.NickName] = user
}

func (cs *ChatServer) handleConnection(conn net.Conn) {
	conn.Write([]byte("> Welcome to my chat server, what's your name?\n"))

	cs.readUserName(conn)
	//cs.HandleChatMessages()
}

func (cs *ChatServer) readUserName(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("error on read: %s\n", err)
		}

		name := string(buffer[:n-1])
		if cs.checkUniqueUserName(name) {
			conn.Write([]byte("> username is unique, moving on...\n"))
			user := NewUser(name, conn)
			cs.addUser(user)

			// we have a valid user name
			// TODO send the last 10 lines of chat in the chat-room,
			cs.SendChatHistory(user)
			// TODO broadcast that the user has connected to the other clients
			cs.BoardcastToAll(user)
			// TODO send a list of users that are currently connected to the just-connected client.
			cs.SendUsersList(user)
			break
		} else {
			conn.Write([]byte("> username already taken, please pick a different one\n"))
		}
	}

}

func (cs *ChatServer) SendMessage(name string, conn net.Conn, msg string) {
	chatMsg := &ChatMessage{
		//Timestamp: time.Now(),
		FromUser: name,
		Message:  msg,
	}
	cs.ChatMessages = append(cs.ChatMessages, chatMsg)
	conn.Write([]byte(msg))
}

func (cs *ChatServer) SendChatHistory(user *User) {
	fmt.Println("Sending chat history")
	for _, msg := range cs.ChatMessages {
		user.Conn.Write([]byte(fmt.Sprintf(">> %s: %s\n", msg.FromUser, msg.Message)))
	}
}

func (cs *ChatServer) BoardcastToAll(user *User) {
	allUsers := cs.Users
	for _, u := range allUsers {
		if u.NickName != user.NickName {
			msg := fmt.Sprintf("%s has joined the chat*\n", user.NickName)
			cs.SendMessage("server", u.Conn, fmt.Sprintf("> [%s] *%s", time.Now(), msg))
		}
	}

}

func (cs *ChatServer) SendUsersList(user *User) {
	allUsers := make([]string, 0)
	for name, _ := range cs.Users {
		allUsers = append(allUsers, name)
	}

	user.Conn.Write([]byte(fmt.Sprintf("> You are connected with users: %s\n", allUsers)))

}

func (cs *ChatServer) HandleChatmessages(user *User) {
}

func NewUser(name string, conn net.Conn) *User {
	return &User{
		NickName: name,
		Conn:     conn,
	}
}
