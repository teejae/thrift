package funky

import (
	"github.com/teejae/go-thrift/thrift"
	"shared"
)

// Structs can also be exceptions, if they are nasty.
// 
// Attributes:
//  - what
//  - why
type InvalidOperation struct {
	thrift.TException
	What *int32
	Why *string
}
func NewInvalidOperation() *InvalidOperation {
	return &InvalidOperation{}
}
func (s *InvalidOperation) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		case 1:
			if ftype == thrift.TTYPE_I32 {
				_v22 := iprot.ReadI32()
				s.What = &_v22
			} else {
				thrift.SkipType(iprot, ftype)
			}
		case 2:
			if ftype == thrift.TTYPE_STRING {
				_v23 := iprot.ReadString()
				s.Why = &_v23
			} else {
				thrift.SkipType(iprot, ftype)
			}
		default:
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}
func (s *InvalidOperation) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("InvalidOperation")
	if s.What != nil {
		oprot.WriteFieldBegin("what", thrift.TTYPE_I32, 1)
		oprot.WriteI32(*s.What)
		oprot.WriteFieldEnd()
	}
	if s.Why != nil {
		oprot.WriteFieldBegin("why", thrift.TTYPE_STRING, 2)
		oprot.WriteString(*s.Why)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

// Attributes:
//  - query_params
//  - required_params
//  - ordered_params
type DataStruct struct {
	Query_params *map[*string]*string
	Required_params *map[*string]*bool
	Ordered_params *[]*string
}
func NewDataStruct() *DataStruct {
	return &DataStruct{}
}
func (s *DataStruct) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		case 1:
			if ftype == thrift.TTYPE_MAP {
				_, _, size := iprot.ReadMapBegin() 
				m := make(map[*string]*string, size)
				for i := int32(0); i < size; i++ {
					key := iprot.ReadString()
					val := iprot.ReadString()
					m[&key] = &val
				}
				s.Query_params = &m
				iprot.ReadMapEnd()
			} else {
				thrift.SkipType(iprot, ftype)
			}
		case 2:
			if ftype == thrift.TTYPE_SET {
				_, size := iprot.ReadSetBegin()
				m := make(map[*string]*bool, size)
				for i := int32(0); i < size; i++ {
					e := iprot.ReadString()
					b := true
					m[&e] = &b
				}
				s.Required_params = &m
				iprot.ReadSetEnd()
			} else {
				thrift.SkipType(iprot, ftype)
			}
		case 3:
			if ftype == thrift.TTYPE_LIST {
				_, size := iprot.ReadListBegin()
				l := make([]*string, size)
				for i := int32(0); i < size; i++ {
					e := iprot.ReadString()
					l = append(l, &e)
				}
				s.Ordered_params = &l
				iprot.ReadListEnd()
			} else {
				thrift.SkipType(iprot, ftype)
			}
		default:
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}
func (s *DataStruct) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("DataStruct")
	if s.Query_params != nil {
		oprot.WriteFieldBegin("query_params", thrift.TTYPE_MAP, 1)
		oprot.WriteMapBegin(thrift.TTYPE_STRING, thrift.TTYPE_STRING, int32(len(*s.Query_params)))
		for k, v := range(*s.Query_params) {
			oprot.WriteString(*k)
			oprot.WriteString(*v)			
		}
		oprot.WriteMapEnd()
		oprot.WriteFieldEnd()
	}
	if s.Required_params != nil {
		oprot.WriteFieldBegin("required_params", thrift.TTYPE_SET, 2)
		oprot.WriteSetBegin(thrift.TTYPE_STRING, int32(len(*s.Required_params)))
		for k := range(*s.Required_params) {
			oprot.WriteString(*k)			
		}
		oprot.WriteSetEnd()
		oprot.WriteFieldEnd()
	}
	if s.Ordered_params != nil {
		oprot.WriteFieldBegin("ordered_params", thrift.TTYPE_LIST, 3)
		oprot.WriteListBegin(thrift.TTYPE_STRING, int32(len(*s.Ordered_params)))
		for _, v := range(*s.Ordered_params) {
			oprot.WriteString(*v)			
		}
		oprot.WriteListEnd()
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}
// Ahh, now onto the cool part, defining a service. Services just need a name
// and can optionally inherit from another service using the extends keyword.
type Sanitizer interface {
	shared.SharedService
	// A method definition looks like C code. It has a return type, arguments,
	// and optionally a list of exceptions that it may throw. Note that argument
	// lists and exception lists are specified using the exact same syntax as
	// field lists in struct or exception definitions.
	Ping() thrift.TException
	// Parameters:
	//  - request_params
	Sanitize_params(request_params DataStruct) (*DataStruct, thrift.TException)
}
// Ahh, now onto the cool part, defining a service. Services just need a name
// and can optionally inherit from another service using the extends keyword.
type SanitizerClient struct {
	*shared.SharedServiceClient
	iprot thrift.TProtocol
	oprot thrift.TProtocol
	seqid int32
}
func NewSanitizerClient(iprot, oprot thrift.TProtocol) *SanitizerClient {
	return &SanitizerClient{
		iprot: iprot,
		oprot: oprot,
		SharedServiceClient: shared.NewSharedServiceClient(iprot, oprot),
	}
}
// A method definition looks like C code. It has a return type, arguments,
// and optionally a list of exceptions that it may throw. Note that argument
// lists and exception lists are specified using the exact same syntax as
// field lists in struct or exception definitions.
func (s *SanitizerClient) Ping() thrift.TException {
	s.seqid++
	s.Send_Ping()
	return s.Recv_Ping()
}
func (s *SanitizerClient) Send_Ping() {
	s.oprot.WriteMessageBegin("ping", thrift.TMESSAGETYPE_CALL, s.seqid)
	args := NewPing_args()
	args.Write(s.oprot)
	s.oprot.WriteMessageEnd()
	s.oprot.GetTransport().Flush()
}

