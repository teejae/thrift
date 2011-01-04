package transport

import (
	"os"
)

type TFileTransport struct {
	name string
	file *os.File
}


func NewTFileTransport(name string) *TFileTransport {
	return &TFileTransport{name: name}
}

func (t *TFileTransport) Open() {
	t.file, _ = os.Open(t.name, os.O_RDWR, 0)
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
