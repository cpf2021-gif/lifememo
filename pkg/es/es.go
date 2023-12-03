package es

import (
	es8 "github.com/elastic/go-elasticsearch/v8"
)

type (
	Config struct {
		Addresses  []string
		Username   string
		Password   string
		MaxRetries int
	}

	Es struct {
		*es8.Client
	}
)

func NewEs(conf *Config) (*Es, error) {
	c := es8.Config{
		Addresses:  conf.Addresses,
		Username:   conf.Username,
		Password:   conf.Password,
		MaxRetries: conf.MaxRetries,
	}

	client, err := es8.NewClient(c)
	if err != nil {
		return nil, err
	}

	return &Es{client}, nil
}

func MustNewEs(conf *Config) *Es {
	es, err := NewEs(conf)
	if err != nil {
		panic(err)
	}

	return es
}
