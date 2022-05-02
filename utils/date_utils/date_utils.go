package date_utils

import (
	"fmt"
	"time"
)

const (
	UTCDateLayout      = "2006-01-02T15:04:05Z"
	DateyyyymmddLayout = "2006-01-02"
	DateyyyymmLayout   = "2006-01"
	DateyyyyLayout     = "2006"
	hybrisDateLayout   = "2006-01-02T15:04:05z"
	ApiDbLayout        = "2006-01-02 15:04:05"
	SearchDateLayout   = "2006-01-02T15:04:05Z"
	fileDateLayout     = "20060102T150405"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	return GetNow().Format(UTCDateLayout)
}

func GetNowFileLayout() string {
	return GetNow().Format(fileDateLayout)
}

// func GetNowSearchFormat() string {
// 	return GetNow().Format(ApiDbLayout)
// }
func GetNowSearchFormat() string {
	return GetNow().Format(SearchDateLayout)
}

//add or substract hours, if hour is passed as 0 then no change, -1 would substract 1 hour from given date
func GetHybrisDateFormat(d string, hours time.Duration) string {

	t, _ := time.Parse(ApiDbLayout, d)
	if hours != 0 {
		t = t.Add(time.Hour * hours)
	}
	return t.Format(hybrisDateLayout)
}

//add or substract hours, if hour is passed as 0 then no change, -1 would substract 1 hour from given date
func GetDBDateFormat(d string, hours time.Duration) string {

	t, _ := time.Parse(ApiDbLayout, d)
	if hours != 0 {
		t = t.Add(time.Hour * hours)
	}
	return t.Format(UTCDateLayout)
}

func GetCurrentDateWithSameTime(dt string) string {
	t, _ := time.Parse(ApiDbLayout, dt)
	y, m, d := t.Date()
	ct := time.Now().UTC()
	cy, cm, cd := ct.Date()
	if y != cy || m != cm || d != cd {
		t = time.Date(ct.Year(), ct.Month(), ct.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	}
	//fmt.Println("GetCurrentDateWithSameTime", t.Format(ApiDbLayout))
	return t.Format(ApiDbLayout)
}

//get next schedule date, if hour is passed as 0 then no change, -1 would substract 1 hour from given date
func GetNextScheduleDate(dt string, hours time.Duration) string {
	t, _ := time.Parse(ApiDbLayout, dt)
	if hours != 0 {
		t = t.Add(time.Hour * hours)
	}
	fmt.Sprintln("GetNextScheduleDate", t.Format(ApiDbLayout))
	return t.Format(ApiDbLayout)
}

//get next schedule date, if hour is passed as 0 then no change, -1 would substract 1 hour from given date
func GetNextScheduleDateByMins(d string, minutes time.Duration) string {
	t, _ := time.Parse(ApiDbLayout, d)
	if minutes != 0 {
		t = t.Add(time.Minute * minutes)
	}
	return t.Format(ApiDbLayout)
}

func GeApiDBLayoutDateFormat(t time.Time) string {
	return t.Format(ApiDbLayout)
}
