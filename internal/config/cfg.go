package config

import "errors"

var (
	ErrorsNameNotСorrect     = errors.New("Данное имя не подходит")
	ErrorsNameRepeat         = errors.New("Данное имя уже используеться на данном уровне депортаминта")
	ErrorsDepartmentNotExist = errors.New("Департамент не существует")
)
