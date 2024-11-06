package log

import (
	"fmt"
	"io"
	"strings"
	"time"
)

var DriverKey = "cli"

type DefaultLogRecord struct {
	LogLevel    string
	Time        string
	Message     string
	Attributes  string
	Transaction string
}

type DefaultLog struct {
	Now                       func() time.Time
	Format                    string
	TimeFormat                string
	LogLevel                  int
	IsEnclosedIntoTransaction bool
	Transaction               *Transaction
	out                       io.Writer
}

func (l *DefaultLog) Log(message string, attributes Attributes) {
	rec := DefaultLogRecord{
		Time:    l.Now().Format(l.TimeFormat),
		Message: message,
	}
	trans := ""
	if l.IsEnclosedIntoTransaction {
		trans = l.Transaction.UUID.String()
	}
	rec.LogLevel = l.getLogLevelAsString(L_INFO)
	rec.Attributes = getAttributesAsString(attributes)

	// @todo not concurrency safe
	// @todo add conditional logic to skip log messages based on the selected log level
	fmt.Fprintf(l.out, l.Format, rec.LogLevel, trans, rec.Time, rec.Message, rec.Attributes)
}

func (l *DefaultLog) SetLogLevel(lvl int) {
	l.LogLevel = lvl
}

func (l *DefaultLog) SetTransaction(trans *Transaction) {
	// not concurrency safe
	l.IsEnclosedIntoTransaction = true
	l.Transaction = trans
}

func (l *DefaultLog) ResetTransaction() {
	// not concurrency safe
	l.IsEnclosedIntoTransaction = false
}

func (l *DefaultLog) getLogLevelAsString(lvl int) string {
	levels := []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERR",
	}
	return levels[lvl]
}

func getAttributesAsString(attrs Attributes) string {
	strAttributes := []string{}
	for key, attr := range attrs {
		stringer, ok := attr.(Stringer)
		if ok {
			strAttributes = append(strAttributes, fmt.Sprintf("%s: %s", key, stringer.String()))
			continue
		}
		strAttributes = append(strAttributes, fmt.Sprintf("%s: %v", key, attr))
	}

	return strings.Join(strAttributes, ", ")
}

func NewDefaultLog() DefaultLog {
	return DefaultLog{
		Now:        time.Now,
		Format:     "[%s] %s %s %s %s \n",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		LogLevel:   L_INFO,
	}
}

type DefaultLogDriver struct{}

func (d *DefaultLogDriver) IsSelected(keyFromConfig string) bool {
	return DriverKey == keyFromConfig || keyFromConfig == ""
}

func (d *DefaultLogDriver) Configure(rawConfig []byte) error {
	return nil
}

func (d *DefaultLogDriver) NewLog() Log {
	defaultLog := NewDefaultLog()
	return &defaultLog
}
