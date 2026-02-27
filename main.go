package main

import (
	"log"
	"orgStruct/internal/config"
	"orgStruct/internal/domain"
	"orgStruct/internal/repository"
)

func main() {
	gormDbConection, err := config.NewGormDb()
	if err != nil {
		log.Fatalln(err)
	}
	gormDbRepo := repository.NewGormDbRepo(gormDbConection.Conn)
	OL := domain.NewOrganizationLogic(gormDbRepo)
	i := 1
	err = OL.CreateDivision("test", &i)
	if err != nil {
		log.Fatalln(err)
	}
}
