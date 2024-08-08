package timeutil

import (
	"time"
)

// MonthIntervalTimeFromNum return the start date and end date of the given number
// if mon = 0, indicate current month
// if mon = -1, indicate last month
// if mon = 1, indicate next month
func MonthIntervalTimeFromNum(mon int) (start, end string) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start = thisMonth.AddDate(0, mon, 0).Format("2006-01-02")
	end = thisMonth.AddDate(0, mon+1, -1).Format("2006-01-02")
	return start, end
}

// MonthIntervalTimeFromMon return the start date and end date of the given month
func MonthIntervalTimeFromMon(mon string, layout ...string) (start, end string, err error) {
	timeLayout := "2006-01-02"
	if len(layout) > 0 && layout[0] != "" {
		timeLayout = layout[0]
	}
	var conv time.Time
	conv, err = time.ParseInLocation(timeLayout, mon, time.Local)
	if err != nil {
		return "", "", err
	}
	year, month, _ := conv.Date()
	givenMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start = givenMonth.AddDate(0, 0, 0).Format("2006-01-02")
	end = givenMonth.AddDate(0, 1, -1).Format("2006-01-02")
	return start, end, nil
}
