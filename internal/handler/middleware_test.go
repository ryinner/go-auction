package handler

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ryinner/go-auction/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {
	tests := []struct {
		name     string
		incoming string
		wantSame bool // true — должен сохраниться входящий id
	}{
		{name: "generated when missing", incoming: "", wantSame: false},
		{name: "kept when valid", incoming: "abc-123_DEF", wantSame: true},
		{name: "replaced when too long", incoming: strings.Repeat("a", maxRequestIDLen+1), wantSame: false},
		{name: "replaced when unsafe chars", incoming: "bad id\n", wantSame: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctxID string
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctxID, _ = logger.RequestID(r.Context())
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.incoming != "" {
				req.Header.Set(requestIDHeader, tt.incoming)
			}
			rec := httptest.NewRecorder()

			RequestID(next).ServeHTTP(rec, req)

			require.NotEmpty(t, ctxID)
			require.Equal(t, ctxID, rec.Header().Get(requestIDHeader))
			if tt.wantSame {
				require.Equal(t, tt.incoming, ctxID)
			} else {
				require.NotEqual(t, tt.incoming, ctxID)
				require.Len(t, ctxID, 32)
			}
		})
	}
}

func TestAccessLogHasRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, slog.LevelInfo)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	h := RequestID(AccessLog(log)(next))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(requestIDHeader, "req-42")
	h.ServeHTTP(httptest.NewRecorder(), req)

	var line map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &line))
	require.Equal(t, "req-42", line["request_id"])
	require.Equal(t, "/x", line["path"])
	require.EqualValues(t, http.StatusTeapot, line["status"])
}
