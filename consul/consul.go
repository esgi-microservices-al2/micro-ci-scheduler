package consul

import (
	"log"
	"strconv"
	"time"

	"github.com/System-Glitch/goyave/v2"
	"github.com/hashicorp/consul/api"
)

type AuthenticationCredentials struct {
	Token string
	Host  string
}

type Service struct {
	ID        string
	Name      string
	Tags      []string
	CheckTime time.Duration
}

const (
	id   string = "al2.scheduler"
	name string = "al2.scheduler"
)

var (
	client      *api.Client
	clientAgent *api.Agent
	service     Service
)

func SetConfiguration(credential AuthenticationCredentials) {
	// Set config
	var configuration = api.DefaultConfig()
	configuration.Address = credential.Host
	configuration.Token = credential.Token

	// Create client
	client, error := api.NewClient(configuration)

	if error != nil {
		goyave.ErrLogger.Println(error, "Failed to instanciate new consul Client")
		return
	}

	// Set service data
	service.ID = id
	service.Name = name
	service.Tags = []string{"al2.scheduler"}
	service.CheckTime = time.Second * 10
}

func Register() {
	if client == nil {
		goyave.ErrLogger.Println("Failed to unregister the service, no client instanciate")
		return
	}

	var serviceRegisterConfiguration = &api.AgentServiceRegistration{
		ID:   service.Name,
		Name: service.Name,
		Check: &api.AgentServiceCheck{
			TTL: strconv.FormatFloat(service.CheckTime.Seconds(), 'f', -1, 64),
		},
	}

	clientAgent = client.Agent()

	var error = clientAgent.ServiceRegister(serviceRegisterConfiguration)

	if error != nil {
		goyave.ErrLogger.Println(error, "Failed cannot register service")
		return
	}

	service.updateTTL()
}

func UnRegiter() {
	if client == nil {
		goyave.ErrLogger.Println("Failed to unregister the service, no client instanciate")
		return
	}

	if clientAgent == nil {
		goyave.ErrLogger.Println("Failed to unregister the service, no agent instanciate")
		return
	}

	var error = clientAgent.ServiceDeregister(service.ID)

	if error != nil {
		goyave.ErrLogger.Println(error, "Failed to deregister service")
	}
}

// update TTL with clock
func (s *Service) updateTTL() {
	ticker := time.NewTicker(s.CheckTime / 2)
	for range ticker.C {
		s.update()
	}
}

// update TTL
func (s *Service) update() {
	if agentErr := clientAgent.UpdateTTL("service:"+s.ID, "", "pass"); agentErr != nil {
		log.Print(agentErr)

	}
}
