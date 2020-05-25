package main

import (
	_ "micro-ci-scheduler/http/request"
	"micro-ci-scheduler/http/route"

	"github.com/System-Glitch/goyave/v2"
	"fmt"
	// Import the approriate GORM dialect for the database you're using.
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
)

func main() {
	// This is the entry point of your application.
	// Most applications don't need more than this, but
	// if you are running multiple services, such as a
	// websocket server, you'll need to run Goyave in a routine.
	// See: https://system-glitch.github.io/goyave/guide/advanced/multi-services.html

	goyave.RegisterStartupHook(func() {
		fmt.Println("Started")
	})

	goyave.Start(route.Register)
}
