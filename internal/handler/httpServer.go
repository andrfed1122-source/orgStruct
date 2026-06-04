package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orgStruct/internal/config"
	"orgStruct/internal/domain"
	"orgStruct/internal/ifs"
	"orgStruct/internal/repository"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CreateEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at,omitempty"`
}

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id,omitempty"`
}

type CreateDepartmentResponse struct {
	Department repository.Department `json:"department"`
}

type CreateEmployeeResponse struct {
	Employee Employee `json:"employee"`
}

type Employee struct {
	ID           int        `json:"id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	DepartmentID int        `json:"department_id"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type DepartmentHandler struct {
	logic ifs.Logic
}

func NewDepartmentHandler(logic ifs.Logic) *DepartmentHandler {
	return &DepartmentHandler{logic: logic}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v\n", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Creating department: %s\n", req.Name)

	dep, err := h.logic.CreateDepartmen(req.Name, req.ParentID)
	if err != nil {
		log.Printf("Failed to create department: %v\n", err)
		http.Error(w, "Failed to create department: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Department created successfully: ID=%d, Name=%s\n", dep.ID, dep.Name)

	response := CreateDepartmentResponse{
		Department: repository.Department{
			ID:        dep.ID,
			Name:      dep.Name,
			ParentID:  dep.ParentID,
			CreatedAt: dep.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DepartmentHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	departmentID, err := strconv.Atoi(vars["id"])
	if err != nil || departmentID <= 0 {
		log.Printf("Invalid department ID: %v\n", err)
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var req CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v\n", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.FullName == "" {
		log.Println("full_name is required")
		http.Error(w, "full_name is required", http.StatusBadRequest)
		return
	}
	if req.Position == "" {
		log.Println("position is required")
		http.Error(w, "position is required", http.StatusBadRequest)
		return
	}

	log.Printf("Creating employee: %s in department %d\n", req.FullName, departmentID)

	err = h.logic.CreateEmployee(req.FullName, req.Position, departmentID)
	if err != nil {
		log.Printf("Failed to create employee: %v\n", err)
		http.Error(w, "Failed to create employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Employee created successfully: %s\n", req.FullName)

	now := time.Now()
	response := CreateEmployeeResponse{
		Employee: Employee{
			ID:           0,
			FullName:     req.FullName,
			Position:     req.Position,
			DepartmentID: departmentID,
			HiredAt:      req.HiredAt,
			CreatedAt:    now,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DepartmentHandler) GetDepartmentEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	departmentID, err := strconv.Atoi(vars["id"])
	if err != nil || departmentID <= 0 {
		log.Printf("Invalid department ID: %v\n", err)
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching employees for department: ID=%d\n", departmentID)

	dept, employees, _, err := h.logic.InfoDeportament(departmentID, nil, nil)
	if err != nil {
		log.Printf("Failed to fetch department info: %v\n", err)
		http.Error(w, "Failed to fetch department: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if dept.Name == "" {
		log.Printf("Department not found: ID=%d\n", departmentID)
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"department": dept,
		"employees":  employees,
	})
}

func SetupRoutes(handler *DepartmentHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/departments/", handler.CreateDepartment).Methods("POST")
	router.HandleFunc("/departments/{id}/employees/", handler.CreateEmployee).Methods("POST")
	router.HandleFunc("/departments/{id}/employees/", handler.GetDepartmentEmployees).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := getHTML()
		fmt.Fprint(w, html)
	}).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"Endpoint not found"}`)
	})

	return router
}

func getHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Org Structure</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            color-scheme: dark;
            --bg: #050505;
            --panel: rgba(255, 255, 255, 0.02);
            --accent: #ffffff;
            --text: #f3f3f3;
            --muted: rgba(255, 255, 255, 0.56);
            --danger: #ff5a5a;
            --border: rgba(255, 255, 255, 0.1);
            --border-soft: rgba(255, 255, 255, 0.05);
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        html, body {
            width: 100%;
            height: 100%;
        }

        body {
            background: radial-gradient(circle at top, rgba(255, 255, 255, 0.03), transparent 28%), linear-gradient(180deg, #070707 0%, #040404 100%);
            color: var(--text);
            font-family: 'JetBrains Mono', monospace;
            font-size: 13px;
            letter-spacing: -0.01em;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
        }

        button, input, textarea, select {
            font: inherit;
            color: inherit;
        }

        button {
            outline: none;
            cursor: pointer;
        }

        input, textarea {
            outline: none;
        }

        .container {
            display: flex;
            flex-direction: column;
            width: 100%;
            height: 100%;
            padding: 20px;
        }

        .header {
            margin-bottom: 24px;
        }

        .header-title {
            font-size: 18px;
            font-weight: 700;
            letter-spacing: 0.05em;
            text-transform: uppercase;
            margin-bottom: 6px;
        }

        .header-subtitle {
            color: var(--muted);
            font-size: 12px;
        }

        .content {
            flex: 1;
            overflow-y: auto;
            display: grid;
            grid-template-columns: 1fr 1fr 1fr;
            gap: 16px;
            margin-bottom: 16px;
        }

        .card {
            background: var(--panel);
            border: 1px solid var(--border);
            border-radius: 6px;
            padding: 16px;
            display: flex;
            flex-direction: column;
            gap: 12px;
        }

        .card-title {
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.1em;
            border-bottom: 1px solid var(--border-soft);
            padding-bottom: 8px;
        }

        .form-group {
            display: flex;
            flex-direction: column;
            gap: 4px;
        }

        .form-label {
            font-size: 11px;
            color: var(--muted);
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .form-input {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid var(--border-soft);
            border-radius: 4px;
            padding: 8px 10px;
            font-size: 12px;
            transition: all 0.2s ease;
        }

        .form-input:focus {
            background: rgba(255, 255, 255, 0.04);
            border-color: var(--accent);
        }

        .form-input::placeholder {
            color: rgba(255, 255, 255, 0.3);
        }

        .btn {
            background: rgba(255, 255, 255, 0.08);
            border: 1px solid var(--border);
            border-radius: 4px;
            padding: 10px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.1em;
            transition: all 0.2s ease;
            margin-top: 4px;
        }

        .btn:hover {
            background: rgba(255, 255, 255, 0.12);
            border-color: var(--accent);
            transform: translateY(-1px);
        }

        .btn:active {
            transform: translateY(0);
        }

        .feedback {
            font-size: 11px;
            padding: 8px;
            border-radius: 4px;
            display: none;
            margin-top: 8px;
        }

        .success {
            background: rgba(76, 175, 80, 0.15);
            border: 1px solid rgba(76, 175, 80, 0.3);
            color: #4caf50;
        }

        .error {
            background: rgba(255, 90, 90, 0.15);
            border: 1px solid rgba(255, 90, 90, 0.3);
            color: #ff5a5a;
        }

        .endpoints {
            background: var(--panel);
            border: 1px solid var(--border);
            border-radius: 6px;
            padding: 16px;
        }

        .endpoints-title {
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.1em;
            margin-bottom: 12px;
            border-bottom: 1px solid var(--border-soft);
            padding-bottom: 8px;
        }

        .endpoint-item {
            padding: 10px 0;
            border-bottom: 1px solid var(--border-soft);
            font-size: 11px;
            display: flex;
            gap: 12px;
            align-items: center;
        }

        .endpoint-item:last-child {
            border-bottom: none;
        }

        .method {
            background: rgba(255, 255, 255, 0.08);
            border: 1px solid var(--border);
            border-radius: 3px;
            padding: 4px 8px;
            font-weight: 600;
            min-width: 50px;
            text-align: center;
            font-size: 10px;
        }

        .path {
            font-family: 'JetBrains Mono', monospace;
            color: var(--muted);
        }

        .employee-item {
            padding: 6px 0;
            border-bottom: 1px solid rgba(255,255,255,0.05);
        }

        .employee-name {
            color: #ffffff;
            font-weight: 600;
            font-size: 11px;
        }

        .employee-position {
            color: rgba(255,255,255,0.56);
            font-size: 10px;
        }

        @media (max-width: 1400px) {
            .content {
                grid-template-columns: 1fr 1fr;
            }
        }

        @media (max-width: 900px) {
            .content {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="header-title">Organization Structure</div>
            <div class="header-subtitle">REST API для управления структурой организации</div>
        </div>

        <div class="content">
            <div class="card">
                <div class="card-title">Create Department</div>
                <form onsubmit="createDepartment(event)" style="display: flex; flex-direction: column; gap: 8px;">
                    <div class="form-group">
                        <label class="form-label">Department Name *</label>
                        <input type="text" id="deptName" class="form-input" placeholder="e.g., IT" required>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Parent Department ID</label>
                        <input type="number" id="parentId" class="form-input" placeholder="Leave empty if root">
                    </div>
                    <button type="submit" class="btn">Create Department</button>
                    <div class="feedback success" id="deptSuccess">✓ Department created</div>
                    <div class="feedback error" id="deptError"></div>
                </form>
            </div>

            <div class="card">
                <div class="card-title">Create Employee</div>
                <form onsubmit="createEmployee(event)" style="display: flex; flex-direction: column; gap: 8px;">
                    <div class="form-group">
                        <label class="form-label">Department ID *</label>
                        <input type="number" id="empDeptId" class="form-input" placeholder="e.g., 1" required>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Full Name *</label>
                        <input type="text" id="empName" class="form-input" placeholder="John Doe" required>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Position *</label>
                        <input type="text" id="empPosition" class="form-input" placeholder="Developer" required>
                    </div>
                    <button type="submit" class="btn">Add Employee</button>
                    <div class="feedback success" id="empSuccess">✓ Employee added</div>
                    <div class="feedback error" id="empError"></div>
                </form>
            </div>

            <div class="card">
                <div class="card-title">View Employees</div>
                <form onsubmit="getEmployees(event)" style="display: flex; flex-direction: column; gap: 8px;">
                    <div class="form-group">
                        <label class="form-label">Department ID *</label>
                        <input type="number" id="viewDeptId" class="form-input" placeholder="e.g., 1" required>
                    </div>
                    <button type="submit" class="btn">View Employees</button>
                    <div id="employeesList" style="margin-top: 8px; font-size: 11px; max-height: 200px; overflow-y: auto;"></div>
                    <div class="feedback error" id="viewError"></div>
                </form>
            </div>
        </div>

        <div class="endpoints">
            <div class="endpoints-title">Available Endpoints</div>
            <div class="endpoint-item">
                <span class="method">POST</span>
                <span class="path">/departments/</span>
            </div>
            <div class="endpoint-item">
                <span class="method">POST</span>
                <span class="path">/departments/{id}/employees/</span>
            </div>
            <div class="endpoint-item">
                <span class="method">GET</span>
                <span class="path">/departments/{id}/employees/</span>
            </div>
        </div>
    </div>

    <script>
        function createDepartment(e) {
            e.preventDefault();
            var name = document.getElementById('deptName').value;
            var parentId = document.getElementById('parentId').value;

            fetch('/departments/', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: name,
                    parent_id: parentId ? parseInt(parentId) : null
                })
            }).then(function(response) {
                if (response.ok) {
                    document.getElementById('deptSuccess').style.display = 'block';
                    document.getElementById('deptError').style.display = 'none';
                    document.getElementById('deptName').value = '';
                    document.getElementById('parentId').value = '';
                    setTimeout(function() {
                        document.getElementById('deptSuccess').style.display = 'none';
                    }, 3000);
                } else {
                    return response.text().then(function(text) {
                        throw new Error(text || 'Failed to create department');
                    });
                }
            }).catch(function(err) {
                document.getElementById('deptError').textContent = '✗ ' + err.message;
                document.getElementById('deptError').style.display = 'block';
                setTimeout(function() {
                    document.getElementById('deptError').style.display = 'none';
                }, 5000);
            });
        }

        function createEmployee(e) {
            e.preventDefault();
            var deptId = document.getElementById('empDeptId').value;
            var name = document.getElementById('empName').value;
            var position = document.getElementById('empPosition').value;

            fetch('/departments/' + deptId + '/employees/', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    full_name: name,
                    position: position
                })
            }).then(function(response) {
                if (response.ok) {
                    document.getElementById('empSuccess').style.display = 'block';
                    document.getElementById('empError').style.display = 'none';
                    document.getElementById('empName').value = '';
                    document.getElementById('empPosition').value = '';
                    document.getElementById('empDeptId').value = '';
                    setTimeout(function() {
                        document.getElementById('empSuccess').style.display = 'none';
                    }, 3000);
                } else {
                    return response.text().then(function(text) {
                        throw new Error(text || 'Failed to add employee');
                    });
                }
            }).catch(function(err) {
                document.getElementById('empError').textContent = '✗ ' + err.message;
                document.getElementById('empError').style.display = 'block';
                setTimeout(function() {
                    document.getElementById('empError').style.display = 'none';
                }, 5000);
            });
        }

        function getEmployees(e) {
            e.preventDefault();
            var deptId = document.getElementById('viewDeptId').value;
            var listDiv = document.getElementById('employeesList');
            var errorDiv = document.getElementById('viewError');

            fetch('/departments/' + deptId + '/employees/', {
                method: 'GET'
            }).then(function(response) {
                if (response.ok) {
                    return response.json();
                } else {
                    return response.text().then(function(text) {
                        throw new Error(text || 'Failed to fetch employees');
                    });
                }
            }).then(function(data) {
                errorDiv.style.display = 'none';

                if (data.employees && data.employees.length > 0) {
                    var html = '<div style="border-top: 1px solid rgba(255,255,255,0.1); padding-top: 8px;">';
                    for (var i = 0; i < data.employees.length; i++) {
                        var emp = data.employees[i];
                        html += '<div class="employee-item">';
                        html += '<div class="employee-name">' + emp.full_name + '</div>';
                        html += '<div class="employee-position">' + emp.position + '</div>';
                        html += '</div>';
                    }
                    html += '</div>';
                    listDiv.innerHTML = html;
                } else {
                    listDiv.innerHTML = '<div style="color: rgba(255,255,255,0.56); padding: 8px 0;">No employees found</div>';
                }
            }).catch(function(err) {
                errorDiv.textContent = '✗ ' + err.message;
                errorDiv.style.display = 'block';
                listDiv.innerHTML = '';
            });
        }
    </script>
</body>
</html>`
}

func HttpServerStart() {
	cfg := config.NewGormDb()
	gormDbRepo := cfg.MockDB
	OL := domain.NewOrganizationLogic(gormDbRepo)
	handler := NewDepartmentHandler(OL)

	router := SetupRoutes(handler)

	port := ":8080"
	fmt.Println("🚀 Server starting on http://localhost:8080")
	log.Printf("Listening on %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
