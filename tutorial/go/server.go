package main

import (
	"log"
	"shared"
	"tutorial"
	"github.com/teejae/go-thrift/thrift"
	"github.com/teejae/go-thrift/thrift/server"
	"github.com/teejae/go-thrift/thrift/transport"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()

	runServer()
}

func runServer() {
	processor := tutorial.NewCalculatorProcessor(&Server{log:make(map[int32]*shared.SharedStruct)})
	trans := transport.NewTServerSocket(9090)

	s := server.NewTSimpleServer(processor, trans)

	log.Println("Starting server...")
	s.Serve()
	log.Println("Done.")
}

type Server struct{
	log map[int32]*shared.SharedStruct
}

func (s *Server) GetStruct(key int32) (*shared.SharedStruct, thrift.TException) {
	return s.log[key], nil
}

func (s *Server) Ping() thrift.TException {
	log.Println("Ping!")
	return nil
}

func (s *Server) Add(num1, num2 int32) (*int32, thrift.TException) {
	v := num1 + num2
	return &v, nil
}

func (s *Server) Calculate(logid int32, work tutorial.Work) (*int32, thrift.TException) {
	log.Println("Calculate", logid, work)

	var x *tutorial.InvalidOperation
	var val int32
	switch *work.Op {
	case tutorial.Operation_ADD:
		val = *work.Num1 + *work.Num2
	case tutorial.Operation_SUBTRACT:
		val = *work.Num1 - *work.Num2
	case tutorial.Operation_MULTIPLY:
		val = *work.Num1 * *work.Num2
	case tutorial.Operation_DIVIDE:
		if *work.Num2 == 0 {
			val = 0
			x = tutorial.NewInvalidOperation()
			o := int32(*work.Op)
			x.What = &o
			m := "Cannot divide by 0"
			x.Why = &m
		} else {
			val = *work.Num1 / *work.Num2			
		}
	default:
		val = 0
		x = tutorial.NewInvalidOperation()
		o := int32(*work.Op)
		x.What = &o
		m := "Can't find the operation"
		x.Why = &m
	}
	

    log := shared.NewSharedStruct()
    log.Key = &logid
	v := string(val)
    log.Value = &v
    s.log[logid] = log

	return &val, x	
}

func (s *Server) Zip() thrift.TException {
	log.Println("Zip!")
	return nil
}
