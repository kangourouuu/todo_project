package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"todo_project/dto"
	"todo_project/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock test service
type mockTodoService struct {
	mock.Mock
}

func (m *mockTodoService) CreateTodo(t *model.Todo) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *mockTodoService) GetTodoByID(id uint) (*model.Todo, error) {
	args := m.Called(id)
	var result *model.Todo
	if args.Get(0) != nil {
		result = args.Get(0).(*model.Todo)
	}

	return result, args.Error(1)
}

func (m *mockTodoService) GetAllTodos() ([]*model.Todo, error) {
	args := m.Called()
	var result []*model.Todo

	if args.Get(0) != nil {
		result = args.Get(0).([]*model.Todo)
	}

	return result, args.Error(1)
}

func (m *mockTodoService) UpdateTodo(t *model.Todo) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *mockTodoService) DeleteTodo(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock redis client
type mockRedisClient struct {
	mock.Mock
}

func (m *mockRedisClient) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockRedisClient) Set(key string, value interface{}) (string, error) {
	args := m.Called(key, value)
	return args.String(0), args.Error(1)
}

func (m *mockRedisClient) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *mockRedisClient) Delete(key string) (int64, error) {
	args := m.Called(key)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockRedisClient) GetClient() *redis.Client {
	args := m.Called()
	return args.Get(0).(*redis.Client)
}

