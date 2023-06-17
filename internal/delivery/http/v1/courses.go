package v1

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/begenov/courses-service/internal/domain"
	"github.com/begenov/courses-service/pkg/course/api/courses"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) initCoursesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	{
		courses.GET("/:id", h.getCourseByID)
		courses.POST("/create", h.createCourses)
		courses.PUT("/:id/update", h.updateCourses)
		courses.DELETE("/:id/delete", h.deleteCourse)
		courses.GET("/:id/courses", h.getCoursesByStudentID)
		courses.GET("/:id/students", h.getStudentsByCoursId)

	}
}

// type createCourses struct {
// 	Name        string   `json:"name" binding:"required"`
// 	Description string   `json:"description" binding:"required"`
// 	Students    []string `json:"students"  binding:"required"`
// }

// @Summary		Create New Courses
// @Tags			Courses
// @Description	 Create New Courses
// @Accept			json
// @Produce		json
// @Param			account	body		createCourses	true	"Courses"
// @Success		201		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/courses/create [post]
func (h *Handler) createCourses(ctx *gin.Context) {
	var inp courses.Courses

	data, _ := io.ReadAll(ctx.Request.Body)

	err := proto.Unmarshal(data, &inp)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	/*
		if err := ctx.BindJSON(&inp); err != nil {
			newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
			return
		}
	*/
	if err := h.service.Courses.Create(context.Background(), domain.Courses{
		Description: inp.Description,
		Name:        inp.Name,
		Students:    inp.Students,
		CreatedAt:   time.Now(),
	}); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Error when creating a courses")
		return
	}

	ctx.JSON(http.StatusCreated, Resposne{"The course is successfully established"})
}

// @Summary		Get Course By ID
// @Tags			Courses
// @Description	 Create New Courses
// @ModuleID getCourseByID
// @Accept			json
// @Produce		json
// @Param			id path string	true	"course id"
// @Success		200		{object}	domain.Courses
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/courses/{id} [get]
func (h *Handler) getCourseByID(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect course ID format")
		return
	}
	course, err := h.service.Courses.GetByID(context.Background(), id)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Mistake in getting a course")
		return
	}

	fmt.Println(course)

	ctx.JSON(http.StatusOK, course)
}

// type inputCourse struct {
// 	Name        string   `json:"name"`
// 	Description string   `json:"description"`
// 	Students    []string `json:"students"`
// }

// @Summary		Update Course
// @Tags			Courses
// @Description	 Update Course
// @Accept			json
// @Produce		json
// @Param			account	body		inputCourse	true	"course update info"
// @Param			id path string		true	"course id"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/courses/{id}/update [put]
func (h *Handler) updateCourses(ctx *gin.Context) {
	var inp courses.Courses
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect course ID format")
		return
	}

	data, _ := io.ReadAll(ctx.Request.Body)

	if err := proto.Unmarshal(data, &inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	if err := h.service.Courses.Update(context.Background(), domain.Courses{
		Name:        inp.Name,
		Description: inp.Description,
		Students:    inp.Students,
		ID:          id,
	}); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Error when updating a course")
		return
	}

	ctx.JSON(http.StatusOK, Resposne{"Course has been successfully updated"})
}

// @Summary		 Delete Course
// @Tags			Courses
// @Description	 Delete Course
// @Accept			json
// @Produce		json
// @Param			id path string		true	"course id"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/courses/{id}/delete [delete]
func (h *Handler) deleteCourse(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect course ID format")
		return
	}

	if err := h.service.Courses.Delete(context.Background(), id); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Error in getting course by course ID")
		return
	}

	ctx.JSON(http.StatusOK, Resposne{"Student successfully deleted"})
}

// @Summary	  Get Courses By StudentID
// @Tags			Courses
// @Description	 Get Courses By StudentID
// @Accept			json
// @Produce		json
// @Param			id path string		true	"course id"
// @Success		200		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/courses/{id}/courses [get]
func (h *Handler) getCoursesByStudentID(ctx *gin.Context) {
	id := ctx.Param("id")

	courses, err := h.service.Courses.GetCoursesByIdStudent(ctx, id)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Error in getting course by student ID")
		return
	}

	ctx.JSON(http.StatusOK, courses)
}

// @Summary		Get Students By CoursId
// @Tags			Courses
// @Description	 Get Students By CoursId
// @Accept			json
// @Produce		json
// @Param			id path string		true	"course id"
// @Success		200		{object}	string
// @Failure		500		{object}	Resposne
// @Router			/courses/{id}/students [get]
func (h *Handler) getStudentsByCoursId(ctx *gin.Context) {
	param := ctx.Param("id")

	err := h.service.Kafka.SendMessages("students-request", param)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to get information about students")
		return
	}

	responseData := <-h.responseCh
	ctx.Data(http.StatusOK, "application/json", responseData)
}

func (h *Handler) consumeResponseMessages() {
	err := h.service.Kafka.ConsumeMessages("students-response", h.handleResponseMessage)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) handleResponseMessage(message string) {
	h.responseCh <- []byte(message)
}

/*
var (
	api = "http://localhost:8000/api/v1/students/"
)

func (h *Handler) getStudentsByCoursId(ctx *gin.Context) {
	param := ctx.Param("id")

	url := api + param + "/students"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get information about students" + err.Error(),
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
*/
