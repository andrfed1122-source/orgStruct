package ifs

import (
	"orgStruct/internal/repository"
)

type Db interface {
	SelectDepartmentById(id int) (repository.Department, error)
	SelectDepartmentWhereParentID(ParentID *int) ([]repository.Department, error)
	InsertDepartment(Depart repository.Department) repository.Department
	InsertEmployee(Emplo repository.Employee)
	SelectEmployeeWhereDepartmentId(DepartmentId int) ([]repository.Employee, error)
	DeleteDepartment(DepartmentId int)
	DeleteEmployee(Id int)
	UpdateEmployee(employee repository.Employee)
	SelectDepartmentByName(name string, ParentID *int) (repository.Department, error)
	UpdateDepartmentParentID(ParentID int, Id int)
}

type Logic interface {
	CreateDepartmen(name string, ParentID *int) (repository.Department, error)
	CreateEmployee(fullName string, position string, idDepartmen int) error
	InfoDeportament(idDepartmen int, depth *int, includeEmployees *bool) (repository.Department, []repository.Employee, []repository.Department, error)
	InfoChildrenDeportament(idDepartmen *int) ([]repository.Department, error)
	DeleteDeportament(idDepartmen int, mode string, reassignDepartmentId int) error
	UpdateDeportament(idDepartmen int, name string, ParentID int) error
}
