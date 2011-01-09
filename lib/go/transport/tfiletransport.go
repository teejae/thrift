package transport

import (
	"os"
	"thrift"
)

type TFileTransport struct {
	file *os.File
}


func NewTFileTransport(file *os.File) thrift.TTransport {
	return &TFileTransport{file: file}
}

func (t *TFileTransport) Open() os.Error {
	return nil
}

func (t *TFileTransport) Close() {
	t.file.Close()
}

func (t *TFileTransport) IsOpen() bool {
	return (t.file != nil)
}

func (t *TFileTransport) Flush() {}

func (t *TFileTransport) Read(b []byte) (n int, err os.Error) {
	return t.file.Read(b)
}

func (t *TFileTransport) Write(b []byte) (n int, err os.Error) {
	return t.file.Write(b)
}
