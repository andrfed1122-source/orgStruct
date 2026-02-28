package config

import (
	//"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"gorm.io/gorm/clause"
	"time"
)

type GormDb struct {
	Conn *gorm.DB
}

func NewGormDb() *GormDb {
	// Строка подключения к PostgreSQL
	dsn := "host=127.0.0.1 port=5432 user=postgres password=12345 dbname=postgres sslmode=disable TimeZone=UTC"

	// Открываем соединение
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true, // Подготовка запросов для производительности
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Получаем объект SQL DB для настройки пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database object")
	}

	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &GormDb{Conn: db}
}
