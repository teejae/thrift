/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * This Thrift file can be included by other Thrift files that want to share
 * these definitions.
 */

package shared

import (
	"fmt"
	"github.com/teejae/go-thrift/thrift"
)

type SharedStruct struct {
	Key   *int32
	Value *string
}

func NewSharedStruct() *SharedStruct {
	return &SharedStruct{}
}

func (s *SharedStruct) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		case 1: // key
			if ftype == thrift.TTYPE_I32 {
				v := iprot.ReadI32()
				s.Key = &v
			} else {
				thrift.SkipType(iprot, ftype)
			}
		case 2: // value
			if ftype == thrift.TTYPE_STRING {
				v := iprot.ReadString()
				s.Value = &v
			} else {
				thrift.SkipType(iprot, ftype)
			}
		default: // unknown
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}

func (s *SharedStruct) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("SharedStruct")
	if s.Key != nil {
		oprot.WriteFieldBegin("key", thrift.TTYPE_I32, 1)
		oprot.WriteI32(*s.Key)
		oprot.WriteFieldEnd()
	}
	if s.Value != nil {
		oprot.WriteFieldBegin("value", thrift.TTYPE_STRING, 2)
		oprot.WriteString(*s.Value)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}


type SharedService interface {
	GetStruct(key int32) (*SharedStruct, *thrift.TException)
}

type SharedServiceClient struct {
	iprot thrift.TProtocol
	oprot thrift.TProtocol
	seqid int32
}

func NewSharedServiceClient(iprot, oprot thrift.TProtocol) *SharedServiceClient {
	return &SharedServiceClient{iprot: iprot, oprot: oprot}
}

func (s *SharedServiceClient) GetStruct(key int32) (*SharedStruct, thrift.TException) {
	s.send_getStruct(key)
	return s.recv_getStruct()
}

func (c *SharedServiceClient) send_getStruct(key int32) {
	c.oprot.WriteMessageBegin("getStruct", thrift.TMESSAGETYPE_CALL, c.seqid)
	args := new_getStruct_args()
	args.Key = &key
	args.Write(c.oprot)
	c.oprot.WriteMessageEnd()
	c.oprot.GetTransport().Flush()
}

func (c *SharedServiceClient) recv_getStruct() (*SharedStruct, thrift.TException) {
	_, mtype, _ := c.iprot.ReadMessageBegin()
	defer c.iprot.ReadMessageEnd()

	if mtype == thrift.TMESSAGETYPE_EXCEPTION {
		x := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_UNKNOWN, "")
		x.Read(c.iprot)
		return nil, x
	}
	result := new_getStruct_result()
	result.Read(c.iprot)
	if result.Success != nil {
		return result.Success, nil
	}
	x := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_MISSING_RESULT, "getStruct failed: unknown result")
	return nil, x
}

type SharedServiceProcessor struct {
	handler SharedService
}

func NewSharedServiceProcessor(handler SharedService) *SharedServiceProcessor {
	p := &SharedServiceProcessor{handler: handler}
	return p
}

func (p *SharedServiceProcessor) Process(iprot, oprot thrift.TProtocol) (bool, thrift.TException) {
	name, _, seqid := iprot.ReadMessageBegin()
	switch name {
	case "getStruct":
		p.process_GetStruct(seqid, iprot, oprot)
	default:
		thrift.SkipType(iprot, thrift.TTYPE_STRUCT)
		iprot.ReadMessageEnd()
		err := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_UNKNOWN_METHOD, fmt.Sprintf("Unknown function %s", name))
		oprot.WriteMessageBegin(name, TMessageType.EXCEPTION, seqid)
		err.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.GetTransport().Flush()
		return false, err // FIXME: check that this is correct. we may not need this.
	}
	return true, nil
}

func (p *SharedServiceProcessor) process_GetStruct(seqid int32, iprot, oprot thrift.TProtocol) {
	args := new_getStruct_args()
	args.Read(iprot)
	iprot.ReadMessageEnd()
	result := new_getStruct_result()
	result.Success, _ := p.handler.GetStruct(*args.Key)
	oprot.WriteMessageBegin("getStruct", thrift.TMESSAGETYPE_REPLY, seqid)
	result.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.GetTransport().Flush()
}


type getStruct_args struct {
	Key *int32
}

func new_getStruct_args() *getStruct_args {
	return &getStruct_args{}
}

func (s *getStruct_args) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("getStruct_args")
	if s.Key != nil {
		oprot.WriteFieldBegin("key", thrift.TTYPE_I32, 1)
		oprot.WriteI32(*s.Key)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

func (s *getStruct_args) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	defer iprot.ReadStructEnd()

	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		case 1: // key
			if ftype == thrift.TTYPE_I32 {
				v := iprot.ReadI32()
				s.Key = &v
			} else {
				thrift.SkipType(iprot, ftype)
			}
		default: // unknown
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
}

type getStruct_result struct {
	Success *SharedStruct
}

func new_getStruct_result() *getStruct_result {
	return &getStruct_result{}
}

func (s *getStruct_result) Write(oprot thrift.TProtocol) {
	oprot.WriteStructBegin("getStruct_result")
	if s.Success != nil {
		oprot.WriteFieldBegin("success", thrift.TTYPE_STRUCT, 0)
		s.Success.Write(oprot)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

func (s *getStruct_result) Read(iprot thrift.TProtocol) {
	iprot.ReadStructBegin()
	defer iprot.ReadStructEnd()

	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == thrift.TTYPE_STOP {
			break
		}

		switch fid {
		case 0: // key
			if ftype == thrift.TTYPE_STRUCT {
				s.Success = NewSharedStruct()
				s.Success.Read(iprot)
			} else {
				thrift.SkipType(iprot, ftype)
			}
		default: // unknown
			thrift.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
}
