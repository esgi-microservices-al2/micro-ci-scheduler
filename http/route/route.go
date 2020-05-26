package route

import (
	"micro-ci-scheduler/http/controller/hello"
	"micro-ci-scheduler/http/controller/job"
	"micro-ci-scheduler/http/request/echorequest"
	"micro-ci-scheduler/http/request/jobrequest"

	"github.com/System-Glitch/goyave/v2"
	"github.com/System-Glitch/goyave/v2/cors"
	"github.com/System-Glitch/goyave/v2/log"
)

// Routing is an essential part of any Goyave application.
// Routes definition is the action of associating a URI, sometimes having
// parameters, with a handler which will process the request and respond to it.

// Routes are defined in routes registrer functions.
// The main route registrer is passed to "goyave.Start()" and is executed
// automatically with a newly created root-level router.

// Register all the application routes. This is the main route registrer.
func Register(router *goyave.Router) {

	// Applying default CORS settings (allow all methods and all origins)
	// Learn more about CORS options here: https://system-glitch.github.io/goyave/guide/advanced/cors.html
	router.CORS(cors.Default())
	router.Middleware(log.CombinedLogMiddleware())

	// Register your routes here

	// Route without validation
	router.Route("GET", "/hello", hello.SayHi, nil)

	// Route with validation
	router.Route("GET", "/echo", hello.Echo, echorequest.Echo)

	router.Get("/job", job.Index, nil)
	router.Get("/job/{id:[0-9]+}", job.Show, nil)
	router.Post("/job", job.Store, jobrequest.Store)
	router.Put("/job/{id:[0-9]+}", job.Update, jobrequest.Store)
	router.Delete("/job/{id:[0-9]+}", job.Destroy, nil)
}
