package transport

import (
	"io"
)

type TTransport interface {
	io.Reader
	io.Writer
	Open()
	Close()
	IsOpen() bool
	Flush()
}

type TServerTransport interface {
	TTransport
	Listen()
	Accept() TTransport
}
