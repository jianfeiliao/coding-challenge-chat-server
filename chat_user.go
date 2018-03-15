package main

import (
	"net"
)

type User struct {
	NickName string
	Conn     net.Conn
}

func NewUser(name string, conn net.Conn) *User {
	return &User{
		NickName: name,
		Conn:     conn,
	}
}
