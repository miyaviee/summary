package main

// Employees employee id list
type Employees []int64

// FindEmployees find normal employee id list
func FindEmployees() Employees {
	query := `
	SELECT
		id
	FROM
		employees
	WHERE
		department_id != 1
	`

	var employees Employees
	if _, err := dbm.Select(&employees, query); err != nil {
		panic(err)
	}

	return employees
}
