package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/begenov/courses-service/internal/domain"
	"github.com/begenov/courses-service/internal/repository"
	"github.com/begenov/courses-service/pkg/cache"
)

type CoursesService struct {
	repo     repository.Courses
	cache    cache.Cache
	ttlCache time.Duration
}

func NewCoursesService(repo repository.Courses, cache cache.Cache, ttlCache time.Duration) *CoursesService {
	return &CoursesService{
		repo:     repo,
		cache:    cache,
		ttlCache: ttlCache,
	}
}

func (s *CoursesService) Create(ctx context.Context, course domain.Courses) error {
	return s.repo.Create(ctx, course)
}

func (s *CoursesService) GetByID(ctx context.Context, id int) (domain.Courses, error) {
	cachedCourses, err := s.cache.Get(ctx, "course:"+strconv.Itoa(id))

	if err == nil {
		cachedData, ok := cachedCourses.(domain.Courses)
		if ok {
			return cachedData, nil
		}
	}

	student, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return student, err
	}

	if err := s.cache.Set(ctx, "course:"+strconv.Itoa(id), student, s.ttlCache); err != nil {
		log.Printf("error caching course with ID %d:", err)
	}

	return student, nil
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
