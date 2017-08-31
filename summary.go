package main

import "time"

// Summary aggregate employee monthly data
type Summary struct {
	ID                    int64     `db:"id"`
	EmployeeID            int64     `db:"employee_id"`
	Year                  int64     `db:"year"`
	Month                 int64     `db:"month"`
	SiteWorkTime          float64   `db:"site_work_time"`
	HouseWorkTime         float64   `db:"house_work_time"`
	MidnightWorkTime      float64   `db:"midnight_work_time"`
	HolidayWorkTime       float64   `db:"holiday_work_time"`
	PaidVacationCount     float64   `db:"paid_vacation_count"`
	VacationCount         float64   `db:"vacation_count"`
	PassPrice             int64     `db:"pass_price"`
	TransportationExpense int64     `db:"transportation_expense"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

// Find find summary
func (s *Summary) Find() error {
	query := `
	SELECT
		id,
		employee_id,
		year,
		month,
		created_at
	FROM
		summaries
	WHERE
		employee_id = ?
	AND
		year = ?
	AND
		month = ?
	`

	return dbm.SelectOne(&s, query, s.EmployeeID, s.Year, s.Month)
}

// Save save insert or update summary
func (s *Summary) Save() error {
	now := time.Now()
	s.UpdatedAt = now
	if s.ID != 0 {
		_, err := dbm.Update(s)

		return err
	}

	s.CreatedAt = now
	return dbm.Insert(s)
}

// SetWorks set values from work data
func (s *Summary) SetWorks(works Works) {
	for _, v := range works {
		s.VacationCount += v.GetVacationCount()
		s.PaidVacationCount += v.GetPaidVacationCount()
		s.MidnightWorkTime += v.GetMidnightWorkTime()
		s.HolidayWorkTime += v.GetHolidayWorkTime()
		s.SiteWorkTime += v.TotalTime
	}
}

// SetHouseWorks set values from house work data
func (s *Summary) SetHouseWorks(works Works) {
	for _, v := range works {
		s.MidnightWorkTime += v.GetMidnightWorkTime()
		s.HolidayWorkTime += v.GetHolidayWorkTime()
		s.HouseWorkTime += v.TotalTime
	}
}
