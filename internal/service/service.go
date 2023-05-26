package service

import (
	"context"

	"github.com/begenov/courses-service/internal/domain"
	"github.com/begenov/courses-service/internal/repository"
)

type Courses interface {
	Create(ctx context.Context, course domain.Courses) error
	GetByID(ctx context.Context, id int) (domain.Courses, error)
	Update(ctx context.Context, course domain.Courses) error
	Delete(ctx context.Context, id int) error
	GetCoursesByIdStudent(ctx context.Context, studentId string) ([]domain.Courses, error)
}

type Service struct {
	Courses Courses
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Courses: NewCoursesService(repo.Courses),
	}
}
