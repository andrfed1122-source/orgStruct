package handler

import (
	"encoding/json"
	"net/http"
	"orgStruct/internal/config"
	"orgStruct/internal/domain"
	"orgStruct/internal/ifs"
	"orgStruct/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type CreateEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at,omitempty"`
}

// CreateDepartmentRequest структура запроса для создания подразделения
type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id,omitempty"`
}
type CreateDepartmentResponse struct {
	Department repository.Department `json:"department"`
}
type DepartmentHandler struct {
	logic ifs.Logic
}

// NewDepartmentHandler создает новый экземпляр обработчика
func NewDepartmentHandler(logic ifs.Logic) *DepartmentHandler {
	return &DepartmentHandler{logic: logic}
}

// CreateDepartment обработчик для создания подразделения
func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON из тела запроса
	var req CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем функцию создания подразделения
	// Предполагается, что функция CreateDepartmen определена где-то в коде
	dep, err := h.logic.CreateDepartmen(req.Name, req.ParentID)
	if err != nil {
		http.Error(w, "Failed to create department: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем ответ с данными созданного подразделения
	// В реальном приложении здесь нужно получить ID созданного подразделения
	// из базы данных или из результата функции CreateDepartmen
	response := CreateDepartmentResponse{
		Department: repository.Department{
			ID:        dep.ID,
			Name:      dep.Name,
			ParentID:  dep.ParentID,
			CreatedAt: dep.CreatedAt,
		},
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Отправляем ответ
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DepartmentHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID подразделения из URL
	vars := mux.Vars(r)
	departmentID, err := strconv.Atoi(vars["id"])
	if err != nil || departmentID <= 0 {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из тела запроса
	var req CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if req.FullName == "" {
		http.Error(w, "full_name is required", http.StatusBadRequest)
		return
	}
	if req.Position == "" {
		http.Error(w, "position is required", http.StatusBadRequest)
		return
	}

	// Проверяем существование подразделения (опционально)
	exists, err := h.checkDepartmentExists(departmentID)
	if err != nil {
		http.Error(w, "Failed to check department: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	// Вызываем существующую функцию CreateEmployee
	err = CreateEmployee(req.FullName, req.Position, departmentID)
	if err != nil {
		http.Error(w, "Failed to create employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем созданного сотрудника (в реальном приложении нужно получить из БД)
	// Здесь для примера создаем объект с заглушками
	now := time.Now()
	employee := Employee{
		ID:           1, // Здесь должен быть реальный ID из БД
		FullName:     req.FullName,
		Position:     req.Position,
		DepartmentID: departmentID,
		HiredAt:      req.HiredAt,
		CreatedAt:    now,
	}

	// Создаем ответ
	response := CreateEmployeeResponse{
		Employee: employee,
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Отправляем ответ
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// SetupRoutes настраивает маршруты для API
func SetupRoutes(handler *DepartmentHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Регистрируем обработчик для создания подразделения
	mux.HandleFunc("/departments/", handler.CreateDepartment)
	mux.HandleFunc("/departments/{id}/employees/", handler.CreateEmployee).Methods("POST")

	return mux
}

// Пример использования:
func HttpServerStart() {
	cfg := config.NewGormDb()
	gormDbRepo := repository.NewGormDbRepo(cfg.Conn)
	OL := domain.NewOrganizationLogic(gormDbRepo)
	handler := NewDepartmentHandler(OL)

	mux := SetupRoutes(handler)

	// Запускаем сервер на порту 8080
	http.ListenAndServe(":8080", mux)
}
