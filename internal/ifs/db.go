package ifs

import (
	"orgStruct/internal/repository"
)

type Db interface {
	SelectDepartmentById(id int) (repository.Department, error)
	SelectDepartmentWhereParentID(ParentID *int) ([]repository.Department, error)
	InsertDepartment(Depart repository.Department)
	InsertEmployee(Emplo repository.Employee)
	SelectEmployeeWhereDepartmentId(DepartmentId int) ([]repository.Employee, error)
	DeleteDepartment(DepartmentId int)
	DeleteEmployee(Id int)
	UpdateEmployee(employee repository.Employee)
	SelectDepartmentByName(Name string) (repository.Department, error)
	UpdateDepartmentParentID(ParentID int, Id int)
}
