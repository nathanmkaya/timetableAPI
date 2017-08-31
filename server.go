package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nathanmkaya/timetableAPI/controllers"
)

func main() {
	router := gin.Default()
	examController := controllers.NewExamController()

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, fmt.Sprintf("everything is okay"))
	})

	exams := router.Group("/exams")
	{
		exams.POST("/upload", examController.UploadExams)
		exams.GET("/:shift", examController.GetExam)
	}

	classes := router.Group("/classes")
	{
		classes.POST("/upload")
		classes.GET("/:shift")
	}

	router.Run(":3000")
}
