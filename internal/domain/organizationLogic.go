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
func (logic *organizationLogic) CreateDepartmen(name string, ParentID *int) (repository.Department, error) {
	// Убераем пробелы по краям имени и проверяем имя на соблюдение правил
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 200 {
		return repository.Department{}, config.ErrorsNameNotСorrect
	}
	//проверка корневого депортамента на наличие соседий с такимже именем
	if ParentID == nil {
		levelDepartment, err := logic.db.SelectDepartmentWhereParentID(nil)
		if err != nil {
			return repository.Department{}, err
		}
		for _, department := range levelDepartment {
			if department.Name == name {
				return repository.Department{}, config.ErrorsNameRepeat
			}
		}
	} else { //запрос верхнего депортамента, всех его дочерних депортаментов и сравнение имени с новым
		ParentDepartment, err := logic.db.SelectDepartmentById(*ParentID)
		if err != nil {
			return repository.Department{}, err
		}
		if ParentDepartment.Name == "" {
			return repository.Department{}, config.ErrorsDepartmentNotExist
		}
		levelDepartment, err := logic.db.SelectDepartmentWhereParentID(&ParentDepartment.ID)
		if err != nil {
			return repository.Department{}, err
		}
		for _, department := range levelDepartment {
			if department.Name == name {
				return repository.Department{}, config.ErrorsNameRepeat
			}
		}
	}
	//если все условия соблюдены
	newDep := logic.db.InsertDepartment(repository.Department{ParentID: ParentID, Name: name, CreatedAt: time.Now()})
	return newDep, nil
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
	Departmen.Children, err = logic.db.SelectDepartmentWhereParentID(&idDepartmen)
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
			midlChildren, err = logic.db.SelectDepartmentWhereParentID(&departmentC.ID)
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

func (logic *organizationLogic) DeleteDeportament(idDepartmen int, mode string, reassignDepartmentId int) error {
	if mode != "cascade" && mode != "reassign" {
		return config.ErrorsModeNotСorrect
	}
	if mode == "reassign" && reassignDepartmentId < 0 {
		return config.ErrorsReassignDepartmentId
	}
	//проверяем существует ли данное подрозделение
	Departmen, err := logic.db.SelectDepartmentById(idDepartmen)
	if err != nil {
		return err
	}
	if Departmen.Name == "" {
		return config.ErrorsDepartmentNotExist
	}

	//сбор всех дочерних пердпреятий
	Departmen.Children, err = logic.db.SelectDepartmentWhereParentID(&idDepartmen)
	if err != nil {
		return err
	}
	chekChildren := Departmen.Children
	dep := 1
	for {
		if dep == 0 {
			break
		}
		dep = 0
		var midlChildren []repository.Department
		for _, departmentC := range chekChildren {
			midlChildren, err = logic.db.SelectDepartmentWhereParentID(&departmentC.ID)
			if err != nil {
				return err
			}
			clear(chekChildren)
			for _, midlChild := range midlChildren {
				dep++
				Departmen.Children = append(Departmen.Children, midlChild)
				chekChildren = append(chekChildren, midlChild)

			}
		}
	}
	//собераем всех сотрудников
	var ChildrenEmployee = []repository.Employee{}
	Departmen.Employees, err = logic.db.SelectEmployeeWhereDepartmentId(idDepartmen)
	if err != nil {
		return err
	}
	for _, department := range Departmen.Children {
		ChildrenEmployee, err = logic.db.SelectEmployeeWhereDepartmentId(department.ID)
		for _, CE := range ChildrenEmployee {
			Departmen.Employees = append(Departmen.Employees, CE)
		}
	}

	if mode == "cascade" {
		//удаляем всех
		for _, DE := range Departmen.Employees {
			logic.db.DeleteEmployee(DE.ID)
		}
		for _, DC := range Departmen.Children {
			logic.db.DeleteDepartment(DC.ID)
		}
	}
	if mode == "reassign" {
		for _, DE := range Departmen.Employees {
			DE.DepartmentID = reassignDepartmentId
			logic.db.UpdateEmployee(DE)
		}
		for _, DC := range Departmen.Children {
			logic.db.DeleteDepartment(DC.ID)
		}
	}
	return nil
}

func (logic *organizationLogic) UpdateDeportament(idDepartmen int, name string, ParentID int) error {
	//if idDepartmen < 1 && name == "" {
	//	return config.ErrorsNameNotСorrect
	//}
	//var Departmen repository.Department
	//var err error
	//if name != "" {
	//	Departmen, err = logic.db.SelectDepartmentById(idDepartmen)
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	Departmen, err = logic.db.SelectDepartmentByName(name)
	//	if err != nil {
	//		return err
	//	}
	//}
	////проверка на создание кольца депортаментов
	////собираем депортаменты дочернии первого порядка
	//Departmen.Children, err = logic.db.SelectDepartmentWhereParentID(&idDepartmen)
	//if err != nil {
	//	return err
	//}
	//chekChildren := Departmen.Children
	//dep := 1
	//for {
	//	if dep == 0 {
	//		break
	//	}
	//	dep = 0
	//	var midlChildren []repository.Department
	//	for _, departmentC := range chekChildren {
	//		midlChildren, err = logic.db.SelectDepartmentWhereParentID(&departmentC.ID)
	//		if err != nil {
	//			return err
	//		}
	//		clear(chekChildren)
	//		for _, midlChild := range midlChildren {
	//			dep++
	//			Departmen.Children = append(Departmen.Children, midlChild)
	//			chekChildren = append(chekChildren, midlChild)
	//
	//		}
	//	}
	//}
	//for _, DC := range Departmen.Children {
	//	if DC.ID == idDepartmen || DC.Name == name {
	//		return config.ErrorsLoopingDepartment
	//	}
	//}
	//logic.db.UpdateDepartmentParentID(ParentID, Departmen.ID)
	return nil
}