func (s *SanitizerClient) Recv_Ping() thrift.TException {
	_, mtype, _ := s.iprot.ReadMessageBegin()
	defer s.iprot.ReadMessageEnd()
	if mtype == thrift.TMESSAGETYPE_EXCEPTION {
		x := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_UNKNOWN, "")
		x.Read(s.iprot)
		return x
	}
	result := NewPing_result()
	result.Read(s.iprot)
	return nil
}
// Parameters:
//  - request_params
func (s *SanitizerClient) Sanitize_params(request_params DataStruct) (*DataStruct, thrift.TException) {
	s.seqid++
	s.Send_Sanitize_params(request_params)
	return s.Recv_Sanitize_params()
}
func (s *SanitizerClient) Send_Sanitize_params(request_params DataStruct) {
	s.oprot.WriteMessageBegin("sanitize_params", thrift.TMESSAGETYPE_CALL, s.seqid)
	args := NewSanitize_params_args()
	args.Request_params = &request_params
	args.Write(s.oprot)
	s.oprot.WriteMessageEnd()
	s.oprot.GetTransport().Flush()
}

func (s *SanitizerClient) Recv_Sanitize_params() (*DataStruct, thrift.TException) {
	_, mtype, _ := s.iprot.ReadMessageBegin()
	defer s.iprot.ReadMessageEnd()
	if mtype == thrift.TMESSAGETYPE_EXCEPTION {
		x := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_UNKNOWN, "")
		x.Read(s.iprot)
		return nil, x
	}
	result := NewSanitize_params_result()
	result.Read(s.iprot)
	if result.Success != nil {
		return result.Success, nil
	}
	if result.Ouch != nil {
		return nil, result.Ouch
	}
	x := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_MISSING_RESULT, "sanitize_params failed: unknown result")
	return nil, x
}

