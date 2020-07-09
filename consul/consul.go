package consul

import (
	"fmt"
	"log"
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
	id   string = "al2-scheduler"
	name string = "al2-scheduler"
)

var (
	client      *api.Client
	clientAgent *api.Agent
	service     Service
)

func SetConfiguration(credential AuthenticationCredentials) {
	// Set config
	configuration := api.DefaultConfig()
	configuration.Address = credential.Host
	configuration.Token = credential.Token

	// Create client
	var err error
	client, err = api.NewClient(configuration)

	if err != nil {
		goyave.ErrLogger.Println(err, "Failed to instantiate new consul Client")
		return
	}

	// Set service data
	service.ID = id
	service.Name = name
	service.Tags = []string{id, "traefik.enable=true", "traefik.frontend.entryPoints=http", "traefik.frontend.rule=PathPrefix:/scheduler/"}
	service.CheckTime = time.Second * 30
}

func Start() {
	if client == nil {
		goyave.ErrLogger.Println("Failed to unregister the service, no client instantiate")
		return
	}

	var serviceRegisterConfiguration = &api.AgentServiceRegistration{
		ID:   service.Name,
		Name: service.Name,
		Check: &api.AgentServiceCheck{
			TTL: fmt.Sprintf("%fs", service.CheckTime.Seconds()),
		},
	}

	clientAgent = client.Agent()

	var err = clientAgent.ServiceRegister(serviceRegisterConfiguration)

	if err != nil {
		goyave.ErrLogger.Println(err, "Failed cannot register consul service")
		return
	}

	service.updateTTL()
}

func Stop() {
	if client == nil {
		goyave.ErrLogger.Println("Failed to unregister the service, no client instantiate")
		return
	}

	if clientAgent == nil {
		goyave.ErrLogger.Println("Failed to unregister the service, no agent instantiate")
		return
	}

	var err = clientAgent.ServiceDeregister(service.ID)

	if err != nil {
		goyave.ErrLogger.Println(err, "Failed to deregister service")
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