func (m *mockRedisClient) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// test-case create
func TestCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockTodoService)
	mockRc := new(mockRedisClient)
	handler := &TodoHandler{
		todoService: mockSvc,
		redisClient: mockRc,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := dto.CreateTodoRequest{
			// example, replace your fields in your dto
			//Name:  "Khanh",
			//Price: 2.55,
		}

		jsonData, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		// mock test service
		mockSvc.On("CreateTodo", mock.AnythingOfType("*model.Todo")).Return(nil)

		// mock test redis
		mockRc.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("OK", nil)

		r := gin.Default()
		r.POST("/test", handler.CreateTodo)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
		mockRc.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {

		jsonData := []byte(`{name : "Khanh", "price":2.55}`)

		req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		// mock test service
		mockSvc.On("CreateTodo", mock.AnythingOfType("*model.Todo")).Return(nil)

		// mock test redis
		mockRc.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("interface{}")).Return("OK", nil)

		r := gin.Default()
		r.POST("/test", handler.CreateTodo)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("redisErr", func(t *testing.T) {
		resp := dto.CreateTodoRequest{
			// example, replace your fields in your dto
			//Name:  "Khanh",
			//Price: 2.55,
		}

		jsonData, _ := json.Marshal(resp)

		req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		// mock test service
		mockSvc.On("CreateTodo", mock.AnythingOfType("*model.Todo")).Return(nil)

		// mock test redis
		mockRc.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", errors.New("redis error"))

		r := gin.Default()
		r.POST("/test", handler.CreateTodo)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

// test-case getAll
func TestGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSv := new(mockTodoService)
	mockRc := new(mockRedisClient)
	handler := &TodoHandler{
		todoService: mockSv,
		redisClient: mockRc,
	}
	t.Run("success", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetAllTodo").Return([]*model.Todo{
			// example, replace your fields in your model
			//{ID: 1, Name: "Khanh", Price: 2.55},
		},nil)

		mockRc.On("Get", mock.AnythingOfType("string")).Return("OK", nil)

		r := gin.Default()
		r.GET("/test", handler.GetAllTodos)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetAllTodos").Return([]*model.Todo{},nil)

		mockRc.On("Get", mock.AnythingOfType("string")).Return("OK", nil)

		r := gin.Default()
		r.GET("/test", handler.GetAllTodos)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("id valid", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodGet, "/test/1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetTodoByID", uint(1)).Return(&model.Todo{
			// example, replace your fields in your model
			//{ID: 1, Name: "Khanh", Price: 2.55},
		}, nil)

		mockRc.On("Get", mock.AnythingOfType("string")).Return("OK", nil)

		r := gin.Default()
		r.GET("/test/:id", handler.GetTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("id invalid", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodGet, "/test/abc", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/test/:id", handler.GetTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("id not found", func(t *testing.T) {
		mockSv.On("GetTodoByID", uint(999)).Return(nil, errors.New("not found"))
		mockRc.On("Get", "test_999").Return("", errors.New("cache miss"))

		req, _ := http.NewRequest(http.MethodGet, "/test/999", nil)
		w := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/test/:id", handler.GetTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}


func TestUpdate(t *testing.T) {
	mockSv := new(mockTodoService)
	mockRc := new(mockRedisClient)
	handler := &TodoHandler{
		todoService: mockSv,
		redisClient: mockRc,
	}
	t.Run("success", func(t *testing.T) {
		reqBody := dto.UpdateTodoRequest{
			// example, replace your fields in your dto
			//ID:    1,
			//Name:  "Khanh",
			//Price: 3.55,
		}

		jsonData, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPut, "/test/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetTodoByID", uint(1)).Return(&model.Todo{
			// replace your fields in your model
			//ID: 1, Name: "Khanh", Price: 2.55,
		}, nil)

		mockSv.On("UpdateTodo", mock.AnythingOfType("*model.Todo")).Return(nil)

		r := gin.Default()
		r.PUT("/test/:id", handler.UpdateTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		jsonData := []byte(`{id : 1, "name" : "Giang Khanh", "price" : 3.55}`)

		req, _ := http.NewRequest(http.MethodPut, "/test/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetTodoByID", uint(1)).Return(&model.Todo{
			// example, replace your fields in your model
			//{ID: 1, Name: "Khanh", Price: 2.55},
		}, nil)

		mockSv.On("UpdateTodo", mock.AnythingOfType("*model.Todo")).Return(errors.New("invalid input"))

		r := gin.Default()
		r.PUT("/test/:id", handler.UpdateTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		jsonData := []byte(`{"id" : 1, "name" : "", "price" : 3.55}`)

		req, _ := http.NewRequest(http.MethodPut, "/test/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetTodoByID", uint(1)).Return(&model.Todo{
			// example, replace your fields in your model
			//{ID: 1, Name: "Khanh", Price: 2.55},
		}, nil)

		mockSv.On("UpdateTodo", mock.AnythingOfType("*model.Todo")).Return(errors.New("invalid input"))

		r := gin.Default()
		r.PUT("/test/:id", handler.UpdateTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("id not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/test/2", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetTodoByID", uint(2)).Return(nil, errors.New("id not found"))

		// mockSv.On("DeleteTodo", uint(2)).Return(errors.New("id not found"))

		r := gin.Default()
		r.DELETE("/test/:id", handler.UpdateTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSv := new(mockTodoService)
	mockRc := new(mockRedisClient)
	handler := &TodoHandler{
		todoService: mockSv,
		redisClient: mockRc,
	}
	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/test/1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetTodoByID", uint(1)).Return(&model.Todo{
			// example, replace your fields in your model
			//{ID: 1, Name: "Khanh", Price: 2.55},
		}, nil)

		mockSv.On("DeleteTodo", uint(1)).Return(nil)

		mockRc.On("Delete", mock.AnythingOfType("string")).Return(1, nil)

		r := gin.Default()
		r.DELETE("/test/:id", handler.DeleteTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("id invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/test/abc", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// mockSv.On("GetTodoByID", uint(1)).Return(&model.Todo{
		// 	ID: 1, Name: "Khanh",Price: 2.55,
		// }, nil)

		mockSv.On("DeleteTodo", uint(1)).Return(nil)

		mockRc.On("Delete", mock.AnythingOfType("string")).Return(1, nil)

		r := gin.Default()
		r.DELETE("/test/:id", handler.DeleteTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("id not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/test/2", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSv.On("GetTodoByID", uint(2)).Return(nil, errors.New("id not found"))

		// mockSv.On("DeleteTodo", uint(2)).Return(errors.New("id not found"))

		r := gin.Default()
		r.DELETE("/test/:id", handler.DeleteTodo)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

