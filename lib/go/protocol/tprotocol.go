package protocol

type TProtocol interface {
	WriteMessageBegin(name string, ttype TType, seq int32)
	WriteMessageEnd()
	WriteStructBegin(name string)
	WriteStructEnd()
	WriteFieldBegin(name string, ttype TType, id int32)
	WriteFieldEnd()
	WriteFieldStop()
	WriteMapBegin(ktype TType, vtype TType, size int32)
	WriteMapEnd()
	WriteListBegin(etype TType, size int32)
	WriteListEnd()
	WriteSetBegin(etype TType, size int32)
	WriteSetEnd()
	WriteBool(bool)
	WriteByte(byte)
	WriteI16(int16)
	WriteI32(int32)
	WriteI64(int64)
	WriteDouble(float64)

	ReadMessageBegin() (name string, ttype TType, seq int32)
	ReadMessageEnd()
	ReadStructBegin() (name string)
	ReadStructEnd()
	ReadFieldBegin() (name string, ttype TType, id int32)
	ReadFieldEnd()
	ReadMapBegin() (ktype TType, vtype TType, size int32)
	ReadMapEnd()
	ReadListBegin() (etype TType, size int32)
	ReadListEnd()
	ReadSetBegin() (etype TType, size int32)
	ReadSetEnd()
	ReadBool() bool
	ReadByte() byte
	ReadI16() int16
	ReadI32() int32
	ReadI64() int64
	ReadDouble() float64
	ReadString() string
}
