package repository

import "time"

type Department struct {
	ID        int       `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;size:200;not null"`
	ParentID  *int      `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`

	// Связи
	Parent    *Department  `gorm:"foreignKey:ParentID;references:ID"`     // Родительский отдел
	Children  []Department `gorm:"foreignKey:ParentID;references:ID"`     // Дочерние отделы
	Employees []Employee   `gorm:"foreignKey:DepartmentID;references:ID"` // Сотрудники отдела
}

// Employee модель для таблицы employee
type Employee struct {
	ID           int        `gorm:"primaryKey;column:id"`
	DepartmentID int        `gorm:"column:parent_id;not null"`
	FullName     string     `gorm:"column:full_name;size:200;not null"`
	Position     string     `gorm:"column:position;size:200;not null"`
	HiredAt      *time.Time `gorm:"column:hired_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`

	// Связи
	Department Department `gorm:"foreignKey:DepartmentID;references:ID"`
}

func (Department) TableName() string {
	return "department"
}

// TableName задает имя таблицы для модели Employee
func (Employee) TableName() string {
	return "employee"
}
