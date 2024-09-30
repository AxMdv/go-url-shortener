package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	err := InitLogger()
	require.NoError(t, err)
	Log.Info("Example message:",
		zap.String("Key 1", "Value 1"),
		zap.String("Key 2", "Value 2"),
	)
}
