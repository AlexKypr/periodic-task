package ptask

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newMockPTask(prd string, tz *time.Location, t1 time.Time, t2 time.Time) *PTask {
	return NewPTask(prd, tz, t1, t2)
}

func TestGetPTasks(t *testing.T) {
	tests := []struct {
		ptask *PTask
		want  []string
		err   error
	}{
		{newMockPTask("1h", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2021, 10, 9, 15, 15, 0, 0, time.UTC)), []string{"20211009T110000Z", "20211009T120000Z", "20211009T130000Z", "20211009T140000Z", "20211009T150000Z"}, nil},
		{newMockPTask("1d", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2021, 10, 14, 15, 15, 0, 0, time.UTC)), []string{"20211010T000000Z", "20211011T000000Z", "20211012T000000Z", "20211013T000000Z", "20211014T000000Z"}, nil},
		{newMockPTask("1w", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2021, 10, 9, 15, 15, 0, 0, time.UTC)), nil, errors.New("invalid period value")},
		{newMockPTask("1mo", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2022, 1, 9, 15, 15, 0, 0, time.UTC)), []string{"20211101T000000Z", "20211201T000000Z", "20220101T000000Z"}, nil},
		{newMockPTask("1y", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2024, 11, 9, 15, 15, 0, 0, time.UTC)), []string{"20220101T000000Z", "20230101T000000Z", "20240101T000000Z"}, nil},
		{newMockPTask("1h", time.UTC, time.Date(2021, 10, 11, 10, 29, 0, 0, time.UTC), time.Date(2021, 10, 9, 15, 15, 0, 0, time.UTC)), []string{}, nil},
	}
	for _, test := range tests {
		got, err := test.ptask.GetPTasks()
		require.Equal(t, test.want, got)
		require.Equal(t, test.err, err)
	}
}

func TestGetPTaskFromURLQueries(t *testing.T) {
	happytests := []struct {
		url  string
		want *PTask
	}{
		{"/ptlist?period=1h&tz=UTC&t1=20211009T102900Z&t2=20211009T151500Z", newMockPTask("1h", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2021, 10, 9, 15, 15, 0, 0, time.UTC))},
		{"/ptlist?period=1d&tz=UTC&t1=20211009T102900Z&t2=20211014T151500Z", newMockPTask("1d", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2021, 10, 14, 15, 15, 0, 0, time.UTC))},
		{"/ptlist?period=1mo&tz=UTC&t1=20211009T102900Z&t2=20220109T151500Z", newMockPTask("1mo", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2022, 1, 9, 15, 15, 0, 0, time.UTC))},
		{"/ptlist?period=1y&tz=UTC&t1=20211009T102900Z&t2=20241109T151500Z", newMockPTask("1y", time.UTC, time.Date(2021, 10, 9, 10, 29, 0, 0, time.UTC), time.Date(2024, 11, 9, 15, 15, 0, 0, time.UTC))},
		{"/ptlist?period=1h&tz=UTC&t1=20211011T102900Z&t2=20211009T151500Z", newMockPTask("1h", time.UTC, time.Date(2021, 10, 11, 10, 29, 0, 0, time.UTC), time.Date(2021, 10, 9, 15, 15, 0, 0, time.UTC))},
	}
	badtests := []struct {
		url string
	}{
		{"/ptlist?period=1w&tz=UTC&t1=20211009T102900Z&t2=20211009T151500Z"},
		{"/ptlist?period=1h&tz=UTC&t1=2021009T102900Z&t2=20211009T151500Z"},
		{"/ptlist?period=1h&tz=UTC&t1=20211009T102900Z&t2=20211009151500Z"},
		{"/ptlist?period=1h&tz=UTCS&t1=20211009T102900Z&t2=20211009T151500Z"},
	}

	l, err := zap.NewProduction()
	require.NoError(t, err)

	for _, test := range happytests {
		r, err := http.NewRequest(http.MethodGet, test.url, nil)
		require.NoError(err)

		got, err := GetPTaskFromURLQueries(r, l)
		require.NoError(t, err)
		require.Equal(t, test.want, got)
	}

	for _, test := range badtests {
		r, err := http.NewRequest(http.MethodGet, test.url, nil)
		require.NoError(err)

		got, err := GetPTaskFromURLQueries(r, l)
		require.Empty(t, got)
		require.Error(t, err)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		input string
		want  error
	}{
		{input: "1h", want: nil},
		{input: "1d", want: nil},
		{input: "1w", want: errors.New("invalid period value")},
		{input: "1mo", want: nil},
		{input: "1y", want: nil},
	}
	for _, test := range tests {
		got := period(test.input).isValid()
		require.Equal(t, test.want, got)
	}
}

func TestAddPeriod(t *testing.T) {
	tests := []struct {
		inputTime   time.Time
		inputPeriod period
		want        time.Time
		err         error
	}{
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), hour, time.Date(2020, 10, 9, 22, 50, 15, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), day, time.Date(2020, 10, 10, 21, 50, 15, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), month, time.Date(2020, 11, 9, 21, 50, 15, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), year, time.Date(2021, 10, 9, 21, 50, 15, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), period("1w"), time.Time{}, errors.New("invalid period value")},
	}
	for _, test := range tests {
		got, err := addPeriod(test.inputTime, test.inputPeriod)
		require.Equal(t, test.want, got)
		require.Equal(t, test.err, err)
	}
}

func TestRoundUp(t *testing.T) {
	tests := []struct {
		inputTime   time.Time
		inputPeriod period
		want        time.Time
		err         error
	}{
		{time.Date(2020, 10, 9, 21, 10, 15, 0, time.UTC), hour, time.Date(2020, 10, 9, 22, 0, 0, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), hour, time.Date(2020, 10, 9, 22, 0, 0, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), day, time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), month, time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), year, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), nil},
		{time.Date(2020, 10, 9, 21, 50, 15, 0, time.UTC), period("1w"), time.Time{}, errors.New("invalid period value")},
	}
	for _, test := range tests {
		got, err := roundUp(test.inputTime, test.inputPeriod)
		require.Equal(t, test.want, got)
		require.Equal(t, test.err, err)
	}
}
