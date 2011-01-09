package thrift

type TType uint32

const (
	TTYPE_STOP   TType = 0
	TTYPE_VOID   = 1
	TTYPE_BOOL   = 2
	TTYPE_BYTE   = 3
	TTYPE_DOUBLE = 4
	TTYPE_I16    = 6
	TTYPE_I32    = 8
	TTYPE_I64    = 10
	TTYPE_STRING = 11
	TTYPE_STRUCT = 12
	TTYPE_MAP    = 13
	TTYPE_SET    = 14
	TTYPE_LIST   = 15
)

const (
	TMESSAGETYPE_CALL      = 1
	TMESSAGETYPE_REPLY     = 2
	TMESSAGETYPE_EXCEPTION = 3
	TMESSAGETYPE_ONEWAY    = 4
)


type TProtocol interface {
	WriteMessageBegin(name string, ttype TType, seq int32)
	WriteMessageEnd()
	WriteStructBegin(name string)
	WriteStructEnd()
	WriteFieldBegin(name string, ttype TType, id int16)
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
	WriteString(string)

	ReadMessageBegin() (name string, ttype TType, seq int32)
	ReadMessageEnd()
	ReadStructBegin() (name string)
	ReadStructEnd()
	ReadFieldBegin() (name string, ttype TType, id int16)
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

	// utility
	GetTransport() TTransport
}

func SkipType(p TProtocol, ttype TType) {
	switch ttype {
	case TTYPE_STOP:
		return

	case TTYPE_BOOL:
		p.ReadBool()
	case TTYPE_BYTE:
		p.ReadByte()
	case TTYPE_I16:
		p.ReadI16()
	case TTYPE_I32:
		p.ReadI32()
	case TTYPE_I64:
		p.ReadI64()
	case TTYPE_DOUBLE:
		p.ReadDouble()
	case TTYPE_STRING:
		p.ReadString()
	case TTYPE_STRUCT:
		p.ReadStructBegin()
		for {
			_, ttype, _ = p.ReadFieldBegin()
			if ttype == TTYPE_STOP {
				break
			}
			SkipType(p, ttype)
			p.ReadFieldEnd()
		}
		p.ReadStructEnd()
	case TTYPE_MAP:
		ktype, vtype, size := p.ReadMapBegin()
		for i := int32(0); i < size; i++ {
			SkipType(p, ktype)
			SkipType(p, vtype)
		}
		p.ReadMapEnd()
	case TTYPE_SET:
		etype, size := p.ReadSetBegin()
		for i := int32(0); i < size; i++ {
			SkipType(p, etype)
		}
		p.ReadSetEnd()
	case TTYPE_LIST:
		etype, size := p.ReadListBegin()
		for i := int32(0); i < size; i++ {
			SkipType(p, etype)
		}
		p.ReadListEnd()
	}
}
