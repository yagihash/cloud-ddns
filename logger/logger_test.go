package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Run("DefaultWithStdout", func(t *testing.T) {
		_, sync, err := New("stdout")
		assert.NoError(t, err)
		defer sync()
	})

	t.Run("ErrorWithNonWritableLog", func(t *testing.T) {
		_, sync, err := New("/dev/urandom")
		assert.NoError(t, err)
		defer sync()
	})

	t.Run("ErrorWithBlankConfig", func(t *testing.T) {
		LogConfigDefault = zap.Config{}
		_, _, err := New("stdout")
		assert.Error(t, err)
	})
}
