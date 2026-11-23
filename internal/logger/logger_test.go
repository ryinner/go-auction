package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func decode(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &m))
	return m
}

func TestLoggerAddsRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, slog.LevelInfo)

	ctx := WithRequestID(context.Background(), "req-1")
	log.InfoContext(ctx, "hello")

	require.Equal(t, "req-1", decode(t, &buf)["request_id"])
}

func TestLoggerWithoutRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, slog.LevelInfo)

	log.InfoContext(context.Background(), "hello")

	require.NotContains(t, decode(t, &buf), "request_id")
}

func TestLoggerKeepsRequestIDAfterWith(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, slog.LevelInfo).With("handler", "health").WithGroup("g")

	ctx := WithRequestID(context.Background(), "req-2")
	log.InfoContext(ctx, "hello", "k", "v")

	got := decode(t, &buf)
	require.Equal(t, "health", got["handler"])
	// после WithGroup атрибут записи попадает в группу — главное, что он не потерялся
	require.Equal(t, map[string]any{"k": "v", "request_id": "req-2"}, got["g"])
}

func TestLoggerRespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, slog.LevelWarn)

	log.InfoContext(context.Background(), "skip me")

	require.Zero(t, buf.Len())
}
