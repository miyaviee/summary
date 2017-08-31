package main

import (
	"fmt"
	"time"
)

// Works work list
type Works []Work

// Work work
type Work struct {
	Year        int64   `db:"year"`
	Month       int64   `db:"month"`
	Day         int64   `db:"day"`
	HolidayType int64   `db:"holiday_type"`
	StartTime   float64 `db:"start_time"`
	EndTime     float64 `db:"end_time"`
	BreakTime   float64 `db:"break_time"`
	TotalTime   float64 `db:"total_time"`
	WorkType    int64   `db:"work_type"`
}

// GetVacationCount get vacation count
func (w *Work) GetVacationCount() float64 {
	if w.WorkType == 2 {
		return 1
	}

	return 0
}

// GetPaidVacationCount get paid vacation count
func (w *Work) GetPaidVacationCount() float64 {
	if w.WorkType == 4 {
		return 1
	}

	if w.WorkType == 5 || w.WorkType == 6 {
		return 0.5
	}

	return 0
}

// GetMidnightWorkTime get midnight work time
func (w *Work) GetMidnightWorkTime() float64 {
	if 22 < w.StartTime {
		return w.TotalTime
	}

	if 22 < w.EndTime {
		return w.EndTime - 22
	}

	return 0
}

// GetHolidayWorkTime get holiday work time
func (w *Work) GetHolidayWorkTime() float64 {
	tomorrowWorkTime := w.getTomorrowWorkTime()

	holidayWorkTime := 0.0
	date := fmt.Sprintf("%04d%02d%02d", w.Year, w.Month, w.Day)
	t, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	t = t.AddDate(0, 0, 1)
	if t.Weekday() == 0 || t.Weekday() == 6 {
		holidayWorkTime += tomorrowWorkTime
	}

	if w.HolidayType == 0 {
		return holidayWorkTime
	}

	holidayWorkTime += w.TotalTime - tomorrowWorkTime

	return holidayWorkTime
}

func (w *Work) getTomorrowWorkTime() float64 {
	if 24 < w.StartTime {
		return w.TotalTime
	}

	if 24 < w.EndTime {
		return w.EndTime - 24
	}

	return 0
}

// FindSiteWorks find site work
func FindSiteWorks(params ...interface{}) Works {
	var works Works
	dbm.Select(&works, getQueryFindWorks("works"), params...)

	return works
}

// FindHouseWorks find house work
func FindHouseWorks(params ...interface{}) Works {
	var works Works
	dbm.Select(&works, getQueryFindWorks("house_works"), params...)

	return works
}

func getQueryFindWorks(tableName string) string {
	query := `
	SELECT
		year,
		month,
		day,
		holiday_type,
		IFNULL(start_time, 0) as start_time,
		IFNULL(end_time, 0) as end_time,
		IFNULL(break_time, 0) as break_time,
		IFNULL(total_time, 0) as total_time,
		work_type
	FROM
		%s
	WHERE
		employee_id = ?
	AND
		year = ?
	AND
		month = ?
	`

	return fmt.Sprintf(query, tableName)
}
