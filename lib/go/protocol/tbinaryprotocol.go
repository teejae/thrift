package protocol

import (
	"encoding/binary"
	"io"
	"thrift/transport"
)

const (
	VERSION_MASK = 0xffff0000
	VERSION_1    = 0x80010000
	TYPE_MASK    = 0x000000ff
)

type TBinaryProtocol struct {
	transport   transport.TTransport
	strictRead  bool
	strictWrite bool
}

func NewTBinaryProtocol(transport transport.TTransport, strictRead, strictWrite bool) *TBinaryProtocol {
	return &TBinaryProtocol{transport: transport, strictRead: strictRead, strictWrite: strictWrite}
}

func (p *TBinaryProtocol) WriteMessageBegin(name string, ttype TType, seq int32) {
	if p.strictWrite {
		p.WriteI32(int32(VERSION_1 | ttype))
		p.WriteString(name)
		p.WriteI32(seq)
	} else {
		p.WriteString(name)
		p.WriteByte(byte(ttype))
		p.WriteI32(seq)
	}
}

func (p *TBinaryProtocol) WriteMessageEnd()             {}
func (p *TBinaryProtocol) WriteStructBegin(name string) {}
func (p *TBinaryProtocol) WriteStructEnd()              {}
func (p *TBinaryProtocol) WriteFieldBegin(name string, ttype TType, id int16) {
	p.WriteByte(byte(ttype))
	p.WriteI16(id)
}
func (p *TBinaryProtocol) WriteFieldEnd() {}
func (p *TBinaryProtocol) WriteFieldStop() {
	p.WriteByte(byte(TTYPE_STOP))
}
func (p *TBinaryProtocol) WriteMapBegin(ktype TType, vtype TType, size int32) {
	p.WriteByte(byte(ktype))
	p.WriteByte(byte(vtype))
	p.WriteI32(size)
}
func (p *TBinaryProtocol) WriteMapEnd() {}

func (p *TBinaryProtocol) WriteListBegin(etype TType, size int32) {
	p.WriteByte(byte(etype))
	p.WriteI32(size)
}

func (p *TBinaryProtocol) WriteListEnd() {}

func (p *TBinaryProtocol) WriteSetBegin(etype TType, size int32) {
	p.WriteByte(byte(etype))
	p.WriteI32(size)
}
func (p *TBinaryProtocol) WriteSetEnd() {}
func (p *TBinaryProtocol) WriteBool(b bool) {
	if b {
		p.WriteByte(1)
	} else {
		p.WriteByte(0)
	}
}

func (p *TBinaryProtocol) WriteByte(b byte) {
	binary.Write(p.transport, binary.BigEndian, b)
}

func (p *TBinaryProtocol) WriteI16(i int16) {
	binary.Write(p.transport, binary.BigEndian, i)
}

func (p *TBinaryProtocol) WriteI32(i int32) {
	binary.Write(p.transport, binary.BigEndian, i)
}

func (p *TBinaryProtocol) WriteI64(i int64) {
	binary.Write(p.transport, binary.BigEndian, i)
}

func (p *TBinaryProtocol) WriteDouble(d float64) {
	binary.Write(p.transport, binary.BigEndian, d)
}

func (p *TBinaryProtocol) WriteString(s string) {
	p.WriteI32(int32(len(s)))
	p.transport.Write([]byte(s))
}

func (p *TBinaryProtocol) ReadMessageBegin() (name string, ttype TType, seq int32) {
	size := p.ReadI32()
	if size < 0 {
		version := uint32(size) & VERSION_MASK
		if version != VERSION_1 {
			// make an error
			// throw new TProtocolError(TProtocolError.BAD_VERSION, "Bad version in readMessageBegin");
		}
		name = p.ReadString()
		ttype = TType(uint32(size) & TYPE_MASK)
		seq = p.ReadI32()
	} else {
		if p.strictRead {
			// throw new TProtocolError(TProtocolError.BAD_VERSION, "Missing version in readMessageBegin, old client?");
		}
		name_buf := make([]byte, size)
		io.ReadFull(p.transport, name_buf)
		name = string(name_buf)
		ttype = TType(p.ReadByte())
		seq = p.ReadI32()
	}
	return
}

func (p *TBinaryProtocol) ReadMessageEnd() {}
func (p *TBinaryProtocol) ReadStructBegin() (name string) {
	name = ""
	return
}

func (p *TBinaryProtocol) ReadStructEnd() {}

func (p *TBinaryProtocol) ReadFieldBegin() (name string, ttype TType, id int16) {
	ttype = TType(p.ReadByte())
	if ttype == TTYPE_STOP {
		id = 0
	} else {
		id = p.ReadI16()
	}
	return
}
func (p *TBinaryProtocol) ReadFieldEnd() {}

func (p *TBinaryProtocol) ReadMapBegin() (ktype TType, vtype TType, size int32) {
	ktype = TType(p.ReadByte())
	vtype = TType(p.ReadByte())
	size = p.ReadI32()
	return
}

func (p *TBinaryProtocol) ReadMapEnd() {}

func (p *TBinaryProtocol) ReadListBegin() (etype TType, size int32) {
	etype = TType(p.ReadByte())
	size = p.ReadI32()
	return
}

func (p *TBinaryProtocol) ReadListEnd() {}

func (p *TBinaryProtocol) ReadSetBegin() (etype TType, size int32) {
	etype = TType(p.ReadByte())
	size = p.ReadI32()
	return
}

func (p *TBinaryProtocol) ReadSetEnd() {}

func (p *TBinaryProtocol) ReadBool() bool {
	b := p.ReadByte()
	if b == 0 {
		return false
	}
	return true
}

func (p *TBinaryProtocol) ReadByte() (b byte) {
	binary.Read(p.transport, binary.BigEndian, &b)
	return
}

func (p *TBinaryProtocol) ReadI16() (i int16) {
	binary.Read(p.transport, binary.BigEndian, &i)
	return
}

func (p *TBinaryProtocol) ReadI32() (i int32) {
	binary.Read(p.transport, binary.BigEndian, &i)
	return
}

func (p *TBinaryProtocol) ReadI64() (i int64) {
	binary.Read(p.transport, binary.BigEndian, &i)
	return
}

func (p *TBinaryProtocol) ReadDouble() (d float64) {
	binary.Read(p.transport, binary.BigEndian, &d)
	return
}

func (p *TBinaryProtocol) ReadString() (s string) {
	size := p.ReadI32()
	s_buf := make([]byte, size)
	io.ReadFull(p.transport, s_buf)
	s = string(s_buf)
	return
}
