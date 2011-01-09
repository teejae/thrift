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
	"os"
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
	GetStruct(key int32) *SharedStruct
}

type SharedServiceClient struct {
	iprot thrift.TProtocol
	oprot thrift.TProtocol
	seqid int32
}

func NewSharedServiceClient(iprot, oprot thrift.TProtocol) *SharedServiceClient {
	return &SharedServiceClient{iprot: iprot, oprot: oprot}
}

func (s *SharedServiceClient) GetStruct(key int32) *SharedStruct {
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

func (c *SharedServiceClient) recv_getStruct() *SharedStruct {
	_, mtype, _ := c.iprot.ReadMessageBegin()
	defer c.iprot.ReadMessageEnd()

	if mtype == thrift.TMESSAGETYPE_EXCEPTION {
		// FIXME: finish this
		// x = TApplicationException()
		// x.read(self._iprot)
		// self._iprot.readMessageEnd()
		// raise x	
	}
	result := new_getStruct_result()
	result.Read(c.iprot)
	if result.Success != nil {
		return result.Success
	}
	// FIXME
	// raise TApplicationException(TApplicationException.MISSING_RESULT, "getStruct failed: unknown result");
	panic("getStruct failed: unknown result")
}

type processFunc func(seqid int32, iprot, oprot thrift.TProtocol)
type processMap map[string]processFunc

type SharedServiceProcessor struct {
	handler SharedService
	pMap    processMap
}

func NewSharedServiceProcessor(handler SharedService) *SharedServiceProcessor {
	// map services
	pMap := make(processMap)
	p := &SharedServiceProcessor{handler: handler, pMap: pMap}
	pMap["getStruct"] = func(seqid int32, iprot, oprot thrift.TProtocol) {
		p.process_GetStruct(seqid, iprot, oprot)
	}

	return p
}

func (p *SharedServiceProcessor) Process(iprot, oprot thrift.TProtocol) (bool, os.Error) {
	name, _, seqid := iprot.ReadMessageBegin()
	if f := p.pMap[name]; f != nil {
		f(seqid, iprot, oprot)
		return true, nil
	}
	thrift.SkipType(iprot, thrift.TTYPE_STRUCT)
	iprot.ReadMessageEnd()
	err := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_UNKNOWN_METHOD, fmt.Sprintf("Unknown function %s", name))
	return false, err
}

func (p *SharedServiceProcessor) process_GetStruct(seqid int32, iprot, oprot thrift.TProtocol) {
	args := new_getStruct_args()
	args.Read(iprot)
	iprot.ReadMessageEnd()
	result := new_getStruct_result()
	result.Success = p.handler.GetStruct(*args.Key)
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
