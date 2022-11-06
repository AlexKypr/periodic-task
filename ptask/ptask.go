package ptask

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	tsFormat = "20060102T150405Z"
)

type period string

const (
	hour  period = "1h"
	day   period = "1d"
	month period = "1mo"
	year  period = "1y"
)

func (p period) isValid() error {
	switch p {
	case hour, day, month, year:
		return nil
	}
	return errors.New("invalid period value")
}

type PTask struct {
	period period
	tz     *time.Location
	t1     time.Time
	t2     time.Time
}

func GetPTaskFromURLQueries(r *http.Request, l *log.Logger) (*PTask, error) {
	l.Println("Evaluating url queries")

	prd := r.URL.Query().Get("period")
	if err := period(prd).isValid(); err != nil {
		err = fmt.Errorf("unsupported period: %w", err)
		return nil, err
	}

	loc := r.URL.Query().Get("tz")
	location, err := time.LoadLocation(loc)
	if err != nil {
		err = fmt.Errorf("unsupported location: %w", err)
		return nil, err
	}

	tstr := r.URL.Query().Get("t1")
	t1, err := time.Parse(tsFormat, tstr)
	t1 = t1.In(location)
	if err != nil {
		err = fmt.Errorf("unsupported t1: %w", err)
		return nil, err
	}

	tstr = r.URL.Query().Get("t2")
	t2, err := time.Parse(tsFormat, tstr)
	t2 = t2.In(location)
	if err != nil {
		err = fmt.Errorf("unsupported t2: %w", err)
		return nil, err
	}
	return NewPTask(prd, location, t1, t2), nil
}

func NewPTask(prd string, tz *time.Location, t1, t2 time.Time) *PTask {
	return &PTask{
		period: period(prd),
		tz:     tz,
		t1:     t1,
		t2:     t2,
	}
}

func (p *PTask) GetPTasks() ([]string, error) {
	t1Limit := roundUp(p.t1, p.period)
	ts := []string{}
	runningTs := t1Limit
	for runningTs.Before(p.t2) {
		ts = append(ts, runningTs.UTC().Format(tsFormat))
		runningTs = addPeriod(runningTs, p.period)
	}
	return ts, nil
}

func roundUp(t time.Time, p period) time.Time {
	switch p {
	case year:
		t = t.AddDate(1, -int(t.Month()), 0)
		fallthrough
	case month:
		t = t.AddDate(0, 1, -t.Day())
		fallthrough
	case day:
		t = t.AddDate(0, 0, 1)
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case hour:
		t = t.Round(time.Hour)
	}
	return t
}

func addPeriod(ts time.Time, p period) time.Time {
	switch p {
	case hour:
		ts = ts.Add(time.Hour)
	case day:
		ts = ts.AddDate(0, 0, 1)
	case month:
		ts = ts.AddDate(0, 1, 0)
	case year:
		ts = ts.AddDate(1, 0, 0)
	}
	return ts
}
