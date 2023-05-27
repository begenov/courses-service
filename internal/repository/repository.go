package repository

import (
	"context"
	"database/sql"

	"github.com/begenov/courses-service/internal/domain"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Courses interface {
	Create(ctx context.Context, course domain.Courses) error
	GetByID(ctx context.Context, id int) (domain.Courses, error)
	Update(ctx context.Context, course domain.Courses) error
	Delele(ctx context.Context, id int) error
	GetCoursesByIdStudent(ctx context.Context, studentId string) ([]domain.Courses, error)
}

type Repository struct {
	Courses Courses
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Courses: NewCoursesRepo(db),
	}
}
