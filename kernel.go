package main

import (
	"micro-ci-scheduler/cron"
	_ "micro-ci-scheduler/http/request"
	"micro-ci-scheduler/http/route"
	"micro-ci-scheduler/rabbit"

	"github.com/System-Glitch/goyave/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	// This is the entry point of your application.
	// Most applications don't need more than this, but
	// if you are running multiple services, such as a
	// websocket server, you'll need to run Goyave in a routine.
	// See: https://system-glitch.github.io/goyave/guide/advanced/multi-services.html

	goyave.Logger.Println("Starting HTTP server...")
	goyave.RegisterStartupHook(func() {
		rabbit.Connect()
		cron.Start()
		goyave.Logger.Println("Ready.")
	})

	goyave.Start(route.Register)
	cron.Stop()
	rabbit.Stop()
}
