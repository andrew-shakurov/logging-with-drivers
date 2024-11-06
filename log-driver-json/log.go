package logdriverjson

import (
	"encoding/json"
	"io"

	"example.com/log"
)

const DriverKey = "json-file"

type JSONLogRecord struct {
	LogLvel    string                   `json:"sevirity"`
	Time       int                      `json:"time"`
	Message    string                   `json:"message"`
	Attributes []JSONLogRecordAttribute `json:"attributes"`
}

type JSONLogRecordAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type JSONLog struct {
	config                    Config
	out                       io.Writer
	LogLevel                  int
	IsEnclosedIntoTransaction bool
	Transaction               *log.Transaction
}

func (l *JSONLog) Log(message string, attributes log.Attributes) {
	// @todo make concurency safe
	record := &JSONLogRecord{}
	record.Message = message
	for key, attr := range attributes {
		record.Attributes = append(record.Attributes, JSONLogRecordAttribute{
			Key:   key,
			Value: attr.String(),
		})
	}
	// no meaningful way to handle, supress posible error
	encRecord, _ := json.Marshal(record)
	l.out.Write(encRecord)
}

func (l *JSONLog) SetLogLevel(lvl int) {
	l.LogLevel = lvl
}

func (l *JSONLog) SetTransaction(trans *log.Transaction) {
	// not concurrency safe
	l.IsEnclosedIntoTransaction = true
	l.Transaction = trans
}

func (l *JSONLog) ResetTransaction() {
	// not concurrency safe
	l.IsEnclosedIntoTransaction = false
}