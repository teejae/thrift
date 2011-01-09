package server

import (
	"thrift"
	"thrift/protocol"
)

type TSimpleServer struct {
	processor thrift.TProcessor
	transport thrift.TServerTransport
}

func NewTSimpleServer(processor thrift.TProcessor, transport thrift.TServerTransport) thrift.TServer {
	return thrift.TServer(&TSimpleServer{processor: processor, transport: transport})
}

func (s *TSimpleServer) Serve() {
	s.transport.Listen()
	for {
		clientTransport := s.transport.Accept()
		// potentially wrap factories here
		itrans := clientTransport
		otrans := clientTransport
		iprot := protocol.NewTBinaryProtocol(itrans, true, true)
		oprot := protocol.NewTBinaryProtocol(otrans, true, true)

		for {
			_, err := s.processor.Process(iprot, oprot)
			if err != nil {
				// might just be that we're done!
				break
			}
		}
		itrans.Close()
		otrans.Close()
	}
}
