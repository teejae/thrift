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
	"thrift/protocol"
)

type SharedStruct struct {
	_key        int32
	_value      string
	_isSetKey   bool
	_isSetValue bool
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
				s.SetKey(key)
			} else {
				protocol.SkipType(iprot, ftype)
			}
		case 2: // value
			if ftype == protocol.TTYPE_STRING {
				value := iprot.ReadString()
				s.SetValue(value)
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
	if s._isSetKey {
		oprot.WriteFieldBegin("key", protocol.TTYPE_I32, 1)
		oprot.WriteI32(s._key)
		oprot.WriteFieldEnd()
	}
	if s._isSetValue {
		oprot.WriteFieldBegin("value", protocol.TTYPE_STRING, 2)
		oprot.WriteString(s._value)
		oprot.WriteFieldEnd()
	}
	oprot.WriteFieldStop()
	oprot.WriteStructEnd()
}


func (s *SharedStruct) GetKey() int32 {
	return s._key
}

func (s *SharedStruct) SetKey(key int32) {
	s._key = key
	s._isSetKey = true
}

func (s *SharedStruct) GetValue() string {
	return s._value
}

func (s *SharedStruct) SetValue(value string) {
	s._value = value
	s._isSetValue = true
}

type SharedService interface {
	getStruct(key int32) SharedStruct
}
