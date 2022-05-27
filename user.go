package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// NewUser 创建新用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	//启动监听
	go user.ListenMessage()

	return user
}

// Online 用户上线
func (this *User) Online() {

	//用户上线，将用户加入到map中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "用户已上线")
}

// Offline 用户下线
func (this *User) Offline() {

	//用户下线，将用户从map中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播当前用户下线消息
	this.server.BroadCast(this, "用户已下线")
}

// DoMessage 广播当前消息
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// ListenMessage 监听用户发送的消息
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		//写入消息，处理err
		_, err := this.conn.Write([]byte(msg + "\n"))
		if err != nil {
			return
		}
	}
}
