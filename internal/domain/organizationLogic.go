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
func (logic *organizationLogic) CreateDepartmen(name string, ParentID *int) error {
	// Убераем пробелы по краям имени и проверяем имя на соблюдение правил
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 200 {
		return config.ErrorsNameNotСorrect
	}
	//проверка корневого депортамента на наличие соседий с такимже именем
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
	} else { //запрос верхнего депортамента, всех его дочерних депортаментов и сравнение имени с новым
		ParentDepartment, err := logic.db.SelectDepartmentById(*ParentID)
		if err != nil {
			return err
		}
		if ParentDepartment.Name == "" {
			return config.ErrorsDepartmentNotExist
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
	//если все условия соблюдены
	logic.db.InsertDepartment(repository.Department{ParentID: ParentID, Name: name, CreatedAt: time.Now()})
	return nil
}

// создание сотрудника
func (logic *organizationLogic) CreateEmployee(fullName string, position string, idDepartmen int) error {
	//проверяем имя и должность на коректность
	fullName = strings.TrimSpace(fullName)
	position = strings.TrimSpace(position)
	if fullName == "" || len(fullName) > 200 || position == "" || len(position) > 200 {
		return config.ErrorsNameNotСorrect
	}
	//проверяем существует ли данное подрозделение
	Departmen, err := logic.db.SelectDepartmentById(idDepartmen)
	if err != nil {
		return err
	}
	if Departmen.Name == "" {
		return config.ErrorsDepartmentNotExist
	}
	//добовляем сотрудника
	logic.db.InsertEmployee(repository.Employee{FullName: fullName, Position: position, CreatedAt: time.Now(), DepartmentID: idDepartmen})
	return nil
}

func (logic *organizationLogic) InfoDeportament(idDepartmen int, depth *int, includeEmployees *bool) (repository.Department, []repository.Employee, []repository.Department, error) {
	//проверяем переменные на не допустимые значения, и при конфликте оставляем константы
	incEmployees := true
	dep := 1
	if includeEmployees != nil {
		incEmployees = *includeEmployees
	}
	if depth != nil && *depth > 0 && *depth < 6 {
		dep = *depth
	}
	//проверяем существует ли данное подрозделение
	Departmen, err := logic.db.SelectDepartmentById(idDepartmen)
	if err != nil {
		return repository.Department{}, nil, nil, err
	}
	if Departmen.Name == "" {
		return repository.Department{}, nil, nil, config.ErrorsDepartmentNotExist
	}
	Departmen.Employees, err = logic.db.SelectEmployeeWhereDepartmentId(idDepartmen)
	if err != nil {
		return repository.Department{}, nil, nil, err
	}
	Departmen.Children, err = logic.InfoChildrenDeportament(&idDepartmen)
	if err != nil {
		return repository.Department{}, nil, nil, err
	}
	chekChildren := Departmen.Children
	for {
		if dep < 2 {
			break
		}
		var midlChildren []repository.Department
		for _, departmentC := range chekChildren {
			midlChildren, err = logic.InfoChildrenDeportament(&departmentC.ID)
			if err != nil {
				return repository.Department{}, nil, nil, err
			}
			clear(chekChildren)
			for _, midlChild := range midlChildren {
				Departmen.Children = append(Departmen.Children, midlChild)
				chekChildren = append(chekChildren, midlChild)

			}
		}
		dep--
	}
	if incEmployees {
		return Departmen, Departmen.Employees, Departmen.Children, nil
	}
	return Departmen, nil, Departmen.Children, nil
}

func (logic *organizationLogic) InfoChildrenDeportament(idDepartmen *int) ([]repository.Department, error) {
	children, err := logic.db.SelectDepartmentWhereParentID(idDepartmen)
	if err != nil {
		return nil, err
	}
	return children, nil
}
