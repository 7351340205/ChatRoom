package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户 map+读写锁
	OnlineMap map[string]*User
	maplock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//发送消息至全体在线user
		this.maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.maplock.Unlock()
	}
}

//处理消息
func (this *Server) Handler(conn net.Conn) {
	//当前连接的业务
	fmt.Println("连接建立成功,尝试接收数据,IP:", conn.RemoteAddr())

	user := NewUser(conn)

	//用户上线，将用户加入到map中
	this.maplock.Lock()
	this.OnlineMap[user.Name] = user
	this.maplock.Unlock()

	//广播当前用户上线消息
	this.BroadCast(user, "用户已上线")

	//阻塞
	select {}

}

//启动服务器
func (this *Server) Start() {
	//socket Listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close Listen socket
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	go this.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		go this.Handler(conn)

	}

}
