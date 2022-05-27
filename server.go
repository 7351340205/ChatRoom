package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户 map+读写锁
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// NewServer 创建Server
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
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// Handler 处理消息
func (this *Server) Handler(conn net.Conn) {
	//当前连接的业务
	fmt.Println("连接建立成功,尝试接收数据,IP:", conn.RemoteAddr())

	user := NewUser(conn, this)

	user.Online()

	//接收客户端信息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read err:", err)
				return
			}

			msg := string(buf[:n-1])

			user.DoMessage(msg)

		}
	}()

	//阻塞
	select {}

}

// Start 启动服务器
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

	//启动消息监听goroutine
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
