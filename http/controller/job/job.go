package job

import (
	"micro-ci-scheduler/database/model"
	"net/http"
	"strconv"

	"github.com/System-Glitch/goyave/v2"
	"github.com/System-Glitch/goyave/v2/database"
)

func Index(response *goyave.Response, request *goyave.Request) {
	jobs := []model.Job{}
	result := database.GetConnection().Find(&jobs)
	if response.HandleDatabaseError(result) {
		response.JSON(http.StatusOK, jobs)
	}
}

func Show(response *goyave.Response, request *goyave.Request) {
	job := model.Job{}
	id, _ := strconv.ParseUint(request.Params["id"], 10, 64)
	result := database.GetConnection().First(&job, id)
	if response.HandleDatabaseError(result) {
		response.JSON(http.StatusOK, job)
	}
}

func Store(response *goyave.Response, request *goyave.Request) {
	job := model.Job{
		Name:           request.String("name"),
		CronExpression: request.String("cronexpression"),
		IdProject:      request.Integer("idproject"),
	}
	if err := database.GetConnection().Create(&job).Error; err != nil {
		response.Error(err)
	} else {
		response.JSON(http.StatusCreated, map[string]uint{"id": job.ID})
	}
}

func Update(response *goyave.Response, request *goyave.Request) {
	id, _ := strconv.ParseUint(request.Params["id"], 10, 64)
	job := model.Job{}
	db := database.GetConnection()
	result := db.Select("id").First(&job, id)
	if response.HandleDatabaseError(result) {
		if err := db.Model(&job).Update("name", request.String("name")).Error; err != nil {
			response.Error(err)
		}
	}
}

func Destroy(response *goyave.Response, request *goyave.Request) {
	id, _ := strconv.ParseUint(request.Params["id"], 10, 64)
	job := model.Job{}
	db := database.GetConnection()
	result := db.Select("id").First(&job, id)
	if response.HandleDatabaseError(result) {
		if err := db.Delete(&job).Error; err != nil {
			response.Error(err)
		}
	}
}
