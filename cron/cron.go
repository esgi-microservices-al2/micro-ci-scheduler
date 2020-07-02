package cron

import (
	"encoding/json"
	"io/ioutil"
	"micro-ci-scheduler/database/model"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/System-Glitch/goyave/v2"
	"github.com/System-Glitch/goyave/v2/config"
	"github.com/System-Glitch/goyave/v2/database"
	"github.com/robfig/cron"
)

var c *cron.Cron
var mu sync.Mutex

func Start() {
	goyave.Logger.Println("Starting cron")
	mu.Lock()
	defer mu.Unlock()
	c = cron.New()

	jobs := []model.Job{}
	if err := database.GetConnection().Find(&jobs).Error; err != nil {
		goyave.ErrLogger.Println(err)
		return
	}
	for _, j := range jobs {
		if err := c.AddFunc(j.CronExpression, job(j.IdProject)); err != nil {
			json, _ := json.Marshal(j)
			goyave.ErrLogger.Printf("Cron AddFunc: %s | %s\n", err.Error(), string(json))
		}
	}
	go c.Start()
}

func Stop() {
	goyave.Logger.Println("Stopping cron")
	mu.Lock()
	c.Stop()
	c = nil
	mu.Unlock()
}

func Restart() {
	goyave.Logger.Println("Restarting cron")
	Stop()
	Start()
}

func job(jobID int) func() {
	return func() {
		goyave.Logger.Printf("Execute %d\n", jobID)
		id := strconv.Itoa(jobID)
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := client.PostForm(config.GetString("service-project-host")+"/api/v1/builds/"+id, url.Values{})
		if err != nil {
			goyave.ErrLogger.Println(err, "Job #"+id)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				goyave.ErrLogger.Printf("ERROR Job #%d: status %d\n", jobID, resp.StatusCode)
			} else {
				goyave.ErrLogger.Printf("ERROR Job #%d: status %d, %q\n", jobID, resp.StatusCode, string(body))
			}
		}
	}
}
