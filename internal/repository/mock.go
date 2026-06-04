package repository

import (
	"errors"
	"sync"
	"time"
)

type MockDB struct {
	mu          sync.RWMutex
	departments map[int]*Department
	employees   map[int]*Employee
	depCounter  int
	empCounter  int
}

func NewMockDB() *MockDB {
	return &MockDB{
		departments: make(map[int]*Department),
		employees:   make(map[int]*Employee),
		depCounter:  1,
		empCounter:  1,
	}
}

func (m *MockDB) SelectDepartmentById(id int) (Department, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	dep, exists := m.departments[id]
	if !exists {
		return Department{}, errors.New("department not found")
	}
	return *dep, nil
}

func (m *MockDB) SelectDepartmentWhereParentID(ParentID *int) ([]Department, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []Department
	for _, dep := range m.departments {
		if (ParentID == nil && dep.ParentID == nil) || (ParentID != nil && dep.ParentID != nil && *dep.ParentID == *ParentID) {
			result = append(result, *dep)
		}
	}
	return result, nil
}

func (m *MockDB) InsertDepartment(dep Department) Department {
	m.mu.Lock()
	defer m.mu.Unlock()

	dep.ID = m.depCounter
	m.depCounter++
	dep.CreatedAt = time.Now()
	m.departments[dep.ID] = &dep
	return dep
}

func (m *MockDB) InsertEmployee(emp Employee) {
	m.mu.Lock()
	defer m.mu.Unlock()

	emp.ID = m.empCounter
	m.empCounter++
	emp.CreatedAt = time.Now()
	m.employees[emp.ID] = &emp
}

func (m *MockDB) SelectEmployeeWhereDepartmentId(DepartmentId int) ([]Employee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []Employee
	for _, emp := range m.employees {
		if emp.DepartmentID == DepartmentId {
			result = append(result, *emp)
		}
	}
	return result, nil
}

func (m *MockDB) DeleteDepartment(DepartmentId int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.departments, DepartmentId)
}

func (m *MockDB) DeleteEmployee(Id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.employees, Id)
}

func (m *MockDB) UpdateEmployee(employee Employee) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.employees[employee.ID] = &employee
}

func (m *MockDB) SelectDepartmentByName(name string, ParentID *int) (Department, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, dep := range m.departments {
		if dep.Name == name {
			if (ParentID == nil && dep.ParentID == nil) || (ParentID != nil && dep.ParentID != nil && *dep.ParentID == *ParentID) {
				return *dep, nil
			}
		}
	}
	return Department{}, errors.New("department not found")
}

func (m *MockDB) UpdateDepartmentParentID(ParentID int, Id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if dep, exists := m.departments[Id]; exists {
		dep.ParentID = &ParentID
	}
}
