package v1

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/begenov/courses-service/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initCoursesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	{
		courses.GET("/:id", h.getCourseByID)
		courses.POST("/create", h.createCourses)
		courses.PUT("/:id/update", h.updateCourses)
		courses.DELETE("/:id/delete", h.deleteCourse)
		courses.GET("/:id/courses", h.getCoursesByIdStudent)
		courses.GET("/:id/students", h.getStudentsByCoursId)

	}
}

type createCourses struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Students    []string `json:"students"  binding:"required"`
}

func (h *Handler) createCourses(ctx *gin.Context) {
	var inp createCourses

	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect input data format",
		})
		return
	}

	if err := h.service.Courses.Create(context.Background(), domain.Courses{
		Description: inp.Description,
		Name:        inp.Name,
		Students:    inp.Students,
		CreatedAt:   time.Now(),
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error when creating a courses",
		})
		return
	}
	log.Println(inp, "oii")
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "The course is successfully established",
	})
}

func (h *Handler) getCourseByID(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect course ID format",
		})
		return
	}
	course, err := h.service.Courses.GetByID(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Mistake in getting a course",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"course": course,
	})
}

type inputCourse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Students    []string `json:"students"`
}

func (h *Handler) updateCourses(ctx *gin.Context) {
	var inp inputCourse
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect course ID format",
		})
		return
	}
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect input data format",
		})
		return
	}

	if err := h.service.Courses.Update(context.Background(), domain.Courses{
		Name:        inp.Name,
		Description: inp.Description,
		Students:    inp.Students,
		ID:          id,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error when updating a course",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course has been successfully",
	})
}

func (h *Handler) deleteCourse(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect course ID format",
		})
		return
	}

	if err := h.service.Courses.Delete(context.Background(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error in getting course by course ID",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Student successfully deleted",
	})
}

func (h *Handler) getCoursesByIdStudent(ctx *gin.Context) {
	id := ctx.Param("id")

	courses, err := h.service.Courses.GetCoursesByIdStudent(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error in getting course by student ID",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}

var (
	api = "http://localhost:8000/api/v1/students/"
)

func (h *Handler) getStudentsByCoursId(ctx *gin.Context) {
	param := ctx.Param("id")
	url := api + param + "/students"
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get information about students",
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get information about students",
		})
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the body of the answer",
		})
		return
	}

	var students domain.Response
	if err := json.Unmarshal(body, &students); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to decode response body",
		})
		return
	}
	if len(students.Students) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Students not found.",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}
