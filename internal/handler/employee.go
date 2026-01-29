package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"question5updation/internal/model"
	"question5updation/internal/response"
	"question5updation/internal/storage"
	"strconv"
)

func GetEmployees(r *http.Request, _ map[string]string) response.APIResponse {
	employees, err := storage.ReadEmployees()
	if err != nil {
		return response.APIResponse{Status: 500, Err: response.ErrStorage}
	}

	isActiveParam := r.URL.Query().Get("is_active")
	if isActiveParam != "" {
		value, err := strconv.ParseBool(isActiveParam)
		if err != nil {
			return response.APIResponse{Status: 400, Err: errors.New("is_active must be true or false")}
		}

		filtered := []model.Employee{}
		for _, e := range employees {
			if e.IsActive != nil && *e.IsActive == value {
				filtered = append(filtered, e)
			}
		}
		employees = filtered
	}

	return response.APIResponse{Status: 200, Data: employees}
}

func GetEmployeeByName(r *http.Request, params map[string]string) response.APIResponse {
	name := params["name"]
	emp, err := storage.GetEmployeeByName(name)
	if err != nil {
		return response.APIResponse{Status: 500, Err: response.ErrStorage}
	}
	if emp == nil {
		return response.APIResponse{Status: 404, Err: response.ErrEmployeeMissing}
	}
	return response.APIResponse{Status: 200, Data: emp}
}

func CreateEmployee(r *http.Request, _ map[string]string) response.APIResponse {
	var emp model.Employee
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&emp); err != nil {
		return response.APIResponse{Status: 400, Err: response.ErrInvalidJSON}
	}

	if emp.Name == "" || emp.Age <= 0 || emp.Address == "" {
		return response.APIResponse{Status: 400, Err: response.ErrMissingField}
	}

	if emp.IsActive == nil {
		active := true
		emp.IsActive = &active
	}

	existing, err := storage.GetEmployeeByName(emp.Name)
	if err != nil {
		return response.APIResponse{Status: 500, Err: fmt.Errorf("storage error: %w", err)}
	}

	if existing != nil {
		return response.APIResponse{Status: 409, Err: response.ErrEmployeeExists}
	}

	id, err := storage.InsertEmployee(emp)
	if err != nil {
		return response.APIResponse{Status: 500, Err: response.ErrStorage}
	}
	emp.ID = int(id)

	return response.APIResponse{Status: 201, Data: emp}
}

func UpdateEmployee(r *http.Request, params map[string]string) response.APIResponse {
	name := params["name"]
	var emp model.Employee
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&emp); err != nil {
		return response.APIResponse{Status: 400, Err: response.ErrInvalidJSON}
	}

	if emp.Age <= 0 || emp.Address == "" || emp.IsActive == nil {
		return response.APIResponse{Status: 400, Err: response.ErrMissingField}
	}

	err := storage.UpdateEmployeeByName(name, emp)
	if err != nil {
		if err.Error() == "not found" {
			return response.APIResponse{Status: 404, Err: response.ErrEmployeeMissing}
		}
		return response.APIResponse{Status: 500, Err: response.ErrStorage}
	}

	return response.APIResponse{Status: 200, Data: emp}
}

func DeleteEmployee(r *http.Request, params map[string]string) response.APIResponse {
	name := params["name"]
	err := storage.DeleteEmployeeByName(name)
	if err != nil {
		if err.Error() == "not found" {
			return response.APIResponse{Status: 404, Err: response.ErrEmployeeMissing}
		}
		return response.APIResponse{Status: 500, Err: response.ErrStorage}
	}
	return response.APIResponse{Status: 200}
}