type SanitizerProcessor struct {
	handler Sanitizer
	parentProcessor *shared.SharedServiceProcessor
}
func NewSanitizerProcessor(handler Sanitizer) *SanitizerProcessor {
	p := &SanitizerProcessor{handler: handler, parentProcessor: shared.NewSharedServiceProcessor(handler)}
	return p
}
func (p *SanitizerProcessor) Process(iprot, oprot thrift.TProtocol) (bool, thrift.TException) {
	name, ttype, seqid := iprot.ReadMessageBegin()
	return p.ProcessMessage(name, ttype, seqid, iprot, oprot)
}
func (p *SanitizerProcessor) ProcessMessage(name string, ttype thrift.TType, seqid int32, iprot, oprot thrift.TProtocol) (bool, thrift.TException) {
	switch name {
	case "ping":
		p.process_Ping(seqid, iprot, oprot)
	case "sanitize_params":
		p.process_Sanitize_params(seqid, iprot, oprot)
	default:
		return p.parentProcessor.ProcessMessage(name, ttype, seqid, iprot, oprot)
	}
	return true, nil
}
func (p *SanitizerProcessor) process_Ping(seqid int32, iprot, oprot thrift.TProtocol) {
	args := NewPing_args()
	args.Read(iprot)
	iprot.ReadMessageEnd()
	result := NewPing_result()
	p.handler.Ping()
	oprot.WriteMessageBegin("ping", thrift.TMESSAGETYPE_REPLY, seqid)
	result.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.GetTransport().Flush()
}
func (p *SanitizerProcessor) process_Sanitize_params(seqid int32, iprot, oprot thrift.TProtocol) {
	args := NewSanitize_params_args()
	args.Read(iprot)
	iprot.ReadMessageEnd()
	result := NewSanitize_params_result()
	result.Success, _ = p.handler.Sanitize_params(*args.Request_params)
	oprot.WriteMessageBegin("sanitize_params", thrift.TMESSAGETYPE_REPLY, seqid)
	result.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.GetTransport().Flush()
}

// HELPER FUNCTIONS AND STRUCTURES

type Ping_args struct {
}
func NewPing_args() *Ping_args {
	return &Ping_args{}
}
func (s *Ping_args) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		default:
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}
func (s *Ping_args) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("ping_args")
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

type Ping_result struct {
}
func NewPing_result() *Ping_result {
	return &Ping_result{}
}
func (s *Ping_result) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		default:
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}
func (s *Ping_result) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("ping_result")
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

// Attributes:
//  - request_params
type Sanitize_params_args struct {
	Request_params *DataStruct
}
func NewSanitize_params_args() *Sanitize_params_args {
	return &Sanitize_params_args{}
}
func (s *Sanitize_params_args) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		case 1:
			if ftype == thrift.TTYPE_STRUCT {
				s.Request_params = NewDataStruct()
				s.Request_params.Read(iprot)
			} else {
				thrift.SkipType(iprot, ftype)
			}
		default:
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}
func (s *Sanitize_params_args) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("sanitize_params_args")
	if s.Request_params != nil {
		oprot.WriteFieldBegin("request_params", thrift.TTYPE_STRUCT, 1)
		s.Request_params.Write(oprot)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

// Attributes:
//  - success
//  - ouch
type Sanitize_params_result struct {
	Success *DataStruct
	Ouch *InvalidOperation
}
func NewSanitize_params_result() *Sanitize_params_result {
	return &Sanitize_params_result{}
}
func (s *Sanitize_params_result) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		case 0:
			if ftype == thrift.TTYPE_STRUCT {
				s.Success = NewDataStruct()
				s.Success.Read(iprot)
			} else {
				thrift.SkipType(iprot, ftype)
			}
		case 1:
			if ftype == thrift.TTYPE_STRUCT {
				s.Ouch = NewInvalidOperation()
				s.Ouch.Read(iprot)
			} else {
				thrift.SkipType(iprot, ftype)
			}
		default:
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}
func (s *Sanitize_params_result) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("sanitize_params_result")
	if s.Success != nil {
		oprot.WriteFieldBegin("success", thrift.TTYPE_STRUCT, 0)
		s.Success.Write(oprot)
		oprot.WriteFieldEnd()
	}
	if s.Ouch != nil {
		oprot.WriteFieldBegin("ouch", thrift.TTYPE_STRUCT, 1)
		s.Ouch.Write(oprot)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}
