package v2

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"

	"todo_project/dto"
	"todo_project/model"
	"todo_project/service"
	"todo_project/internal/redis"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService service.TodoService
	redisClient redis.IRedis
}

func NewTodoHandler(todoService service.TodoService, redisClient redis.IRedis) *TodoHandler {
	return &TodoHandler{ 
	todoService: todoService,
	redisClient: redisClient }
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	obj := &model.Todo{ // ne : not equal
		Name: req.Name, // ne : not equal
		Description: req.Description,
	}

	if err := h.todoService.CreateTodo(obj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	resp := dto.TodoResponse{
		ID: obj.ID,
		Name: obj.Name,
		Description: obj.Description,
	}

	if h.redisClient != nil {
		// convert struct to json
		respJSON, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Can not convert data")
		} else {
			redisKey := "todo_" + strconv.Itoa(resp.ID)
			str, err := h.redisClient.Set(redisKey, string(respJSON))
			if err != nil {
				fmt.Println("Can not set data to redis cache")
			} else {
				fmt.Println("Successfully cached data " + str )
			}
		}
	} else {
		fmt.Println("Redis is not ready")
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	obj, err := h.todoService.GetTodoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	resp := &dto.TodoResponse{
		ID: obj.ID,
		Name: obj.Name,
		Description: obj.Description,
	}
	if h.redisClient != nil {
		_, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Can not convert data")
		} else {
			redisKey := "todo_" + strconv.Itoa(resp.ID)
			str, err := h.redisClient.Get(redisKey)
			if err != nil {
				fmt.Println("No data found on cache")
			} else {
				err := json.Unmarshal([]byte(str), resp)
				if err != nil {
					logrus.Info("Failed to Unmarshal")
				}
				fmt.Println("Data cached in Redis - Key:", redisKey, "Result: ", str)
			}
		}
	} else {
		fmt.Println("Redis is not ready")
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	list, err := h.todoService.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get todos"})
		return
	}

	var resp []dto.TodoResponse
	for _, s := range list {
		resp = append(resp, dto.TodoResponse{
			ID: s.ID,
			Name: s.Name,
			Description: s.Description,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	existingTodo, err := h.todoService.GetTodoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	todo := &model.Todo{
		ID: existingTodo.ID,
		Name: req.Name,
		Description: req.Description,
	}

	if err := h.todoService.UpdateTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	resp := dto.TodoResponse{
		ID: todo.ID,
		Name: todo.Name,
		Description: todo.Description,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = h.todoService.GetTodoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	err = h.todoService.DeleteTodo(uint(id))
	if err != nil {
		if err.Error() == "id not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	if h.redisClient != nil {
		redisKey := "todo_" + strconv.Itoa(id)
		deletedCount, err := h.redisClient.Delete(redisKey)
		if err != nil {
			fmt.Println("❌ Failed to delete from Redis cache:", err.Error())
		} else {
			if deletedCount > 0 {
				fmt.Println("🗑️ Successfully deleted from Redis cache - Key:", redisKey, "Count:", deletedCount)
			} else {
				fmt.Println("ℹ️ Key not found in Redis cache:", redisKey)
			}
		}
	} else {
		fmt.Println("⚠️ Redis client is not available - cache not deleted")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
