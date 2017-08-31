package controllers

import (
	"fmt"
	"net/http"

	"io"
	"os"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nathanmkaya/timetableAPI/models"
	"github.com/nathanmkaya/timetableAPI/parser"
)

type ExamController struct{}

func NewExamController() *ExamController {
	return &ExamController{}
}

func (ec ExamController) GetExam(c *gin.Context) {
	shift := c.Params.ByName("shift")
	name := c.Query("name")
	names, ok := c.GetQueryArray("name")
	if ok {
		courses, err := models.GetCourses(shift, names)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("%s", err.Error()))
			return
		}
		c.JSON(http.StatusOK, courses)
	} else {
		course, err := models.GetCourse(shift, name)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("%s", err.Error()))
		}
		c.JSON(http.StatusOK, course)
	}
}

func (ec ExamController) UploadExams(c *gin.Context) {
	file, err := c.FormFile("exam")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	src, err := file.Open()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file open err: %s", err.Error()))
		return
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("exam.xlsx")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Create file err: %s", err.Error()))
		return
	}
	defer dst.Close()

	// Copy
	io.Copy(dst, src)

	filename, err := filepath.Abs(dst.Name())
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Create file err: %s", err.Error()))
		return
	}

	parser.ParseExams(filename)
}
