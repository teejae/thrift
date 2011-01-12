package funky

import (
	"github.com/teejae/go-thrift/thrift"
)

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
