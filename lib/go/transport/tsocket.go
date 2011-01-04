package transport

import (
	"net"
	"os"
)

type TSocket struct {
	host string
	port int
	conn *net.TCPConn
}

func NewTSocket(host string, port int) *TSocket {
	return &TSocket{host: host, port: port}
}

func newTSocketFromConnection(conn *net.TCPConn) *TSocket {
	return &TSocket{conn: conn}
}

func (s *TSocket) Open() {
	raddr, _ := net.ResolveTCPAddr(s.host + ":" + string(s.port))
	conn, _ := net.DialTCP(raddr.Network(), nil, raddr)
	s.conn = conn
}

func (s *TSocket) Close() {
	s.conn.Close()
}

func (s *TSocket) IsOpen() bool {
	return (s.conn != nil)
}

func (s *TSocket) Flush() {}

func (s *TSocket) Read(b []byte) (n int, err os.Error) {
	return s.conn.Read(b)
}

func (s *TSocket) Write(b []byte) (n int, err os.Error) {
	return s.conn.Write(b)
}

type TServerSocket struct {
	TSocket
	port     int
	listener *(net.TCPListener)
}

func NewTServerSocket(port int) *TServerSocket {
	return &TServerSocket{port: port}
}

func (s *TServerSocket) Listen() {
	laddr, _ := net.ResolveTCPAddr(":" + string(s.port))
	listener, _ := net.ListenTCP(laddr.Network(), laddr)
	s.listener = listener
}

func (s *TServerSocket) Accept() *TSocket {
	conn, _ := s.listener.AcceptTCP()
	return newTSocketFromConnection(conn)
}
