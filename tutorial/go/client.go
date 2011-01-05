package main

import (
	"log"
	"os"
	"shared"
	"thrift/protocol"
	"thrift/transport"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()

	f, _ := os.Open("shared_struct.tob", os.O_RDONLY, 0)

	t := transport.NewTFileTransport(f)
	t.Open()
	log.Println("file: ", t)
	p := protocol.NewTBinaryProtocol(t, true, true)

	s := shared.NewSharedStruct()
	s.Read(p)

	t.Close()
	log.Println("struct: ", s)

	log.Println("writing")
	f, _ = os.Open("shared_struct_out.tob", os.O_CREATE|os.O_WRONLY, 0777)
	t = transport.NewTFileTransport(f)
	log.Println("file: ", t)
	t.Open()
	p = protocol.NewTBinaryProtocol(t, true, true)
	s.Write(p)

	t.Close()

	tsock := transport.NewTSocket("127.0.0.1", 9090)
	tsock.Open()
	log.Println("tsock ", tsock)
	p = protocol.NewTBinaryProtocol(tsock, true, true)
	c := shared.NewSharedServiceClient(p, p)
	s = c.GetStruct(1234)
	tsock.Close()
	log.Println("getStruct -> ", s)
}
