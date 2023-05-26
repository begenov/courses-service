package service

import (
	"context"
	"fmt"
	"time"

	"github.com/begenov/courses-service/internal/domain"
	"github.com/begenov/courses-service/internal/repository"
)

type CoursesService struct {
	repo repository.Courses
}

func NewCoursesService(repo repository.Courses) *CoursesService {
	return &CoursesService{
		repo: repo,
	}
}

func (s *CoursesService) Create(ctx context.Context, course domain.Courses) error {
	return s.repo.Create(ctx, course)
}

func (s *CoursesService) GetByID(ctx context.Context, id int) (domain.Courses, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CoursesService) Update(ctx context.Context, course domain.Courses) error {
	c, err := s.repo.GetByID(ctx, course.ID)
	if err != nil {
		return err
	}

	if course.Name == "" {
		course.Name = c.Name
	}

	if len(course.Students) == 0 {
		course.Students = c.Students
	}

	if course.Description == "" {
		course.Description = c.Description
	}

	course.UpdatedAt = time.Now()

	fmt.Println(course, "service")
	return s.repo.Update(ctx, course)
}

func (s *CoursesService) Delete(ctx context.Context, id int) error {
	return s.repo.Delele(ctx, id)
}

func (s *CoursesService) GetCoursesByIdStudent(ctx context.Context, studentId string) ([]domain.Courses, error) {
	return s.repo.GetCoursesByIdStudent(ctx, studentId)
}
