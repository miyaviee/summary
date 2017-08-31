package main

import (
	"database/sql"
	"os"
	"sync"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var dbm *gorp.DbMap

func main() {
	initDb()
	defer dbm.Db.Close()

	date, err := getTargetDate()
	if err != nil {
		panic(err)
	}

	year := int64(date.Year())
	month := int64(date.Month())

	wg := sync.WaitGroup{}
	semaphore := make(chan int, 8)
	employees := FindEmployees()
	for _, employeeID := range employees {
		wg.Add(1)
		go func(employeeID, year, month int64) {
			defer wg.Done()
			semaphore <- 1
			s := Summary{EmployeeID: employeeID, Year: year, Month: month}
			s.Find()

			params := []interface{}{employeeID, year, month}

			works := FindSiteWorks(params...)
			s.SetWorks(works)

			houseWorks := FindHouseWorks(params...)
			s.SetHouseWorks(houseWorks)

			t := FindTransportationExpense(params...)
			s.TransportationExpense += t.Cost

			p := FindPassPrice(params...)
			s.PassPrice += p.Cost

			err := s.Save()
			if err != nil {
				panic(err)
			}
			<-semaphore
		}(employeeID, year, month)
	}
	wg.Wait()
}

func getTargetDate() (time.Time, error) {
	if len(os.Args) != 2 {
		return time.Now(), nil
	}

	return time.ParseInLocation("200601", os.Args[1], time.Local)
}

func initDb() {
	db, err := sql.Open("mysql", os.Getenv("DB_INFO")+"?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		panic(err)
	}

	dbm = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbm.AddTableWithName(Summary{}, "summaries").SetKeys(true, "ID")
}
