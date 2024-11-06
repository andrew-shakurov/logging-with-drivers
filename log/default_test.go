package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Run("a log message gets printed", func(t *testing.T) {
		message := "abc"
		timeNow := time.Time{}

		buff := bytes.Buffer{}
		log := NewDefaultLog()
		log.out = &buff
		log.Now = func() time.Time { return timeNow }

		log.Log(message, nil)

		abc := buff.String()
		print(abc)
		expected := "[INFO]  0001-01-01T00:00:00Z abc  \n"
		assert.Equal(t, expected, buff.String())
	})
}
