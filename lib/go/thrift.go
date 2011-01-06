package thrift

import (
	"os"
	"thrift/protocol"
)

type TProcessor interface {
	Process(iprot, oprot protocol.TProtocol) (bool, os.Error)
}

type TExceptionType int

type TException struct {
	os.Error
	Message    *string
	Type       *TExceptionType
	structName string
}

func (e *TException) Write(oprot protocol.TProtocol) {
	oprot.WriteStructBegin(e.structName)
	if e.Message != nil {
		oprot.WriteFieldBegin("message", protocol.TTYPE_STRING, 1)
		oprot.WriteString(*e.Message)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldBegin("type", protocol.TTYPE_I32, 2)
	oprot.WriteI32(int32(*e.Type))
	oprot.WriteFieldEnd()

	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

func (e *TException) Read(iprot protocol.TProtocol) {
	iprot.ReadStructBegin()
	defer iprot.ReadStructEnd()

	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == protocol.TTYPE_STOP {
			break
		}

		switch fid {
		case 1: // message
			if ftype == protocol.TTYPE_STRING {
				v := iprot.ReadString()
				e.Message = &v
			} else {
				protocol.SkipType(iprot, ftype)
			}
		case 2: // type
			if ftype == protocol.TTYPE_I32 {
				v := iprot.ReadI32()
				t := TExceptionType(v)
				e.Type = &t
			} else {
				protocol.SkipType(iprot, ftype)
			}
		default: // unknown
			protocol.SkipType(iprot, ftype)
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

type TApplicationException TException

func NewTApplicationException(exceptionType TApplicationExceptionType, message string) *TApplicationException {
	t := TExceptionType(exceptionType)
	return &TApplicationException{Type: &t, Message: &message, structName: "TApplicationException"}
}

type TTransportExceptionType TExceptionType

const (
	TTRANSPORT_EXCEPTION_TYPE_UNKNOWN      = 0
	TTRANSPORT_EXCEPTION_TYPE_NOT_OPEN     = 1
	TTRANSPORT_EXCEPTION_TYPE_ALREADY_OPEN = 2
	TTRANSPORT_EXCEPTION_TYPE_TIMED_OUT    = 3
	TTRANSPORT_EXCEPTION_TYPE_END_OF_FILE  = 4
)

type TTransportException TException

func NewTTransportException(exceptionType TTransportExceptionType, message string) *TTransportException {
	t := TExceptionType(exceptionType)
	return &TTransportException{Type: &t, Message: &message, structName: "TTransportException"}
}
