package storage

import (
	"database/sql"
	"fmt"
	"question5updation/internal/model"
)

func ReadEmployees() ([]model.Employee, error) {
	rows, err := DB.Query("SELECT id, name, age, address, is_active FROM srihari_employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []model.Employee
	for rows.Next() {
		var e model.Employee
		err := rows.Scan(&e.ID, &e.Name, &e.Age, &e.Address, &e.IsActive)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func GetEmployeeByName(name string) (*model.Employee, error) {
	var e model.Employee
	err := DB.QueryRow("SELECT id, name, age, address, is_active FROM srihari_employees WHERE name=?", name).
		Scan(&e.ID, &e.Name, &e.Age, &e.Address, &e.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func InsertEmployee(emp model.Employee) (int64, error) {

	result, err := DB.Exec(
		"INSERT INTO srihari_employees (name, age, address, is_active) VALUES (?, ?, ?, ?)",
		emp.Name, emp.Age, emp.Address, emp.IsActive,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateEmployeeByName(name string, emp model.Employee) error {
	res, err := DB.Exec("UPDATE srihari_employees SET age=?, address=?, is_active=? WHERE name=?",
		emp.Age, emp.Address, *emp.IsActive, name)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func DeleteEmployeeByName(name string) error {
	res, err := DB.Exec("DELETE FROM srihari_employees WHERE name=?", name)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}
