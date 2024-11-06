package log

import (
	"github.com/google/uuid"
)

const (
	L_DEBUG = iota
	L_INFO
	L_WARN
	L_ERR
)

type Stringer interface {
	String() string
}

type Attributes map[string]interface{}

type Transaction struct {
	UUID       uuid.UUID
	Attributes Attributes
}

type Log interface {
	Log(message string, attributes Attributes)
	SetLogLevel(int)
	SetTransaction(*Transaction)
	ResetTransaction()
}
