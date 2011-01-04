package protocol

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
