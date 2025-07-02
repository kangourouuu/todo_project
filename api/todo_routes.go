package api

import (
	v2 "todo_project/api/v2"
	"todo_project/repository"
	"todo_project/service"
	"todo_project/internal/redis"
	"todo_project/internal"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.RouterGroup, redisClient redis.IRedis) {
	todoRepo := repository.NewTodoRepository(internal.GormSqlClient.GetDB())
	todoService := service.NewTodoService(todoRepo)
	todoHandler := v2.NewTodoHandler(todoService, redisClient)
	r.POST("/todo", todoHandler.CreateTodo)
	r.GET("/todo/:id", todoHandler.GetTodo)
	r.GET("/todo", todoHandler.GetAllTodos)
	r.PUT("/todo/:id", todoHandler.UpdateTodo)
	r.DELETE("/todo/:id", todoHandler.DeleteTodo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/openapi.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

}

//@Summary Lấy tất cả todo
//@Description Trả về thông tin danh sách todo
//@Tags todo
//@Produce json 
//@Param id path int true "todo ID" 
//@Success 200 {object} model.Todo 
//@Router /todo [get] 
func GetAllTodos(c *gin.Context) {}

// @Summary Lấy todo theo ID
// @Description Trả về thông tin todo theo ID
// @Tags todo
// @Produce json
// @Param id path int true "todo ID"
// @Success 200 {object} model.Todo
// @Router /todo/{id} [get]
func GetTodo(c *gin.Context) {}

// @Summary Thêm mới todo
// @Description Thêm một đôi todo mới vào hệ thống
// @Tags todo
// @Accept json
// @Produce json
// @Param data body model.Todo true "todo info"
// @Success 201 {object} model.Todo
// @Router /todo [post]
func CreateTodo(c *gin.Context) {}

// @Summary Xoá todo
// @Description Xoá một đôi todo theo ID
// @Tags todo
// @Produce json
// @Param id path int true "todo ID"
// @Success 204
// @Router /todo/{id} [delete]
func DeleteTodo(c *gin.Context) {}

// @Summary Cập nhật todo
// @Description Cập nhật thông tin một đôi todo
// @Tags todo
// @Accept json
// @Produce json
// @Param id path int true "todo ID"
// @Param data body model.Todo true "todo info"
// @Success 200 {object} model.Todo
// @Router /todo/{id} [put]
func UpdateTodo(c *gin.Context) {}