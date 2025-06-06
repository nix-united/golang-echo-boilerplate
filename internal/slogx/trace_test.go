package slogx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTrace struct {
	Trace string `json:"trace"`
	Index int64  `json:"index"`
}

type testLog struct {
	Level string    `json:"level"`
	Msg   string    `json:"msg"`
	Key   string    `json:"key"`
	Trace testTrace `json:"trace"`
}

func TestTrace(t *testing.T) {
	tracer := NewTraceStarter(func() (uuid.UUID, error) {
		return uuid.MustParse("11111111-1111-1111-1111-111111111111"), nil
	})

	buffer := new(bytes.Buffer)

	logger := slog.New(newTraceHandler(slog.NewJSONHandler(buffer, nil)))

	ctx, err := tracer.Start(t.Context())
	require.NoError(t, err)

	logger.Info("Message")
	logger.InfoContext(ctx, "First message with context")
	logger.InfoContext(ctx, "Second message with context")
	logger.InfoContext(t.Context(), "Message with context but without trace")

	wantLogs := []testLog{
		{
			Level: "INFO",
			Msg:   "Message",
		},
		{
			Level: "INFO",
			Msg:   "First message with context",
			Trace: testTrace{
				Trace: "11111111-1111-1111-1111-111111111111",
				Index: 1,
			},
		},
		{
			Level: "INFO",
			Msg:   "Second message with context",
			Trace: testTrace{
				Trace: "11111111-1111-1111-1111-111111111111",
				Index: 2,
			},
		},
		{
			Level: "INFO",
			Msg:   "Message with context but without trace",
		},
	}

	lines := strings.Split(strings.TrimSuffix(buffer.String(), "\n"), "\n")
	require.Len(t, lines, len(wantLogs))

	for i, line := range lines {
		var gotLog testLog
		err := json.Unmarshal([]byte(line), &gotLog)
		require.NoError(t, err)

		assert.Equal(t, wantLogs[i], gotLog)
	}
}
