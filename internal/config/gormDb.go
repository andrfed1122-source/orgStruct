package config

import (
	"gorm.io/gorm"
	"orgStruct/internal/repository"
)

type GormDb struct {
	Conn   *gorm.DB
	MockDB *repository.MockDB
}

func NewGormDb() *GormDb {
	// Используем in-memory mock БД для разработки (без необходимости CGO и GCC)
	mockDB := repository.NewMockDB()
	
	return &GormDb{
		Conn:   &gorm.DB{},
		MockDB: mockDB,
	}
}
