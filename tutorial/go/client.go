package main

import (
	"log"
	"tutorial"
	"github.com/teejae/go-thrift/thrift/protocol"
	"github.com/teejae/go-thrift/thrift/transport"
)

func main() {
	// Make socket
	tsock := transport.NewTSocket("localhost", 9090)

	// Wrap in protocol
	p := protocol.NewTBinaryProtocol(tsock, true, true)
	
	// Create a client to use the protocol encoder
	client := tutorial.NewCalculatorClient(p, p)
	
	// Connect!
	tsock.Open()
	
	client.Ping()
	log.Println("ping()");
	
	work := tutorial.NewWork()
	
	work.Op = new(tutorial.Operation)
	*work.Op = tutorial.Operation_DIVIDE
	n1 := int32(1)
	work.Num1 = &n1
	n2 := int32(0)
	work.Num2 = &n2
	
	quotient, err := client.Calculate(1, *work)
	if err != nil {
		log.Println(err.String())
	} else {
		log.Println("You know how to divide by 0!")
		log.Println("quotient", *quotient)
	}
	
	op := tutorial.Operation_SUBTRACT
	work.Op = &op
	n1 = int32(15)
	work.Num1 = &n1
	n2 = int32(10)
	work.Num2 = &n2
	
	diff, err := client.Calculate(1, *work)
	if err != nil{
		log.Println("InvalidOperation: ", err)		
	}
	
	log.Println("15-10=", *diff)
	
	s, err := client.GetStruct(1)
	if err != nil {
		// panic(err)
	}
	
	log.Println("Check log: ", *s.Value)
	
	// Close!
	tsock.Close()
}
