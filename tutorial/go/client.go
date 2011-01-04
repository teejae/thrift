package main

import (
	"log"
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

	t := transport.NewTFileTransport("shared_struct.tob")
	t.Open()
	log.Println("file: ", t)
	p := protocol.NewTBinaryProtocol(t, true, true)

	s := shared.NewSharedStruct()
	s.Read(p)

	t.Close()

	log.Println("struct: ", s)
}
