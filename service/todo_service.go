package service

import (
	"todo_project/model"
	"todo_project/repository"
)

type TodoService interface {
	CreateTodo(todo *model.Todo) error
	GetTodoByID(id uint) (*model.Todo, error)
	GetAllTodos() ([]*model.Todo, error)
	UpdateTodo(todo *model.Todo) error
	DeleteTodo(id uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(todo *model.Todo) error {
	return s.repo.Create(todo)
}

func (s *todoService) GetTodoByID(id uint) (*model.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *todoService) GetAllTodos() ([]*model.Todo, error) {
	return s.repo.FindAll()
}

func (s *todoService) UpdateTodo(todo *model.Todo) error {
	return s.repo.Update(todo)
}

func (s *todoService) DeleteTodo(id uint) error {
	return s.repo.Delete(id)
}