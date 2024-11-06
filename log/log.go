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
	Debug(message string, attributes Attributes)
	Info(message string, attributes Attributes)
	Warning(message string, attributes Attributes)
	Error(message string, attributes Attributes)

	SetLogLevel(int)
	SetTransaction(*Transaction)
	ResetTransaction()

	Log(message string, attributes Attributes)
	SetMessageLogLevelOfLog(int)
}
