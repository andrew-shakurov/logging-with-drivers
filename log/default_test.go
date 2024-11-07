package log

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Run("a log message gets printed", func(t *testing.T) {
		log, buff, messageAwaitingWg := getLogStreamingToBuff()

		message := "abc"
		log.Log(message, nil)

		expected := "[INFO]  0001-01-01T00:00:00Z abc  \n"

		messageAwaitingWg.Wait()
		assert.Equal(t, expected, buff.String())
	})

	t.Run("an attribute is printed, when provided", func(t *testing.T) {
		log, buff, messageAwaitingWg := getLogStreamingToBuff()

		attributeKey := "userId"
		attributeVal := 123

		log.Log("", map[string]interface{}{attributeKey: attributeVal})

		expected := "[INFO]  0001-01-01T00:00:00Z  userId: 123 \n"

		messageAwaitingWg.Wait()
		assert.Equal(t, expected, buff.String())
	})

	t.Run("two attributes are printed, when provided", func(t *testing.T) {
		log, buff, messageAwaitingWg := getLogStreamingToBuff()

		attributeKey := "userId"
		attributeVal := 123
		anotherAttributeKey := "httpMethod"
		anotherValue := "GET"

		log.Log("", map[string]interface{}{
			attributeKey:        attributeVal,
			anotherAttributeKey: anotherValue,
		})

		expected := "[INFO]  0001-01-01T00:00:00Z  userId: 123, httpMethod: GET \n"

		messageAwaitingWg.Wait()
		assert.Equal(t, expected, buff.String())
	})

	t.Run("transaction is printed, when enabled", func(t *testing.T) {
		log, buff, messageAwaitingWg := getLogStreamingToBuff()
		transUUIDString := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		trans := &Transaction{
			UUID: uuid.Must(uuid.Parse(transUUIDString)),
		}
		log.SetTransaction(trans)

		log.Log("", nil)

		expected := "[INFO] 6ba7b810-9dad-11d1-80b4-00c04fd430c8 0001-01-01T00:00:00Z   \n"

		messageAwaitingWg.Wait()
		assert.Equal(t, expected, buff.String())
	})

	t.Run("transaction and its attributes are printed when ts is enabled and attrs are specified", func(t *testing.T) {
		log, buff, messageAwaitingWg := getLogStreamingToBuff()
		transUUIDString := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		trans := &Transaction{
			UUID:       uuid.Must(uuid.Parse(transUUIDString)),
			Attributes: map[string]interface{}{"abc": 123},
		}
		log.SetTransaction(trans)

		log.Log("", nil)

		expected := "[INFO] 6ba7b810-9dad-11d1-80b4-00c04fd430c8 abc: 123 0001-01-01T00:00:00Z   \n"

		messageAwaitingWg.Wait()
		assert.Equal(t, expected, buff.String())
	})

	t.Run("an INFO log message is completely omited, when log level set to ERROR", func(t *testing.T) {
		log, buff, messageAwaitingWg := getLogStreamingToBuff()
		log.LogLevel = L_ERR

		message := "abc"
		log.Log(message, nil)

		expected := ""

		messageAwaitingWg.Wait()
		assert.Equal(t, expected, buff.String())
	})

	t.Run("a few concurrently called Log() produce messages, printed sequentially", func(t *testing.T) {
		log, buff, messageAwaitingWg := getLogStreamingToBuff()

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				log.Log("abc", nil)
				wg.Done()
			}()
		}

		expected := `[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
[INFO]  0001-01-01T00:00:00Z abc  
`
		wg.Wait()
		messageAwaitingWg.Wait()
		assert.Equal(t, expected, buff.String())
	})
}

func getLogStreamingToBuff() (*DefaultLog, *bytes.Buffer, *sync.WaitGroup) {
	someFixedPointInTime := time.Time{}
	log := NewDefaultLog()
	log.Now = func() time.Time { return someFixedPointInTime }
	outCh := make(chan string)
	log.outCh = outCh
	messageAwaitingWG := &sync.WaitGroup{}
	log.messageAwaitingWG = messageAwaitingWG
	buff := &bytes.Buffer{}

	go func() {
		for message := range outCh {
			buff.Write([]byte(message))
			messageAwaitingWG.Done()
		}
	}()

	return &log, buff, messageAwaitingWG
}
