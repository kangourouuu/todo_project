# ✅ Todo App - Task Management API

A lightweight yet production-ready RESTful API built in Go for managing tasks with full CRUD support. Designed with scalability in mind and integrated with PostgreSQL, Redis, Docker, and Swagger for API documentation.

📦 Ideal for learning and bootstrapping real-world backend systems.

---

## 📚 Table of Contents

- 🎯 Overview  
- ⚙️ Features  
- 🔧 Requirements  
- 🚀 Installation & Run  
- 🐳 Docker Workflow  
- 📘 API Documentation (Swagger)  
- 📁 Project Structure

---

## 🎯 Overview

**Todo App** is a backend service that provides CRUD functionality for managing tasks (`name`, `description`). Built using the Go programming language and the Gin framework, the project integrates Docker, PostgreSQL, Redis, and Swagger to support a full-stack backend experience.

---

## ⚙️ Features

- 🛠️ CRUD operations for managing todo tasks  
- 🧱 PostgreSQL integration for data persistence  
- ⚡ Redis caching layer  
- 🐳 Docker and Docker Compose support  
- 📄 Swagger for API documentation  
- ✅ Clean project structure with separation of concerns

---

## 🔧 Requirements

- Go 1.20+  
- Docker & Docker Compose  
- Swagger CLI (`swag`)  
- Git

---

## 🚀 Installation & Run

### ▶️ Run locally without Docker
```bash
# Clone the repo
git clone https://github.com/kangourouuu/todo_project.git
cd todo_project

# Install all Go dependencies
go mod tidy

# Run the project
go run main.go
```
## 🐳 Docker Workflow
### 🔨 Basic Docker commands:
```bash
# Build and start containers
docker-compose up --build

# Stop containers (recommended when stable)
docker-compose stop

# Start stopped containers again
docker-compose start

# Remove containers + network + volumes
docker-compose down -v
```
## 📘 Swagger API Docs
```bash
📦 Install Swagger CLI:

go install github.com/swaggo/swag/cmd/swag@latest

🔄 Upgrade to latest version to fix LeftDelim/RightDelim error:

go get github.com/swaggo/swag@latest

🚀 Generate Swagger docs:

swag init

Swagger UI will be available at: http://localhost:<your-port>/swagger/index.html
```
## 📁 Project Structure
```bash
todo-app/
├── main.go                 # Entry point
├── Dockerfile              # Service container definition
├── docker-compose.yml      # Service orchestration
├── api/                    # Route handlers
│   └── v2/                 # CRUD endpoints for Todo
├── configs/                # Configuration files
├── model/                  # Todo model definitions
├── repository/             # Data access layer
├── service/                # Business logic
├── docs/                   # Swagger-generated docs
├── internal/
│   ├── db/                 # PostgreSQL connection setup
│   └── redis/              # Redis setup
└── common/                 # Utility functions (logger, response, etc.)
```
## 🚀 API Endpoints
```bash
| Method | Endpoint            | Description          |
| ------ | ------------------- | -------------------- |
| GET    | `/api/v2/todo`     | Get all todos        |
| GET    | `/api/v2/todo/:id` | Get todo by ID       |
| POST   | `/api/v2/todo`     | Create a new todo    |
| PUT    | `/api/v2/todo/:id` | Update existing todo |
| DELETE | `/api/v2/todo/:id` | Delete a todo        |
```
---