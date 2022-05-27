package main

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {

	fmt.Println("连接建立成功,尝试接收数据")
	fmt.Println("IP:", conn.RemoteAddr())

	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		//fmt.Println("got", n, "bytes.")
		buf = append(buf, tmp[:n]...)

		for _, val := range buf {
			if val == 10 {
				//fmt.Println(buf[index])
				//buf[index] = nil
				//fmt.Println(buf[index])

				var strr string = string(buf[:])
				fmt.Println("message:", strr)
				buf = nil
			}
		}

		//var strr string = string(buf[:])
		//fmt.Println("message:", strr)

	}

	fmt.Println("连接断开,数据接收中止")
	//var str string = string(buf[:])
	//
	//fmt.Println("total size:", len(buf))
	//fmt.Println("message:", str)

}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		go this.Handler(conn)

	}

}
