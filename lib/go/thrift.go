package thrift

import (
	"os"
)

type TProcessor interface {
	Process(iprot, oprot TProtocol) (bool, *TException)
}

type TExceptionType int

type TException struct {
	os.Error
	Message    *string
	Type       *TExceptionType
	structName string
}

func newTException(exceptionType TExceptionType, message string, structName string) *TException {
	return &TException{Message: &message, Type: &exceptionType, structName: structName}
}

func (e *TException) Write(oprot TProtocol) {
	oprot.WriteStructBegin(e.structName)
	if e.Message != nil {
		oprot.WriteFieldBegin("message", TTYPE_STRING, 1)
		oprot.WriteString(*e.Message)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldBegin("type", TTYPE_I32, 2)
	oprot.WriteI32(int32(*e.Type))
	oprot.WriteFieldEnd()

	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

func (e *TException) Read(iprot TProtocol) {
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
				e.Message = &v
			} else {
				SkipType(iprot, ftype)
			}
		case 2: // type
			if ftype == TTYPE_I32 {
				v := iprot.ReadI32()
				t := TExceptionType(v)
				e.Type = &t
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

func NewTApplicationException(exceptionType TApplicationExceptionType, message string) *TException {
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

func NewTTransportException(exceptionType TTransportExceptionType, message string) *TException {
	t := TExceptionType(exceptionType)
	return newTException(t, message, "TTransportException")
}
