package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetTsArray(t *testing.T) {
	tests := []struct {
		method   string
		url      string
		httpCode int
	}{
		{"GET", "http://127.0.0.1:8080/ptlist?period=1y&tz=Europe/Athens&t1=20180214T204603Z&t2=20211115T123456Z", 200},
		{"GET", "http://127.0.0.1:8080/ptlist?period=1y&tz=Europe/Athens&t1=20180214T204603Z&t2=2021115T123456Z", 400},
		{"POST", "http://127.0.0.1:8080/ptlist?period=1y&tz=Europe/Athens&t1=20180214T204603Z&t2=20211115T123456Z", 400},
	}

	l, err := zap.NewProduction()
	require.NoError(t, err)

	for _, test := range tests {
		r, err := http.NewRequest(test.method, test.url, nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		th := NewTaskHandler(l)
		th.ServeHTTP(recorder, r)

		require.Equal(t, test.httpCode, recorder.Code)
	}
}
