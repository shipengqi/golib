package timeutil

import (
	"fmt"
	"time"
)

// MonthIntervalTimeFromNow return the start date and end date of the month
// if mon = 0, indicate current month
// if mon = -1, indicate last month
// if mon = 1, indicate next month
func MonthIntervalTimeFromNow(mon int) (start, end string) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start = thisMonth.AddDate(0, mon, 0).Format("2006-01-02")
	end = thisMonth.AddDate(0, mon+1, -1).Format("2006-01-02")
	return start, end
}

// MonthIntervalTimeWithGivenMon return the start date and end date of the given month
func MonthIntervalTimeWithGivenMon(mon string, layout ...string) (start, end string) {
	timeLayout := "2006-01-02"
	if len(layout) > 0 && layout[0] != "" {
		timeLayout = layout[0]
	}
	conv, err := time.ParseInLocation(timeLayout, mon, time.Local)
	if err != nil {
		fmt.Println(err)
	}
	year, month, _ := conv.Date()
	givenMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start = givenMonth.AddDate(0, 0, 0).Format("2006-01-02")
	end = givenMonth.AddDate(0, 1, -1).Format("2006-01-02")
	return start, end
}
