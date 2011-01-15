package thrift

import (
	"os"
)

type TProcessor interface {
	Process(iprot, oprot TProtocol) (bool, TException)
}

type TExceptionType int

type TException interface {
	os.Error
	Type() *TExceptionType
	Write(oprot TProtocol)
	Read(iprot TProtocol)
}

func newTException(exceptionType TExceptionType, message string, structName string) TException {
	return &TExceptionImpl{message: &message, eType: &exceptionType, structName: structName}
}

type TExceptionImpl struct {
	message    *string
	eType      *TExceptionType
	structName string
}

func (e *TExceptionImpl) String() string {
	return *e.message
}

func (e *TExceptionImpl) Type() *TExceptionType {
	return e.eType
}

func (e *TExceptionImpl) Write(oprot TProtocol) {
	oprot.WriteStructBegin(e.structName)
	if e.message != nil {
		oprot.WriteFieldBegin("message", TTYPE_STRING, 1)
		oprot.WriteString(*e.message)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldBegin("type", TTYPE_I32, 2)
	oprot.WriteI32(int32(*e.eType))
	oprot.WriteFieldEnd()

	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

func (e *TExceptionImpl) Read(iprot TProtocol) {
	iprot.ReadStructBegin()
	defer iprot.ReadStructEnd()

	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == TTYPE_STOP {
			break
		}

		switch fid {
		case 1: // message
			if ftype == TTYPE_STRING {
				v := iprot.ReadString()
				e.message = &v
			} else {
				SkipType(iprot, ftype)
			}
		case 2: // type
			if ftype == TTYPE_I32 {
				v := iprot.ReadI32()
				t := TExceptionType(v)
				e.eType = &t
			} else {
				SkipType(iprot, ftype)
			}
		default: // unknown
			SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
}

type TApplicationExceptionType TExceptionType

const (
	TAPPLICATION_EXCEPTION_UNKNOWN              TApplicationExceptionType = 0
	TAPPLICATION_EXCEPTION_UNKNOWN_METHOD       = 1
	TAPPLICATION_EXCEPTION_INVALID_MESSAGE_TYPE = 2
	TAPPLICATION_EXCEPTION_WRONG_METHOD_NAME    = 3
	TAPPLICATION_EXCEPTION_BAD_SEQUENCE_ID      = 4
	TAPPLICATION_EXCEPTION_MISSING_RESULT       = 5
	TAPPLICATION_EXCEPTION_INTERNAL_ERROR       = 6
	TAPPLICATION_EXCEPTION_PROTOCOL_ERROR       = 7
)

func NewTApplicationException(exceptionType TApplicationExceptionType, message string) TException {
	t := TExceptionType(exceptionType)
	return newTException(t, message, "TApplicationException")
}

type TTransportExceptionType TExceptionType

const (
	TTRANSPORT_EXCEPTION_TYPE_UNKNOWN      = 0
	TTRANSPORT_EXCEPTION_TYPE_NOT_OPEN     = 1
	TTRANSPORT_EXCEPTION_TYPE_ALREADY_OPEN = 2
	TTRANSPORT_EXCEPTION_TYPE_TIMED_OUT    = 3
	TTRANSPORT_EXCEPTION_TYPE_END_OF_FILE  = 4
)

func NewTTransportException(exceptionType TTransportExceptionType, message string) TException {
	t := TExceptionType(exceptionType)
	return newTException(t, message, "TTransportException")
}
