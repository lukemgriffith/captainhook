package configparser

import (
	"errors"
	"github.com/lukemgriffith/captainhook"
)

type EndpointService struct {
	Config *Config
}

func (e *EndpointService) Endpoint(name string) (*captainhook.Endpoint, error) {

	if len(*e.Config.GetEndpoints()) == 0 {
		return nil, errors.New("No endpoints configured")
	}

	for _, endpoint := range *e.Config.GetEndpoints() {
		if endpoint.Name == name {
			return &endpoint, nil
		}
	}

	return nil, errors.New("Unable to find endpoint by name")
}

func (e *EndpointService) Endpoints() (*[]captainhook.Endpoint, error) {
	return e.Config.GetEndpoints(), nil
}

func (e *EndpointService) CreateEndpoint() error {
	return errors.New("Unable to create endpoint, defined from static config")
}

func (e *EndpointService) DeleteEndpoint() error {
	return errors.New("Unable to create endpoint, defined from static config")
}
