package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	serverIp   string
	serverPort int
	name       string
	conn       net.Conn
	flag       int //当前client的模式
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
		flag:       999,
	}

	return client
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名")
	fmt.Scanln(&client.name)

	sendMsg := "rename|" + client.name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

// 处理server回应的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	//一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) PublicChat() {
	var chatMsg string
	for chatMsg != "exit" {
		fmt.Println(">>>>请输入聊天内容，exit退出>>>>")
		fmt.Scanln(&chatMsg)
		if len(chatMsg) == 0 {
			continue
		}
		sendMsg := chatMsg + "\n"
		_, err := client.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("conn.Write err:", err)
			break
		}

	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式")
			break
		case 3:
			//更新用户名
			client.UpdateName()
			break
		}
	}
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(ServerIp, ServerPort)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.serverIp, client.serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
	}
	client.conn = conn

	//单独开启一个goroutine去处理server的回执消息
	go client.DealResponse()

	//启动客户端的业务
	client.Run()
}
