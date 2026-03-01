package config

import "errors"

var (
	ErrorsNameNotСorrect       = errors.New("Данное имя не подходит")
	ErrorsNameRepeat           = errors.New("Данное имя уже используеться на данном уровне депортаминта")
	ErrorsDepartmentNotExist   = errors.New("Департамент не существует")
	ErrorsModeNotСorrect       = errors.New("Данный мод не коректный")
	ErrorsReassignDepartmentId = errors.New("ReassignDepartmentId меньше нуля")
	ErrorsLoopingDepartment    = errors.New("Вы пытаетесь сделать депортамент дочерним от его дочерних депортаментов")
)
