package main

import "fmt"

// Expense expsnse cost
type Expense struct {
	Cost int64 `db:"cost"`
}

// FindTransportationExpense find transportation_expenses total cost
func FindTransportationExpense(params ...interface{}) Expense {
	var expense Expense
	dbm.SelectOne(&expense, getQueryFindExpenses("transportation_expenses"), params...)

	return expense
}

// FindPassPrice find pass_prices total cost
func FindPassPrice(params ...interface{}) Expense {
	var expense Expense
	dbm.SelectOne(&expense, getQueryFindExpenses("pass_prices"), params...)

	return expense
}

func getQueryFindExpenses(tableName string) string {
	query := `
	SELECT
		SUM(cost) AS cost
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
