package svc

import (
	"lifememo/application/moment/mq/internal/config"
	"lifememo/pkg/es"
)

type ServiceContext struct {
	Config config.Config
	Es     *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Es: es.MustNewEs(&es.Config{
			Addresses: c.Es.Addresses,
			Username:  c.Es.Username,
			Password:  c.Es.Password,
		}),
	}
}
