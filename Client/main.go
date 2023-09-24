package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	serverIp   string
	serverPort int
	name       string
	conn       net.Conn
}

var ServerIp string
var ServerPort int

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "设置服务器IP地址(默认127.0.0.1)")
	flag.IntVar(&ServerPort, "port", 8888, "设置服务器端口(默认8888)")
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		serverIp:   serverIp,
		serverPort: serverPort,
	}

	return client
}

func main() {
	flag.Parse()
	client := NewClient(ServerIp, ServerPort)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.serverIp, client.serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
	}
	client.conn = conn

	select {}
}
