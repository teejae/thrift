package thrift

import (
	"io"
	"os"
)

type TTransport interface {
	io.Reader
	io.Writer
	Open() os.Error
	Close()
	IsOpen() bool
	Flush()
}

type TServerTransport interface {
	TTransport
	Listen()
	Accept() TTransport
}
