package service

import (
	"context"
	"testing"
	"time"

	"github.com/begenov/courses-service/internal/domain"
	"github.com/begenov/courses-service/pkg/cache"

	repoMock "github.com/begenov/courses-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCoursesService_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockCourses(ctl)

	service := NewCoursesService(repo, &cache.MemoryCache{}, 10*time.Minute)

	ctx := context.Background()

	course := domain.Courses{
		Name:        "test",
		Description: "test",
		CreatedAt:   time.Now(),
		Students:    []string{"1"},
	}

	repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := service.Create(ctx, course)
	require.NoError(t, err)
}

func TestCoursesService_CreateError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockCourses(ctl)

	service := NewCoursesService(repo, &cache.MemoryCache{}, 10*time.Minute)

	ctx := context.Background()

	course := domain.Courses{
		Name:        "test",
		Description: "test",
	}

	repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := service.Create(ctx, course)
	require.Nil(t, err)
}

func TestCoursesService_Delete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockCourses(ctl)

	service := NewCoursesService(repo, &cache.MemoryCache{}, 10*time.Minute)

	ctx := context.Background()
	repo.EXPECT().Delele(ctx, gomock.Any()).Return(nil)

	err := service.Delete(ctx, 1)
	require.NoError(t, err)
}

func TestCoursesService_GetCoursesByIdStudent(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockCourses(ctl)

	service := NewCoursesService(repo, &cache.MemoryCache{}, 10*time.Minute)

	ctx := context.Background()

	id := "1"

	courses := []domain.Courses{
		{ID: 1, Name: "Course 1"},
		{ID: 2, Name: "Course 2"},
	}

	repo.EXPECT().GetCoursesByIdStudent(ctx, gomock.Any()).Return(courses, nil)

	result, err := service.GetCoursesByIdStudent(ctx, id)
	require.NoError(t, err)
	require.Equal(t, courses, result)
}

func TestCoursesService_Update(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockCourses(ctl)

	service := NewCoursesService(repo, &cache.MemoryCache{}, 10*time.Minute)

	ctx := context.Background()
	course := domain.Courses{
		ID:   1,
		Name: "test",
	}

	couseRes := domain.Courses{
		ID:          1,
		Name:        "test",
		Description: "test",
		CreatedAt:   time.Now(),
		Students:    []string{"1"},
	}

	repo.EXPECT().GetByID(ctx, gomock.Any()).Return(couseRes, nil)
	repo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	service.Update(ctx, course)
}
