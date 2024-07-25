package timeutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonthIntervalTimeFromNow(t *testing.T) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	start, end := MonthIntervalTimeWithGivenMon(thisMonth)
	s, e := MonthIntervalTimeFromNow(0)
	assert.Equal(t, start, s)
	assert.Equal(t, end, e)
}

func TestMonthIntervalTimeWithGivenMon(t *testing.T) {
	tests := []struct {
		mon   string
		start string
		end   string
	}{
		{"2023-07-11", "2023-07-01", "2023-07-31"},
		{"2023-07-01", "2023-07-01", "2023-07-31"},
		{"2023-07-31", "2023-07-01", "2023-07-31"},
		{"2023-12-01", "2023-12-01", "2023-12-31"},
		{"2023-12-21", "2023-12-01", "2023-12-31"},
		{"2023-12-31", "2023-12-01", "2023-12-31"},
		{"2024-01-01", "2024-01-01", "2024-01-31"},
		{"2024-01-11", "2024-01-01", "2024-01-31"},
		{"2024-01-31", "2024-01-01", "2024-01-31"},
	}

	for _, v := range tests {
		t.Run(fmt.Sprintf("mon %s", v.mon), func(t *testing.T) {
			start, end := MonthIntervalTimeWithGivenMon(v.mon)
			assert.Equal(t, v.start, start)
			assert.Equal(t, v.end, end)
		})
	}
}

func TestMonthIntervalTimeWithGivenMonTimeParam(t *testing.T) {
	tests := []struct {
		mon    string
		start  string
		end    string
		layout string
	}{
		{"2023-07-11 15:04:05", "2023-07-01", "2023-07-31", "2006-01-02 15:04:05"},
		{"2023-07-01 12:04:05", "2023-07-01", "2023-07-31", "2006-01-02 15:04:05"},
		{"2023-07-31 00:04:05", "2023-07-01", "2023-07-31", "2006-01-02 15:04:05"},
		{"2023-07-31 23:04:05", "2023-07-01", "2023-07-31", "2006-01-02 15:04:05"},
	}

	for _, v := range tests {
		t.Run(fmt.Sprintf("mon %s", v.mon), func(t *testing.T) {
			start, end := MonthIntervalTimeWithGivenMon(v.mon, v.layout)
			assert.Equal(t, v.start, start)
			assert.Equal(t, v.end, end)
		})
	}
}
