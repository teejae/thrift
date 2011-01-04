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

package tutorial

import (
	"thrift/protocol"
)

type SharedStruct struct {
	Key        *int32
	Value      string
	isSetKey   bool
	isSetValue bool
}

func NewSharedStruct() *SharedStruct {
	return &SharedStruct{}
}

func (s *SharedStruct) Read(iprot protocol.TProtocol) {
	iprot.ReadStructBegin()
	for {
		_, ftype, fid := iprot.ReadFieldBegin()
		if ftype == protocol.TTYPE_STOP {
			break
		}

		switch fid {
		case 1: // key
			if ftype == protocol.TTYPE_I32 {
				key := iprot.ReadI32()
				s.Key = &key
			} else {
				protocol.SkipType(iprot, ftype)
			}
		case 2: // value
			if ftype == protocol.TTYPE_STRING {
				s.Value = iprot.ReadString()
			} else {
				protocol.SkipType(iprot, ftype)
			}
		default: // unknown
			protocol.SkipType(iprot, ftype)
		}
		iprot.ReadFieldEnd()
	}
	iprot.ReadStructEnd()
}

func (s *SharedStruct) Write(oprot protocol.TProtocol) {
	oprot.WriteStructBegin("SharedStruct")
	if s.isSetKey {
		oprot.WriteFieldBegin("key", protocol.TTYPE_I32, 1)
		oprot.WriteI32(*s.Key)
		oprot.WriteFieldEnd()
	}
	if s.isSetValue {
		oprot.WriteFieldBegin("value", protocol.TTYPE_STRING, 2)
		oprot.WriteString(s.Value)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}

type SharedService interface {
	getStruct(key int32) SharedStruct
}
