package handlers

import (
	"encoding/json"
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

func (p period) IsValid() error {
	switch p {
	case hour, day, month, year:
		return nil
	}
	return errors.New("invalid period value")
}

type errorDto struct {
	Status string `json:"status"`
	Desc   string `json:"desc"`
}

func newErrorDTO(status, desc string) errorDto {
	return errorDto{
		Status: status,
		Desc:   desc,
	}
}

func toJSON(rw http.ResponseWriter, data interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(data)
}

type ptlist struct {
	period period
	tz     *time.Location
	t1     time.Time
	t2     time.Time
}

func newPtlist(period period, tz *time.Location, t1, t2 time.Time) *ptlist {
	return &ptlist{
		period: period,
		tz:     tz,
		t1:     t1,
		t2:     t2,
	}
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

func (p *ptlist) tsArray() ([]string, error) {
	t1Limit := roundUp(p.t1, p.period)
	ts := []string{}
	runningTs := t1Limit
	for runningTs.Before(p.t2) {
		ts = append(ts, runningTs.UTC().Format(tsFormat))
		runningTs = addPeriod(runningTs, p.period)
	}
	return ts, nil
}

type PtlistQuery struct {
	l *log.Logger
}

func NewPtlistQuery(l *log.Logger) *PtlistQuery {
	return &PtlistQuery{
		l: l,
	}
}

func (p *PtlistQuery) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getTsArray(rw, r)
		return
	}

	// catch http methods that are not supported
	rw.WriteHeader(http.StatusMethodNotAllowed) // maybe replace it with 400
}

func (p *PtlistQuery) evaluateURLQueries(r *http.Request) (*ptlist, error) {
	p.l.Println("Evaluating url queries")

	prd := r.URL.Query().Get("period")
	if err := period(prd).IsValid(); err != nil {
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
	return newPtlist(period(prd), location, t1, t2), nil
}

func (p *PtlistQuery) getTsArray(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Getting timestamp array")

	ptlist, err := p.evaluateURLQueries(r)
	if err != nil {
		p.l.Println(err.Error())
		errDto := newErrorDTO("error", err.Error())
		toJSON(rw, errDto, http.StatusNotFound)
		return
	}

	tsarray, err := ptlist.tsArray()
	if err != nil {
		p.l.Println(err.Error())
		errDto := newErrorDTO("error", err.Error())
		toJSON(rw, errDto, http.StatusNotFound)
		return
	}

	toJSON(rw, tsarray, http.StatusOK)
}
