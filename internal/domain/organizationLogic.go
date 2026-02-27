package domain

import (
	"orgStruct/internal/config"
	"orgStruct/internal/ifs"
	"orgStruct/internal/repository"
	"strings"
	"time"
)

type organizationLogic struct {
	db ifs.Db
}

func NewOrganizationLogic(db ifs.Db) *organizationLogic {
	return &organizationLogic{db: db}
}

// Создание Депортамента
func (logic *organizationLogic) CreateDivision(name string, ParentID *int) error {
	// Убераем пробелы по краям имени и проверяем имя на соблюдение правил
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 200 {
		return config.ErrorsNameNotСorrect
	}
	if ParentID == nil {
		levelDepartment, err := logic.db.SelectDepartmentWhereParentID(nil)
		if err != nil {
			return err
		}
		for _, department := range levelDepartment {
			if department.Name == name {
				return config.ErrorsNameRepeat
			}
		}
	} else {
		ParentDepartment, err := logic.db.SelectDepartmentById(*ParentID)
		if err != nil {
			return err
		}

		levelDepartment, err := logic.db.SelectDepartmentWhereParentID(&ParentDepartment.ID)
		if err != nil {
			return err
		}
		for _, department := range levelDepartment {
			if department.Name == name {
				return config.ErrorsNameRepeat
			}
		}
	}
	logic.db.InsertDepartment(repository.Department{ParentID: ParentID, Name: name, CreatedAt: time.Now()})
	return nil
}
