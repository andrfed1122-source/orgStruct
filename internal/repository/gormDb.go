package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type gormDbRepo struct {
	db *gorm.DB
}

func NewGormDbRepo(db *gorm.DB) *gormDbRepo {
	return &gormDbRepo{db: db}
}

// получение всех депортаминтов
func (repo *gormDbRepo) SelectDepartment() error {
	var departments []Department
	err := repo.db.Find(&departments)
	if err.Error != nil {
		return err.Error
	}
	fmt.Println(departments)
	return nil
}

// получение всех сотрудников
func (repo *gormDbRepo) SelectEmployee() error {
	var employee []Employee
	err := repo.db.Find(&employee)
	if err.Error != nil {
		return err.Error
	}
	fmt.Println(employee)
	return nil
}

// получение конкретного дипортамента
func (repo *gormDbRepo) SelectDepartmentById(id int) (Department, error) {
	var department Department
	err := repo.db.Where("ID = ?", id).Find(&department)
	if err.Error != nil {
		return Department{}, err.Error
	}
	return department, nil
}

func (repo *gormDbRepo) SelectDepartmentByName(name string, ParentID *int) (Department, error) {
	var department Department
	err := repo.db.Where("Name = ? AND ParentID = ?", name, ParentID).Find(&department)
	if err.Error != nil {
		return Department{}, err.Error
	}
	return department, nil
}

// получение всех дочерних дипортаментов
func (repo *gormDbRepo) SelectDepartmentWhereParentID(ParentID *int) ([]Department, error) {
	var departments []Department
	if ParentID == nil {
		err := repo.db.Where("parent_id is null").Find(&departments)
		if err.Error != nil {
			return nil, err.Error
		}
	} else {
		err := repo.db.Where("parent_id = ?", *ParentID).Find(&departments)
		if err.Error != nil {
			return nil, err.Error
		}
	}
	return departments, nil
}

func (repo *gormDbRepo) InsertDepartment(Depart Department) Department {
	//var maxID int
	//repo.db.Raw("SELECT MAX(id) FROM department").Scan(&maxID)
	//Depart.ID = maxID
	repo.db.Create(&Depart)
	return Depart
}
func (repo *gormDbRepo) InsertEmployee(Emplo Employee) {
	repo.db.Create(&Emplo)
}

func (repo *gormDbRepo) SelectEmployeeWhereDepartmentId(DepartmentId int) ([]Employee, error) {
	var employees []Employee
	err := repo.db.Where("parent_id = ?", DepartmentId).Find(&employees)
	if err.Error != nil {
		return nil, err.Error
	}
	return employees, nil
}

func (repo *gormDbRepo) DeleteDepartment(DepartmentId int) {
	repo.db.Where("ID = ?", DepartmentId).Delete(&Department{})
}

func (repo *gormDbRepo) DeleteEmployee(Id int) {
	repo.db.Where("ID = ?", Id).Delete(&Employee{})
}

func (repo *gormDbRepo) DeleteEmployeeWhereIdPe(Id int) {
	repo.db.Where("ID = ?", Id).Delete(&Employee{})
}

func (repo *gormDbRepo) UpdateEmployee(employee Employee) {
	repo.db.Model(employee).Where("ID = ?", employee.ID).Updates(&employee)
}

func (repo *gormDbRepo) UpdateDepartmentParentID(ParentID int, Id int) {
	repo.db.Model(Department{}).Where("ID = ?", Id).Update("parent_id", ParentID)
}
