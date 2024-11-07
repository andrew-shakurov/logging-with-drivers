package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var DriverKey = "cli"

type defaultLogRecord struct {
	logLevel    string
	time        string
	message     string
	attributes  string
	transaction string
}

type DefaultLog struct {
	now                       func() time.Time
	format                    string
	timeFormat                string
	logLevel                  int
	isEnclosedIntoTransaction bool
	transaction               *Transaction
	messageLogLevelOfLog      int
	outCh                     chan string
	messageAwaitingWG         *sync.WaitGroup
}

func (l *DefaultLog) log(message string, attributes Attributes, messageSevirity int) {
	if messageSevirity < l.logLevel {
		return
	}

	rec := defaultLogRecord{
		time:    l.now().Format(l.timeFormat),
		message: message,
	}
	trans := ""
	if l.isEnclosedIntoTransaction {
		trans = l.transaction.UUID.String()
		if len(l.transaction.Attributes) > 0 {
			trans += " " + getAttributesAsString(l.transaction.Attributes)
		}
	}
	rec.logLevel = l.getLogLevelAsString(L_INFO)
	rec.attributes = getAttributesAsString(attributes)

	l.messageAwaitingWG.Add(1)
	l.outCh <- fmt.Sprintf(l.format, rec.logLevel, trans, rec.time, rec.message, rec.attributes)
}

func (l *DefaultLog) Debug(message string, attributes Attributes) {
	l.log(message, attributes, L_DEBUG)
}

func (l *DefaultLog) Info(message string, attributes Attributes) {
	l.log(message, attributes, L_INFO)
}

func (l *DefaultLog) Warning(message string, attributes Attributes) {
	l.log(message, attributes, L_WARN)
}

func (l *DefaultLog) Error(message string, attributes Attributes) {
	l.log(message, attributes, L_ERR)
}

func (l *DefaultLog) Log(message string, attributes Attributes) {
	l.log(message, attributes, l.messageLogLevelOfLog)
}

func (l *DefaultLog) SetMessageLogLevelOfLog(lvl int) {
	l.messageLogLevelOfLog = lvl
}

func (l *DefaultLog) SetLogLevel(lvl int) {
	l.logLevel = lvl
}

func (l *DefaultLog) SetTransaction(trans *Transaction) {
	// not concurrency safe
	l.isEnclosedIntoTransaction = true
	l.transaction = trans
}

func (l *DefaultLog) ResetTransaction() {
	// not concurrency safe
	l.isEnclosedIntoTransaction = false
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

func newDefaultLog() DefaultLog {
	return DefaultLog{
		now:                  time.Now,
		format:               "[%s] %s %s %s %s \n",
		timeFormat:           "2006-01-02T15:04:05Z07:00",
		logLevel:             L_INFO,
		messageLogLevelOfLog: L_INFO,
	}
}

type DefaultLogDriver struct {
	outCh           chan string
	MessageAwaiting *sync.WaitGroup
}

func (d *DefaultLogDriver) IsSelected(keyFromConfig string) bool {
	return DriverKey == keyFromConfig || keyFromConfig == ""
}

func (d *DefaultLogDriver) Configure(rawConfig []byte) error {
	return nil
}

func (d *DefaultLogDriver) NewLog() Log {
	defaultLog := newDefaultLog()
	defaultLog.outCh = d.outCh
	defaultLog.messageAwaitingWG = d.MessageAwaiting
	return &defaultLog
}

func (d *DefaultLogDriver) Close() {
	d.MessageAwaiting.Wait()
}

func newDefaultLogDriver() *DefaultLogDriver {
	outCh := make(chan string)
	messageAwaitingWG := &sync.WaitGroup{}

	go func() {
		for message := range outCh {
			os.Stdout.Write([]byte(message))
			messageAwaitingWG.Done()
		}
	}()

	return &DefaultLogDriver{
		outCh:           outCh,
		MessageAwaiting: messageAwaitingWG,
	}
}
