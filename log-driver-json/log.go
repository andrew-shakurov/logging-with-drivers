package logdriverjson

import (
	"encoding/json"
	"fmt"
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
	MessageLogLevelOfLog      int
}

func (l *JSONLog) Debug(message string, attributes log.Attributes) {
	l.log(message, attributes, log.L_DEBUG)
}

func (l *JSONLog) Info(message string, attributes log.Attributes) {
	l.log(message, attributes, log.L_INFO)
}

func (l *JSONLog) Warning(message string, attributes log.Attributes) {
	l.log(message, attributes, log.L_WARN)
}

func (l *JSONLog) Error(message string, attributes log.Attributes) {
	l.log(message, attributes, log.L_ERR)
}

func (l *JSONLog) Log(message string, attributes log.Attributes) {
	l.log(message, attributes, l.MessageLogLevelOfLog)
}

func (l *JSONLog) SetMessageLogLevelOfLog(lvl int) {
	l.MessageLogLevelOfLog = lvl
}

func (l *JSONLog) log(message string, attributes log.Attributes, messageSevirity int) {
	if messageSevirity < l.LogLevel {
		return
	}

	// @todo make concurency safe
	record := &JSONLogRecord{}
	record.Message = message
	for key, attr := range attributes {
		jsonAttr := JSONLogRecordAttribute{
			Key: key,
		}
		stringer, ok := attr.(log.Stringer)
		if ok {
			jsonAttr.Value = stringer.String()
			record.Attributes = append(record.Attributes, jsonAttr)
			continue
		}

		jsonAttr.Value = fmt.Sprintf("%v", attr)
		record.Attributes = append(record.Attributes, jsonAttr)
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
