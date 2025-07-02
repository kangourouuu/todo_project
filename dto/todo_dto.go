package dto

type CreateTodoRequest struct {
	Name string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateTodoRequest struct {
	ID int `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type TodoResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}
