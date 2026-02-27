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
	err := repo.db.Find(&department, id)
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

func (repo *gormDbRepo) InsertDepartment(Depart Department) {
	repo.db.Create(&Depart)
}
